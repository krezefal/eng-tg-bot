package main

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/app"
	"github.com/krezefal/eng-tg-bot/pkg/log"
)

const serviceName = "eng-tg-bot"

// TODO: add linters
// TODO: add unit tests
// TODO: add metrics
func main() {
	ctx := context.Background()
	zerolog := log.For(serviceName)

	app, err := app.New(ctx, zerolog)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("error init app")
	}

	// TODO: context cancellation?
	app.Start()

	zerolog.Info().Msg("shutdown app")
}
