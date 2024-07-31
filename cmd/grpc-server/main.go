package main

import (
	"context"
	"log"

	"github.com/fatih/color"

	"github.com/neracastle/chat-server/internal/app"
)

func main() {
	ap := app.NewApp(context.Background())

	err := ap.Start()
	if err != nil {
		log.Fatal(color.RedString("failed to start app: %v", err))
	}
}
