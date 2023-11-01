package messenger

import "time"

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
	*OptInEvent
}

type PostbackEvent struct {
	PostbackEvent map[string]any `json:"postback"`
}

type MessageEvent struct {
	Message Message `json:"message"`
}

type OptInEvent struct {
	OptIn OptIn `json:"optin"`
}

type OptInType string

const (
	OptInTypeNotificationMessage OptInType = "notification_message"
)

type NotificationMessageFreq string

const (
	NotificationMessageFreqDaily   string = "DAILY"
	NotificationMessageFreqWeekly  string = "WEEKLY"
	NotificationMessageFreqMonthly string = "MONTHLY"
)

type NotificationMessageStatus string

const (
	NotificationMessageStatusStop   NotificationMessageStatus = "STOP NOTIFICATIONS"
	NotificationMessageStatusResume NotificationMessageStatus = "RESUME NOTIFICATIONS"
)

type UserTokenStatus string

const (
	UserTokenStatusRefreshed    UserTokenStatus = "REFRESHED"
	UserTokenStatusNotRefreshed UserTokenStatus = "NOT_REFRESHED"
)

type OptIn struct {
	Type                      OptInType                 `json:"type"`
	Payload                   string                    `json:"payload"`
	NotificationMessageToken  string                    `json:"notification_messages_token"`
	NotificationMessageFreq   NotificationMessageFreq   `json:"notification_messages_frequency,omitempty"`
	NotificationTimeZone      string                    `json:"notification_messages_timezone,omitempty"`
	TokenExpiryTimestamp      *time.Time                `json:"token_expiry_timestamp,omitempty"`
	UserTokenStatus           UserTokenStatus           `json:"user_token_status,omitempty"`
	NotificationMessageStatus NotificationMessageStatus `json:"notification_messages_status,omitempty"`
	Title                     string                    `json:"title,omitempty"`
}
