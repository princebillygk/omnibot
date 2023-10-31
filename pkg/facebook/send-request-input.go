package facebook

type SendRequestInputBody struct {
	Recipient Recipient      `json:"recipient"`
	Message   map[string]any `json:"message"`
}

type Recipient struct {
	ID string `json:"id"`
}
