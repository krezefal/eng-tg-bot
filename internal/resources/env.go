package resources

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"

	"github.com/krezefal/eng-tg-bot/pkg/log"
)

type Env struct {
	Token   string        `envconfig:"TOKEN" required:"true"`
	Timeout time.Duration `envconfig:"POLLING_TIMEOUT" default:"10s"`
	DSN     string        `envconfig:"DB_DSN" required:"true"`
}

func init() {
	if err := gotenv.Load(); err != nil {
		log.Logger.Error().Err(err).Msg("dotenv load error")
	}
}

func (r *Resources) initEnv() error {
	var e Env
	err := envconfig.Process("", &e)
	if err != nil {
		return err
	}

	r.Env = &e
	log.Logger.Info().Msg("init env success")

	return nil
}
