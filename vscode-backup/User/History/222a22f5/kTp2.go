package handler

import (
	"net/http"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func GetCoinRoutesExchangeAccounts(wl wlog.Logger, cc *coinroutesapi.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
