package app

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/repository/postgres"
	"github.com/krezefal/eng-tg-bot/internal/resources"
	"github.com/krezefal/eng-tg-bot/internal/transport/telegram"
	"github.com/krezefal/eng-tg-bot/internal/usecase"
)

type App struct {
	logger *zerolog.Logger
	tgSrv  *telegram.Server
}

func New(ctx context.Context, logger *zerolog.Logger) (*App, error) {
	resources := resources.MustGet()

	logger.Info().Msg("initializing application...")

	pref := tele.Settings{
		Token:  resources.Env.Token,
		Poller: &tele.LongPoller{Timeout: resources.Env.Timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("error init telebot: %w", err)
	}

	logger.Info().Msg("telebot instance initialized")

	userRepo := postgres.NewUserRepo(resources.Db, logger)
	// TODO: uncomment
	//dictRepo := repository.NewDictionaryRepo(resources.Db, logger)
	//subsRepo := repository.NewSubscriptionRepo(resources.Db, logger)
	//wordStateRepo := repository.NewWordStateRepo(resources.Db, logger)

	onboardUC := usecase.NewOnboardingUsecase(userRepo, logger)
	// TODO: uncomment
	//catalogUC := usecase.NewCatalogUsecase(dictRepo, logger)
	//subscUC := usecase.NewSubscriptionUsecase(userRepo, subsRepo, logger)
	//learningUC := usecase.NewLearningUsecase(subsRepo, wordStateRepo, logger)
	//reviewUC := usecase.NewReviewUsecase(wordStateRepo, logger)

	handlers := telegram.NewHandler(
		onboardUC,
		nil,
		nil,
		nil,
		nil,
		logger,
	)

	tgSrv := telegram.NewServer(b, logger)
	tgSrv.InitRoutes(ctx, handlers)

	app := &App{
		logger: logger,
		tgSrv:  tgSrv,
	}

	logger.Info().Msg("application initialized")

	return app, nil
}

func (a *App) Start() {
	a.logger.Info().Msg("starting bot")
	a.tgSrv.Start()
}
