package floodgate

import "context"

type service struct {
	OriginError   error  `json:"error"`
	StatusMessage string `json:"status"`
	MessageText   string `json:"message"`
	check         CheckerFunc
}

// Gets the state error. Nil if there is no error.
func (s service) Err() error {
	return s.OriginError
}

// Gets the state status. Returns "healthy" or "error"
func (s service) Status() string {
	return s.StatusMessage
}

// Gets the message. Error message if there is error. "ok" if there is no error
func (s service) Message() string {
	return s.MessageText
}

func (s *service) SetError(err error) {
	s.OriginError = err
}

func (s *service) SetStatus(status string) {
	s.StatusMessage = status
}

func (s *service) SetMessage(message string) {
	s.MessageText = message
}

func (s service) Check(ctx context.Context) error {
	return s.check(ctx)
}
