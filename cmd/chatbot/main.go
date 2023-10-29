package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/princebillygk/omnibot/internal/controller/messenger"
	"github.com/princebillygk/omnibot/internal/services/users"
	"github.com/princebillygk/omnibot/internal/utility"
	"github.com/princebillygk/omnibot/pkg/facebook"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := utility.GetEnv[int64]("PORT", 3000)
	pageAccessToken := utility.MustGetEnv[string]("PAGE_ACCESS_TOKEN")
	mongoURI := utility.MustGetEnv[string]("MONGO_URI")

	ctx := context.TODO()

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
	userService := users.NewService(db)

	err = userService.UpdateLastActivityTime(ctx, "1234", time.Time)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/chat/messenger", messenger.New(facebook.NewPageService(pageAccessToken)).HandleWebhook)
	fmt.Printf("Running http server at port %d\n", port)

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
