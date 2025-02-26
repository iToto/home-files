// Package render provide helper functions to render json reponses,
// it also provide a custom error type suitable for error marshaling.
package render

import (
	"boilerplate-go-api/internal/wlog"
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

// JSONErr writes the json-encoded error message to the response.
func JSONErr(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error, status int) {
	var e *Error
	if !errors.As(err, &e) {
		err = &Error{Err: err}
	}
	JSON(ctx, wl, w, err, status)
}

// InternalError logs the error and returns a 500 internal server error.
func InternalError(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Error(err)
	JSONErr(ctx, wl, w, ErrInternal, http.StatusInternalServerError)
}

// AdminInternalError log and writes the json-encoded error message to the response
// with a 500 internal server error.
func AdminInternalError(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Error(err)
	JSONErr(ctx, wl, w, err, http.StatusInternalServerError)
}

// NotFound writes the json-encoded error message to the response
// with a 404 not found status code.
func NotFound(ctx context.Context, wl wlog.Logger, w http.ResponseWriter, err error) {
	wl.Debug(err.Error())
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
