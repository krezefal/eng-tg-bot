package telebot

import (
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/logic"
	"github.com/krezefal/eng-tg-bot/internal/models"
	telebotServer "github.com/krezefal/eng-tg-bot/internal/server/telebot"
)

var _ telebotServer.Handlers = (*Handlers)(nil)

type Handlers struct {
	logic  models.Logic
	logger *zerolog.Logger
}

func New(logic *logic.MainLogic, logger *zerolog.Logger) *Handlers {
	return &Handlers{logic, logger}
}

func (h *Handlers) Start(c tele.Context) error {
	h.logic.Start(c)

	return nil
}

func (h *Handlers) Repeat(c tele.Context) error {
	h.logic.Repeat(c)

	return nil
}

func (h *Handlers) List(c tele.Context) error {
	h.logic.List(c)

	return nil
}

func (h *Handlers) Dict(c tele.Context) error {
	h.logic.Dict(c)

	return nil
}

func (h *Handlers) Subscribe(c tele.Context) error {
	h.logic.Subscribe(c)

	return nil
}

func (h *Handlers) Help(c tele.Context) error {
	h.logic.Help(c)

	return nil
}

func (h *Handlers) Remove(c tele.Context) error {
	h.logic.Remove(c)

	return nil
}
