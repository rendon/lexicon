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
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

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

type PrintMode int

const (
	ShortDef PrintMode = iota
	FullDef
)

const entrySeparator = "---"

func defineName(lexicon *db.Lexicon, name string, waitTime time.Duration) error {
	res, err := lexicon.Find(name)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			// Wait time so as not to overload the service and get throttled/blocked.
			// TODO: Improve this hack for define-batch.
			time.Sleep(waitTime)

			res, err = saveDefinition(lexicon, name)
			if err != nil {
				return err
			}
			printLexeme(res, NewWord, ShortDef)
			return nil
		} else {
			return err
		}
	}
	printLexeme(res, SavedWord, ShortDef)

	return nil
}

func formatCognates(cognates []api.Cognate) string {
	var res string

	for i, cognate := range cognates {
		res += cognate.Label + " " + strings.Join(util.QuoteStrings(cognate.Targets), ",")
		if i+1 < len(cognates) {
			res += " | "
		}
	}
	return res
}

func printLexeme(lexeme *db.Lexeme, status WordStatus, printMode PrintMode) {
	var out = new(strings.Builder)
	_, _ = fmt.Fprintf(out, "\n")
	label := labelName(status)

	title := color.New(color.FgGreen, color.Bold)
	_, _ = fmt.Fprintf(out, "%s", title.Sprintf("\n%s", lexeme.Name))
	_, _ = fmt.Fprintf(out, "\t%s", color.RedString(label))

	_, _ = fmt.Fprintf(out, "%s", title.Sprintf("\n%s\n", strings.Repeat("=", len(lexeme.Name))))

	_, _ = fmt.Fprintf(out, "Added on %s\n", util.FormatTime(lexeme.CreatedAt))

	var lex api.Lexeme
	if err := json.Unmarshal([]byte(lexeme.Definition), &lex); err != nil {
		log.Printf("Unable to parse definition: %s", err)
		return
	}

	subtitle := color.New(color.FgBlue)
	for i, e := range lex.Entries {
		_, _ = fmt.Fprintf(out, "\n")

		gf := e.GrammaticalFunction
		if len(gf) > 0 {
			_, _ = fmt.Fprintf(out, "%s\n", subtitle.Sprintf("%s — %s", e.Headword.Text, gf))

			for _, sd := range e.ShortDefinitions {
				_, _ = fmt.Fprintf(out, "• %s\n", sd)
			}
		} else if len(e.Cognates) > 0 {
			cognates := formatCognates(e.Cognates)
			_, _ = fmt.Fprintf(out, "%s\n", subtitle.Sprintf("%s — %s", e.Headword.Text, cognates))
		}
		if printMode == FullDef {
			for _, d := range e.Definitions {
				printDefinition(out, d)
			}

			if len(e.Quotes) > 0 {
				_, _ = fmt.Fprintf(out, "\n%s\n", color.BlueString("Quotes"))
				for _, q := range e.Quotes {
					printQuote(out, q)
				}
			}

			if i+1 < len(lex.Entries) {
				_, _ = fmt.Fprintf(out, "%s\n", strings.Repeat("—", 80))
			}
		}
		_, _ = fmt.Fprintf(out, "%s\n", entrySeparator)
	}

	fmt.Print(abridgeOutput(out))
}

// abridgeOutput shortens the output so that it can fit nicely on a screen w/o scrolling.
func abridgeOutput(builder *strings.Builder) string {
	lines := strings.Split(builder.String(), "\n")
	var out strings.Builder
	for i, line := range lines {
		if line == entrySeparator {
			// Arbitrary number
			// TODO: Make this value configurable.
			if i > 25 {
				break
			}
			continue
		}
		_, _ = fmt.Fprintf(&out, "%s\n", line)
	}
	return out.String()
}

