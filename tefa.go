package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"math/rand"

	"github.com/Masterminds/sprig/v3"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/google/uuid"
)

type tefa struct {
	*gofakeit.Faker
	idx int
	tp  *template.Template
}

func (f *tefa) Index() int {
	return f.idx
}

func (f *tefa) Execute(w io.Writer, n int) error {
	for i := 0; i < n; i++ {
		f.idx = i
		if err := f.tp.Execute(w, f); err != nil {
			return err
		}
	}

	return nil
}

func newTefa(vars map[string]string, templateFiles ...string) (*tefa, error) {
	funcs := sprig.FuncMap()
	funcs["csv"] = escapeCsv
	funcs["lines"] = readlines
	funcs["scan"] = scanlines
	funcs["any"] = anyOf
	funcs["nth"] = nth
	funcs["tick"] = tick
	funcs["bool"] = randomBool
	funcs["shuffle"] = shuffle
	funcs["islice"] = interfaceSlice
	funcs["mapf"] = mapf
	funcs["kv"] = func(k string) string {
		return vars[k]
	}
	funcs["atoi"] = strconv.Atoi
	funcs["uuidv7"] = uuid.NewV7

	tp, err := template.New(templateFiles[0]).Funcs(funcs).ParseFiles(templateFiles...)

	if err != nil {
		return nil, err
	}

	tefa := &tefa{
		Faker: gofakeit.NewFaker(source.NewCrypto(), true),
		tp:    tp,
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

func scanlines(filepath string) (chan string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	lines := make(chan string, 1000)
	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	return lines, nil
}

func anyOf(arr []string) string {
	rnd := rand.Intn(len(arr))
	return arr[rnd]
}

func nth(n int, arr []string) string {
	return arr[n]
}

func tick(n int) chan int {
	c := make(chan int)
	go func() {
		var i int
		for i < n {
			c <- i
			i++
		}

		close(c)
	}()

	return c
}

func shuffle(arr []string) []string {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})

	return arr
}

func randomBool(p float32) bool {
	return rand.Float32() < p
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func mapf(f string, s interface{}) []string {
	sli := interfaceSlice(s)
	ret := make([]string, len(sli))
	for i, v := range sli {
		ret[i] = fmt.Sprintf(f, v)
	}

	return ret
}
