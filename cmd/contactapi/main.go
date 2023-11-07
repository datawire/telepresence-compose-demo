package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/knlambert/telepresence-compose-demo/pkg/dao"
    _ "github.com/lib/pq" // initialize the Postgres driver
    "github.com/sethvargo/go-envconfig"
)

var d *dao.UserDAO

type AppConfig struct {
    Host     string `env:"DB_HOST"`
    Username string `env:"DB_USERNAME"`
    Password string `env:"DB_PASSWORD"`
    Dataset  string `env:"DATASET,default=local"`
    Port     int    `env:"PORT,default=8080"`
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

    err = d.InitContacts(appConfig.Dataset)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to init db: %w", err))
    }

    router := httprouter.New()
    router.GET("/contacts/:userId", GetContactFromUserId)

    router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        w.Write([]byte("OK"))
    })

    log.Printf("Contact API Listening on %d", appConfig.Port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", appConfig.Port), router))
}

func GetContactFromUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    userID := ps.ByName("userId")

    contact, err := d.GetContacts(ps.ByName("userId"))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error: failed to get contactapi for dao %s, %s", userID, err)
        return
    }

    output, err := json.Marshal(contact)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error: failed to marshal contactapi %s", err)
        return
    }

    log.Printf("Info: Returned contactapi for %s", userID)
    w.WriteHeader(http.StatusOK)
    w.Write(output)
}
