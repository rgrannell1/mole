package main

import (
	"database/sql"
	"log"
)

// ReadSqlite reads data
func ReadSqlite(query string, db *sql.DB, scan Scanner) (chan []interface{}, chan error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan []interface{}, 100_000)
	fatals := make(chan error)

	go func() {
		defer close(out)
		defer close(fatals)
		defer rows.Close()

		for rows.Next() {
			entry, err := scan(rows)

			if err != nil {
				fatals <- err
				return
			}

			out <- entry
		}

		err = rows.Err()
		if err != nil {
			fatals <- err
			return
		}
	}()

	return out, fatals
}

// ReadSqlite2 reads data
func ReadSqlite2(query string, db *sql.DB, scan Scanner2) (chan map[string]interface{}, chan error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan map[string]interface{}, 100_000)
	fatals := make(chan error)

	go func() {
		defer close(out)
		defer close(fatals)
		defer rows.Close()

		for rows.Next() {
			entry, err := scan(rows)

			if err != nil {
				fatals <- err
				return
			}

			out <- entry
		}

		err = rows.Err()
		if err != nil {
			fatals <- err
			return
		}
	}()

	return out, fatals
}
