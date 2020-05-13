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

func GetDoor(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	delay, err := impl.GetDoorDelay(uhppoted.GetDoorDelayRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	})

	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-door", fmt.Sprintf("Error retrieving delay for door %d from device %v", door, deviceID)),
			err
	}

	if delay == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-door", fmt.Sprintf("Error retrieving delay for door %d from device %v", door, deviceID)),
			fmt.Errorf("No response returned to request for all cards from device %v", deviceID)
	}

	control, err := impl.GetDoorControl(uhppoted.GetDoorControlRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	})

	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-door-control", fmt.Sprintf("Error retrieving door control for device %v, door %d", deviceID, door)),
			err
	}

	if control == nil {
		return http.StatusOK, nil, nil
	}

	reply := struct {
		Delay        uint8                 `json:"delay"`
		ControlState uhppoted.ControlState `json:"control"`
	}{
		Delay:        delay.Delay,
		ControlState: control.Control,
	}

	return http.StatusOK, &struct {
		Door interface{} `json:"door"`
	}{
		Door: reply,
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
