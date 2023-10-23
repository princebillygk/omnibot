package messenger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/princebillygk/se-job-aggregator-chatbot/pkg/facebook"
)

var msngrVerfToken string

func init() {
	var ok bool
	msngrVerfToken, ok = os.LookupEnv("MESSENGER_VERIFY_TOKEN")
	if !ok {
		panic("Messenger Verification Webhook doesn't exists")
	}
}

// Messenger is a controller for messaging services
type Messenger struct {
	pgSrvc *facebook.PageService
}

func New(pgSrvc *facebook.PageService) *Messenger {
	return &Messenger{pgSrvc}
}

func (c Messenger) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.handleNotification(w, r)
	case "GET":
		c.handleVerification(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c Messenger) handleNotification(w http.ResponseWriter, r *http.Request) {
	var body *Notification
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", body)

	if body.Object != "page" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, entry := range body.Entry {
		switch e := entry.Messaging[0]; {
		case e.MessageEvent != nil:
			err = c.handleMessage(w, &MessageInput{
				MessageEvent: e.MessageEvent,
				EventProps:   &e.EventProps,
			})
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (c Messenger) handleVerification(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mode, token, challenge := query.Get("hub.mode"), query.Get("hub.verify_token"), query.Get("hub.challenge")

	if mode != "subscribe" && token != msngrVerfToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(challenge))
	defer log.Println("Verified webhook callback url!")
}

type MessageInput struct {
	*MessageEvent
	*EventProps
}

func (m Messenger) handleMessage(w http.ResponseWriter, input *MessageInput) error {
	w.WriteHeader(http.StatusOK)
	err := m.pgSrvc.SendMsg(input.Sender.ID, fmt.Sprintf("Message received with love %s", input.Message.Text))
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
