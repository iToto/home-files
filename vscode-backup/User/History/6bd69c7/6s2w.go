package wlog

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

// BasicLogger is basic version of a wlog.Logger. It provides methods to log
// at different severity levels as well as methods to add metadata to the logger.
type BasicLogger struct {
	zlog zerolog.Logger
}

// Debug logs a Debug level message.
func (bl BasicLogger) Debug(msg string) {
	bl.zlog.Debug().Msg(msg)
}

// Debugf logs a Debug level message with formatting.
func (bl BasicLogger) Debugf(msg string, v ...interface{}) {
	bl.zlog.Debug().Msgf(msg, v...)
}

// Info logs an Info level message.
func (bl BasicLogger) Info(msg string) {
	bl.zlog.Info().Msg(msg)
}

// Infof logs an Info level message with formatting.
func (bl BasicLogger) Infof(msg string, v ...interface{}) {
	bl.zlog.Info().Msgf(msg, v...)
}

// Error logs an Error level message.
func (bl BasicLogger) Error(err error) {
	bl.zlog.Error().Msg(err.Error())
}

// WithStr returns the logger with added key-value strings metadata.
func (bl BasicLogger) WithStr(key string, value string) Logger {
	bl.zlog = bl.zlog.With().Str(key, value).Logger()
	return bl
}

// NewBasicLogger initializes a new BasicLogger with
// config values parsed from the runtime environment.
func NewBasicLogger() (Logger, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return NewBasicLoggerWithConfig(cfg)
}

// NewBasicLoggerWithConfig is the same as NewBasicLogger
// but can be used to specify a custom Logger config.
func NewBasicLoggerWithConfig(cfg *Config) (Logger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	zlog, err := zLogFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return BasicLogger{zlog: zlog}, nil
}
