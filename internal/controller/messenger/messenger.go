package messenger

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/princebillygk/omnibot/internal/config"
	"github.com/princebillygk/omnibot/internal/services/users"
	"github.com/princebillygk/omnibot/internal/utility"
	"github.com/princebillygk/omnibot/pkg/facebook"
	"github.com/princebillygk/omnibot/pkg/facebook/template"
)

// TODO: Handle duplicacy of request

var msngrVerfToken string
var appSecret string

func init() {
	msngrVerfToken = utility.MustGetEnv[string]("MESSENGER_VERIFY_TOKEN")
	appSecret = utility.MustGetEnv[string]("APP_SECRET")
}

// Messenger is a controller for messaging services
type Messenger struct {
	pgSrvc  *facebook.PageService
	usrSrvc *users.Service

	logger *config.Logger
}

func New(pgSrvc *facebook.PageService, usrSrvc *users.Service, logger *config.Logger) *Messenger {
	return &Messenger{pgSrvc, usrSrvc, logger}
}

// HandleWebhook routes webhook request to the actual handler depending on the request method
func (m Messenger) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		m.handleNotification(w, r)
	case "GET":
		m.handleWebhookVerification(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m Messenger) verifyRequestSignature(r *http.Request, payload []byte) bool {
	sign := r.Header.Get("X-Hub-Signature-256")

	if sign == "" {
		return false
	}

	givenHash, ok := strings.CutPrefix(sign, "sha256=")
	if !ok {
		return false
	}

	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write(payload)

	expectedHash := hex.EncodeToString(h.Sum(nil))

	if givenHash != expectedHash {
		return false
	}
	return true
}

// handleNotification handle notifications sent from messenger webhooks
func (m Messenger) handleNotification(w http.ResponseWriter, r *http.Request) {
	var body *Notification
	data, err := io.ReadAll(r.Body)
	fmt.Println(string(data))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := m.verifyRequestSignature(r, data)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Object != "page" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, entry := range body.Entry {
		// fmt.Printf("%#v", entry.Messaging)
		switch e := entry.Messaging[0]; {
		case e.MessageEvent != nil:
			err = m.handleMessageNotification(r.Context(), w, e.MessageEvent, &e.EventProps)
		case e.OptInEvent != nil:
			err = m.handleOptInEvent(r.Context(), w, e.OptInEvent, &e.EventProps)
		}

		if err != nil {
			m.logger.LogError(err)
		}
	}
}

// handleWebhookVerification handles messenger webhook verification
func (m Messenger) handleWebhookVerification(w http.ResponseWriter, r *http.Request) {
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

func (m Messenger) handleMessageNotification(ctx context.Context, w http.ResponseWriter, me *MessageEvent, props *EventProps) error {
	w.WriteHeader(http.StatusOK)
	switch me.Message.Text {
	case "subscribe":
		return m.pgSrvc.SendOneTimeNotificationRequest(props.Sender.ID, "Subscribe", "subscribe")
	case "buttons":
		return m.pgSrvc.SendFromButtonTemplate(
			props.Sender.ID,
			fmt.Sprintf("Message received with love %s", me.Message.Text),
			[]template.Button{
				template.URLButton{
					Title: "My Portfolio",
					URL:   "https://princebillygk.github.io/",
				},
				template.PostbackButton{
					Title:   "Poke me",
					Payload: "message",
				},
				template.PhoneNumberButton{
					Title:       "Call me",
					PhoneNumber: "01521432424",
				},
			},
		)
	default:
		return m.pgSrvc.SendTextMessage(props.Sender.ID, "Prince Billy Graham Karmoker")
	}
}

func (m Messenger) handleOptInEvent(ctx context.Context, w http.ResponseWriter, oe *OptInEvent, props *EventProps) error {
	fmt.Printf("%#v", oe)
	return nil
}
