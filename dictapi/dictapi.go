// dictapi implements an client for dictionaryapi.com.
package dictapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lexicon/types"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// GetDefinitionResult is a struct representation of the data returned by dictionaryapi.com.
// Generated with https://transform.tools/json-to-go using the JSON response from
// https://dictionaryapi.com/api/v3/references/collegiate/json/accept?key={key}
// JSON documentation: https://dictionaryapi.com/products/json

type DMeta struct {
	ID        string   `json:"id"`
	UUID      string   `json:"uuid"`
	Sort      string   `json:"sort"`
	Src       string   `json:"src"`
	Section   string   `json:"section"`
	Stems     []string `json:"stems"`
	Offensive bool     `json:"offensive"`
}

type MDef struct {
	Vd   string            `json:"vd"`
	Sseq [][][]interface{} `json:"sseq"`
}

type DQuote struct {
	T  string `json:"t"`
	Aq struct {
		Auth   string `json:"auth"`
		Source string `json:"source"`
		Aqdate string `json:"aqdate"`
	} `json:"aq"`
}
type DPrs struct {
	Mw    string `json:"mw"`
	Sound struct {
		Audio string `json:"audio"`
		Ref   string `json:"ref"`
		Stat  string `json:"stat"`
	} `json:"sound,omitempty"`
	L string `json:"l,omitempty"`
}

type DHwi struct {
	Hw  string `json:"hw"`
	Prs []DPrs `json:"prs"`
}

type DCxs struct {
	Cxl   string `json:"cxl"`
	Cxtis []struct {
		Cxt string `json:"cxt"`
	} `json:"cxtis"`
}

type DEntry struct {
	Meta     DMeta           `json:"meta"`
	Hwi      DHwi            `json:"hwi"`
	Cxs      []DCxs          `json:"cxs"`
	Fl       string          `json:"fl"`
	Def      []MDef          `json:"def"`
	Quotes   []DQuote        `json:"quotes,omitempty"`
	Et       [][]interface{} `json:"et,omitempty"`
	Date     string          `json:"date,omitempty"`
	Shortdef []string        `json:"shortdef"`
	Ins      []struct {
		If string `json:"if"`
	} `json:"ins,omitempty"`
	Uros []struct {
		Ure string `json:"ure"`
		Prs []struct {
			Mw    string `json:"mw"`
			Sound struct {
				Audio string `json:"audio"`
				Ref   string `json:"ref"`
				Stat  string `json:"stat"`
			} `json:"sound"`
		} `json:"prs"`
		Fl string `json:"fl"`
	} `json:"uros,omitempty"`
	Dros []struct {
		Drp string `json:"drp"`
		Vrs []struct {
			Vl string `json:"vl"`
			Va string `json:"va"`
		} `json:"vrs"`
		Def []struct {
			Sseq [][][]interface{} `json:"sseq"`
		} `json:"def"`
	} `json:"dros,omitempty"`
}

type GetDefinitionResult []DEntry

func getDictionaryApiKey() string {
	return os.Getenv("DICTIONARY_API_KEY")
}

var DefNotFound = errors.New("found no definitions")

func wodNotFound(date string) error {
	return fmt.Errorf("word of the day for '%s' not found", date)
}

func Define(name string) (*types.Definition, error) {
	key := getDictionaryApiKey()
	if len(key) == 0 {
		return nil, errors.New("missing API key")
	}

	// url.QueryEscape() vs url. PathEscape().
	// See https://stackoverflow.com/q/2678551/526189
	u := fmt.Sprintf(
		`https://dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s`,
		url.PathEscape(name), url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Got %s: %s", res.Status, body)
		return nil, fmt.Errorf("service returned %s: %s", res.Status, body)
	}

	ss := parseSpellingSuggestions(body)
	if len(ss) > 0 {
		message := "The word you've entered isn't in the dictionary. Spelling suggestions:\n"
		for _, e := range ss {
			message += fmt.Sprintf("• %s\n", e)
		}
		return nil, errors.New(message)
	}

	var data GetDefinitionResult
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, DefNotFound
	}
	var definition types.Definition
	for _, entry := range data {
		definition.Entries = append(definition.Entries, parseEntry(entry))
	}
	return &definition, nil
}

var client *http.Client

func post(u, name string) error {
	cookie := os.Getenv("MERRIAM_WEBSTER_COOKIE")
	if cookie == "" {
		return errors.New("missing Merriam-Webster cookie")
	}

	if client == nil {
		client = &http.Client{}
	}

	payload := fmt.Sprintf("word=%s&type=d", url.QueryEscape(name))
	req, err := http.NewRequest("POST", u, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	req.Header.Set("Accept", "text/javascript")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", cookie)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("service returned %s: %s", res.Status, body)
	}
	return nil
}

func Save(name string) error {
	u := "https://www.merriam-webster.com/lapi/v1/wordlist/save"
	return post(u, name)
}