func printDefinition(out *strings.Builder, d api.Definition) {
	_, _ = fmt.Fprintf(out, "\n")
	if len(d.VerbDivider) > 0 {
		_, _ = fmt.Fprintf(out, "%s\n", d.VerbDivider)
	}
	for _, s := range d.Senses {
		_, _ = fmt.Fprintf(out, "%s\n", s.Text)

		if len(s.UsageNotes) > 0 {
			for _, u := range s.UsageNotes {
				_, _ = fmt.Fprintf(out, "  •%q\n", u)
			}
		}
		if len(s.VerbalIllustrations) > 0 {
			for _, i := range s.VerbalIllustrations {
				_, _ = fmt.Fprintf(out, "  • %q\n", i)
			}
			_, _ = fmt.Fprintf(out, "\n")
		}
	}
}

func printQuote(out *strings.Builder, q api.Quote) {
	_, _ = fmt.Fprintf(out, "  %q\n", q.Text)
	_, _ = fmt.Fprintf(out, "  %s, %s, %s\n\n", q.Source, q.Author, q.PublicationDate)
}

func labelName(status WordStatus) string {
	if status == NewWord {
		return "\tNew"
	}
	return ""
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

		input := strings.ToLower(strings.TrimSpace(line))
		if strings.TrimSpace(input) == "!exit" {
			break
		}

		// Skip empty lines
		if len(input) == 0 {
			continue
		}

		if err := defineName(lexicon, input, 0); err != nil {
			log.Printf("Unable to define %q: %s", input, err)
		}
	}
}

// defineBatch reads words from a file and defines all words in it. If the words contain a timestamp
// the createdAt and updatedAt timestamps are set to such timestamp. This command is useful for
// importing words from other sources while still keeping the original dates.
func defineBatch(lexicon *db.Lexicon) error {
	if len(os.Args) < 3 {
		return errors.New("missing file name")
	}
	f, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var failed []string
	lines := strings.Split(strings.TrimSpace(string(buf)), "\n")
	for _, line := range lines {
		log.Printf("%q", line)
		tokens := strings.Split(line, ",")
		name := strings.ToLower(tokens[0])
		waitTime := time.Millisecond * time.Duration(100+rand.Intn(2000))
		if err := defineName(lexicon, name, waitTime); err != nil {
			log.Printf("Define name for %q failed with error: %s", name, err)
			failed = append(failed, line)
			continue
		}

		// The word has a timestamp, we'll update the database timestamps
		if len(tokens) == 2 {
			timestamp, err := time.Parse("2006-01-02 15:04:05", tokens[1])
			if err != nil {
				log.Printf("Unable to parse timestamp: %s", err)
				return err
			}

			var lexeme = db.Lexeme{
				Name:      name,
				UpdatedAt: timestamp,
				CreatedAt: timestamp,
			}
			if err := lexicon.UpdateTimestamps(lexeme); err != nil {
				log.Printf("Unable to update %s's timestamps: %s", name, err)
			}
		}
	}

	log.Println()
	log.Printf("Successful definitions: %d", len(lines)-len(failed))
	log.Printf("Failed definitions: %d", len(failed))
	for _, f := range failed {
		log.Printf("• %s", f)
	}
	return nil
}

func wod(_ *db.Lexicon) error {
	now := time.Now()
	date := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	if len(os.Args) >= 4 && os.Args[2] == "--date" && len(os.Args[3]) > 0 {
		date = os.Args[3]
	}

	res, err := api.GetWod(date)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", res.Date, res.Word)
	return nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	lexicon, err := db.NewLexicon()
	if err != nil {
		log.Fatalf("Failed to set up database: %s", err)
	}
	defer lexicon.Close()

	if len(os.Args) <= 1 {
		// Launch the lexicon in interactive mode
		interactive(lexicon)
		return
	}

	command := os.Args[1]
	if command == "define-batch" {
		if err := defineBatch(lexicon); err != nil {
			log.Fatalf("define-batch failed with error: %q", err)
		}
	} else if command == "wod" {
		if err := wod(lexicon); err != nil {
			log.Fatalf("wod failed with error: %q", err)
		}
	}
}
