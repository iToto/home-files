package reportsvc

import (
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/signalSourceDAL"
	"yield-mvp/internal/signallogger"
	"yield-mvp/internal/strategyDAL"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/userDAL"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/coinroutespriceconsumer"
	"yield-mvp/pkg/signalapi"

	"github.com/jmoiron/sqlx"
)

type reportService struct {
	db       *sqlx.DB
	ethPrice *coinroutespriceconsumer.Consumer
	btcPrice *coinroutespriceconsumer.Consumer
}

func New(
	s *signalapi.Client,
	c *coinroutesapi.Client,
	d *sqlx.DB,
	t *tradelogger.DataLogger,
	sl *signallogger.DataLogger,
	e exchangeDAL.DAL,
	ssd signalSourceDAL.DAL,
	sd strategyDAL.DAL,
	ud userDAL.DAL,
	ep *coinroutespriceconsumer.Consumer,
	bp *coinroutespriceconsumer.Consumer,
) (SVC, error) {
	return &signalService{
		db:          d,
		sc:          s,
		cc:          c,
		dl:          t,
		sl:          sl,
		exdal:       e,
		signalDAL:   ssd,
		strategyDAL: sd,
		userDAL:     ud,
		ethPrice:    ep,
		btcPrice:    bp,
	}, nil
}
