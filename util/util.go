package util

import (
	"encoding/json"
	"fmt"
)

// Min returns the minimum of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// QuoteStrings wraps in string with double quotes and returns the result.
func QuoteStrings(strs []string) []string {
	var res []string
	for _, s := range strs {
		res = append(res, fmt.Sprintf("%q", s))
	}
	return res
}

// Serialize converts interface object to a JSON document.
func Serialize(i interface{}) ([]byte, error) {
	doc, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
