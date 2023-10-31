package template

type ButtonType string

const (
	ButtonTypeWebURL   ButtonType = "web_url"
	ButtonTypePostBack ButtonType = "postback"
	ButtonPhoneNumber  ButtonType = "phone_number"
)

type Button interface {
	GetButtonObject() map[string]any
}

// URLButton represents button object with web url
type URLButton struct {
	Title string
	URL   string
}

func (w URLButton) GetButtonObject() map[string]any {
	return map[string]any{
		"type":  ButtonTypeWebURL,
		"title": w.Title,
		"url":   w.URL,
	}
}

// PostbackButton respresents a postback event in webhook
type PostbackButton struct {
	Title   string
	Payload string
}

func (p PostbackButton) GetButtonObject() map[string]any {
	return map[string]any{
		"type":    ButtonTypePostBack,
		"title":   p.Title,
		"payload": p.Payload,
	}
}

// PhoneNumberButton represent a call button to a specified phone number
type PhoneNumberButton struct {
	Title       string
	PhoneNumber string
}

func (p PhoneNumberButton) GetButtonObject() map[string]any {
	return map[string]any{
		"type":    ButtonPhoneNumber,
		"title":   p.Title,
		"payload": p.PhoneNumber,
	}
}
