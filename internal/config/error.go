package config

type ApplicationError struct {
	HttpStatus   int
	Message      string
	DebugMessage string
}

func (ae ApplicationError) Error() string {
	return ae.Message
}
