package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
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
		return nil, nil
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
		return nil, nil
	}

	if !response.Deleted {
		return nil, errors.Errorf(err, deviceID, "delete-cards", fmt.Sprintf("Failed to delete all cards for device %v", deviceID))
	}

	return nil, nil
}

func GetCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByID(deviceID, cardNumber)
	if err != nil {
		warn(ctx, deviceID, "get-card-by-id", err)
		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
		return
	}

	if card == nil {
		http.Error(w, "Card record does not exist", http.StatusNotFound)
		return
	}

	response := struct {
		Card types.Card `json:"card"`
	}{
		Card: *card,
	}

	reply(ctx, w, response)
}

func PutCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, deviceID, "put-card", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		From  *types.Date `json:"start-date"`
		To    *types.Date `json:"end-date"`
		Doors []bool      `json:"doors"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		warn(ctx, deviceID, "put-card", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if body.From == nil {
		warn(ctx, deviceID, "put-card", fmt.Errorf("Missing 'start-date'"))
		http.Error(w, "Invalid request: missing 'start-date'", http.StatusBadRequest)
		return
	}

	if body.To == nil {
		warn(ctx, deviceID, "put-card", fmt.Errorf("Missing 'end-date'"))
		http.Error(w, "Invalid request: missing 'end-date'", http.StatusBadRequest)
		return
	}

	card := types.Card{
		CardNumber: cardNumber,
		From:       *body.From,
		To:         *body.To,
		Doors:      body.Doors,
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).PutCard(deviceID, card)
	if err != nil {
		warn(ctx, deviceID, "put-card", err)
		http.Error(w, "Error adding/updating card", http.StatusInternalServerError)
		return
	}

	if !result.Succeeded {
		warn(ctx, deviceID, "put-card", fmt.Errorf("Request failed"))
		http.Error(w, "Error adding/updating card", http.StatusInternalServerError)
		return
	}
}

func DeleteCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(deviceID, cardNumber)
	if err != nil {
		warn(ctx, deviceID, "delete-card-by-id", err)
		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
		return
	}

	if !result.Succeeded {
		warn(ctx, deviceID, "delete-card", fmt.Errorf("Request failed"))
		http.Error(w, "Error deleting card", http.StatusInternalServerError)
		return
	}
}
