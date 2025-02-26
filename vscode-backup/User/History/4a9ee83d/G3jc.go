package reportsvc

import (
	"context"
	"yield-mvp/internal/emailer"
	"yield-mvp/internal/orderDAL"
	"yield-mvp/internal/strategyDAL"
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
	od       orderDAL.DAL
	sd       strategyDAL.DAL
	ethPrice *coinroutespriceconsumer.Consumer
	btcPrice *coinroutespriceconsumer.Consumer
	e        emailer.Client
	path     string
}

func New(
	d *sqlx.DB,
	od orderDAL.DAL,
	sd strategyDAL.DAL,
	ep *coinroutespriceconsumer.Consumer,
	bp *coinroutespriceconsumer.Consumer,
	e emailer.Client,
	path string,
) (SVC, error) {
	return &reportService{
		db:       d,
		od:       od,
		sd:       sd,
		ethPrice: ep,
		btcPrice: bp,
		e:        e,
		path:     path,
	}, nil
}
