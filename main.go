package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lexicon/dictapi"
	"lexicon/lexapi"
	"lexicon/lexdb"
	"lexicon/types"
	"lexicon/util"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	_ "github.com/mattn/go-sqlite3"
)

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02 15:04:05"

const (
	newEntry      = 1
	existingEntry = 2
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

func getDefinition(name string, dictionary types.Dictionary) (*types.Lexeme, int, error) {
	res, err := dictionary.Find(name)
	if err != nil {
		if !errors.Is(err, types.NotFound) {
			return nil, 0, err
		}

		def, err := dictapi.Define(name)
		if err != nil {
			return nil, 0, err
		}

		defstr, err := util.Serialize(def)
		if err != nil {
			log.Printf("Unable to serialize %v: %s", def, err)
			return nil, 0, err
		}
		lex := types.Lexeme{
			Name:       name,
			Definition: string(defstr),
			Source:     DictionaryApi,
		}

		return &lex, newEntry, err
	}
	return res, existingEntry, nil
}

func defineName(name string, dictionary types.Dictionary) error {
	def, status, err := getDefinition(name, dictionary)
	if err != nil {
		return err
	}

	if status == newEntry {
		if err := dictionary.Save(def); err != nil {
			log.Printf("Unable to save %v: %s", def, err)
			return err
		}
	}

	printLexeme(def, status, ShortDef)
	return nil
}

func formatCognates(cognates []types.Cognate) string {
	var res string

	for i, cognate := range cognates {
		res += cognate.Label + " " + strings.Join(util.QuoteStrings(cognate.Targets), ",")
		if i+1 < len(cognates) {
			res += " | "
		}
	}
	return res
}

func getPronunciations(headword types.Headword) string {
	if len(headword.Pronunciations) == 0 {
		return ""
	}
	var prons []string
	for _, p := range headword.Pronunciations {
		prons = append(prons, p.Text)
	}
	return fmt.Sprintf("    (%s)", strings.Join(prons, ","))
}

func printLexeme(lexeme *types.Lexeme, nameStatus int, printMode PrintMode) {
	var out = new(strings.Builder)
	_, _ = fmt.Fprintf(out, "\n")
	label := labelName(nameStatus)

	title := color.New(color.FgGreen, color.Bold)
	_, _ = fmt.Fprintf(out, "%s", title.Sprintf("\n%s", lexeme.Name))
	_, _ = fmt.Fprintf(out, "\t%s", color.RedString(label))

	_, _ = fmt.Fprintf(out, "%s", title.Sprintf("\n%s\n", strings.Repeat("=", len(lexeme.Name))))

	_, _ = fmt.Fprintf(out, "Added on %s\n", util.FormatDateTime(lexeme.CreatedAt))

	var lex types.Definition
	if err := json.Unmarshal([]byte(lexeme.Definition), &lex); err != nil {
		log.Printf("Unable to parse definition: %s -> %s", err, lexeme.Definition)
		return
	}

	subtitle := color.New(color.FgBlue)
	for i, e := range lex.Entries {
		_, _ = fmt.Fprintf(out, "\n")

		gf := e.GrammaticalFunction
		if len(gf) > 0 {
			hw := e.Headword.Text
			prons := getPronunciations(e.Headword)
			_, _ = fmt.Fprintf(out, "%s\n", subtitle.Sprintf("%s — %s%s", hw, gf, prons))

			for _, sd := range e.ShortDefinitions {
				_, _ = fmt.Fprintf(out, "• %s\n", sd)
			}
		} else if len(e.Cognates) > 0 {
			cognates := formatCognates(e.Cognates)
			_, _ = fmt.Fprintf(out, "%s\n", subtitle.Sprintf("%s — %s", e.Headword.Text, cognates))
		}
		if printMode == FullDef {
			for _, d := range e.Defs {
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

func printDefinition(out *strings.Builder, d types.Def) {
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

func printQuote(out *strings.Builder, q types.Quote) {
	_, _ = fmt.Fprintf(out, "  %q\n", q.Text)
	_, _ = fmt.Fprintf(out, "  %s, %s, %s\n\n", q.Source, q.Author, q.PublicationDate)
}

func labelName(nameStatus int) string {
	if nameStatus == newEntry {
		return "\tNew"
	}
	return ""
}

// interactive launches an interactive session where the user can define as many words as needed.
func interactive(dictionary types.Dictionary) {
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

		if err := defineName(input, dictionary); err != nil {
			log.Printf("Unable to define %q: %s", input, err)
		}
	}
}

// defineBatch reads words from a file and defines all words in it. If the words contain a timestamp
// the createdAt and updatedAt timestamps are set to such timestamp. This command is useful for
// importing words from other sources while still keeping the original dates.
func defineBatch(dictionary types.Dictionary) error {
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
		tokens := strings.Split(line, ",")
		name := strings.ToLower(tokens[0])
		def, nameStatus, err := getDefinition(name, dictionary)
		if err != nil {
			log.Printf("Unable to define %q: %s", name, err)
			failed = append(failed, line)
			continue
		}

		// The word comes with a timestamp, we'll update the timestamps accordingly.
		if len(tokens) == 2 {
			timestamp, err := time.Parse(dateTimeFormat, tokens[1])
			if err != nil {
				log.Printf("Unable to parse timestamp: %s", err)
				return err
			}
			def.UpdatedAt = &timestamp
			def.CreatedAt = &timestamp
		}

		if err := dictionary.Save(def); err != nil {
			log.Printf("Unable to save %q: %s", name, err)
			failed = append(failed, line)
			continue
		}

		if nameStatus == newEntry {
			// Wait to avoid throttling.
			waitTime := time.Millisecond * time.Duration(100+rand.Intn(2000))
			time.Sleep(waitTime)
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

func getWod(date string) (*types.Wod, error) {
	u := fmt.Sprintf("https://rafaelrendon.io/wod/%s", url.PathEscape(date))
	res, err := http.Get(u)
	if err != nil {
		log.Printf("HTTP call to %s failed with error: %s", u, err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("word of the day for '%s' not found", date)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Unable to read response's body: %s", err)
		return nil, err
	}

	var wod types.Wod
	if err := json.Unmarshal(buf, &wod); err != nil {
		log.Printf("Unable to unmarshal response's body: %s", err)
		return nil, err
	}

	return &wod, nil
}

func printWod(date string) error {
	res, err := getWod(date)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", res.Date, res.Word)
	return nil
}

func wod() error {
	// Default date range (today's date)
	startDate := time.Now()
	endDate := startDate

	// Parse one or two dates
	var dates []time.Time
	for i := 2; i < 4 && i < len(os.Args); i++ {
		d, err := time.Parse(dateFormat, os.Args[i])
		if err != nil {
			return err
		}
		dates = append(dates, d)
	}

	if len(dates) > 0 {
		startDate = dates[0]
		endDate = dates[len(dates)-1]
	}

	for d := startDate; d.Unix() <= endDate.Unix(); d = d.Add(time.Hour * 24) {
		if err := printWod(util.FormatDate(d)); err != nil {
			return err
		}
	}
	return nil
}

func migrateToApi(lexicon *lexdb.Lexicon, dictionary types.Dictionary) error {
	lexemes, err := lexicon.All()
	if err != nil {
		return err
	}

	// register words with the API
	for _, lex := range lexemes {
		if err := dictionary.Save(lex); err != nil {
			log.Printf("Unable to create lexeme: %s", err)
		}
	}

	return nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	var dictionary types.Dictionary

	// TODO: read config from toml file.
	if os.Getenv("DATA_SOURCE_TYPE") == "API" {
		ac, err := lexapi.NewDictionary()
		if err != nil {
			log.Fatalf("Failed to set up API client: %s", err)
		}
		dictionary = ac
	} else {
		d, err := lexdb.NewDictionary()
		if err != nil {
			log.Fatalf("Failed to set up database: %s", err)
		}
		dictionary = d
	}
	defer dictionary.Close()

	if len(os.Args) <= 1 {
		// Launch the lexicon in interactive mode
		interactive(dictionary)
		return
	}

	command := os.Args[1]
	if command == "define-batch" {
		if err := defineBatch(dictionary); err != nil {
			log.Fatalf("define-batch failed with error: %q", err)
		}
	} else if command == "wod" {
		if err := wod(); err != nil {
			log.Fatalf("wod failed with error: %q", err)
		}
	} else if command == "migrate-to-api" {
		lex, err := lexdb.NewLexicon()
		if err != nil {
			log.Fatalf("Failed to set up database: %s", err)
		}
		if err := migrateToApi(lex, dictionary); err != nil {
			log.Fatalf("migrate-to-api failed with error: %s", err)
		}
	}
}
