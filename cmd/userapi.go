package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // initialize the Postgres driver
    "github.com/sethvargo/go-envconfig"
)

var (
    dao *DAO
)

func main() {
    // Loading env vars.
    var appConfig AppConfig
    if err := envconfig.Process(context.Background(), &appConfig); err != nil {
        log.Fatal(err)
    }

    var err error
    dao, err = newDAO(appConfig)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to create DAO: %w", err))
    }

    err = dao.Init(appConfig.Dataset)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to init db: %w", err))
    }

    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        users, err := dao.ListUsers()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("Error: failed to list users %s", err)
            return
        }

        output, err := json.Marshal(users)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("Error: failed to unmarshal users %s", err)
            return
        }

        log.Printf("Info: Return %d users", len(users))
        w.WriteHeader(http.StatusOK)
        w.Write(output)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    port := 8080
    log.Printf("Listening on %d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}

type AppConfig struct {
    Host     string `env:"DB_HOST"`
    Username string `env:"DB_USERNAME"`
    Password string `env:"DB_PASSWORD"`
    Dataset  string `env:"DATASET,default=local"`
}

func newDAO(config AppConfig) (*DAO, error) {
    chain := fmt.Sprintf(
        "host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode=disable",
        config.Host, 5432, config.Username, config.Password, "postgres",
    )

    db, err := sqlx.Open("postgres", chain)
    if err != nil {
        return nil, fmt.Errorf("failed to open chain: %w", err)
    }

    return &DAO{
        db: db,
    }, nil
}

type DAO struct {
    db *sqlx.DB
}

func (d *DAO) Init(dataset string) error {
    schema := `CREATE TABLE IF NOT EXISTS users (
        id              SERIAL PRIMARY KEY,
        name text NOT NULL UNIQUE
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
        datasets := map[string][]string{
            "staging": {"kevin", "jake", "guillaume", "nick", "thomas", "jose"},
            "local":   {"local-dummy-user"},
        }

        for _, name := range datasets[dataset] {
            users = append(users, fmt.Sprintf("('%s')", name))
        }

        _, err = d.db.Exec("INSERT INTO users (name) VALUES" + strings.Join(users, ", "))
        if err != nil {
            return fmt.Errorf("failed to insert dataset: %w", err)
        }
    }

    return nil
}

type User struct {
    ID   int    `db:"id" json:"id"`
    Name string `db:"name" json:"name"`
}

func (d *DAO) ListUsers() ([]*User, error) {
    var users []*User
    err := d.db.Select(&users, "SELECT id, name FROM users")
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    return users, nil
}
