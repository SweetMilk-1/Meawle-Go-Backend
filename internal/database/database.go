package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// Database представляет подключение к SQLite базе данных
type Database struct {
	DB *sql.DB
}

// New создает новое подключение к SQLite базе данных
func New(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s", dbPath)

	return &Database{DB: db}, nil
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Execute выполняет SQL запрос без возврата данных (CREATE, INSERT, UPDATE, DELETE)
func (d *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return result, nil
}

// Query выполняет SQL запрос с возвратом данных (SELECT)
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	return rows, nil
}

// QueryRow выполняет SQL запрос с возвратом одной строки
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

// RunMigrations выполняет миграции базы данных
func (d *Database) RunMigrations(migrationsPath string) error {
	driver, err := sqlite3.WithInstance(d.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Выполняем миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}
