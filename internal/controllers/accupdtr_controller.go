package controllers

import (
	"context"
	"database/sql"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/accupdter"
	"github.com/vladwithcode/lex_app/internal/db"
	"github.com/vladwithcode/lex_app/internal/fetchers"
	"github.com/vladwithcode/lex_app/internal/readers"
)

type AccUpdterOpts accupdter.AccUpdterOpts

type AccordUpdaterCtl struct {
	ctx   context.Context
	appDb *sql.DB

	generalUpdater *accupdter.GeneralUpdater
}

func NewAccordUpdaterCtl() *AccordUpdaterCtl {
	return &AccordUpdaterCtl{
		generalUpdater: accupdter.NewGeneralUpdater(&accupdter.GenUpdterConf{
			Region:          internal.RegionDefault,
			ReadFn:          readers.NewReader(internal.RegionDefault),
			FetchFn:         fetchers.NewFetcher(internal.RegionDefault),
			SearchStartDate: time.Now(),
			MaxSearchBack:   0,
		}),
	}
}

func (ctl *AccordUpdaterCtl) FindUpdates(
	caseKeys []string,
	searchStartDate time.Time,
	maxSearchBack int,
	exhaustSearch bool,
) ([]*accupdter.UpdatedAccord, error) {
	return ctl.generalUpdater.FindUpdates(caseKeys, searchStartDate, maxSearchBack, exhaustSearch)
}

func (ctl *AccordUpdaterCtl) Update(caseKeys []string, searchStartDate time.Time, maxSearchBack int, exhaustSearch bool) ([]string, error) {
	return ctl.generalUpdater.Update(caseKeys, searchStartDate, maxSearchBack, exhaustSearch)
}

func (ctl *AccordUpdaterCtl) FindCasesAndUpdate(
	searchStartDate time.Time,
	maxSearchBack int,
	exhaustSearch bool,
	findOpts *db.FindCaseOptions,
) (notFoundKeys []string, err error) {
	var caseKeys []string
	var cases []*db.LexCase
	if findOpts == nil {
		cases, err = db.FindFilteredCases(ctl.ctx, ctl.appDb, findOpts)
	} else {
		cases, err = db.FindAllCases(ctl.ctx, ctl.appDb)
	}
	if err != nil {
		return nil, err
	}
	for _, c := range cases {
		caseKeys = append(caseKeys, c.GetCaseKey())
	}
	return ctl.generalUpdater.Update(caseKeys, searchStartDate, maxSearchBack, exhaustSearch)
}

func (ctl *AccordUpdaterCtl) Startup(ctx context.Context, db *sql.DB) {
	ctl.ctx = ctx
	ctl.appDb = db

	ctl.generalUpdater.SetStore(accupdter.NewDefaultCaseStore(ctx, db))
}
