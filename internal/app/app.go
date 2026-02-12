package app

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/logic"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v4"

	telebotHandlers "github.com/krezefal/eng-tg-bot/internal/handlers/telebot"
	"github.com/krezefal/eng-tg-bot/internal/resources"
	telebotServer "github.com/krezefal/eng-tg-bot/internal/server/telebot"
)

type App struct {
	logger     *zerolog.Logger
	telebotSrv *telebotServer.TelebotServer
}

func New(ctx context.Context, l *zerolog.Logger) (*App, error) {
	log.Info().Msg("initializing application")

	resources := resources.Get()

	pref := tele.Settings{
		Token:  resources.Env.Token,
		Poller: &tele.LongPoller{Timeout: resources.Env.Timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, errors.Wrap(err, "error init telebot")
	}

	log.Info().Msg("telebot instance initialized")

	logic := logic.New(nil)
	handlers := telebotHandlers.New(logic, l)
	telebotSrv := telebotServer.New(b, l)

	telebotSrv.InitRoutes(ctx, handlers)

	log.Info().Msg("telebot routes initialized")

	app := &App{
		logger:     l,
		telebotSrv: telebotSrv,
	}

	log.Info().Msg("application initialized")

	return app, nil
}

func (a *App) Start() {
	a.logger.Info().Msg("starting bot")
	a.telebotSrv.Start()
}
