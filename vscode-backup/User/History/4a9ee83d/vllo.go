package reportsvc

import (
	"context"
	"yield-mvp/internal/orderDAL"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutespriceconsumer"

	"github.com/jmoiron/sqlx"
)

type SVC interface {
	GenerateTradeReport(
		ctx context.Context,
		wl wlog.Logger,
	) error
}

type reportService struct {
	db       *sqlx.DB
	od       *orderDAL.DAL
	ethPrice *coinroutespriceconsumer.Consumer
	btcPrice *coinroutespriceconsumer.Consumer
}

func New(
	d *sqlx.DB,
	ep *coinroutespriceconsumer.Consumer,
	bp *coinroutespriceconsumer.Consumer,
) (SVC, error) {
	return &reportService{
		db:       d,
		ethPrice: ep,
		btcPrice: bp,
	}, nil
}
