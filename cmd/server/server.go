package main

import (
	"log/slog"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8000", nil); err != nil {
		slog.Error("Failed to start HTTP server",
			"port", 8000,
			"error", err)
	}
}
