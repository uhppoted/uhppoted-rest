package errors

import (
	"fmt"
)

type RESTError struct {
	Err       error  `json:"-"`
	DeviceID  uint32 `json:"-"`
	RequestID string `json:"request-id,omitempty"`
	Tag       string `json:"tag"`
	Status    int    `json:"-"`
	Message   string `json:"message"`
	Debug     string `json:"debug,omitempty"`
}

func (e *RESTError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewRESTError(tag, msg string) error {
	return &RESTError{
		RequestID: "",
		Tag:       tag,
		Message:   msg,
	}
}
