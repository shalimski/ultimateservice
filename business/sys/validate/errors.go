package validate

import "errors"

// ErrInvalidID occurs when an ID is not a valid form
var ErrInvalidID = errors.New("Id is not valid")

// ErrorResponse is the form used for API responses from failures
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

// RequestError to pass an error during the request  
type RequestError struct {
	Err    error
	Status int
	Fields error
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}
