package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/util"
)

type repository struct {
	db *sql.DB
}

func GetDataSourceName() string {
	host := util.GetEnv("POSTGRES_HOST", "localhost")
	dbUser := util.GetEnv("DB_USER", "nobody")
	dbPassword := util.GetEnv("DB_PASSWORD", "nobody")
	dbName := util.GetEnv("DB_NAME", "go-active-learning")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbUser, dbPassword, dbName,
	)
}

func New() (*repository, error) {
	db, err := sql.Open("postgres", GetDataSourceName())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	return &repository{db: db}, nil
}

func (r *repository) Ping() error {
	return r.db.Ping()
}

func (r *repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	} else {
		return nil
	}
}
