package facebook

type SendRequestInputBody struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}
