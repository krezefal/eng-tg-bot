package models

import tele "gopkg.in/telebot.v4"

type Logic interface {
	Start(c tele.Context) error
	Repeat(c tele.Context) error

	List(c tele.Context) error
	Dict(c tele.Context) error
	Subscribe(c tele.Context) error

	Help(c tele.Context) error
	Remove(c tele.Context) error
}
