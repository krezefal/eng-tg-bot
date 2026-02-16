package resources

import (
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/krezefal/eng-tg-bot/pkg/log"
)

type Resources struct {
	Env *Env
	Db  *sql.DB
}

func MustGet() *Resources {
	r := &Resources{}

	if err := r.initEnv(); err != nil {
		log.Logger.Fatal().Err(err).Msg("init env error")
	}

	var group errgroup.Group
	group.Go(func() error {
		if err := r.initDb(); err != nil {
			return fmt.Errorf("init db: %w", err)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Logger.Fatal().Err(err).Msg("init resources error")
	}

	return r
}
