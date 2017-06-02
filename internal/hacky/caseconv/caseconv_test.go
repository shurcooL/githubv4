package caseconv_test

import (
	"fmt"

	"github.com/shurcooL/githubql/internal/hacky/caseconv"
)

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

func ExampleMixedCapsToUnderscoreSep() {
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("StringURLAppend"))
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("URLFrom"))
	fmt.Println(caseconv.MixedCapsToUnderscoreSep("SetURLHTML"))

	// Output:
	// string_URL_append
	// URL_from
	// set_URL_HTML
}
