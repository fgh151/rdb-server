package utils

import (
	"log"
	"os"
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var (
	levelMap = map[logrus.Level]sentry.Level{
		logrus.TraceLevel: sentry.LevelDebug,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.PanicLevel: sentry.LevelFatal,
	}
)

type Converter func(entry *logrus.Entry, event *sentry.Event, hub *sentry.Hub)

type Option func(h *Hook)

type Hook struct {
	hub       *sentry.Hub
	levels    []logrus.Level
	tags      map[string]string
	extra     map[string]interface{}
	converter Converter
}

func NewSentryHook(dsn string, levels []logrus.Level, options ...Option) Hook {

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Environment:      os.Getenv("SENTRY_ENVIRONMENT"),
	})

	if err != nil {
		log.Fatal(err)
	}

	h := Hook{
		levels:    levels,
		hub:       sentry.CurrentHub(),
		converter: DefaultConverter,
	}

	for _, option := range options {
		option(&h)
	}

	return h
}

func (hook Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook Hook) Fire(entry *logrus.Entry) error {
	event := sentry.NewEvent()
	for k, v := range hook.extra {
		event.Extra[k] = v
	}
	for k, v := range hook.tags {
		event.Tags[k] = v
	}

	hook.converter(entry, event, hook.hub)

	sentry.CaptureEvent(event)

	return nil
}

func DefaultConverter(entry *logrus.Entry, event *sentry.Event, hub *sentry.Hub) {
	event.Level = levelMap[entry.Level]
	event.Message = entry.Message

	for k, v := range entry.Data {
		event.Extra[k] = v
	}

	if err, ok := entry.Data[logrus.ErrorKey].(error); ok {
		exception := sentry.Exception{
			Type:  reflect.TypeOf(err).String(),
			Value: err.Error(),
		}

		client := hub.Client()
		if client != nil && client.Options().AttachStacktrace {
			exception.Stacktrace = sentry.ExtractStacktrace(err)
		}

		event.Exception = []sentry.Exception{exception}
	}
}
