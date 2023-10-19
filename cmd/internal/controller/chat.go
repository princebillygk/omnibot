package controller

import (
	"fmt"
	"net/http"
)

// Chat is a controller for messaging services
type Chat struct {
}

func (c Chat) HandleMessenger(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mode, token, challenge := query.Get("hub.mode"), query.Get("hub.verify_token"), query.Get("hub.challenge")

	if mode == "subscribe" && token == "myToken" {
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s %s %s\n", mode, token, challenge)
	}
	w.WriteHeader(http.StatusForbidden)
}
