package database

import (
	"context"
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/x6r/asunder/internal/config"
)

// Entry is a struct the holds all necessary informations for a TOTP entry.
type Entry struct {
	bun.BaseModel `bun:"table:asunder"`

	ID       int    `bun:",pk,autoincrement"`
	Username string `survey:"username"`
	Issuer   string `survey:"issuer"`
	Secret   string `survey:"secret"`
}

// InitDB returns the databse or creates it if it does not exist.
func InitDB() (*bun.DB, error) {
	if err := os.MkdirAll(config.PathAsunder, 0755); err != nil {
		return nil, err
	}
	sqlite, err := sql.Open(sqliteshim.ShimName, config.PathDB)
	if err != nil {
		return nil, err
	}
	sqlite.SetMaxOpenConns(1)
	db := bun.NewDB(sqlite, sqlitedialect.New())
	ctx := context.Background()
	if _, err := db.NewCreateTable().IfNotExists().Model((*Entry)(nil)).Exec(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
