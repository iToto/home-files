package handler

import (
	"net/http"
	"reflect"
	"yield-mvp/internal/service/exchangesvc"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/render"
)

func GenerateExchangeReport(
	wl wlog.Logger,
	exchangeService exchangesvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "exchange")
		data, err := exchangeService.GenereateReport(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.HTMLTable(ctx, wl, w, convertStructToMap(data), http.StatusOK)
	}
}

func convertStructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return result
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		result[field.Name] = fieldValue.Interface()
	}

	return result
}
