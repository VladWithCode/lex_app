package accupdter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/fetchers"
	"github.com/vladwithcode/lex_app/internal/readers"
)

var (
	ErrNoCaseKeys  = errors.New("no caseKeys were supplied")
	ErrFatalSearch = errors.New("search for case updates failed")
	ErrNoUpdates   = errors.New("found no updates for the provided parameters")
	ErrNilStore    = errors.New("the configured store is nil. But a store dependant method was called")
	ErrFailSave    = errors.New("failed to save accords")
)

type CaseTypesMap map[internal.CaseType][]string

type GenUpdterConf struct {
	// Refers to the unique id used in the provided store
	// to identify individual cases (independant of CaseId and CaseType)
	Store           CaseStore
	Region          internal.Region
	MaxSearchBack   int
	SearchStartDate time.Time
	FetchFn         func(time.Time, internal.CaseType) (*[]byte, error)
	ReadFn          func(*[]byte) (*readers.CaseTable, error)

	ctx context.Context
	db  *sql.DB
}

type GeneralUpdater struct {
	conf *GenUpdterConf
}

func NewGeneralUpdater(conf *GenUpdterConf) *GeneralUpdater {
	if conf.Region == "" {
		conf.Region = internal.RegionDefault
	}
	if conf.Store == nil {
		conf.Store = NewDefaultCaseStore(conf.ctx, conf.db)
	}

	if conf.FetchFn == nil {
		conf.FetchFn = fetchers.NewFetcher(conf.Region)
	}
	if conf.ReadFn == nil {
		conf.ReadFn = readers.NewReader(conf.Region)
	}
	if conf.ctx == nil {
		conf.ctx = context.Background()
	}

	return &GeneralUpdater{conf}
}

func (updter *GeneralUpdater) Update(
	caseKeys []string,
	startSearchDate time.Time,
	maxSearchBack int,
	exhaustSearch bool,
) (notFoundKeys []string, err error) {
	if len(caseKeys) == 0 {
		return nil, ErrNoCaseKeys
	}
	if startSearchDate == (time.Time{}) {
		startSearchDate = updter.conf.SearchStartDate
	}
	if maxSearchBack < 0 {
		maxSearchBack = updter.conf.MaxSearchBack
	}

	caseTypesMap := genCaseTypeMap(caseKeys)

	updatedAccords := []*UpdatedAccord{}
	searchErrors := []error{}
	// Refers to the searches per caseType not per caseId
	pendingSearch := len(caseTypesMap)

	updates := make(chan []*UpdatedAccord)
	complete := make(chan error)

	for cType, cIds := range caseTypesMap {
		go updter.getUpdates(&getUpdatesParams{
			updates:       updates,
			complete:      complete,
			caseType:      cType,
			caseIds:       cIds,
			startDate:     startSearchDate,
			daysBack:      maxSearchBack,
			exhaustSearch: exhaustSearch,
		})
	}

	for pendingSearch > 0 {
		select {
		case updt := <-updates:
			updatedAccords = append(updatedAccords, updt...)
		case err := <-complete:
			pendingSearch--
			if err != nil {
				searchErrors = append(searchErrors, err)
			}
		}
	}

	if len(updatedAccords) == 0 {
		for i, sErr := range searchErrors {
			fmt.Printf("Search Err[%d]: %v\n", i, sErr)
		}

		return nil, ErrNoUpdates
	}

	store := updter.getStore()
	if store == nil {
		return nil, ErrNilStore
	}

	err = store.Save(updatedAccords)
	if err != nil {
		return nil, ErrFailSave
	}

	foundMap := map[string]bool{}
	notFoundKeys = []string{}

	for _, acc := range updatedAccords {
		foundMap[acc.CaseKey] = true
	}

	for _, k := range caseKeys {
		if !foundMap[k] {
			notFoundKeys = append(notFoundKeys, k)
		}
	}

	return
}

