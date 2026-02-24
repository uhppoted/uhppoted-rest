package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"

	"github.com/uhppoted/uhppoted-rest/errors"
	"github.com/uhppoted/uhppoted-rest/lib"
)

func GetDoor(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)
	door := ctx.Value(lib.Door).(uint8)

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
			fmt.Errorf("no response returned to request for all cards from device %v", deviceID)
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
		Delay        uint8              `json:"delay"`
		ControlState types.ControlState `json:"control"`
	}{
		Delay:        delay.Delay,
		ControlState: control.Control,
	}

	return http.StatusOK, &struct {
		Door any `json:"door"`
	}{
		Door: reply,
	}, nil
}

func SetDoorDelay(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)
	door := ctx.Value(lib.Door).(uint8)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-delay", "Error reading request"),
			err
	}

	body := struct {
		Delay *uint8 `json:"delay"`
	}{}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-delay", "Error parsing request"),
			err
	}

	if body.Delay == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-delay", "Missing/invalid door delay"),
			fmt.Errorf("missing/invalid door delay value in request body (%s)", string(blob))
	}

	if err := impl.SetDoorDelay(deviceID, door, *body.Delay); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-delay", "Error setting device door delay"),
			err
	}

	return http.StatusOK, &struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: *body.Delay,
	}, nil
}

func SetDoorControl(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)
	door := ctx.Value(lib.Door).(uint8)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-control", "Error reading request"),
			err
	}

	body := struct {
		Control *types.ControlState `json:"control"`
	}{}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-control", "Error parsing request"),
			err
	}

	if body.Control == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-control", "Missing/invalid door control"),
			fmt.Errorf("missing/invalid door control value in request body (%s)", string(blob))
	}

	if err := impl.SetDoorControl(deviceID, door, *body.Control); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-control", "Error setting device door control"),
			err
	}

	return http.StatusOK, &struct {
		Control types.ControlState `json:"control"`
	}{
		Control: *body.Control,
	}, nil
}

func SetDoorPasscodes(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)
	door := ctx.Value(lib.Door).(uint8)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-passcodes", "Error reading request"),
			err
	}

	body := struct {
		Passcodes []uint32 `json:"passcodes"`
	}{
		Passcodes: []uint32{},
	}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-door-passcodes", "Error parsing request"),
			err
	}

	if err := impl.SetDoorPasscodes(deviceID, door, body.Passcodes...); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-door-passcodes", "Error setting controller door passcodes"),
			err
	}

	return http.StatusOK, &struct {
		Passcodes []uint32 `json:"passcodes"`
	}{
		Passcodes: body.Passcodes,
	}, nil
}

