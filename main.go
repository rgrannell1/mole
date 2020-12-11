package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"github.com/docopt/docopt-go"
)

// findChromeHistory finds the location of the Chrome history file
func findChromeHistory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	candidates := []string{
		path.Join(home, ".config/google-chrome/Default/History"),
	}

	for _, candidate := range candidates {
		exists, err := fileExists(candidate)

		if err != nil {
			return "", err
		}

		if exists {
			return candidate, nil
		}
	}

	return "", errors.New("no candidate locations matched")
}

// fileExists check if a file exists
func fileExists(fpath string) (bool, error) {
	info, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		return false, nil
	}

	return !info.IsDir(), nil
}

func copy(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Scanner takes an SQL rows object
type Scanner func(rows *sql.Rows) ([]interface{}, error)

// Scanner2 takes an SQL rows object
type Scanner2 func(rows *sql.Rows) (map[string]interface{}, error)

func scanTableNames(rows *sql.Rows) ([]interface{}, error) {
	var name string
	err := rows.Scan(&name)

	if err != nil {
		return nil, err
	}

	return []interface{}{name}, err
}

func scan(rows *sql.Rows) (map[string]interface{}, error) {
	types, _ := rows.ColumnTypes()
	columns, _ := rows.Columns()

	row := make([]interface{}, len(types))
	valuePtrs := make([]interface{}, len(types))

	for ith := range types {
		valuePtrs[ith] = &row[ith]
	}

	err := rows.Scan(valuePtrs...)

	if err != nil {
		return nil, err
	}

	mapped := make(map[string]interface{})

	for ith, col := range columns {
		mapped[col] = row[ith]
	}

	return mapped, err
}

// emitHistory prints Chrome history as JSON
func emitHistory(opts docopt.Opts, fpath string) error {
	copyLocation := fpath + ".copy"
	err := copy(fpath, copyLocation)

	if err != nil {
		return err
	}

	knownTables := map[string]bool{
		"downloads":               true,
		"downloads_slices":        true,
		"downloads_url_chains":    true,
		"keyword_search_terms":    true,
		"meta":                    true,
		"segment_usage":           true,
		"segments":                true,
		"sqlite_sequence":         true,
		"typed_url_sync_metadata": true,
		"urls":                    true,
		"visit_source":            true,
		"visits":                  true,
	}

	// -- open the database
	db, err := sql.Open("sqlite3", copyLocation)
	if err != nil {
		return err
	}
	defer db.Close()

	tableName, optErr := opts.String("<tablename>")
	if optErr != nil {
		return optErr
	}
	if !knownTables[tableName] {
		return errors.New(tableName + " not supported.")
	}

	queries := map[string]string{
		"downloads":               "SELECT * FROM downloads",
		"downloads_slices":        "SELECT * FROM downloads_slices",
		"downloads_url_chains":    "SELECT * FROM downloads_url_chains",
		"keyword_search_terms":    "SELECT hidden,last_visit_time,normalized_term,term,title,typed_count,url,visit_count FROM keyword_search_terms LEFT JOIN urls ON keyword_search_terms.url_id = urls.id",
		"meta":                    "SELECT * FROM meta",
		"segment_usage":           "SELECT * FROM segment_usage LEFT JOIN segments ON segment_usage.segment_id = segments.id",
		"segments":                "SELECT name,url,title,visit_count,typed_count,last_visit_time FROM segments LEFT JOIN urls ON segments.url_id = urls.id",
		"sqlite_sequence":         "SELECT * FROM sqlite_sequence",
		"typed_url_sync_metadata": "SELECT * FROM typed_url_sync_metadata",
		"urls":                    "SELECT * FROM urls",
		"visit_source":            "SELECT * FROM visit_source",
		"visits":                  "SELECT * FROM visits",
	}

	// -- it does, print all data out. TODO use querystring instead.
	rows, readErrors := ReadSqlite2(queries[tableName], db, scan)

	for readErr := range readErrors {
		if readErr != nil {
			return readErr
		}
	}

	for elem := range rows {
		jsonBytes, err := json.Marshal(elem)
		if err != nil {
			return err
		}

		fmt.Println(string(jsonBytes))
	}

	return nil
}

// Mole is the core app
func Mole(args docopt.Opts) error {
	fpath, err := findChromeHistory()

	if err != nil {
		return err
	}

	return emitHistory(args, fpath)
}

// main is the CLI wrapper for Mole
func main() {
	usage := `Mole
Usage:
	mole [--db <string>]kb
	mole ls <tablename> [--db <string>]

Options:
	--db <string>    the Chrome database path. If not provided,

Author:
	Róisín Grannell <r.grannell2@gmail.com>
	`

	opts, _ := docopt.ParseDoc(usage)

	err := Mole(opts)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
