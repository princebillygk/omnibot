package config

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type httpLogger struct {
	handler http.Handler
}

func (h httpLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Recover
	defer func() {
		if r := recover(); r != nil {
			sentry.CurrentHub().Recover(r)
			sentry.Flush(time.Second * 5)
		}
	}()
	h.handler.ServeHTTP(w, r)
}

// Logger is used to keep logs of application
type Logger struct {
	Db                  *mongo.Database       // Required
	SentryOptions       *sentry.ClientOptions // Optional
	isSentryInitialized bool
}

// Init different services required for logger. It is required to initialize
// logger before using any service of logger
func (l *Logger) Init() error {
	if l.SentryOptions != nil {
		if err := sentry.Init(*l.SentryOptions); err != nil {
			return err
		}
		l.isSentryInitialized = true
	}

	if l.Db == nil {
		panic("A Mongo database instance is required for logger!")
	}
	return nil
}

// Shutdown closes and cleanup services that are intialized for logger
func (l *Logger) Shutdown() {
	if l.isSentryInitialized {
		sentry.Flush(time.Second * 10)
	}
}

// WrapHTTPHandler wraps regular http handler with logger
func (l *Logger) WrapHTTPHandler(h http.Handler) http.Handler {
	return httpLogger{h}
}
