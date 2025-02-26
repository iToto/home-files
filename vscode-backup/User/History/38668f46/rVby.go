// Package render provide helper functions to render json reponses,
// it also provide a custom error type suitable for error marshaling.
package render

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"yield-mvp/internal/wlog"
)

// JSON writes the json-encoded message to the response.
func JSON(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		wl.Error(err)
	}
}

// HTMLTable writes a slice of structs in a HTML table to the response.
func HTMLTable(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, data []interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if len(data) == 0 {
		return
	}

	// Generate HTML table
	table := "<table style=\"border-collapse: collapse;\">"
	table += "<tr>"
	// Get the field names from the first struct
	fields := reflect.ValueOf(data[0]).Elem().Type()
	wl.Debugf("fields: %v", fields)

	numFields := fields.NumField()
	wl.Debugf("numFields: %v", numFields)

	for i := 0; i < numFields; i++ {
		field := fields.Field(i)
		table += "<th style=\"border: 1px solid black;\">" + field.Name + "</th>"
	}
	table += "</tr>"

	// Iterate over the data and generate table rows
	for _, item := range data {
		table += "<tr>"
		values := reflect.ValueOf(item).Elem()
		for i := 0; i < values.NumField(); i++ {
			value := values.Field(i)
			table += "<td style=\"border: 1px solid black;\">" + fmt.Sprintf("%v", value.Interface()) + "</td>"
		}
		table += "</tr>"
	}

	table += "</table>"

	// Write the HTML table to the response
	if _, err := w.Write([]byte(table)); err != nil {
		wl.Error(err)
	}
	JSONErr(ctx, wl, w, ErrNotFound, http.StatusNotFound)
}

// Unauthorized writes the json-encoded error message to the response
// with a 401 unauthorized status code.
func Unauthorized(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Info(err.Error())
	JSONErr(ctx, wl, w, ErrUnauthorized, http.StatusUnauthorized)
}

// Forbidden writes the json-encoded error message to the response
// with a 403 forbidden status code.
func Forbidden(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Info(err.Error())
	JSONErr(ctx, wl, w, ErrForbidden, http.StatusForbidden)
}

// BadRequest writes the json-encoded error message to the response
// with a 400 bad request status code.
func BadRequest(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Debug(err.Error())
	JSONErr(ctx, wl, w, err, http.StatusBadRequest)
}

// Conflict writes the json-encoded error message to the response
// with a 409 conflict status code.
func Conflict(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Info(err.Error())
	JSONErr(ctx, wl, w, err, http.StatusConflict)
}

// TooManyRequests writes the json-encoded error message to the response
// with a 429 too many requests status code.
func TooManyRequests(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Info(err.Error())
	JSONErr(ctx, wl, w, err, http.StatusTooManyRequests)
}

// UpgradeRequired responds with a 412 status codes and includes the
// the a friendly user message in the response body.
func UpgradeRequired(ctx context.Context, wl wlog.Logger, w http.ResponseWriter) {
	body := map[string]string{
		"message": "Please download the latest version of our app to continue.",
	}

	wl.Info("forcing client to upgrade")
	JSON(ctx, wl, w, body, http.StatusPreconditionFailed)
}

// ImagePNG writes the image to the response.
func ImagePNG(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	if b == nil {
		return
	}
	n, err := w.Write(b)
	if err != nil {
		wl.Error(err)
	} else if n != len(b) {
		wl.Error(errors.New("failed to write all bytes"))
	}
}

// ACKPushEvent acknowledges a PubSub events
func ACKPushEvent(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	if err != nil {
		wl.Error(err)
	}
	w.WriteHeader(http.StatusOK) // only 200 seems to be accepted by the pubsub emulator
}

// NACKPushEvent un-acknowledges a PubSub event
// use this when you want to retry the event
func NACKPushEvent(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	if err != nil {
		InternalError(ctx, wl, w, err)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
