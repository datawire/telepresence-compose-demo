package dao

import (
    "fmt"
    "strings"

    "github.com/jmoiron/sqlx"
)

type UserDAO struct {
    db *sqlx.DB
}

func NewUserDAO(host, username, password string) (*UserDAO, error) {
    db, err := DB(host, username, password)
    if err != nil {
        return nil, fmt.Errorf("failed to create DB: %w", err)
    }

    return &UserDAO{
        db: db,
    }, nil
}

func (d *UserDAO) InitUsers(dataset string) error {
    schema := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE
    );`

    _, err := d.db.Exec(schema)
    if err != nil {
        return fmt.Errorf("failed to init schema: %w", err)
    }

    var count int
    err = d.db.Get(&count, "SELECT count(*) FROM users")
    if err != nil {
        return fmt.Errorf("failed to count users: %w", err)
    }

    if count == 0 {
        var users = make([]string, 0)
        datasets := map[string][][]any{
            "staging": {{1, "kevin"}, {2, "jake"}, {3, "guillaume"}, {4, "nick"}, {5, "thomas"}, {6, "jose"}},
            "local":   {{1, "local-dummy-dao"}},
        }

        for _, values := range datasets[dataset] {
            users = append(users, fmt.Sprintf("(%d, '%s')", values[0], values[1]))
        }

        _, err = d.db.Exec("INSERT INTO users (id, name) VALUES" + strings.Join(users, ", "))
        if err != nil {
            return fmt.Errorf("failed to insert dataset: %w", err)
        }
    }

    return nil
}

type User struct {
    ID       int      `db:"id" json:"id"`
    Name     string   `db:"name" json:"name"`
    Contacts *Contact `json:"contact"`
}

func (d *UserDAO) ListUsers() ([]*User, error) {
    var users []*User
    err := d.db.Select(&users, "SELECT id, name FROM users")
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    return users, nil
}
