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

`lines`: read lines from a file.

`csv`: escape the input so it's safe to put in a csv cell.

`nth`: get nth element from a slice.

`any`: get a random element from a slice.

`islice`: convert any slice (`[]int`, `[]string`, etc.) to a `[]interface{}`
