package messenger

import (
	"encoding/json"
	"net/http"

	messenger "github.com/princebillygk/se-job-aggregator-chatbot/internal/controller/messenger/inputs"
)

type MessageInput struct {
	*messenger.MessageEvent
	*messenger.EventProps
}

type MessageResponse struct {
	Text string `json:"text"`
}

func handleMessage(w http.ResponseWriter, input *MessageInput) error {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(MessageResponse{
		Text: "Message Received!",
	})
	return err
}
