package utilities

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

type LoggerParams struct {
	Mode         string
	SentryParams struct {
		Dsn         string
		Environment string
	}
}

func InitLogger(params LoggerParams) (Logger, func(time.Duration) bool, error) {
	switch params.Mode {
	case "stdout":
		return &StdoutLogger{}, nil, nil
	case "sentry":
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              params.SentryParams.Dsn,
			TracesSampleRate: 1.0,
			Environment:      params.SentryParams.Environment,
		})
		return &sentryLogger{}, sentry.Flush, err
	default:
		return &NullLogger{}, nil, nil
	}
}

type Logger interface {
	LogMessageln(a ...any)
	LogMessagef(format string, a ...any)
	LogError(err error)
}

type NullLogger struct{}

func (l *NullLogger) LogMessageln(a ...any) {}

func (l *NullLogger) LogMessagef(format string, a ...any) {}

func (l *NullLogger) LogError(err error) {}

type StdoutLogger struct{}

func (l *StdoutLogger) LogMessageln(a ...any) {
	fmt.Println(a...)
}

func (l *StdoutLogger) LogMessagef(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (l *StdoutLogger) LogError(err error) {
	fmt.Println(err)
}

type sentryLogger struct{}

func (l *sentryLogger) LogMessageln(a ...any) {
	sentry.CaptureMessage(fmt.Sprintln(a...))
}

func (l *sentryLogger) LogMessagef(format string, a ...any) {
	sentry.CaptureMessage(fmt.Sprintf(format, a...))
}

func (l *sentryLogger) LogError(err error) {
	sentry.CaptureException(err)
}
