package telebot

import (
	"context"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

// TODO: add comments
type Handlers interface {
	Start(c tele.Context) error
	Repeat(c tele.Context) error

	List(c tele.Context) error
	Dict(c tele.Context) error
	Subscribe(c tele.Context) error

	Help(c tele.Context) error
	Remove(c tele.Context) error
}

type TelebotServer struct {
	telebot *tele.Bot
	logger  *zerolog.Logger
}

func New(telebot *tele.Bot, logger *zerolog.Logger) *TelebotServer {
	return &TelebotServer{telebot, logger}
}

func (t *TelebotServer) InitRoutes(ctx context.Context, handlers Handlers) {
	t.telebot.Handle("/start", handlers.Start)
	t.telebot.Handle("/repeat", handlers.Repeat)

	t.telebot.Handle("/list", handlers.List)
	t.telebot.Handle("/dict", handlers.Dict)
	t.telebot.Handle("/subscribe", handlers.Subscribe)

	t.telebot.Handle("/help", handlers.Help)
	t.telebot.Handle("/removeAll", handlers.Remove)
}

func (t *TelebotServer) Start() {
	t.telebot.Start()
}
