package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"

    "github.com/knlambert/telepresence-compose-demo/pkg/dao"
    _ "github.com/lib/pq" // initialize the Postgres driver
    "github.com/sethvargo/go-envconfig"
)

var (
    d *dao.UserDAO
)

type AppConfig struct {
    Host          string `env:"DB_HOST"`
    Username      string `env:"DB_USERNAME"`
    Password      string `env:"DB_PASSWORD"`
    Dataset       string `env:"DATASET,default=local"`
    Port          int    `env:"PORT,default=8080"`
    ContactAPIURL string `env:"CONTACT_API_URL"`
}

func main() {
    // Loading env vars.
    var appConfig AppConfig
    if err := envconfig.Process(context.Background(), &appConfig); err != nil {
        log.Fatal(err)
    }

    var err error
    d, err = dao.NewUserDAO(appConfig.Host, appConfig.Username, appConfig.Password)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to create DAO: %w", err))
    }

    err = d.InitUsers(appConfig.Dataset)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to init db: %w", err))
    }

    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        users, err := d.ListUsers()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("Error: failed to list users %s", err)
            return
        }

        for i := range users {
            resp, err := http.Get(fmt.Sprintf("%s/contacts/%d", strings.TrimSuffix(appConfig.ContactAPIURL, "/"), users[i].ID))
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("Error: failed to get contact for user %d: %v", users[i].ID, err)
                return
            }

            body, err := io.ReadAll(resp.Body)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("Error: failed to read contact for user %d: %v", users[i].ID, err)
                return
            }

            contact := &dao.Contact{}
            err = json.Unmarshal(body, contact)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Printf("Error: failed to decode contact for user %d: %v", users[i].ID, err)
                return
            }

            users[i].Contacts = contact
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

    log.Printf("Listening on %d", appConfig.Port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", appConfig.Port), nil))
}
