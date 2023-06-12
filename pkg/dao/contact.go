package dao

import (
    "fmt"
    "strings"
)

func (d *UserDAO) InitContacts(dataset string) error {
    schema := `CREATE TABLE IF NOT EXISTS contacts (
        user_id  INTEGER NOT NULL,
        country TEXT NOT NULL
    );`

    _, err := d.db.Exec(schema)
    if err != nil {
        return fmt.Errorf("failed to init schema: %w", err)
    }

    var count int
    err = d.db.Get(&count, "SELECT count(*) FROM contacts")
    if err != nil {
        return fmt.Errorf("failed to count contacts: %w", err)
    }

    if count == 0 {
        var contacts = make([]string, 0)
        datasets := map[string][][]any{
            "staging": {{1, "CANADA"}, {2, "US"}, {3, "CANADA"}, {4, "US"}, {5, "SWEDEN"}, {6, "CANADA"}},
            "local":   {{1, "US"}},
        }

        for _, values := range datasets[dataset] {
            contacts = append(contacts, fmt.Sprintf("(%d, '%s')", values[0], values[1]))
        }

        _, err = d.db.Exec("INSERT INTO contacts (user_id, country) VALUES" + strings.Join(contacts, ", "))
        if err != nil {
            return fmt.Errorf("failed to insert dataset: %w", err)
        }
    }

    return nil
}

type Contact struct {
    Country string `db:"country" json:"country"`
}

func (d *UserDAO) GetContacts(userId string) (*Contact, error) {
    contact := Contact{}

    err := d.db.Get(&contact, "SELECT country FROM contacts WHERE user_id = $1 LIMIT 1", userId)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    return &contact, nil
}
