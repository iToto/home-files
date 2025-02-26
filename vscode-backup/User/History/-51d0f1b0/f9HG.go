package reportsvc

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/davecgh/go-spew/spew"
)

const (
	pathToCSV   = "reports/"
	csvFileName = "yield-report.csv"
)

func (rs *reportService) GenerateTradeReport(
	ctx context.Context,
	wl wlog.Logger,
) error {
	var reportData = make(map[string][]*entities.OrderReportRecord)

	// get all active strategies
	strategies, err := rs.sd.GetActiveStrategies(ctx, wl)
	if err != nil {
		return err
	}

	if len(strategies) < 1 {
		wl.Info("no strategies found for report")
		return fmt.Errorf("no strategies found for report")
	}

	wl.Debugf("strategies is of len %d", len(strategies))
	spew.Dump(strategies)

	// get report data for each strategy
	for _, strategy := range strategies {
		data, err := rs.od.GetOrderReportForStrategy(ctx, wl, strategy.Name)
		if err != nil {
			return err
		}

		reportData[strategy.Name] = data
	}

	// spew.Dump(reportData)
	csvFile, err := os.Create(pathToCSV + csvFileName)
	if err != nil {
		wl.Error(err)
		return err
	}

	headerSet := false
	// generate CSV for each strategy
	for _, s := range reportData {
		w := csv.NewWriter(csvFile)
		var headers []string

		// write headers for strategy csv
		if !headerSet {
			headers = s[0].GetHeaders()
			err = w.Write(headers)
			if err != nil {
				wl.Error(err)
				return err
			}
			headerSet = true
		}

		// calculate P&L
		wl.Debugf("calculating PL for: %s", s[0].Strategy)
		err = rs.calculatePLForStrategy(ctx, wl, s)
		if err != nil {
			wl.Error(err)
			return err
		}

		// write values to csv
		for _, row := range s {
			err = w.Write(row.GetValues())
			if err != nil {
				wl.Error(err)
				return err
			}
		}

		w.Flush()

	}

	// email csv to list of recipients
	recipients := []string{"jordan@yieldtechnologies.com"}
	subject := fmt.Sprintf("Weekly: yield-mvp trade report for %s", time.Now().Local().Format(time.RFC1123))
	message := fmt.Sprintf("Attached is the CSV for all trades made up until %s", time.Now().Format(time.RFC1123))
	pathToCSV := pathToCSV + csvFileName

	err = rs.e.SendEmailReport(ctx, wl, recipients, subject, message, pathToCSV)
	if err != nil {
		wl.Infof("error trying to email report: %w", err)
	}

	return nil
}

func (rs *reportService) calculatePLForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	data []*entities.OrderReportRecord,
) error {
	ethPrice, err := rs.ethPrice.GetPrice()
	if err != nil {
		wl.Error(err)
		return err
	}

	btcPrice, err := rs.btcPrice.GetPrice()
	if err != nil {
		wl.Error(err)
		return err
	}

	// loop through report data and populate PnL property
	for i := len(data) - 1; i > 0; i-- {
		// buy/long: avg-price of i-1 / avg-price i - 1
		// sell/short: -1 * avg-price of i-1 / avg-price i - 1
		// neutral: 0
		p, _ := strconv.ParseFloat(data[i-1].AvgPrice, 32)
		q, _ := strconv.ParseFloat(data[i].AvgPrice, 32)
		if data[i].Signal == "long" {
			data[i].PNL = strconv.FormatFloat(p/q-1, 'g', 2, 64)
		} else if data[i].Signal == "short" {
			data[i].PNL = strconv.FormatFloat(-1*(p/q-1), 'g', 2, 64)
		} else if data[i].Signal == "neutral" {
			data[i].PNL = "0"
		} else {
			wl.Infof("invalid signal: %s", data[i].Signal)
		}
	}

	// calculate PL on last trade using current price
	var p float64
	if data[0].Coin == "btc" {
		p = btcPrice
	} else {
		p = ethPrice
	}

	wl.Debugf("current price for %s = %f", data[0].Coin, p)

	q, _ := strconv.ParseFloat(data[0].AvgPrice, 32)
	if data[0].Signal == "long" {
		data[0].PNL = strconv.FormatFloat(p/q-1, 'g', 2, 64)
	} else if data[0].Signal == "short" {
		data[0].PNL = strconv.FormatFloat(-1*(p/q-1), 'g', 2, 64)
	} else if data[0].Signal == "neutral" {
		data[0].PNL = "0"
	} else {
		wl.Infof("invalid signal: %s", data[0].Signal)
	}

	return nil
}
