package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Accord struct {
	Id      string    `json:"id" db:"id"`
	ForCase string    `json:"forCase" db:"for_case"`
	Content string    `json:"content" db:"content"`
	Date    time.Time `json:"date" db:"date"`
	rawData string
}

func NewAccord(caseId string) *Accord {
	return &Accord{
		Id:      uuid.Must(uuid.NewV7()).String(),
		ForCase: caseId,
	}
}

func (a *Accord) GetRawData() string {
	return a.rawData
}

func FindAllAccordsForCase(ctx context.Context, appDb *sql.DB, caseId string) ([]*Accord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := appDb.QueryContext(
		ctx,
		"SELECT id, for_case, content, date, raw_data FROM accords WHERE for_case = $1",
		caseId,
	)
	if err != nil {
		return nil, err
	}

	accords := []*Accord{}
	for rows.Next() {
		a := Accord{}
		var rd sql.NullString
		rows.Scan(
			&a.Id,
			&a.ForCase,
			&a.Content,
			&a.Date,
			&rd,
		)
		if rd.Valid {
			v, _ := rd.Value()
			a.rawData = v.(string)
		}
		accords = append(accords, &a)
	}

	return accords, nil
}

func FindLatestAccordForCase(ctx context.Context, appDb *sql.DB, caseId string) (*Accord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := appDb.QueryRowContext(
		ctx,
		"SELECT id, for_case, content, max(date), raw_data FROM accords WHERE for_case = $1",
		caseId,
	)

	accord := Accord{}
	var rd sql.NullString
	row.Scan(
		&accord.Id,
		&accord.ForCase,
		&accord.Content,
		&accord.Date,
		&rd,
	)
	if rd.Valid {
		v, _ := rd.Value()
		accord.rawData = v.(string)
	}

	return &accord, nil
}

func InsertAccord(ctx context.Context, appDb *sql.DB, accord *Accord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := appDb.ExecContext(
		ctx,
		"INSERT INTO accords (id, for_case, content, date, raw_data) VALUES @Id, @ForCase, @Content, @Date, @RawData",
		sql.Named("Id", accord.Id),
		sql.Named("ForCase", accord.ForCase),
		sql.Named("Content", accord.Content),
		sql.Named("Date", accord.Date),
		sql.Named("RawData", accord.rawData),
	)
	if err != nil {
		return err
	}

	return nil
}
