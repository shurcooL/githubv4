package githubql

import (
	"bytes"
	"io"
	"reflect"
	"sort"

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
// E.g., map[string]interface{}{"a": Int(123), "b": NewBoolean(true)} -> "$a:Int!$b:Boolean".
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
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			// TODO: Support t.Elem() being a pointer, if needed. Probably want to do this recursively.
			s += "[" + t.Elem().Name() + "!]" // E.g., "[IssueState!]".
		case reflect.Ptr:
			// Pointer is an optional type, so no "!" at the end.
			s += t.Elem().Name() // E.g., "Int".
		default:
			name := t.Name()
			if name == "string" { // HACK: Workaround for https://github.com/shurcooL/githubql/issues/12.
				name = "ID"
			}
			// Value is a required type, so add "!" to the end.
			s += name + "!" // E.g., "Int!".
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
	querifyType(&buf, reflect.TypeOf(v), false)
	return buf.String()
}

func querifyType(w io.Writer, t reflect.Type, inline bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		querifyType(w, t.Elem(), false)
	case reflect.Struct:
		// Special handling of scalar struct types.
		if t == dateTimeType || t == gitTimestampType || t == uriType || t == x509CertificateType {
			return
		}
		if !inline {
			io.WriteString(w, "{")
		}
		sep := false
		for i := 0; i < t.NumField(); i++ {
			if !sep {
				sep = true
			} else {
				io.WriteString(w, ",")
			}
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					io.WriteString(w, value)
				} else {
					io.WriteString(w, caseconv.MixedCapsToLowerCamelCase(f.Name))
				}
			}
			querifyType(w, f.Type, inlineField)
		}
		if !inline {
			io.WriteString(w, "}")
		}
	}
}

var (
	dateTimeType        = reflect.TypeOf(DateTime{})
	gitTimestampType    = reflect.TypeOf(GitTimestamp{})
	uriType             = reflect.TypeOf(URI{})
	x509CertificateType = reflect.TypeOf(X509Certificate{})
)
