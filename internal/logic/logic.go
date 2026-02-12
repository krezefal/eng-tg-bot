package logic

import (
	"strings"

	tele "gopkg.in/telebot.v4"
)

type Storage interface {
	Start(c tele.Context) error
	Repeat(c tele.Context) error

	List(c tele.Context) error
	Dict(c tele.Context) error
	Subscribe(c tele.Context) error

	Help(c tele.Context) error
	Remove(c tele.Context) error
}

type MainLogic struct {
	storage Storage
}

func New(storage Storage) *MainLogic {
	return &MainLogic{storage}
}

func (l *MainLogic) Start(c tele.Context) error {
	return c.Send("Start Message")
}

func (l *MainLogic) Repeat(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /repeat <vocabulary>")
	}
	return c.Send(strings.ToUpper(args[0]))
}

func (l *MainLogic) List(c tele.Context) error {
	return c.Send("List")
}

func (l *MainLogic) Dict(c tele.Context) error {
	return c.Send("Dict")
}

func (l *MainLogic) Subscribe(c tele.Context) error {
	return c.Send("Subscribe")
}

func (l *MainLogic) Help(c tele.Context) error {
	return c.Send("Help")
}

func (l *MainLogic) Remove(c tele.Context) error {
	return c.Send("Remove")
}
