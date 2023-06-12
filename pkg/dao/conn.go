package dao

import (
    "fmt"

    "github.com/jmoiron/sqlx"
)

var _db *sqlx.DB

func DB(host, username, password string) (*sqlx.DB, error) {
    if _db == nil {
        chain := fmt.Sprintf(
            "host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode=disable",
            host, 5432, username, password, "postgres",
        )

        var err error
        _db, err = sqlx.Open("postgres", chain)

        if err != nil {
            return nil, fmt.Errorf("failed to open chain: %w", err)
        }
    }
    return _db, nil
}
