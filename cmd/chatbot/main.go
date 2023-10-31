package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/sentry-go"
	"github.com/princebillygk/omnibot/internal/config"
	"github.com/princebillygk/omnibot/internal/controller/messenger"
	"github.com/princebillygk/omnibot/internal/services/users"
	"github.com/princebillygk/omnibot/internal/utility"
	"github.com/princebillygk/omnibot/pkg/facebook"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	port := utility.GetEnv[int64]("PORT", 3000)

	pageAccessToken := utility.MustGetEnv[string]("PAGE_ACCESS_TOKEN")
	mongoURI := utility.MustGetEnv[string]("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("omnibot")
	usrSrvc := users.NewService(db)

	sentryDsn := utility.MustGetEnv[string]("SENTRY_DSN")
	logger := config.Logger{
		Db:            db,
		SentryOptions: &sentry.ClientOptions{Dsn: sentryDsn},
	}

	err = logger.Init()
	if err != nil {
		panic(err)
	}

	// Setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/chat/messenger", messenger.New(facebook.NewPageService(pageAccessToken), usrSrvc, &logger).HandleWebhook)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: logger.WrapHTTPHandler(mux),
	}

	sigs, exit := make(chan os.Signal, 1), make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Cleanup
	go func() {
		<-sigs
		fmt.Println("Shutting down all services...")
		logger.Shutdown()
		_ = server.Shutdown(ctx)
		fmt.Println("Exiting!")
		exit <- true
	}()

	fmt.Printf("Running http server at port %d\n", port)
	server.ListenAndServe()
	<-exit
}
