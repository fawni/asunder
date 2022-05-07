package database

import (
	"context"
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/x6r/asunder/internal/common"
)

type DB = bun.DB

// Entry is a struct the holds all necessary informations for a TOTP entry.
type Entry struct {
	bun.BaseModel `bun:"table:asunder"`

	ID       int    `bun:",pk,autoincrement"`
	Username string `survey:"username"`
	Issuer   string `survey:"issuer"`
	Secret   string `survey:"secret"`
}

// InitDB returns the databse or creates it if it does not exist.
func InitDB() (*DB, error) {
	if err := os.MkdirAll(common.PathAsunder, 0755); err != nil {
		return nil, err
	}
	sqlite, err := sql.Open(sqliteshim.ShimName, common.PathDB)
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

func GetEntries(db *DB, key []byte) ([]Entry, error) {
	var entries []Entry
	var ctx = context.Background()
	err := db.NewSelect().Model(&entries).OrderExpr("id ASC").Scan(ctx)
	if err != nil {
		return []Entry{}, err
	}
	for i, entry := range entries {
		entries[i].Username, err = Decrypt(key, entry.Username)
		if err != nil {
			return []Entry{}, err
		}
		entries[i].Issuer, err = Decrypt(key, entry.Issuer)
		if err != nil {
			return []Entry{}, err
		}
		entries[i].Secret, err = Decrypt(key, entry.Secret)
		if err != nil {
			return []Entry{}, err
		}
	}
	return entries, nil
}
