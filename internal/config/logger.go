package config

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type httpLogger struct {
	handler http.Handler
}

func (h httpLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("Request received")
	h.handler.ServeHTTP(w, r)
	log.Println("Request is: ", time.Since(start))
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error occurred: ===================>")
		}
	}()
}

type Logger struct {
	Db                  *mongo.Database
	SentryOptions       *sentry.ClientOptions
	isSentryInitialized bool
}

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

func (l *Logger) Shutdown() {
	if l.isSentryInitialized {
		sentry.Flush(time.Second * 10)
	}
}

func (l *Logger) WrapHTTPHandler(h http.Handler) http.Handler {
	return httpLogger{h}
}
