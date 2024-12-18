package controllers

import (
	"context"
	"database/sql"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/accupdter"
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

func (ctl *AccordUpdaterCtl) Update(caseKeys []string, searchStartDate time.Time, maxSearchBack int) ([]string, error) {
	return ctl.generalUpdater.Update(caseKeys, searchStartDate, maxSearchBack)
}

func (ctl *AccordUpdaterCtl) Startup(ctx context.Context, db *sql.DB) {
	ctl.ctx = ctx
	ctl.appDb = db

	ctl.generalUpdater.SetStore(accupdter.NewDefaultCaseStore(ctx, db))
}
