package resources

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/krezefal/eng-tg-bot/pkg/log"
)

const dbPingTimeout = 5 * time.Second

func (r *Resources) initDb() error {
	const op = "resources.initDb"

	db, err := sql.Open("postgres", r.Env.DSN)
	if err != nil {
		return fmt.Errorf("%s: open db: %w", op, err)
	}

	// TODO: adjust db conn settings (idle conn, max conn, etc...)

	// TODO: move to readiness probe in app.go
	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("%s: ping db: %w", op, err)
	}

	r.Db = db
	log.Logger.Info().Msg("init db connection success")

	return nil
}
