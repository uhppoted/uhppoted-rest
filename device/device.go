package device

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"log"
	"net/http"
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

func debug(ctx context.Context, serialNumber uint32, operation string, r *http.Request) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-12d %-20s %v\n", serialNumber, operation, *r)
}

func warn(ctx context.Context, serialNumber uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", serialNumber, operation, err)
}
