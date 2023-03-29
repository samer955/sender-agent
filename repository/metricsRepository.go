package repository

import "database/sql"

type MetricsRepository interface {
}

type MetricsRepositoryImpl struct {
	db *sql.DB
}
