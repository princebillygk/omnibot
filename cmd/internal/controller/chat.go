package controller

import (
	"log"
	"net/http"
	"os"
)

var msngrVerfToken string

func init() {
	var ok bool
	msngrVerfToken, ok = os.LookupEnv("MESSENGER_WEBHOOK_INTEGRATION_TOKEN")
	if !ok {
		panic("Messenger Verification Webhook doesn't exists")
	}
}

// Chat is a controller for messaging services
type Chat struct {
}

func (c Chat) HandleMessenger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)

	query := r.URL.Query()
	mode, token, challenge := query.Get("hub.mode"), query.Get("hub.verify_token"), query.Get("hub.challenge")

	if mode == "subscribe" && token == msngrVerfToken {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusForbidden)
}
