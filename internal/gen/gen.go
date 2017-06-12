// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/shurcooL/githubql/internal/hacky/caseconv"
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	schema, err := loadSchema()
	if err != nil {
		return err
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err := t.Execute(&buf, schema)
		if err != nil {
			return err
		}
		out, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		fmt.Println("writing", filename)
		err = ioutil.WriteFile(filename, out, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadSchema() (schema interface{}, err error) {
	f, err := os.Open(filepath.Join("internal", "gen", "schema.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&schema)
	return schema, err
}

// Filename -> Template.
var templates = map[string]*template.Template{
	"enum.go": t(`package githubql
{{range .data.__schema.types}}{{if and (eq .kind "ENUM") (not (internal .name))}}
{{template "enum" .}}
{{end}}{{end}}


{{- define "enum" -}}
// {{.name}} represents {{.description | endSentence}}
type {{.name}} string

// {{.description | fullSentence}}
const ({{range .enumValues}}
	{{$.name}}{{.name | enumIdentifier}} {{$.name}} = {{.name | quote}} // {{.description | fullSentence}}{{end}}
)
{{- end -}}
`),
}

func t(text string) *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"internal": func(s string) bool { return strings.HasPrefix(s, "__") },
		"quote":    strconv.Quote,
		"enumIdentifier": func(s string) string {
			return caseconv.UnderscoreSepToMixedCaps(strings.ToLower(s))
		},
		"endSentence": func(s string) string {
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return strings.ToLower(s[0:1]) + s[1:]
		},
		"fullSentence": func(s string) string {
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
	}).Parse(text))
}
