package resources

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
	"gorm.io/gorm"
)

type Resources struct {
	Env *Env
	Db  *gorm.DB
}

type Env struct {
	Token   string        `envconfig:"TOKEN" required:"true"`
	Timeout time.Duration `envconfig:"POLLING_TIMEOUT" default:"10s"`
}

func init() {
	if err := gotenv.Load(); err != nil {
		log.Error().Err(err).Msg("dotenv load error")
	}
}

func Get() *Resources {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal().Err(err).Msg("envconfig error")
	}

	return &Resources{
		Env: &env,
	}
}