func Remove(name string) error {
	u := "https://www.merriam-webster.com/lapi/v1/wordlist/delete"
	return post(u, name)
}

// parseSpellingSuggestions tries to parse spelling suggestions, which is an array of strings. If
// successful, it returns the list of suggestions, otherwise, it returns an empty array.
func parseSpellingSuggestions(buf []byte) []string {
	var suggestions []string
	if err := json.Unmarshal(buf, &suggestions); err != nil {
		return []string{}
	}
	return suggestions
}

func parseEntry(entry DEntry) types.Entry {
	return types.Entry{
		Meta:                parseMeta(entry.Meta),
		Headword:            parseHeadword(entry.Hwi),
		Cognates:            parseCognates(entry.Cxs),
		GrammaticalFunction: entry.Fl,
		ShortDefinitions:    entry.Shortdef,
		Defs:                parseDefinitions(entry.Def),
		Quotes:              parseQuotes(entry.Quotes),
	}
}

func parseCognates(cxs []DCxs) []types.Cognate {
	var cognates []types.Cognate
	for _, ref := range cxs {
		cognates = append(cognates, parseCognate(ref))
	}
	return cognates
}

func parseCognate(ref DCxs) types.Cognate {
	var targets []string
	for _, t := range ref.Cxtis {
		targets = append(targets, t.Cxt)
	}

	return types.Cognate{Label: ref.Cxl, Targets: targets}
}

func parseMeta(meta DMeta) types.Meta {
	return types.Meta{
		ID:        meta.ID,
		UUID:      meta.UUID,
		Sort:      meta.Sort,
		Source:    meta.Src,
		Section:   meta.Section,
		Stems:     meta.Stems,
		Offensive: meta.Offensive,
	}
}

func parseHeadword(hwi DHwi) types.Headword {
	return types.Headword{
		Text:           hwi.Hw,
		Pronunciations: parsePronunciations(hwi.Prs),
	}
}

func parsePronunciations(prs []DPrs) []types.Pronunciation {
	var res []types.Pronunciation
	for _, p := range prs {
		res = append(res, types.Pronunciation{
			Text:  p.Mw,
			Sound: p.Sound.Audio,
		})
	}
	return res
}

func parseDefinitions(defs []MDef) []types.Def {
	var res []types.Def
	for _, def := range defs {
		res = append(res, parseDefinition(def))
	}
	return res
}

// parseDefinition parses definitions.
func parseDefinition(def MDef) types.Def {
	return types.Def{
		VerbDivider: def.Vd,
		Senses:      parseSenses(def.Sseq),
	}
}

// Ref: https://dictionaryapi.com/products/json#sec-2.dt
func parseSenses(sseq [][][]interface{}) []types.Sense {
	var senses []types.Sense
	for _, a := range sseq {
		for _, b := range a {
			for _, c := range b {
				if _, ok := c.(map[string]interface{}); ok {
					var sense types.Sense
					if s, isMap := c.(map[string]interface{}); isMap {
						if s["dt"] == nil {
							continue
						}
						dt, isArray := s["dt"].([]interface{})
						if !isArray {
							continue
						}
						sense.Text = extractText(dt, "text")
						sense.UsageNotes = parseUsageNotes(dt)
						sense.VerbalIllustrations = parseVerbalIllustrations(dt)
					}
					senses = append(senses, sense)
				}
			}
		}
	}
	return senses
}

// extractText returns the first string in an array next to the "text" string.
func extractText(i interface{}, key string) string {
	if arr, isArray := i.([]interface{}); isArray {
		if len(arr) >= 2 {
			e0, ok0 := arr[0].(string)
			e1, ok1 := arr[1].(string)
			if ok0 && ok1 && e0 == key {
				return e1
			}
		}
		for _, a := range arr {
			if text := extractText(a, key); text != "" {
				return text
			}
		}
	}

	if mp, isMap := i.(map[string]interface{}); isMap {
		for k, v := range mp {
			if s, ok := v.(string); ok && k == key {
				return s
			}
			if text := extractText(v, key); text != "" {
				return text
			}
		}
	}
	return ""
}

func parseUsageNotes(dt []interface{}) []string {
	var notes []string
	for _, e := range dt {
		if x, isArray := e.([]interface{}); isArray && len(x) >= 2 && x[0] == "uns" {
			notes = append(notes, extractText(x, "text"))
		}
	}
	return notes
}

func parseVerbalIllustrations(dt []interface{}) []string {
	var res []string

	for _, e := range dt {
		if x, isArray := e.([]interface{}); isArray && len(x) >= 2 && x[0] == "vis" {
			res = append(res, extractText(x, "t"))
		}
	}
	return res
}

func parseQuotes(quotes []DQuote) []types.Quote {
	var res []types.Quote
	for _, q := range quotes {
		res = append(res, types.Quote{
			Text:            q.T,
			Author:          q.Aq.Auth,
			Source:          q.Aq.Source,
			PublicationDate: q.Aq.Aqdate,
		})
	}
	return res
}
