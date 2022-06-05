package main

import (
	"context"
	"jira-bot/internal/app"
	"log"
)

func main() {
	if err := app.NewApplication().Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
