package main

import (
	"fmt"
	"net/http"

	"github.com/princebillygk/se-job-aggregator-chatbot/cmd/internal/controller/messenger"
	"github.com/princebillygk/se-job-aggregator-chatbot/cmd/internal/utility"
)

func main() {
	port := utility.GetEnv[int64]("PORT", 3000)
	http.HandleFunc("/chat/messenger", messenger.New().HandleWebhook)
	fmt.Printf("Running http server at port %d\n", port)

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
