package types

import (
	"errors"
	"time"
)

var NotFound = errors.New("Not fond")

// Quote represents a quote.
type Quote struct {
	Text            string `json:"text,omitempty"`            // t
	Author          string `json:"author,omitempty"`          // auth
	Source          string `json:"source,omitempty"`          // source
	PublicationDate string `json:"publicationDate,omitempty"` // aqdate
}

type Sense struct {
	Number              string   `json:"number,omitempty"`              // sn
	Text                string   `json:"text,omitempty"`                // dt
	UsageNotes          []string `json:"usageNotes,omitempty"`          // uns
	VerbalIllustrations []string `json:"verbalIllustrations,omitempty"` // vis
}

// Def describes a definition in a lexicon entry.
type Def struct {
	VerbDivider string  `json:"verbDivider,omitempty"` // vd
	Senses      []Sense `json:"senses,omitempty"`      // sseq
}

type Meta struct {
	ID        string   `json:"id,omitempty"`
	UUID      string   `json:"uuid,omitempty"`
	Sort      string   `json:"sort,omitempty"`
	Source    string   `json:"source,omitempty"`
	Section   string   `json:"section,omitempty"`
	Stems     []string `json:"stems,omitempty"`
	Offensive bool     `json:"offensive,omitempty"`
}

type Pronunciation struct {
	Text  string `json:"text,omitempty"`
	Sound string `json:"sound,omitempty"`
}

type Headword struct {
	Text           string          `json:"text,omitempty"`
	Pronunciations []Pronunciation `json:"pronunciations,omitempty"`
}

// Cognate when a headword is a less common spelling of another word with there
// same meaning, there will be a cognate cross-reference pointing to the headwordwith
// with the more common spelling.
type Cognate struct {
	Label   string   `json:"label"`
	Targets []string `json:"targets"`
}

// Entry represents a meaning intended or conveyed.
type Entry struct {
	Meta                Meta      `json:"meta,omitempty"`
	Headword            Headword  `json:"headword,omitempty"`
	Cognates            []Cognate `json:"cognates"`
	GrammaticalFunction string    `json:"grammaticalFunction,omitempty"`
	ShortDefinitions    []string  `json:"shortDefinitions,omitempty"`
	Defs                []Def     `json:"defs,omitempty"`
	Quotes              []Quote   `json:"quotes,omitempty"`
	Etymology           []string  `json:"etymology,omitempty"`
}

// Lexeme represents a linguistic unit.
type Definition struct {
	Entries []Entry `json:"entries,omitempty"`
}

// Wod represents a Word of the Day.
type Wod struct {
	Date string `json:"date"`
	Word string `json:"word"`
}

type Lexeme struct {
	Name       string     `db:"name" json:"name"`
	Definition string     `db:"definition" json:"definition"`
	Source     string     `db:"source" json:"source"`
	CreatedAt  *time.Time `db:"createdAt" json:"created_at"`
	UpdatedAt  *time.Time `db:"updatedAt" json:"updated_at"`
}

// Dictionary defines the operations that every dictionary must implement.
type Dictionary interface {
	Find(name string) (*Lexeme, error)
	Save(lexeme *Lexeme) error
	Close() error
}
