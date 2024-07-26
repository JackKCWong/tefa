# Tefa

A utility for TEmplating with FAke data. The main usecase I have in mind is to generate random csv data that I can import into a database for testing.

## Get started

```zsh
# install
go install github.com/JackKCWong/tefa@latest

# get help
tefa -h
Usage:
  tefa [options] <template files path> [flags]

Flags:
  -r, --repeat int       number of times to repeat the template (default 1)
  -h, --help            help for tefa
  -o, --output string   output file


## e.g.
tefa <(<<_EOF
{{.Name}}
_EOF
)
```

It combines [Go template](https://pkg.go.dev/text/template) with [gofakeit](https://github.com/brianvoe/gofakeit) and [sprig](https://masterminds.github.io/sprig/). You can check the respective documents to see what you can use.

## Additional functions / methods

`.Index`: When a template is repeated multiple times, this function can be used to generate a sequence of numbers starting from 0.

`lines(filepath string)`: read all lines from a file.

`scan(filepath string)`: read a file line by line.

`csv(val string)`: escape the input so it's safe to put in a csv cell.

`nth([]string)`: get nth element from a `[]string`.

`any([]string)`: get a random element from a `[]string`.

`shuffle([]string)`: shuffle a `[]string`.

`islice`: convert any slice (`[]int`, `[]string`, etc.) to a `[]interface{}`

`bool(p float32)`: a true/false value with a given probability.

`kv(k string)`: return a value defined by `-D` flag.

`atoi(string)`: convert string to int.

`mapf(f string, s any)`: convert `[]any` to `[]string` using `printf`.

`uuidv7`: generate a UUID v7 which can be time sorted. [Ref.](https://pkg.go.dev/github.com/google/uuid#NewV7)

`ulid`: Universally Unique Lexicographically Sortable Identifier. [Ref.](https://pkg.go.dev/github.com/oklog/ulid/v2#Make)

`sqid`: A [Sqids](https://pkg.go.dev/github.com/sqids/sqids-go#section-readme) by encoding `.Index` and `time.Now()`
