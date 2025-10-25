package repositories

import "database/sql"

// Database определяет интерфейс для работы с базой данных
type Database interface {
	Execute(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}