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

var (
	InvalidDeviceID    = fmt.Errorf("%w: Missing device ID", uhppoted.BadRequest)
	InvalidCard        = fmt.Errorf("%w: Missing/invalid card", uhppoted.BadRequest)
	InvalidCardNumber  = fmt.Errorf("%w: Missing/invalid card number", uhppoted.BadRequest)
	InvalidDoorID      = fmt.Errorf("%w: Missing/invalid door ID", uhppoted.BadRequest)
	InvalidDoorDelay   = fmt.Errorf("%w: Missing/invalid door delay", uhppoted.BadRequest)
	InvalidDoorControl = fmt.Errorf("%w: Missing/invalid door control", uhppoted.BadRequest)
	InvalidEventID     = fmt.Errorf("%w: Missing/invalid event ID", uhppoted.BadRequest)
	InvalidDate        = fmt.Errorf("%w: Missing/invalid date", uhppoted.BadRequest)
	InvalidDateTime    = fmt.Errorf("%w: Missing/invalid date/time", uhppoted.BadRequest)
	RequestFailed      = fmt.Errorf("%w: Request failed", uhppoted.InternalServerError)
)

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
