package main

import (
	"io"
	"strings"
	"text/template"
)

type tefa struct {
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
	fa := &tefa{te: te}

	_, err := te.Funcs(template.FuncMap{
		"csv": func(str string) string {
			return escapeCsv(str)
		},
		"seq": func() int64 {
			return fa.Seq()
		},
	}).Parse(str)

	if err != nil {
		return nil, err
	}

	return fa, nil
}

// escapeCsv escapes a string for use in a CSV file.
func escapeCsv(str string) string {
	// Add double quotes around the string if it contains double quotes.
	if strings.ContainsAny(str, "\"") {
		str = "\"" + strings.ReplaceAll(str, "\"", "\"\"") + "\""
	}

	return str
}
