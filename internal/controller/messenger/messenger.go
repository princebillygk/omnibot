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

	"github.com/princebillygk/omnibot/internal/services/users"
	"github.com/princebillygk/omnibot/internal/utility"
	"github.com/princebillygk/omnibot/pkg/facebook"
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
}

func New(pgSrvc *facebook.PageService, usrSrvc *users.Service) *Messenger {
	return &Messenger{pgSrvc, usrSrvc}
}

func (c Messenger) HandleWebhook(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		c.handleNotification(w, r)
	case "GET":
		c.handleWebhookVerification(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c Messenger) handleNotification(w http.ResponseWriter, r *http.Request) {
	var body *Notification
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := c.verifyRequestSignature(r, data)
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
		switch e := entry.Messaging[0]; {
		case e.MessageEvent != nil:
			err = c.handleMessage(r.Context(), w, &MessageInput{
				MessageEvent: e.MessageEvent,
				EventProps:   &e.EventProps,
			})
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (c Messenger) handleWebhookVerification(w http.ResponseWriter, r *http.Request) {
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

func (m Messenger) handleMessage(ctx context.Context, w http.ResponseWriter, input *MessageInput) error {
	w.WriteHeader(http.StatusOK)
	err := m.pgSrvc.SendMsg(input.Sender.ID, fmt.Sprintf("Message received with love %s", input.Message.Text))
	panic("wer234r23radff")
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (c Messenger) verifyRequestSignature(r *http.Request, payload []byte) bool {
	sign := r.Header.Get("X-Hub-Signature-256")
	fmt.Println("Signature", sign)

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
