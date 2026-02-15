package telebot

import (
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

type TelebotServer struct {
	telebot *tele.Bot
	logger  *zerolog.Logger
}

func New(telebot *tele.Bot, logger *zerolog.Logger) *TelebotServer {
	return &TelebotServer{telebot, logger}
}

func (t *TelebotServer) Start() {
	t.telebot.Start()
}
