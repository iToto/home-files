// Package wlog provides facilities to wrap an existing log package so that we can decouple
// the app's ability to log from a specific log package. It currently wraps zerolog under
// the hood.
package wlog

import (
	"context"
	"fmt"
	"os"
	"time"
	"yield-mvp/internal/ycontext"

	"github.com/rs/zerolog"
)

// Log item keys
const (
	LogKeyUserID    = "user_id"
	LogKeyPrn       = "prn"
	LogKeyChain     = "chain"
	LogKeyEventID   = "event_id"
	LogKeyEventType = "event_type"
	LogKeyRewardID  = "reward_id"
	LogKeyTraceID   = "logging.googleapis.com/trace"
)

func init() {
	// renames level to severity for GCP
	zerolog.LevelFieldName = "severity"
}

// Logger is an interface that provides standard log methods.
type Logger interface {
	// Debug logs a Debug level message.
	Debug(msg string)
	// Debugf logs a Debug level message with formatting.
	Debugf(msg string, v ...interface{})
	// Info logs an Info level message.
	Info(msg string)
	// Infof logs an Info level message with formatting.
	Infof(msg string, v ...interface{})
	// Error logs an Error level message.
	Error(err error)
	// WithStr returns the logger with added key-value strings metadata.
	WithStr(key string, value string) Logger
}

func zLogFromConfig(cfg *Config) (zerolog.Logger, error) {
	// create a new zerologger
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// config pretty logs
	if cfg.PrettyLogs {
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMilli})
	}

	// we want all available precision
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// set min log level
	level, err := zerolog.ParseLevel(cfg.MinLogLevel)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("error setting min log level: %w", err)
	}

	l = l.Level(level)

	return l, nil
}

// WithServiceRequest extracts any request context vars
// and adds them to the logger as metadata. This should be used at the top level
// when handling requests in a service.
// Currently this adds the following metadata to the logger:
//	- requestID
//	- callerID
//	- serviceName
//	- userID
//	- traceID
func WithServiceRequest(ctx context.Context, l Logger, serviceName string) Logger {
	if requestID, ok := ctx.Value(ycontext.ContextKeyRequestIDHeader).(string); ok {
		l = l.WithStr("requestID", requestID)
	}
	if callerID, ok := ctx.Value(ycontext.ContextKeyCallerIDHeader).(string); ok {
		l = l.WithStr("callerID", callerID)
	}
	if userID, ok := ctx.Value(ycontext.ContextKeyUserID).(string); ok {
		l = l.WithStr(LogKeyUserID, userID)
	}
	if traceID, ok := ctx.Value(ycontext.ContextKeyTraceIDHeader).(string); ok {
		l = l.WithStr(LogKeyTraceID, traceID)
	}

	l = l.WithStr("serviceName", serviceName)

	return l
}

// WithUserID adds the user id to the logger
func WithUserID(l Logger, userID string) Logger {
	return l.WithStr(LogKeyUserID, userID)
}

// WithPRN adds the prn to the logger
func WithPRN(l Logger, prn string) Logger {
	return l.WithStr(LogKeyPrn, prn)
}

func WithChain(l Logger, chain string) Logger {
	return l.WithStr()
}

// WithEventID adds the event id to the logger
func WithEventID(l Logger, eventID string) Logger {
	return l.WithStr(LogKeyEventID, eventID)
}

// WithEventType adds the event type to the logger
func WithEventType(l Logger, eventType string) Logger {
	return l.WithStr(LogKeyEventType, eventType)
}

// WithRewardID adds the user reward id to the logger
func WithRewardID(l Logger, userRewardID string) Logger {
	return l.WithStr(LogKeyRewardID, userRewardID)
}
