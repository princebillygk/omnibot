package messenger

type NotificationInput struct {
	Entry  []Entry `json:"entry"`
	Object string  `json:"object"`
}

type Message struct {
	Mid  string `json:"mid,omitempty"`
	Text string `json:"text,omitempty"`
}

type Recipient struct {
	ID string `json:"id,omitempty"`
}

type Sender struct {
	ID string `json:"id,omitempty"`
}

type WebhookEvent interface {
	MessageEvent | PostbackEvent
}

type WebhookEventBase struct {
	Recipient Recipient `json:"recipient"`
	Sender    Sender    `json:"sender"`
	Timestamp int64     `json:"timestamp"`
}

type MessageEvent struct {
	WebhookEventBase
	Message Message `json:"message"`
}

type PostbackEvent struct {
	WebhookEventBase
	PostbackEvent map[string]any `json:'postback'`
}

type Entry struct {
	ID        string         `json:"id"`
	Messaging []WebhookEvent `json:"messaging"`
	Time      int64          `json:"time"`
}
