package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/orderlogger"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/render"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

func UpdateOpenOrders(
	wl wlog.Logger,
	cc *coinroutesapi.Client,
	db *sqlx.DB,
	ol *orderlogger.DataLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "order")

		// query all open orders in DB
		orders, err := getOpenOrders(ctx, db)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		if len(orders) < 1 {
			wl.Infof("no orders to process")
			render.JSON(ctx, wl, w, nil, http.StatusOK)
			return
		}

		// loop through and update info for each
		for _, o := range orders {
			wl.Debugf("cur order: %+v", o)

			// get latest state of from CR
			resp, err := cc.GetClientOrder(ctx, o.ClientOrderId)
			if err != nil {
				wl.Error(fmt.Errorf("unable to get order details from coinroutes: %w", err))
				continue
			}

			// update properties that need to be persisted in the DB
			if resp.FinishedAt != "" {
				ft, err := time.Parse(time.RFC3339Nano, resp.FinishedAt)
				if err != nil {
					wl.Error(fmt.Errorf("unable to parse finished time: %w", err))
					continue
				}
				o.FinishedAt = null.NewTime(ft, true)
			}

			o.AvgPrice = resp.AvgPrice
			o.Status = entities.OrderStatusType(resp.OrderStatus)
			o.ExecutedQty = resp.ExecutedQty

			// if status == closed or cancelled, log to BQ
			if o.Status == entities.Closed || o.Status == entities.Cancelled {
				if resp.FinishedAt == "" {
					wl.Debugf("finishedAt was empty for closed or cancelled order: %+v", resp)
				}

				err = ol.Log(ctx, wl, resp)
				if err != nil {
					wl.Infof("unable to log order to BQ: %w", err)
					continue
				}
			}

			// update record in DB (update if we successfully logged in BQ)
			err = updateOrder(ctx, &o, db)
			if err != nil {
				wl.Error(fmt.Errorf("unable to update order %s: %w", o.ClientOrderId, err))
				continue
			}
			wl.Debugf("updated order %s", o.ClientOrderId)

		}

		render.JSON(ctx, wl, w, nil, http.StatusOK)

	}
}

func getOpenOrders(ctx context.Context, db *sqlx.DB) ([]entities.Order, error) {
	var openOrders []entities.Order

	queryOrderStatus := "open"
	query := `SELECT 
		id,
		client_order_id,
		strategy,
		status,
		currency_pair,
		avg_price,
		executed_qty,
		finished_at,
		created_at,
		updated_at,
		deleted_at
	FROM mvp_order
	WHERE status = $1`

	rows, err := db.Queryx(query, queryOrderStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(
			&order.ID,
			&order.ClientOrderId,
			&order.Strategy,
			&order.Status,
			&order.CurrencyPair,
			&order.AvgPrice,
			&order.ExecutedQty,
			&order.FinishedAt,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan items: %w", err)
		}
		openOrders = append(openOrders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close getOpenOrders rows: %w", err)
	}

	return openOrders, nil
}

func updateOrder(ctx context.Context, order *entities.Order, db *sqlx.DB) error {
	var query string
	if null.Time.IsZero(order.FinishedAt) {
		query = `UPDATE mvp_order SET
			status = :status,
			avg_price = :avg_price,
			executed_qty = :executed_qty,
			updated_at = NOW()
		WHERE
			client_order_id = :client_order_id
		`
	} else {
		query = `UPDATE mvp_order SET
			status = :status,
			avg_price = :avg_price,
			executed_qty = :executed_qty,
			finished_at = :finished_at,
			updated_at = NOW()
		WHERE
			client_order_id = :client_order_id
		`
	}
	_, err := db.NamedQuery(query, order)
	if err != nil {
		return err
	}

	return nil
}
