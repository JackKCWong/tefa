package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"text/template"

	"math/rand"

	"github.com/Masterminds/sprig/v3"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

type tefa struct {
	*gofakeit.Faker
	seq  int64
	tp   *template.Template
	Data chan any
}

func (f *tefa) Seq() int64 {
	f.seq++
	return f.seq
}

func (f *tefa) Execute(w io.Writer, i int) error {
	go func() {
		for j := 0; j < i; j++ {
			f.Data <- f
		}
		close(f.Data)
	}()

	return f.tp.Execute(w, f)
}

func newTefa(preTemplate, mainTemplate string) (*tefa, error) {
	funcs := sprig.FuncMap()
	funcs["csv"] = escapeCsv
	funcs["lines"] = readlines
	funcs["any"] = anyOf

	templateBody := bytes.Buffer{}
	templateBody.WriteString(preTemplate)
	templateBody.WriteString("{{range .Data}}")
	templateBody.WriteString(mainTemplate)
	templateBody.WriteString("{{end}}")

	tp, err := template.New("").Funcs(funcs).Parse(templateBody.String())

	if err != nil {
		return nil, err
	}

	tefa := &tefa{
		Faker: gofakeit.NewFaker(source.NewCrypto(), true),
		tp:    tp,
		Data:  make(chan any),
	}

	return tefa, nil
}

// escapeCsv escapes a string for use in a CSV file.
func escapeCsv(str string) string {
	// Add double quotes around the string if it contains double quotes, commas, or newlines
	if strings.Contains(str, "\"") || strings.Contains(str, ",") || strings.Contains(str, "\n") {
		str = strings.ReplaceAll(str, "\"", "\"\"")
		str = "\"" + str + "\""
	}

	// Replace all newlines
	str = strings.ReplaceAll(str, "\n", "\\n")

	return str
}

// readlines reads a file line by line and returns a slice of strings.
func readlines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func anyOf(strs []string) string {
	if len(strs) < 1 {
		return ""
	}

	rnd := rand.Intn(len(strs))
	return strs[rnd]
}
