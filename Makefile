BIN_DIR := ./bin

.PHONY: migrator seeder migrate-up migrate-down seed-up seed-down

migrator:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/migrator ./cmd/migrator

seeder:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/seeder ./cmd/seeder

migrate-up: migrator
	$(BIN_DIR)/migrator --up

migrate-down: migrator
	$(BIN_DIR)/migrator --down

seed-up: seeder
	$(BIN_DIR)/seeder --up

seed-down: seeder
	$(BIN_DIR)/seeder --down
