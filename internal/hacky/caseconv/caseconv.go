// Package caseconv implements functions for converting between CamelCase, snake_case and MixedCaps forms for identifier names.
package caseconv

import (
	"strings"
	"unicode"
)

// LowerCamelCaseToMixedCaps converts "clientMutationId" to "ClientMutationID" form.
func LowerCamelCaseToMixedCaps(s string) string {
	return UnderscoreSepToMixedCaps(CamelCaseToUnderscoreSep(s))
}

// MixedCapsToLowerCamelCase converts "ClientMutationID" to "clientMutationId" form.
func MixedCapsToLowerCamelCase(s string) string {
	r := []rune(UnderscoreSepToCamelCase(MixedCapsToUnderscoreSep(s)))
	if len(r) == 0 {
		return ""
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// UnderscoreSepToCamelCase converts "string_URL_append" to "StringUrlAppend" form.
func UnderscoreSepToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// UnderscoreSepToMixedCaps converts "string_URL_append" to "StringURLAppend" form.
func UnderscoreSepToMixedCaps(in string) string {
	var out string
	ss := strings.Split(in, "_")
	for _, s := range ss {
		initialism := strings.ToUpper(s)
		if _, ok := initialisms[initialism]; ok {
			out += initialism
		} else {
			out += strings.Title(s)
		}
	}
	return out
}

// UpperUnderscoreSepToMixedCaps converts "STRING_URL_APPEND" to "StringURLAppend" form.
func UpperUnderscoreSepToMixedCaps(s string) string {
	return UnderscoreSepToMixedCaps(strings.ToLower(s))
}

// CamelCaseToUnderscoreSep converts "StringUrlAppend" to "string_url_append" form.
func CamelCaseToUnderscoreSep(s string) string {
	var out []rune
	var seg []rune
	for _, r := range s {
		if !unicode.IsLower(r) {
			out = addSegment(out, seg)
			seg = nil
		}
		seg = append(seg, unicode.ToLower(r))
	}
	out = addSegment(out, seg)
	return string(out)
}

// MixedCapsToUnderscoreSep converts "StringURLAppend" to "string_URL_append" form.
func MixedCapsToUnderscoreSep(in string) string {
	var out []rune
	var seg []rune
	for _, r := range in {
		if !unicode.IsLower(r) {
			out = addSegment(out, seg)
			seg = nil
		}
		seg = append(seg, unicode.ToLower(r))
	}
	out = addSegment(out, seg)
	// This is awful. You're welcome to make it better.
	// For each supported initialism, do a strings.Replace to fix up bad output like "_u_r_l_" with "_URL_".
	s := "_" + string(out) + "_"
	for initialism := range initialisms {
		ruinedInitialism := ruinedInitialism(initialism)
		s = strings.Replace(s, ruinedInitialism, "_"+initialism+"_", -1)
	}
	return s[1 : len(s)-1]
}

func addSegment(inout, seg []rune) []rune {
	if len(seg) == 0 {
		return inout
	}
	if len(inout) != 0 {
		inout = append(inout, '_')
	}
	initialism := strings.ToUpper(string(seg))
	if _, ok := initialisms[initialism]; ok {
		inout = append(inout, []rune(initialism)...)
	} else {
		inout = append(inout, seg...)
	}
	return inout
}

// ruinedInitialism returns "_u_r_l_" form that CamelCaseToUnderscoreSep generates for initialism "URL".
func ruinedInitialism(initialism string) string {
	var out = []rune{'_'}
	for _, r := range initialism {
		out = append(out, unicode.ToLower(r))
		out = append(out, '_')
	}
	return string(out)
}

// initialisms is the set of initialisms in Go-style Mixed Caps case.
var initialisms = map[string]struct{}{
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XSRF":  {},
	"XSS":   {},
}