func (updter *GeneralUpdater) FindUpdates(
	caseKeys []string,
	startSearchDate time.Time,
	maxSearchBack int,
	exhaustSearch bool,
) (accords []*UpdatedAccord, err error) {
	if len(caseKeys) == 0 {
		return nil, ErrNoCaseKeys
	}
	if startSearchDate == (time.Time{}) {
		startSearchDate = updter.conf.SearchStartDate
	}
	if maxSearchBack < 0 {
		maxSearchBack = updter.conf.MaxSearchBack
	}

	caseTypesMap := genCaseTypeMap(caseKeys)

	accords = []*UpdatedAccord{}
	searchErrors := []error{}
	// Refers to the searches per caseType not per caseId
	pendingSearch := len(caseTypesMap)

	updates := make(chan []*UpdatedAccord)
	complete := make(chan error)

	for cType, cIds := range caseTypesMap {
		go updter.getUpdates(&getUpdatesParams{
			updates:       updates,
			complete:      complete,
			caseType:      cType,
			caseIds:       cIds,
			startDate:     startSearchDate,
			daysBack:      maxSearchBack,
			exhaustSearch: exhaustSearch,
		})
	}

	for pendingSearch > 0 {
		select {
		case updt := <-updates:
			accords = append(accords, updt...)
		case err := <-complete:
			pendingSearch--
			if err != nil {
				searchErrors = append(searchErrors, err)
			}
		}
	}

	if len(accords) == 0 {
		for i, sErr := range searchErrors {
			fmt.Printf("Search Err[%d]: %v\n", i, sErr)
		}

		return nil, ErrNoUpdates
	}

	return
}

type getUpdatesParams struct {
	updates  chan<- []*UpdatedAccord
	complete chan<- error

	caseType      internal.CaseType
	caseIds       []string
	startDate     time.Time
	daysBack      int
	exhaustSearch bool
}

func (updter *GeneralUpdater) getUpdates(updateParams *getUpdatesParams) {
	pendingIds := make([]string, len(updateParams.caseIds))
	copy(pendingIds, updateParams.caseIds)
	y, m, d := updateParams.startDate.Date()
	searchDate := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	sent := 0

	for i := 0; i <= updateParams.daysBack; i++ {
		updatedAccords := []*UpdatedAccord{}
		fmt.Printf("Fetching with date %q for caseType %q [attempt %d]\n", searchDate, updateParams.caseType, i+1)

		data, err := updter.conf.FetchFn(searchDate, updateParams.caseType)
		if err != nil {
			if i == updateParams.daysBack && sent == 0 {
				fatalErr := errors.Join(
					fmt.Errorf("FetchFail: fetch for CaseType %q errored on date %s", updateParams.caseType, searchDate),
					ErrFatalSearch,
					err,
				)

				updateParams.complete <- fatalErr
				return
			}

			searchDate = searchDate.Add(internal.DayBack)
			continue
		}

		caseTable, err := updter.conf.ReadFn(data)
		if err != nil {
			if i == updateParams.daysBack && sent == 0 {
				fatalErr := errors.Join(
					fmt.Errorf("ReadFail: read for CaseType %q errored on date %s", updateParams.caseType, searchDate),
					ErrFatalSearch,
					err,
				)
				updateParams.complete <- fatalErr
				return
			}

			searchDate = searchDate.Add(internal.DayBack)
			continue
		}

		nextPendingIds := []string{}
		for _, cId := range pendingIds {
			caseRow := caseTable.Find(cId)
			if updateParams.exhaustSearch || caseRow == nil {
				nextPendingIds = append(nextPendingIds, cId)

				if caseRow == nil {
					continue
				}
			}
			caseRow.CaseType = string(updateParams.caseType)
			acc := UpdatedAccord{
				CaseKey:  caseRow.GetCaseKey(),
				CaseType: updateParams.caseType,
				CaseId:   caseRow.CaseId,
				Content:  caseRow.Accord,
				Date:     searchDate,
				Nature:   caseRow.Nature,
				OthIds:   caseRow.AllIds,
			}

			updatedAccords = append(updatedAccords, &acc)
		}

		updateParams.updates <- updatedAccords
		searchDate = searchDate.Add(internal.DayBack)
		sent += len(updatedAccords)

		if len(nextPendingIds) == 0 {
			break
		}
	}

	updateParams.complete <- nil
}

func (updter *GeneralUpdater) getStore() CaseStore {
	return updter.conf.Store
}

func (updter *GeneralUpdater) SetStore(st CaseStore) {
	updter.conf.Store = st
}

func genCaseTypeMap(keys []string) CaseTypesMap {
	caseTypesMap := CaseTypesMap{}

	for _, cK := range keys {
		parts := strings.Split(cK, readers.CaseKeySeparator)
		caseId, caseType := parts[0], internal.CaseType(parts[1])

		if _, ok := caseTypesMap[caseType]; ok {
			caseTypesMap[caseType] = append(caseTypesMap[caseType], caseId)
		} else {
			caseTypesMap[caseType] = []string{caseId}
		}
	}

	return caseTypesMap
}
