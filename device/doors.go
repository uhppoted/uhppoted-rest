package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
)

func GetDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	}

	response, err := impl.GetDoorDelay(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-door-delay", fmt.Sprintf("Error retrieving door delay for device %v, door %d", deviceID, door)),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: response.Delay,
	}, nil
}

func SetDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-delay", "Error reading request"),
			err
	}

	body := struct {
		Delay *uint8 `json:"delay"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-delay", "Error parsing request"),
			err
	}

	if body.Delay == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-delay", "Missing/invalid door delay"),
			fmt.Errorf("Missing/invalid door delay value in request body (%s)", string(blob))
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
		Delay:    *body.Delay,
	}

	response, err := impl.SetDoorDelay(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-delay", "Error setting device door delay"),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: response.Delay,
	}, nil
}

func GetDoorControl(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	}

	response, err := impl.GetDoorControl(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-door-control", fmt.Sprintf("Error retrieving door control for device %v, door %d", deviceID, door)),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		ControlState uhppoted.ControlState `json:"control"`
	}{
		ControlState: response.Control,
	}, nil
}

func SetDoorControl(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-control", "Error reading request"),
			err
	}

	body := struct {
		Control *uhppoted.ControlState `json:"control"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-control", "Error parsing request"),
			err
	}

	if body.Control == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-control", "Missing/invalid door control"),
			fmt.Errorf("Missing/invalid door control value in request body (%s)", string(blob))
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
		Control:  *body.Control,
	}

	response, err := impl.SetDoorControl(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-control", "Error setting device door control"),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		Control uhppoted.ControlState `json:"control"`
	}{
		Control: response.Control,
	}, nil
}
