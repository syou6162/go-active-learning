package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/example"
	"time"
)

func CreateEntryTable(db *sql.DB) (sql.Result, error) {
	schema := `
CREATE TABLE IF NOT EXISTS entry (
  "id" SERIAL,
  "url" TEXT NOT NULL,
  "label" INT NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);
CREATE UNIQUE INDEX "url_idx_entry" ON entry ("url");
`
	return db.Exec(schema)
}

func InsertEntry(db *sql.DB, e *example.Example) (sql.Result, error) {
	now := time.Now()
	return db.Exec(`
INSERT INTO entry (url, label, created_at, updated_at) VALUES ($1, $2, $3, $4)
`, e.Url, e.Label, now, now)
}
