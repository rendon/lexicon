package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

var NotFound = errors.New("NOT FOUND")

type Lexeme struct {
	Name       string    `db:"name" json:"name"`
	Definition string    `db:"definition" json:"definition"`
	Source     string    `db:"source" json:"source"`
	CreatedAt  time.Time `db:"createdAt" json:"created_at"`
	UpdatedAt  time.Time `db:"updatedAt" json:"updated_at"`
}

type Lexicon struct {
	db *sql.DB
}

// NewLexicon returns an instance of the Lexicon DAO.
func NewLexicon() (*Lexicon, error) {
	sourceName := getDataSourceName()
	if len(sourceName) == 0 {
		return nil, errors.New("missing data source name")
	}
	db, err := sql.Open("sqlite3", sourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to create database: %s", err)
	}
	return &Lexicon{db: db}, nil
}

func getDataSourceName() string {
	return os.Getenv("DATA_SOURCE_NAME")
}

func (d *Lexicon) Close() {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			log.Printf("Failed to close the database: %s", err)
		}
	}
}

func (d *Lexicon) Find(name string) (*Lexeme, error) {
	rows, err := d.db.Query(`SELECT * FROM lexicon WHERE name = ?`, name)
	if err != nil {
		log.Printf("Unable to query the database: %s", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, NotFound
	}

	return readRecord(rows)
}

func (d *Lexicon) Exists(name string) bool {
	rows, err := d.db.Query(`SELECT * FROM lexicon WHERE name = ?`, name)
	if err != nil {
		log.Printf("Unable to query the database: %s", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (d *Lexicon) SelectRandom() (*Lexeme, error) {
	rows, err := d.db.Query(`SELECT * FROM lexicon ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		log.Printf("Unable to query lexicon table: %s", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, NotFound
	}

	return readRecord(rows)
}

func readRecord(rows *sql.Rows) (*Lexeme, error) {
	var name, def, source string
	var createdAt, updatedAt int64
	if err := rows.Scan(&name, &def, &source, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	return &Lexeme{
		Name:       name,
		Definition: def,
		Source:     source,
		UpdatedAt:  time.UnixMilli(updatedAt),
		CreatedAt:  time.UnixMilli(createdAt),
	}, nil
}

func (d *Lexicon) Save(lexeme *Lexeme) error {
	stmt, err := d.db.Prepare(
		`INSERT INTO lexicon(name, definition, source, createdAt, updatedAt) values(?,?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("unable to insert prepare statement: %s", err)
	}
	defer stmt.Close()

	lexeme.UpdatedAt = time.Now()
	_, err = stmt.Exec(
		lexeme.Name, lexeme.Definition, lexeme.Source,
		lexeme.CreatedAt.UnixMilli(),
		lexeme.UpdatedAt.UnixMilli(),
	)
	if err != nil {
		log.Printf("Unable to insert record: %s", err)
		return err
	}
	return nil
}

// UpdateTimestamps updates an entry's timestamps.
// TODO: Generalize function --> Update()
func (d *Lexicon) UpdateTimestamps(lexeme Lexeme) error {
	stmt, err := d.db.Prepare(`UPDATE lexicon SET createdAt = ?, updatedAt = ? WHERE name = ?`)
	if err != nil {
		return fmt.Errorf("unable to prepare update statement: %s", err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(lexeme.CreatedAt.UnixMilli(), lexeme.UpdatedAt.UnixMilli(), lexeme.Name)
	return err
}

func (d *Lexicon) All() ([]*Lexeme, error) {
	rows, err := d.db.Query("SELECT * FROM lexicon")
	if err != nil {
		log.Printf("Unable to query lexicon table: %s", err)
		return nil, err
	}
	defer rows.Close()

	var all []*Lexeme
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
