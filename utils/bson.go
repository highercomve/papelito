package utils

import (
	"strings"

	cjson "github.com/gibson042/canonicaljson-go"
)

func bsonQuoteMap(m *map[string]interface{}) map[string]interface{} {
	escapedMap := map[string]interface{}{}
	for k, v := range *m {
		nk := BsonQuote(k)
		escapedMap[nk] = v
	}
	return escapedMap
}

func bsonUnquoteMap(m *map[string]interface{}) map[string]interface{} {
	escapedMap := map[string]interface{}{}
	for k, v := range *m {
		nk := BsonUnquote(k)
		escapedMap[nk] = v
	}
	return escapedMap
}

// BsonQuoteMap create a new map of quotes with escaped indexes
func BsonQuoteMap(m *map[string]interface{}) map[string]interface{} {
	b, err := cjson.Marshal(m)
	if err != nil {
		return bsonQuoteMap(m)
	}

	escapedMap := map[string]interface{}{}
	err = cjson.Unmarshal([]byte(BsonQuote(string(b))), &escapedMap)
	if err != nil {
		return bsonQuoteMap(m)
	}

	return escapedMap
}

// BsonUnquoteMap create a new map of quotes with unescaped indexes
func BsonUnquoteMap(m *map[string]interface{}) map[string]interface{} {
	b, err := cjson.Marshal(m)
	if err != nil {
		return bsonUnquoteMap(m)
	}

	escapedMap := map[string]interface{}{}
	err = cjson.Unmarshal([]byte(BsonUnquote(string(b))), &escapedMap)
	if err != nil {
		return bsonUnquoteMap(m)
	}

	return escapedMap
}

// BsonUnquote unquote a string
func BsonUnquote(s string) string {
	return strings.Replace(
		strings.Replace(s, "\uFF2E", ".", -1),
		"\uFFE0", "$", -1,
	)
}

// BsonQuote quote a string
func BsonQuote(s string) string {
	return strings.Replace(
		strings.Replace(s, ".", "\uFF2E", -1),
		"$", "\uFFE0", -1,
	)
}
