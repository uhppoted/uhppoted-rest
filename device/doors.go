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

func GetDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	}

	response, err := impl.GetDoorDelay(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-door-delay", fmt.Sprintf("Error retrieving door delay for device %v, door %d", deviceID, door))
	} else if response == nil {
		return nil, nil
	}

	return &struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: response.Delay,
	}, nil
}

func SetDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("Error reading request (%w)", err), deviceID, "set-door-delay", "Error reading request")
	}

	body := struct {
		Delay *uint8 `json:"delay"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), deviceID, "set-door-delay", "Error parsing request")
	}

	if body.Delay == nil {
		return nil, errors.Errorf(errors.InvalidDoorDelay, deviceID, "set-door-delay", "Missing/invalid door delay")
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
		Delay:    *body.Delay,
	}

	response, err := impl.SetDoorDelay(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "set-door-delay", "Error setting device door delay")
	} else if response == nil {
		return nil, nil
	}

	return &struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: response.Delay,
	}, nil
}

func GetDoorControl(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	}

	response, err := impl.GetDoorControl(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-door-control", fmt.Sprintf("Error retrieving door control for device %v, door %d", deviceID, door))
	} else if response == nil {
		return nil, nil
	}

	return &struct {
		ControlState uhppoted.ControlState `json:"control"`
	}{
		ControlState: response.Control,
	}, nil
}

func SetDoorControl(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("Error reading request (%w)", err), deviceID, "set-door-control", "Error reading request")
	}

	body := struct {
		Control *uhppoted.ControlState `json:"control"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), deviceID, "set-door-control", "Error parsing request")
	}

	if body.Control == nil {
		return nil, errors.Errorf(errors.InvalidDoorControl, deviceID, "set-door-control", "Missing/invalid door control")
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
		Control:  *body.Control,
	}

	response, err := impl.SetDoorControl(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "set-door-control", "Error setting device door control")
	} else if response == nil {
		return nil, nil
	}

	return &struct {
		Control uhppoted.ControlState `json:"control"`
	}{
		Control: response.Control,
	}, nil
	//	states := map[string]uint8{
	//		"normally open":   1,
	//		"normally closed": 2,
	//		"controlled":      3,
	//	}
	//
	//	lookup := map[uint8]string{
	//		1: "normally open",
	//		2: "normally closed",
	//		3: "controlled",
	//	}
	//
	//	deviceID := ctx.Value("device-id").(uint32)
	//	door := ctx.Value("door").(uint8)
	//
	//	blob, err := ioutil.ReadAll(r.Body)
	//	if err != nil {
	//		warn(ctx, deviceID, "set-door-control", err)
	//		http.Error(w, "Error reading request", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	body := struct {
	//		ControlState string `json:"control"`
	//	}{}
	//
	//	err = json.Unmarshal(blob, &body)
	//	if err != nil {
	//		warn(ctx, deviceID, "set-door-control", err)
	//		http.Error(w, "Invalid request format", http.StatusBadRequest)
	//		return
	//	} else if _, ok := states[body.ControlState]; !ok {
	//		warn(ctx, deviceID, "set-door-control", fmt.Errorf("Invalid request control value: '%s'", body.ControlState))
	//		http.Error(w, fmt.Sprintf("Invalid request value: '%s'", body.ControlState), http.StatusBadRequest)
	//		return
	//	}
	//
	//	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceID, door)
	//	if err != nil {
	//		warn(ctx, deviceID, "set-door-control", err)
	//		http.Error(w, "Error setting door control", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(deviceID, door, states[body.ControlState], state.Delay)
	//	if err != nil {
	//		warn(ctx, deviceID, "set-door-control", err)
	//		http.Error(w, "Error setting door control", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	response := struct {
	//		ControlState string `json:"control"`
	//	}{
	//		ControlState: lookup[result.ControlState],
	//	}
	//
	//	reply(ctx, w, response)
}
