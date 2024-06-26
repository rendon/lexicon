// lexapi implements a client for lexapi.io (not a real service).
package lexapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lexicon/dictapi"
	"lexicon/types"
	"lexicon/util"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseURL = "https://rafaelrendon.io"

var client *http.Client

type APIDictionary struct {
	httpc  *http.Client
	apiKey string
}

// NewDictionary return a new client ready to use.
func NewDictionary() (*APIDictionary, error) {
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("API_KEY is missing")
	}
	return &APIDictionary{
		httpc:  &http.Client{Timeout: time.Second * 3},
		apiKey: apiKey,
	}, nil
}

func (a *APIDictionary) Find(name string) (*types.Lexeme, error) {
	u := fmt.Sprintf("%s/lexemes/%s", baseURL, url.PathEscape(name))
	res, err := a.httpc.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, types.NotFound
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned %s: %s", res.Status, body)
	}

	var lexeme types.Lexeme
	if err := json.Unmarshal(body, &lexeme); err != nil {
		log.Printf("Unable to unmarshal %s", body)
		return nil, err
	}
	return &lexeme, nil
}

// createRequest represents a request object for the POST /lexemes/ API.
type createRequest struct {
	Lexeme *types.Lexeme `json:"lexeme"`
}

// Save calls the /lexemes API and saves a new lexeme.
func (a *APIDictionary) Save(lexeme *types.Lexeme) error {
	timestamp := time.Now()
	if lexeme.CreatedAt == nil {
		lexeme.CreatedAt = &timestamp
	}
	if lexeme.UpdatedAt == nil {
		lexeme.UpdatedAt = &timestamp
	}

	payload, err := util.Serialize(createRequest{Lexeme: lexeme})
	if err != nil {
		return err
	}

	resp, err := a.post(payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		if strings.Contains(string(body), "already exists") {
			log.Printf("%s already exists", lexeme.Name)
			return nil
		}

		message := fmt.Sprintf("Service returned response %v status code", resp.StatusCode)
		if err == nil {
			message += fmt.Sprintf(" body: %s", body)
		}
		return errors.New(message)
	}
	return nil
}

func (a *APIDictionary) post(payload []byte) (*http.Response, error) {
	u := fmt.Sprintf("%s/lexemes/", baseURL)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", a.apiKey)

	return a.httpc.Do(req)
}

func (a *APIDictionary) _delete(name string) (*http.Response, error) {
	u := fmt.Sprintf("%s/lexemes/%s", baseURL, url.PathEscape(name))
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", a.apiKey)

	return a.httpc.Do(req)
}

func (a *APIDictionary) Remove(name string) error {
	_, err := a._delete(name)
	if err != nil {
		return err
	}

	if err := dictapi.Remove(name); err != nil {
		return err
	}

	return err
}

// Stats calls the /lexemes/stats API and returns the parsed result.
func (a *APIDictionary) Stats() ([]types.Stat, error) {
	res, err := a.httpc.Get(baseURL + "/stats")
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
	var stats []types.Stat
	if err := json.Unmarshal(body, &stats); err != nil {
		log.Printf("Unable to unmarshal %s", body)
		return nil, err
	}
	return stats, nil
}

// Close closes any open connection to the server.
func (a *APIDictionary) Close() error {
	// No need to close the http client.
	return nil
}
