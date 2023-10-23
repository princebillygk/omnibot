package main

import (
	"fmt"
	"net/http"

	"github.com/princebillygk/se-job-aggregator-chatbot/internal/controller/messenger"
	"github.com/princebillygk/se-job-aggregator-chatbot/internal/utility"
	"github.com/princebillygk/se-job-aggregator-chatbot/pkg/facebook"
)

func main() {
	port := utility.GetEnv[int64]("PORT", 3000)
	pageAccessToken := utility.MustGetEnv[string]("PAGE_ACCESS_TOKEN")

	http.HandleFunc("/chat/messenger", messenger.New(facebook.NewPageService(pageAccessToken)).HandleWebhook)
	fmt.Printf("Running http server at port %d\n", port)

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
