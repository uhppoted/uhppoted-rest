package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
)

func GetCards(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.GetCardsRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.GetCards(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-cards", fmt.Sprintf("Error retrieving cards for device %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-cards", fmt.Sprintf("Error retrieving all cards from device %v", deviceID)),
			fmt.Errorf("No response returned to request for all cards from device %v", deviceID)
	}

	return http.StatusOK, &struct {
		Cards []uint32 `json:"cards"`
	}{
		Cards: response.Cards,
	}, nil
}

func DeleteCards(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
	}

	response, err := impl.DeleteCards(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-cards", fmt.Sprintf("Error deleting all cards for device %v", deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-cards", fmt.Sprintf("Error deleting all cards for device %v", deviceID)),
			fmt.Errorf("No response returned to request to delete all cards on device %v", deviceID)
	}

	if !response.Deleted {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-cards", fmt.Sprintf("Failed to delete all cards for device %v", deviceID)),
			fmt.Errorf("Request to delete all cards on device %v returned %v", deviceID, response.Deleted)
	}

	return http.StatusOK, nil, nil
}

func GetCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	rq := uhppoted.GetCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-card", fmt.Sprintf("Error retrieving card %v from device %v", cardNumber, deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("get-card", fmt.Sprintf("Error retrieving card %v from device %v", cardNumber, deviceID)),
			fmt.Errorf("No response returned to request for card %v from device %v", cardNumber, deviceID)
	}

	return http.StatusOK, struct {
		Card interface{} `json:"card"`
	}{
		Card: response.Card,
	}, nil
}

func PutCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-card", "Error reading request"),
			err
	}

	card := struct {
		CardNumber uint32                `json:"card-number"`
		From       *types.Date           `json:"start-date"`
		To         *types.Date           `json:"end-date"`
		Doors      map[uint8]interface{} `json:"doors"`
	}{}

	if err = json.Unmarshal(blob, &card); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("put-card", fmt.Sprintf("Error parsing request (%v)", err)),
			err
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: uhppoted.DeviceID(deviceID),
		Card: types.Card{
			CardNumber: cardNumber,
			From:       card.From,
			To:         card.To,
			Doors:      map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0},
		},
	}

	for k, v := range card.Doors {
		switch vv := v.(type) {
		case bool:
			if vv {
				rq.Card.Doors[k] = 1
			}
		case int:
			rq.Card.Doors[k] = vv
		case float64:
			rq.Card.Doors[k] = int(vv)
		}
	}

	if _, err = impl.PutCard(rq); err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-card", fmt.Sprintf("Error storing card %v to device %v", cardNumber, deviceID)),
			err
	}

	return http.StatusOK, nil, nil
}

func DeleteCard(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.DeleteCard(rq)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-card", fmt.Sprintf("Error deleting card %v from device %v", cardNumber, deviceID)),
			err
	} else if response == nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-card", fmt.Sprintf("Error deleting card %v from device %v", cardNumber, deviceID)),
			fmt.Errorf("No response returned to request to delete card %v from device %v", cardNumber, deviceID)
	}

	if !response.Deleted {
		return http.StatusInternalServerError,
			errors.NewRESTError("delete-card", fmt.Sprintf("Failed to delete card %v from device %v", cardNumber, deviceID)),
			fmt.Errorf("Request to delete card %v from device %v returned %v", cardNumber, deviceID, response.Deleted)
	}

	return http.StatusOK, nil, nil
}
