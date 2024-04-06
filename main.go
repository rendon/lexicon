package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lexicon/api"
	"lexicon/db"
	"lexicon/util"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WordStatus int

const (
	NewWord WordStatus = iota
	SavedWord
)
const (
	DictionaryApi = "dictionaryapi.com"
)

func defineName(lexicon *db.Lexicon, name string) error {
	res, err := lexicon.Find(name)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			res, err = saveDefinition(lexicon, name)
			if err != nil {
				return err
			}
			printLexeme(res, NewWord)
			return nil
		} else {
			return err
		}
	}
	printLexeme(res, SavedWord)
	return nil
}

func printLexeme(lexeme *db.Lexeme, status WordStatus) {
	label := labelName(status)
	// def := abridgeDefinition(lexeme.Definition)
	fmt.Printf("\n%s\t%s", lexeme.Name, label)
	fmt.Printf("\n%s\n", strings.Repeat("=", len(lexeme.Name)))

	fmt.Printf("Added on %s\n\n", util.FormatTime(lexeme.CreatedAt))

	var lex api.Lexeme
	if err := json.Unmarshal([]byte(lexeme.Definition), &lex); err != nil {
		log.Printf("Unable to parse definition: %s", err)
		return
	}
	for _, e := range lex.Entries {
		fmt.Printf("%s — %s\n", e.Headword.Text, e.GrammaticalFunction)
		for _, sd := range e.ShortDefinitions {
			fmt.Printf("\t%s\n", sd)
		}
		fmt.Println()
	}
	fmt.Println()
}

func labelName(status WordStatus) string {
	if status == NewWord {
		return "\tNew"
	}
	return ""
}

// abridgeDefinition limits the output to 16 lines.
func abridgeDefinition(def string) string {
	lines := strings.Split(def, "\n")
	maxLines := util.Min(len(lines), 20)
	return strings.Join(lines[0:maxLines], "\n")
}

func saveDefinition(lexicon *db.Lexicon, name string) (*db.Lexeme, error) {
	def, err := api.Define(name)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(def)
	if err != nil {
		log.Printf("Unable to serialize lexeme: %s", err)
		return nil, err
	}

	lexeme := &db.Lexeme{
		Name:       name,
		Definition: string(out),
		Source:     DictionaryApi,
		CreatedAt:  time.Now(),
	}
	if err := lexicon.Save(lexeme); err != nil {
		return nil, fmt.Errorf("unable to save %q: %s: ", name, err)
	}
	return lexeme, nil
}

// interactive launches an interactive session where the user can define as many words as needed.
func interactive(lexicon *db.Lexicon) {
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n> ")
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("scanner error: %s", err)
			continue
		}

		input := strings.TrimSpace(line)
		if strings.TrimSpace(input) == "!exit" {
			break
		}

		// Skip empty lines
		if len(input) == 0 {
			continue
		}

		if err := defineName(lexicon, input); err != nil {
			log.Printf("Unable to define %q: %s", input, err)
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	lexicon, err := db.NewLexicon()
	if err != nil {
		log.Fatalf("Failed to set up database: %s", err)
	}
	defer lexicon.Close()

	if len(os.Args) == 1 {
		// Launch the lexicon in interactive mode
		interactive(lexicon)
		return
	}
}
