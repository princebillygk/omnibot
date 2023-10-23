package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (p PageService) SendMsg(senderId string, msg string) error {
	inputBody, err := json.Marshal(&SendRequestInputBody{
		Recipient: Recipient{ID: senderId},
		Message: Message{
			Text: msg,
		},
	})

	if err != nil {
		log.Fatalf("Failed to parse message: %v", err)
	}

	req, err := http.NewRequest("POST", sendAPIURL, bytes.NewReader(inputBody))
	q := req.URL.Query()
	q.Add("access_token", p.accessToken)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalf("Failed to initiate request %v:", err)
	}

	fmt.Println(string(inputBody))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to send message! API Error %v:", err)
	}

	if res.StatusCode != http.StatusOK {
		var output SendRequestOutputBody
		json.NewDecoder(res.Body).Decode(&output)
		return fmt.Errorf("Unable to send message! Client Error: %s ", output.Error.Message)
	}
	resBody, err := io.ReadAll(res.Body)
	log.Println(string(resBody))

	return nil
}
