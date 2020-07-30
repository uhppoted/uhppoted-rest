package device

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func reply(ctx context.Context, w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(b) > 1024 && ctx.Value("compression") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		encoder := gzip.NewWriter(w)
		encoder.Write(b)
		encoder.Flush()
	} else {
		w.Write(b)
	}
}

func authorized(ctx context.Context, cardNumber uint32) bool {
	cards := ctx.Value("authorized-cards").([]string)
	c := fmt.Sprintf("%v", cardNumber)
	for _, re := range cards {
		if ok, err := regexp.MatchString(re, c); ok && err == nil {
			return true
		}
	}

	return false
}