func OpenDoor(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)
	door := ctx.Value(lib.Door).(uint8)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("open-door", "Error reading request"),
			err
	}

	body := struct {
		CardNumber *uint32 `json:"card-number"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("open-door", "Error parsing request"),
			err
	}

	if door < 1 || door > 4 {
		return http.StatusBadRequest,
			errors.NewRESTError("open-door", fmt.Sprintf("Invalid door (%v)", door)),
			fmt.Errorf("missing/invalid door in request (%v)", door)
	}

	if body.CardNumber == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("open-door", "Missing/invalid user ID"),
			fmt.Errorf("missing/invalid user ID in request body (%s)", string(blob))
	}

	if !authorized(ctx, *body.CardNumber) {
		return http.StatusUnauthorized,
			errors.NewRESTError("open-door", fmt.Sprintf("Not authorized for card %v", *body.CardNumber)),
			fmt.Errorf("not authorized for card %v", *body.CardNumber)
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: *body.CardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("open-door", fmt.Sprintf("Card %v not valid for device %v", *body.CardNumber, deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("open-door", fmt.Sprintf("Card %v not valid for device %v", *body.CardNumber, deviceID)),
			fmt.Errorf("function GetCard returned <nil> for card %v, device %v", *body.CardNumber, deviceID)
	} else {
		card := response.Card

		// Check start/end validity dates
		today := types.Date(time.Now())
		if card.From.IsZero() || card.To.IsZero() || today.Before(card.From) || today.After(card.To) {
			return http.StatusUnauthorized,
				errors.NewRESTError("open-door", fmt.Sprintf("Card %v is not valid for %v", card.CardNumber, today)),
				fmt.Errorf("card %v is not valid for %v", card, deviceID)
		}

		// Check door permissions
		if card.Doors[door] < 1 || card.Doors[door] > 254 {
			return http.StatusUnauthorized,
				errors.NewRESTError("open-door", fmt.Sprintf("Card %v is does not have permission for %v, door %v", card.CardNumber, deviceID, door)),
				fmt.Errorf("card %v is not valid for %v, door %v", card, deviceID, door)
		}

		// Check time profile
		if card.Doors[door] >= 2 && card.Doors[door] <= 254 {
			profileID := uint8(card.Doors[door])
			checked := map[uint8]bool{}

			for {
				profile, err := getTimeProfile(impl, deviceID, profileID)
				if err != nil {
					return http.StatusInternalServerError,
						errors.NewRESTError("open-door", fmt.Sprintf("Error retrieving time profile %v associated with card %v, door %v from device %v", profileID, *body.CardNumber, door, deviceID)),
						err
				}

				if profile == nil {
					return http.StatusInternalServerError,
						errors.NewRESTError("open-door", fmt.Sprintf("Failed to retrieve time profile %v associated with card %v, door %v from device %v", profileID, *body.CardNumber, door, deviceID)),
						fmt.Errorf("GetTimeProfile received <nil> response for time profile %v associated with card %v, door %v from device %v", profileID, *body.CardNumber, door, deviceID)
				}

				if err = checkTimeProfile(deviceID, *body.CardNumber, card.Doors[door], *profile); err == nil {
					break
				}

				if profile.LinkedProfileID < 2 || profile.LinkedProfileID > 254 || checked[profile.LinkedProfileID] {
					return http.StatusUnauthorized, errors.NewRESTError("open-door", fmt.Sprintf("%v", err)), err
				}

				checked[profileID] = true
				profileID = profile.LinkedProfileID
			}
		}
	}

	rqq := uhppoted.OpenDoorRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Door:     door,
	}

	result, err := impl.OpenDoor(rqq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("open-door", "Error opening door"),
			err
	} else if response == nil {
		return http.StatusOK, nil, nil
	}

	return http.StatusOK, &struct {
		Opened bool `json:"opened"`
	}{
		Opened: result.Opened,
	}, nil
}

func SetInterlock(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	deviceID := ctx.Value(lib.DeviceID).(uint32)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-interlock", "Error reading request"),
			err
	}

	body := struct {
		Interlock *types.Interlock `json:"interlock"`
	}{}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-interlock", "Error parsing request"),
			err
	}

	if body.Interlock == nil {
		return http.StatusBadRequest,
			errors.NewRESTError("set-interlock", "Missing door interlock mode"),
			fmt.Errorf("missing door interlock mode in request body (%s)", string(blob))
	}

	interlock := *body.Interlock

	if interlock != 0 && interlock != 1 && interlock != 2 && interlock != 3 && interlock != 4 && interlock != 8 {
		return http.StatusBadRequest,
			errors.NewRESTError("set-interlock", "Invalid door interlock mode"),
			fmt.Errorf("invalid door interlock mode in request body (%s)", interlock)
	}

	if err := impl.SetInterlock(deviceID, interlock); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("set-interlock", "Error setting controller door interlock"),
			err
	}

	return http.StatusOK, &struct {
		Interlock types.Interlock `json:"interlock"`
	}{
		Interlock: interlock,
	}, nil
}

func ActivateKeypads(impl uhppoted.IUHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, any, error) {
	controller := ctx.Value(lib.DeviceID).(uint32)

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("activate-keypads", "Error reading request"),
			err
	}

	body := struct {
		Keypads map[uint8]bool `json:"keypads"`
	}{
		Keypads: map[uint8]bool{
			1: false,
			2: false,
			3: false,
			4: false,
		},
	}

	if err := json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("activate-keypads", "Error parsing request"),
			err
	}

	keypads := map[uint8]bool{
		1: body.Keypads[1],
		2: body.Keypads[2],
		3: body.Keypads[3],
		4: body.Keypads[4],
	}

	if err := impl.ActivateKeypads(controller, keypads); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("activate-keypads", "Error activating controller reader keypads"),
			err
	}

	return http.StatusOK, &struct {
		Keypads map[uint8]bool `json:"keypads"`
	}{
		Keypads: keypads,
	}, nil
}

func getTimeProfile(impl uhppoted.IUHPPOTED, deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
	rq := uhppoted.GetTimeProfileRequest{
		DeviceID:  deviceID,
		ProfileID: profileID,
	}

	response, err := impl.GetTimeProfile(rq)
	if err != nil {
		return nil, err
	}

	return &response.TimeProfile, nil
}

func checkTimeProfile(deviceID, cardNumber uint32, profileID uint8, profile types.TimeProfile) error {
	now := types.NewHHmm(time.Now().Hour(), time.Now().Minute())
	today := types.Date(time.Now())

	if profile.From.IsZero() || profile.To.IsZero() || today.Before(profile.From) || today.After(profile.To) {
		return fmt.Errorf("card %v: time profile %v on device %v is not valid for %v", cardNumber, profileID, deviceID, today)
	}

	if !profile.Weekdays[today.Weekday()] {
		return fmt.Errorf("card %v: time profile %v on device %v is not authorized for %v", cardNumber, profileID, deviceID, today.Weekday())
	}

	for _, p := range []uint8{1, 2, 3} {
		if segment, ok := profile.Segments[p]; ok {
			if !segment.Start.After(now) && !segment.End.Before(now) {
				return nil
			}
		}
	}

	return fmt.Errorf("card %v: time profile %v on device %v is not authorized for %v", cardNumber, profileID, deviceID, now)
}
