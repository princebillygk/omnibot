package facebook

type SendRequestOutputBody struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		FbTraceID string `json:"fbtrace_id"`
	} `json:"error",omitempty`
}
