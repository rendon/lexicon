package lexdb

import (
	"database/sql"
	"errors"
	"fmt"
	"lexicon/types"
	"log"
	"os"
	"time"
)

type Lexicon struct {
	db *sql.DB
}

// NewDictionary returns a new Dictionary backed by an SQLite database.
func NewDictionary() (types.Dictionary, error) {
	sourceName := os.Getenv("DATA_SOURCE_NAME")
	if len(sourceName) == 0 {
		return nil, errors.New("missing data source name")
	}
	db, err := sql.Open("sqlite3", sourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to create database: %s", err)
	}
	return &Lexicon{db: db}, nil
}

// Find finds and returns a name in the database or returns error if the name
// does not exist in the database.
func (d *Lexicon) Find(name string) (*types.Lexeme, error) {
	rows, err := d.db.Query(`SELECT * FROM lexicon WHERE name = ?`, name)
	if err != nil {
		log.Printf("Unable to query the database: %s", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, types.NotFound
	}

	return readRecord(rows)
}

func (d *Lexicon) exists(name string) bool {
	rows, err := d.db.Query(`SELECT * FROM lexicon WHERE name = ?`, name)
	if err != nil {
		log.Printf("Unable to query the database: %s", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (d *Lexicon) selectRandom() (*types.Lexeme, error) {
	rows, err := d.db.Query(`SELECT * FROM lexicon ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		log.Printf("Unable to query lexicon table: %s", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, types.NotFound
	}

	return readRecord(rows)
}

// Save add lexeme to the database. Returns error if the operation fails.
func (d *Lexicon) Save(lexeme *types.Lexeme) error {
	timestamp := time.Now()
	if lexeme.CreatedAt == nil {
		lexeme.CreatedAt = &timestamp
	}
	if lexeme.UpdatedAt == nil {
		lexeme.UpdatedAt = &timestamp
	}

	stmt, err := d.db.Prepare(
		`INSERT INTO lexicon(name, definition, source, createdAt, updatedAt) values(?,?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("unable to insert prepare statement: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		lexeme.Name,
		lexeme.Definition,
		lexeme.Source,
		lexeme.CreatedAt.Unix(),
		lexeme.UpdatedAt.Unix(),
	)
	if err != nil {
		log.Printf("Unable to insert record: %s", err)
		return err
	}
	return nil
}

func (d *Lexicon) All() ([]*types.Lexeme, error) {
	rows, err := d.db.Query("SELECT * FROM lexicon")
	if err != nil {
		log.Printf("Unable to query lexicon table: %s", err)
		return nil, err
	}
	defer rows.Close()

	var all []*types.Lexeme
	for rows.Next() {
		lexeme, err := readRecord(rows)
		if err != nil {
			log.Printf("Unable to read record: %s", err)
			continue
		}
		all = append(all, lexeme)
	}
	return all, nil
}

// Close closes the connection to the database.
func (d *Lexicon) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			log.Printf("Failed to close the database: %s", err)
			return err
		}
	}
	return nil
}

func readRecord(rows *sql.Rows) (*types.Lexeme, error) {
	var name, def, source string
	var createdAt, updatedAt int64
	if err := rows.Scan(&name, &def, &source, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	cat := time.Unix(updatedAt, 0)
	uat := time.Unix(createdAt, 0)
	return &types.Lexeme{
		Name:       name,
		Definition: def,
		Source:     source,
		CreatedAt:  &cat,
		UpdatedAt:  &uat,
	}, nil
}
