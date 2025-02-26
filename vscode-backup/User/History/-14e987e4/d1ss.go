package signalsourceDAL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/jmoiron/sqlx"
)

type DAL interface {
	// GetActiveSignalSources will get all signal sources that are active
	GetActiveSignalSources(
		ctx context.Context,
		wl wlog.Logger,
	) ([]*entities.SignalSource, error)

	// GetSignalSourceByID will get a signal source by it's ID
	GetSignalSourceByID(
		ctx context.Context,
		wl wlog.Logger,
		ID string,
	) (*entities.SignalSource, error)
}

type signalsourceDAL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (DAL, error) {
	return &signalsourceDAL{
		db: db,
	}, nil
}

func (ssd *signalsourceDAL) GetActiveSignalSources(
	ctx context.Context,
	wl wlog.Logger,
) ([]*entities.SignalSource, error) {
	var rows []entities.SignalSource

	query := "SELECT id, enabled, type, ip, signal_version, created_at, updated_at, deleted_at FROM mvp_signal_source WHERE enabled = TRUE"
	err := ssd.db.Select(&rows, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active signals found")
		}

		return nil, err
	}

	var sources []*entities.SignalSource
	for _, row := range rows {
		source := row
		sources = append(sources, &source)
	}

	return sources, nil
}

func (ssd *signalsourceDAL) GetSignalSourceByID(
	ctx context.Context,
	wl wlog.Logger,
	ID string,
) (*entities.SignalSource, error) {
	var signal entities.SignalSource

	query := "SELECT id, enabled, type, ip, signal_version, created_at, updated_at, deleted_at FROM mvp_signal_source WHERE id = $1"
	err := ssd.db.Get(&signal, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no signals found for ID: %s", ID)
		}

		return nil, err
	}

	return &signal, nil
}
