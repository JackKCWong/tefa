package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

type tefa struct {
	*gofakeit.Faker
	seq int64
	te  *template.Template
}

func (f *tefa) Seq() int64 {
	f.seq++
	return f.seq
}

func (f *tefa) Execute(w io.Writer) error {
	return f.te.Execute(w, f)
}

func newTefa(str string) (*tefa, error) {
	te := template.New("template")
	tefa := &tefa{
		Faker: gofakeit.NewFaker(source.NewCrypto(), true),
		te:    te,
	}

	funcs := sprig.FuncMap()
	funcs["csv"] = escapeCsv
	_, err := te.Funcs(funcs).Parse(str)

	if err != nil {
		return nil, err
	}

	return tefa, nil
}

// escapeCsv escapes a string for use in a CSV file.
func escapeCsv(str string) string {
	// Add double quotes around the string if it contains double quotes.
	if strings.ContainsAny(str, "\"") {
		str = "\"" + strings.ReplaceAll(str, "\"", "\"\"") + "\""
	}

	return str
}
