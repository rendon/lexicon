package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func getDictionaryApiKey() string {
	return os.Getenv("DICTIONARY_API_KEY")
}

func Define(name string) (*Lexeme, error) {
	key := getDictionaryApiKey()
	if len(key) == 0 {
		return nil, errors.New("missing API key")
	}
	u := fmt.Sprintf(
		`https://dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s`,
		url.QueryEscape(name), key,
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Got non-OK response %s: %s", res.Status, string(buf))
		return nil, fmt.Errorf("service returned %s", res.Status)
	}

	var data GetDefinitionResult
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, errors.New("found no definitions")
	}
	var lexeme Lexeme
	for _, entry := range data {
		lexeme.Entries = append(lexeme.Entries, parseEntry(entry))
	}
	return &lexeme, nil
}

func parseEntry(entry DEntry) Entry {
	return Entry{
		Meta:                parseMeta(entry.Meta),
		Headword:            parseHeadword(entry.Hwi),
		GrammaticalFunction: entry.Fl,
		ShortDefinitions:    entry.Shortdef,
		Definitions:         parseDefinitions(entry.Def),
		Quotes:              parseQuotes(entry.Quotes),
	}
}

func parseMeta(meta DMeta) Meta {
	return Meta{
		ID:        meta.ID,
		UUID:      meta.UUID,
		Sort:      meta.Sort,
		Source:    meta.Src,
		Section:   meta.Section,
		Stems:     meta.Stems,
		Offensive: meta.Offensive,
	}
}

func parseHeadword(hwi DHwi) Headword {
	return Headword{
		Text:           hwi.Hw,
		Pronunciations: parsePronunciations(hwi.Prs),
	}
}

func parsePronunciations(prs []DPrs) []Pronunciation {
	var res []Pronunciation
	for _, p := range prs {
		res = append(res, Pronunciation{
			Text:  p.Mw,
			Sound: p.Sound.Audio,
		})
	}
	return res
}

func parseDefinitions(defs []MDef) []Definition {
	var res []Definition
	for _, def := range defs {
		res = append(res, parseDefinition(def))
	}
	return res
}

// parseDefinition parses definitions.
func parseDefinition(def MDef) Definition {
	return Definition{
		VerbDivider: def.Vd,
		Senses:      parseSenses(def.Sseq),
	}
}

// Ref: https://dictionaryapi.com/products/json#sec-2.dt
func parseSenses(sseq [][][]interface{}) []Sense {
	var senses []Sense
	for _, a := range sseq {
		for _, b := range a {
			for _, c := range b {
				if _, ok := c.(map[string]interface{}); ok {
					var sense Sense
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

func parseQuotes(quotes []DQuote) []Quote {
	var res []Quote
	for _, q := range quotes {
		res = append(res, Quote{
			Text:            q.T,
			Author:          q.Aq.Auth,
			Source:          q.Aq.Source,
			PublicationDate: q.Aq.Aqdate,
		})
	}
	return res
}
