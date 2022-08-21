package web

import "errors"

type shutdownError struct {
	Message string
}

// NewShutdownError causes to signal a graceful shutdown
func NewShutdownError(message string) error {
	return &shutdownError{Message: message}
}

func (s *shutdownError) Error() string {
	return s.Message
}

func IsShutdown(e error) bool {
	var se *shutdownError
	return errors.As(e, &se)
}
