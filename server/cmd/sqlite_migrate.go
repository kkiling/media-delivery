package main

import (
	"database/sql"
	"fmt"

	"github.com/kkiling/goplatform/log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const dialect = "sqlite3"
const stateMigrations = "migrations/sqlite/state"
const mediaDeliveryMigrations = "migrations/sqlite"

func sqliteMigrate(logger log.Logger, sqliteDsn string) error {
	// Настройка Goose
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	db, err := sql.Open(dialect, sqliteDsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Получаем текущую версию базы данных
	current, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}
	logger.Infof("Current database version: %d", current)

	// Применяем все доступные миграции
	if err := goose.Up(db, stateMigrations); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}
	if err := goose.Up(db, mediaDeliveryMigrations); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}
	logger.Infof("New database version: %d", newVersion)
	logger.Infof("Migrations applied successfully")

	return nil
}
