package api

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

type DEntry struct {
	Meta     DMeta           `json:"meta"`
	Hwi      DHwi            `json:"hwi"`
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

// Definition describes a definition in a lexicon entry.
type Definition struct {
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

// Entry represents a meaning intended or conveyed.
// See https://www.merriam-webster.com/dictionary/sense
type Entry struct {
	Meta                Meta         `json:"meta,omitempty"`
	Headword            Headword     `json:"headword,omitempty"`
	GrammaticalFunction string       `json:"grammaticalFunction,omitempty"`
	ShortDefinitions    []string     `json:"shortDefinitions,omitempty"`
	Definitions         []Definition `json:"definitions,omitempty"`
	Quotes              []Quote      `json:"quotes,omitempty"`
	Etymology           []string     `json:"etymology,omitempty"`
}

// Lexeme represents a linguistic unit.
type Lexeme struct {
	Entries []Entry `json:"entries,omitempty"`
}
