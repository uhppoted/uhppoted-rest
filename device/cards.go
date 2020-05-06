package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
)

func GetCards(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetCardsRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetCards(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-cards", fmt.Sprintf("Error retrieving cards for device %v", deviceID))
	} else if response == nil {
		return nil, errors.Errorf(errors.RequestFailed, deviceID, "get-cards", fmt.Sprintf("Error retrieving all cards from device %v", deviceID))
	}

	return &struct {
		Cards []uint32 `json:"cards"`
	}{
		Cards: response.Cards,
	}, nil
}

func DeleteCards(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.DeleteCards(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "delete-cards", fmt.Sprintf("Error deleting all cards for device %v", deviceID))
	} else if response == nil {
		return nil, errors.Errorf(errors.RequestFailed, deviceID, "delete-cards", fmt.Sprintf("Error deleting all cards for device %v", deviceID))
	}

	if !response.Deleted {
		return nil, errors.Errorf(err, deviceID, "delete-cards", fmt.Sprintf("Failed to delete all cards for device %v", deviceID))
	}

	return nil, nil
}

func GetCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	rq := uhppoted.GetCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "get-card", fmt.Sprintf("Error retrieving card %v from device %v", cardNumber, deviceID))
	} else if response == nil {
		return nil, errors.Errorf(errors.RequestFailed, deviceID, "get-card", fmt.Sprintf("Error retrieving card %v from device %v", cardNumber, deviceID))
	}

	return response.Card, nil
}

func PutCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("Error reading request (%w)", err), deviceID, "put-card", "Error reading request")
	}

	body := struct {
		From  *types.Date `json:"start-date"`
		To    *types.Date `json:"end-date"`
		Doors []bool      `json:"doors"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), deviceID, "put-card", "Error parsing request")
	}

	if body.From == nil {
		return nil, errors.Errorf(errors.InvalidDate, deviceID, "put-card", "Missing 'start-date'")
	}

	if body.To == nil {
		return nil, errors.Errorf(errors.InvalidDate, deviceID, "put-card", "Missing 'end-date'")
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Card: types.Card{
			CardNumber: cardNumber,
			From:       *body.From,
			To:         *body.To,
			Doors:      body.Doors,
		},
	}

	if _, err = impl.PutCard(rq); err != nil {
		return nil, errors.Errorf(err, deviceID, "put-card", fmt.Sprintf("Error storing card %v to device %v", cardNumber, deviceID))
	}

	return nil, nil
}

func DeleteCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.DeleteCard(rq)
	if err != nil {
		return nil, errors.Errorf(err, deviceID, "delete-card", fmt.Sprintf("Error deleting card %v from device %v", cardNumber, deviceID))
	} else if response == nil {
		return nil, errors.Errorf(errors.RequestFailed, deviceID, "delete-card", fmt.Sprintf("Error deleting card %v from device %v", cardNumber, deviceID))
	}

	if !response.Deleted {
		return nil, errors.Errorf(err, deviceID, "delete-card", fmt.Sprintf("Failed to delete card %v from device %v", cardNumber, deviceID))
	}

	return nil, nil
}
