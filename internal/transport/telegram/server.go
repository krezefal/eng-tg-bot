package telegram

import (
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

type Server struct {
	bot    *tele.Bot
	logger *zerolog.Logger
}

func NewServer(bot *tele.Bot, parentLogger *zerolog.Logger) *Server {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}
	logger := parentLogger.With().Str("component", "telegram_server").Logger()

	return &Server{bot, &logger}
}

func (t *Server) Start() {
	t.bot.Start()
}
