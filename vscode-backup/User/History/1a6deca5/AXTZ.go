package render_test

import (
	"boilerplate-go-api/internal/wlog"
	"boilerplate-go-api/pkg/render"
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func TestJSONErr(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		inputCode    int
		expectedResp string
	}{
		{
			"json friendly error",
			render.NewErrorStr("fake error"),
			500,
			`{"error":"fake error"}`,
		},
		{
			"non-json friendly error",
			errors.New("fake error"),
			500,
			`{"error":{}}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.JSONErr(context.Background(), wlog.NewNopLogger(), w, tc.input, tc.inputCode)

		assert.Equal(t, w.Code, tc.inputCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestInternalError(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			errors.New("fake error"),
			500,
			`{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.InternalError(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestNotFound(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			errors.New("fake error"),
			404,
			`{"error":"Not Found"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.NotFound(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestUnauthorized(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			errors.New("fake error"),
			401,
			`{"error":"Unauthorized"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.Unauthorized(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestForbidden(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			errors.New("fake error"),
			403,
			`{"error":"Forbidden"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.Forbidden(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestBadRequest(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			validation.Errors{"field1": errors.New("this field is required")},
			400,
			`{"error":{"field1":"this field is required"}}`,
		},
		{
			"subdocuments",
			validation.Errors{
				"field1": errors.New("this field is required"),
				"field2": validation.Errors{"field3": errors.New("this field is required")},
			},
			400,
			`{"error":{"field1":"this field is required","field2":{"field3":"this field is required"}}}`,
		},
		{
			"non-json friendly error",
			errors.New("fake error"),
			400,
			`{"error":{}}`,
		},
		{
			"json friendly error",
			render.NewErrorStr("fake error"),
			400,
			`{"error":"fake error"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.BadRequest(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}

func TestConflict(t *testing.T) {
	testData := []struct {
		name         string
		input        error
		expectedCode int
		expectedResp string
	}{
		{
			"happy path",
			render.NewErrorStr("fake error"),
			409,
			`{"error":"fake error"}`,
		},
	}

	for _, tc := range testData {
		w := httptest.NewRecorder()
		render.Conflict(context.Background(), wlog.NewNopLogger(), w, tc.input)

		assert.Equal(t, w.Code, tc.expectedCode)
		assert.Equal(t, strings.TrimSpace(w.Body.String()), tc.expectedResp)
	}
}
