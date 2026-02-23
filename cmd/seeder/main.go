package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"

	"github.com/krezefal/eng-tg-bot/pkg/log"
)

const serviceName = "seeder"

const (
	flagUpName   = "up"
	flagDownName = "down"
	flagFileName = "file"
	flagHelpName = "help"

	defaultSeedFile = "seeds/random_pool_a2_basic_50.json"
	envDBDSN        = "DB_DSN"
)

var logger = log.For(serviceName)

type seedDictionary struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
	Author      string `json:"author"`
}

type seedWord struct {
	Spelling      string `json:"spelling"`
	Transcription string `json:"transcription"`
	AudioLink     string `json:"audio"`
	RUTranslation string `json:"ru_translation"`
}

type seedData struct {
	Dictionary seedDictionary `json:"dictionary"`
	Words      []seedWord     `json:"words"`
}

func helpFn() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--up | --down] [--file path]\n\n", os.Args[0])
	flag.PrintDefaults()
}

// TODO: current realization is only for dicts with random_pool mode; add
// support of on_schedule mode
func main() {
	up := flag.Bool(flagUpName, false, "apply dictionary seed")
	down := flag.Bool(flagDownName, false, "rollback dictionary seed")
	filePath := flag.String(flagFileName, "", "path to seed JSON file")
	help := flag.Bool(flagHelpName, false, "show usage")
	flag.Usage = helpFn
	flag.Parse()

	if err := run(*help, *up, *down, *filePath); err != nil {
		logger.Fatal().Err(err).Msg("seeder run error")
	}
}

func run(help, up, down bool, filePath string) error {
	if help {
		flag.Usage()
		return nil
	}

	if up == down {
		return fmt.Errorf("set exactly one flag: --%s or --%s", flagUpName, flagDownName)
	}

	if strings.TrimSpace(filePath) == "" {
		logger.Warn().Msg("seed file path wasn't specified, using default path")
		filePath = defaultSeedFile
	}

	seed, err := loadSeed(filePath)
	if err != nil {
		return err
	}

	if err = validateSeed(seed); err != nil {
		return err
	}

	if err = gotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	dsn := strings.TrimSpace(os.Getenv(envDBDSN))
	if dsn == "" {
		return fmt.Errorf("env var %s is empty", envDBDSN)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Warn().Err(closeErr).Msg("db close error")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if up {
		if err = seedUp(ctx, db, seed); err != nil {
			return err
		}
		logger.Info().Str("dictionary", seed.Dictionary.Title).Msg("seed applied")

		return nil
	}

	removed, err := seedDown(ctx, db, seed.Dictionary.Mode, seed.Dictionary.Author)
	if err != nil {
		return err
	}
	logger.Info().
		Str("dictionary", seed.Dictionary.Title).
		Int64("rows_removed", removed).
		Msg("seed rolled back")

	return nil
}

func loadSeed(filePath string) (*seedData, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read seed file: %w", err)
	}

	var seed seedData
	if err = json.Unmarshal(raw, &seed); err != nil {
		return nil, fmt.Errorf("parse seed json: %w", err)
	}

	return &seed, nil
}

func validateSeed(seed *seedData) error {
	if seed == nil {
		return errors.New("seed is nil")
	}

	dict := seed.Dictionary
	if strings.TrimSpace(dict.Title) == "" {
		return errors.New("dictionary.title is required")
	}
	if strings.TrimSpace(dict.Mode) != "random_pool" {
		return errors.New("dictionary.mode must be random_pool")
	}
	if strings.TrimSpace(dict.Author) == "" {
		return errors.New("dictionary.author is required")
	}
	if len(seed.Words) == 0 {
		return errors.New("words must not be empty")
	}

	for i, w := range seed.Words {
		if strings.TrimSpace(w.Spelling) == "" {
			return fmt.Errorf("words[%d].spelling is required", i)
		}
		if strings.TrimSpace(w.RUTranslation) == "" {
			return fmt.Errorf("words[%d].ru_translation is required", i)
		}
	}

	return nil
}

func seedUp(ctx context.Context, db *sql.DB, seed *seedData) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	dictID, err := ensureDictionary(ctx, tx, seed.Dictionary)
	if err != nil {
		return err
	}

	const upsertWordQuery = `
		INSERT INTO dictionary_words (
			dictionary_id,
			batch_id,
			spelling,
			transcription,
			audio,
			ru_translation
		)
		VALUES ($1, NULL, $2, $3, $4, $5)
		ON CONFLICT (dictionary_id, spelling) DO UPDATE
		SET transcription = EXCLUDED.transcription,
			audio = EXCLUDED.audio,
			ru_translation = EXCLUDED.ru_translation,
			batch_id = NULL;
	`

	for _, w := range seed.Words {
		_, err = tx.ExecContext(
			ctx,
			upsertWordQuery,
			dictID,
			strings.TrimSpace(w.Spelling),
			strings.TrimSpace(w.Transcription),
			strings.TrimSpace(w.AudioLink),
			strings.TrimSpace(w.RUTranslation),
		)
		if err != nil {
			return fmt.Errorf("upsert word %q: %w", w.Spelling, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func ensureDictionary(ctx context.Context, tx *sql.Tx, dict seedDictionary) (string, error) {
	const selectQuery = `
		SELECT id
		FROM dictionaries
		WHERE mode = $1 AND author = $2
		ORDER BY created_at ASC
		LIMIT 1;
	`

	var dictID string
	err := tx.QueryRowContext(
		ctx,
		selectQuery,
		strings.TrimSpace(dict.Mode),
		strings.TrimSpace(dict.Author),
	).Scan(&dictID)

	if err == nil {
		return dictID, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("select dictionary: %w", err)
	}

	const insertQuery = `
		INSERT INTO dictionaries (title, description, mode, author)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		strings.TrimSpace(dict.Title),
		strings.TrimSpace(dict.Description),
		strings.TrimSpace(dict.Mode),
		strings.TrimSpace(dict.Author),
	).Scan(&dictID)
	if err != nil {
		return "", fmt.Errorf("insert dictionary: %w", err)
	}

	return dictID, nil
}

func seedDown(ctx context.Context, db *sql.DB, mode, author string) (int64, error) {
	const query = `
		DELETE FROM dictionaries
		WHERE mode = $1 AND author = $2;
	`

	res, err := db.ExecContext(ctx, query, strings.TrimSpace(mode), strings.TrimSpace(author))
	if err != nil {
		return 0, fmt.Errorf("delete seeded dictionaries: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rows, nil
}
