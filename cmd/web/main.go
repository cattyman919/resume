package main

import (
	"log/slog"
	"os"

	"github.com/cattyman919/autocv/internal/web"
)

func main() {
	addr := ":3001"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	webApp, err := web.NewApp()
	if err != nil {
		slog.Error("Failed to create web app", "err", err)
		os.Exit(1)
	}

	slog.Info("AutoCV web server starting", "addr", addr)
	if err := webApp.ListenAndServe(addr); err != nil {
		slog.Error("Server error", "err", err)
		os.Exit(1)
	}
}
