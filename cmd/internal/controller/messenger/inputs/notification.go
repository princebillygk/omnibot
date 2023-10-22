package messenger

type NotificationInput struct {
	Entry  []map[string]any `json:"entry"`
	Object string           `json:"object"`
}

// type NotificationInput struct {
// 	Entry  []Entry `json:"entry`
// 	Object string  `json:"object`
// }

// type Message struct {
// 	Mid  string `json:"mid,omitempty"`
// 	Text string `json:"text,omitempty"`
// }
// type Recipient struct {
// 	ID string `json:"id,omitempty"`
// }
// type Sender struct {
// 	ID string `json:"id,omitempty"`
// }
// type Messaging struct {
// 	Message   Message   `json:"message,omitempty"`
// 	Recipient Recipient `json:"recipient,omitempty"`
// 	Sender    Sender    `json:"sender,omitempty"`
// 	Timestamp int64     `json:"timestamp,omitempty"`
// }
// type Entry struct {
// 	ID        string      `json:"id,omitempty"`
// 	Messaging []Messaging `json:"messaging,omitempty"`
// 	Time      int64       `json:"time,omitempty"`
// }
