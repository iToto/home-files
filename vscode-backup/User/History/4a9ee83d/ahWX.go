package reportsvc

import (
	"yield-mvp/pkg/coinroutespriceconsumer"

	"github.com/jmoiron/sqlx"
)

type reportService struct {
	db       *sqlx.DB
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
