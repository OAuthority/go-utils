// Helpers for parameters in urls; doesn't really do much on its own.
//
// Most of this code is copied from the net/url package and is slightly adapted to ensure that
// it adheres with what MediaWiki expects.
// it is therefore Copyright (2009) of The Go Authors.
// see: https://pkg.go.dev/net/http?tab=licenses

package gowikibot

import (
	"bytes"
	"net/url"
	"sort"
)

type Values map[string]string

// Return the value associated with a specific key.
// if one was not found then we return an empty string.
func (val Values) Get(key string) string {

	if val == nil {
		return ""
	}

	valString, ok := val[key]

	if !ok {
		return ""
	}

	return valString
}

// set a value on a key and replace any exisitng value that may
// exist
func (val Values) Set(key, value string) {
	val[key] = value
}

// Delete the value associated with a key, obviously
func (val Values) Delete(key string) {
	delete(val, key)
}

// Encode all of the values so that they are in URL format, with the exception of the
// token, which is added at the end of the parameters and is not encoded. This ensures
// that if, for some reason, the token has been cut off, that the request is not executed
// see the MediaWiki API documentation for discussion of this.
func (v Values) Encode() string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	keys := v.sortKeys()
	token := false
	for _, k := range keys {
		if k == "token" {
			token = true
			continue
		}
		prefix := url.QueryEscape(k) + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(url.QueryEscape(v[k]))
	}
	if token {
		buf.WriteString("&token=" + url.QueryEscape(v["token"]))
	}
	return buf.String()
}

// Sort our keys into alphabetical order
func (v Values) sortKeys() []string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
