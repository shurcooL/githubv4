package caseconv_test

import (
	"fmt"
	"testing"

	"github.com/shurcooL/githubql/internal/hacky/caseconv"
)

func ExampleLowerCamelCaseToMixedCaps() {
	fmt.Println(caseconv.LowerCamelCaseToMixedCaps("clientMutationId"))

	// Output: ClientMutationID
}

func ExampleMixedCapsToLowerCamelCase() {
	fmt.Println(caseconv.MixedCapsToLowerCamelCase("ClientMutationID"))

	// Output: clientMutationId
}

func ExampleUnderscoreSepToCamelCase() {
	fmt.Println(caseconv.UnderscoreSepToCamelCase("string_URL_append"))

	// Output: StringUrlAppend
}

func ExampleCamelCaseToUnderscoreSep() {
	fmt.Println(caseconv.CamelCaseToUnderscoreSep("StringUrlAppend"))

	// Output: string_URL_append
}

func ExampleUnderscoreSepToMixedCaps() {
	fmt.Println(caseconv.UnderscoreSepToMixedCaps("string_URL_append"))

	// Output: StringURLAppend
}

func ExampleUpperUnderscoreSepToMixedCaps() {
	fmt.Println(caseconv.UpperUnderscoreSepToMixedCaps("STRING_URL_APPEND"))

	// Output: StringURLAppend
}

func ExampleMixedCapsToUnderscoreSep() {
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("StringURLAppend"))
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("URLFrom"))
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("SetURLHTML"))

	// Output:
	// string_URL_append
	// URL_from
	// set_URL_HTML
}

func TestMixedCapsToLowerCamelCase(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "DatabaseID", want: "databaseId"},
		{in: "URL", want: "url"},
		{in: "ID", want: "id"},
		{in: "CreatedAt", want: "createdAt"},
		{in: "Login", want: "login"},
		{in: "ResetAt", want: "resetAt"},
	}
	for _, tc := range tests {
		got := caseconv.MixedCapsToLowerCamelCase(tc.in)
		if got != tc.want {
			t.Errorf("got: %q, want: %q", got, tc.want)
		}
	}
}
