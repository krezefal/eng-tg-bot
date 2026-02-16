package repository

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BaseStorage struct {
	db     *gorm.DB
	logger *zerolog.Logger
}
