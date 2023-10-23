package messenger

type Notification struct {
	Entry  []Entry `json:"entry"`
	Object string  `json:"object"`
}

type Message struct {
	MID  string `json:"mid,omitempty"`
	Text string `json:"text,omitempty"`
}

type Recipient struct {
	ID string `json:"id,omitempty"`
}

type Sender struct {
	ID string `json:"id,omitempty"`
}

type Entry struct {
	ID        string  `json:"id"`
	Messaging []Event `json:"messaging"`
	Time      int64   `json:"time"`
}

type EventProps struct {
	Recipient Recipient `json:"recipient"`
	Sender    Sender    `json:"sender"`
	Timestamp int64     `json:"timestamp"`
}

type Event struct {
	EventProps
	*MessageEvent
	*PostbackEvent
}

type MessageEvent struct {
	Message Message `json:"message,omitempty"`
}

type PostbackEvent struct {
	PostbackEvent map[string]any `json:"postback,omitempty"`
}
