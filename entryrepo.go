package main

import (
	"database/sql"
)

const MaxEntriesPerPage = 10

type EntryRepo struct {
	db *sql.DB
}

func (repo *EntryRepo) AddEntry(name, email, message string) error {
	_, err := repo.db.Exec(`INSERT INTO "entry" ("name", "email", "message") VALUES ($1, $2, $3)`, name, email, message)
	return err
}

func (repo *EntryRepo) CountEntries() (int, error) {
	row := repo.db.QueryRow(`SELECT COUNT(*) FROM "entry"`)
	if row.Err() != nil {
		return 0, row.Err()
	} else {
		var count int
		if err := row.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}
}

func (repo *EntryRepo) ListEntries(page int) ([]Entry, error) {
	var entries []Entry
	if rows, err := repo.db.Query(` SELECT "id", "name", "email", "message", "posted"
                                          FROM "entry"
                                          ORDER BY posted DESC
                                          LIMIT $1 OFFSET $2`,
		MaxEntriesPerPage,
		(page-1)*MaxEntriesPerPage); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var entry Entry
			if err := rows.Scan(&entry.ID, &entry.Name, &entry.Email, &entry.Message, &entry.Posted); err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}
	}
	return entries, nil
}
