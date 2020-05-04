package errors

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"net/http"
)

type IError struct {
	Err      error  `json:"-"`
	DeviceID uint32 `json:"-"`
	Tag      string `json:"-"`
	Code     int    `json:"error-code"`
	Message  string `json:"message"`
}

func (e *IError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func Errorf(err error, deviceID uint32, tag string, msg string) *IError {
	status := http.StatusInternalServerError

	if errors.Is(err, uhppoted.InternalServerError) {
		status = http.StatusInternalServerError
	} else if errors.Is(err, uhppoted.NotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, uhppoted.BadRequest) {
		status = http.StatusBadRequest
	}

	return &IError{
		Err:      err,
		DeviceID: deviceID,
		Tag:      tag,
		Code:     status,
		Message:  msg,
	}
}
