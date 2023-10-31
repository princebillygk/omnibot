package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/princebillygk/omnibot/pkg/facebook/template"
)

const sendAPIURL = "https://graph.facebook.com/v2.6/me/messages"

type PageService struct {
	accessToken string
}

func NewPageService(accessToken string) *PageService {
	return &PageService{accessToken}
}

type FacebookResponse struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		FbtraceID string `json:"fbtrace_id"`
	} `json:"error"`
}

func (p PageService) callSendAPI(senderId string, msg map[string]any) error {
	inputBody, err := json.Marshal(&SendRequestInputBody{
		Recipient: Recipient{ID: senderId},
		Message:   msg,
	})

	if err != nil {
		return fmt.Errorf("Failed to parse message: %v", err)
	}

	req, err := http.NewRequest("POST", sendAPIURL, bytes.NewReader(inputBody))
	q := req.URL.Query()
	q.Add("access_token", p.accessToken)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return fmt.Errorf("Failed to initiate request %v:", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to send message! API Error %v:", err)
	}

	if res.StatusCode != http.StatusOK {
		var output SendRequestErrorBody
		json.NewDecoder(res.Body).Decode(&output)
		return fmt.Errorf("Unable to send message! Client Error: %s ", output.Error.Message)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println(string(inputBody))
	log.Println(string(resBody))
	return nil
}

func (p PageService) SendTextMessage(senderId string, msg string) error {
	return p.callSendAPI(senderId, map[string]any{
		"text": msg,
	})
}

func (p PageService) SendFromButtonTemplate(senderId string, msg string, buttons []template.Button) error {
	buttonObjects := make([]map[string]any, 0, len(buttons))
	for _, b := range buttons {
		buttonObjects = append(buttonObjects, b.GetButtonObject())
	}

	return p.callSendAPI(senderId, map[string]any{
		"attachment": map[string]any{
			"type": "template",
			"payload": map[string]any{
				"template_type": "button",
				"text":          msg,
				"buttons":       buttonObjects,
			},
		},
	})
}
