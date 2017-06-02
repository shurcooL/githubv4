package githubql

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"unicode"

	"github.com/shurcooL/githubql/internal/hacky/caseconv"
)

// WARNING: This file contains hacky (but functional) code. It's very ugly.
//          The goal is to eventually clean up the code here and move it elsewhere,
//          reducing this file to non-existence. But, I'm tackling higher priorities
//          first (such as ensuring the API design will scale and work out), and
//          saving time by deferring this work.

func constructQuery(v interface{}, variables map[string]interface{}) string {
	query := querify(v)
	if variables != nil {
		return "query(" + queryArguments(variables) + ")" + query
	}
	return query
}

func constructMutation(v interface{}, variables map[string]interface{}) string {
	query := querify(v)
	if variables != nil {
		return "mutation(" + queryArguments(variables) + ")" + query
	}
	return "mutation" + query
}

// queryArguments constructs a minified arguments string for variables.
//
// E.g., map[string]interface{}{"A": Int(123), "B": NewBoolean(true)} -> "$A:Int!$B:Boolean".
func queryArguments(variables map[string]interface{}) string {
	sorted := make([]string, 0, len(variables))
	for k := range variables {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	var s string
	for _, k := range sorted {
		v := variables[k]
		s += "$" + k + ":"
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			// Pointer is an optional type, so no "!" at the end.
			s += t.Elem().Name() // E.g., "Int".
		} else {
			// Value is a required type, so add "!" to the end.
			s += t.Name() + "!" // E.g., "Int!".
		}
	}
	return s
}

// querify uses querifyType, which recursively constructs
// a minified query string from the provided struct v.
//
// E.g., struct{Foo Int, Bar *Boolean} -> "{foo,bar}".
func querify(v interface{}) string {
	var buf bytes.Buffer
	querifyType(&buf, reflect.TypeOf(v))
	return buf.String()
}

func querifyType(w io.Writer, t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		querifyType(w, t.Elem())
	case reflect.Struct:
		// Special handling of scalar struct types.
		if t == dateTimeType || t == gitTimestampType || t == uriType || t == x509CertificateType {
			return
		}
		io.WriteString(w, "{")
		var sep bool
		for i := 0; i < t.NumField(); i++ {
			if sep {
				io.WriteString(w, ",")
			} else {
				sep = true
			}
			f := t.Field(i)
			if value, ok := f.Tag.Lookup("graphql"); ok {
				io.WriteString(w, value)
			} else {
				io.WriteString(w, mixedCapsToLowerCamelCase(f.Name))
			}
			querifyType(w, f.Type)
		}
		io.WriteString(w, "}")
	}
}

var (
	dateTimeType        = reflect.TypeOf(DateTime{})
	gitTimestampType    = reflect.TypeOf(GitTimestamp{})
	uriType             = reflect.TypeOf(URI{})
	x509CertificateType = reflect.TypeOf(X509Certificate{})
)

func mixedCapsToLowerCamelCase(s string) string {
	r := []rune(caseconv.UnderscoreSepToCamelCase(caseconv.MixedCapsToUnderscoreSep(s)))
	if len(r) == 0 {
		return ""
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
