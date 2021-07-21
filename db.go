package main

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const dbQuery = "select extension, name from users;"

var excluded = []string{"fax", "test", "ata", "panic"}

type Entry struct {
	Extension string
	User      string
}

func filterAndSort(list []*Entry) []*Entry {
	var out []*Entry

outer:
	for _, e := range list {
		for _, ex := range excluded {
			if strings.Contains(strings.ToLower(e.User), ex) {
				continue outer
			}
		}
		out = append(out, e)
	}

	sort.Slice(out, func(i, j int) bool {
		ii, ierr := strconv.Atoi(out[i].Extension)
		ji, jerr := strconv.Atoi(out[j].Extension)
		if ierr != nil && jerr != nil {
			return out[i].Extension < out[j].Extension
		}
		if ierr != nil {
			return false
		}
		if jerr != nil {
			return true
		}

		return ii < ji
	})
	return out
}

func GetEntries(dsn string) ([]*Entry, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(dbQuery)
	if err != nil {
		return nil, fmt.Errorf("could not query db: %w", err)
	}
	defer rows.Close()

	var entries []*Entry

	for rows.Next() {
		e := new(Entry)
		if err = rows.Scan(&(e.Extension), &(e.User)); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not finish iterating rows: %w", err)
	}

	return filterAndSort(entries), nil
}
