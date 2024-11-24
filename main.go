package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "tefa [options] <template file path>",
	Short: "TEmplating with FAke data.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		round := must(cmd.Flags().GetInt("repeat"))
		if round < 1 {
			return fmt.Errorf("repeat must be greater than 0")
		}

		out, err := mkOutput(must(cmd.Flags().GetString("output")))
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		values := make(map[string]any)
		valuesDefFile := must(cmd.Flags().GetString("values"))
		if valuesDefFile != "" {
			err := loadValues(valuesDefFile, &values)
			if err != nil {
				return err
			}
		}

		valuesFromCli := must(cmd.Flags().GetStringToString("define"))
		for k, v := range inferTypes(valuesFromCli) {
			values[k] = v
		}

		te, err := newTefa(values, args...)
		if err != nil {
			return fmt.Errorf("failed to parse template file: %w", err)
		}

		if err := te.Execute(out, round); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "output file")
	rootCmd.Flags().IntP("repeat", "r", 1, "number of times to repeat the template")
	rootCmd.Flags().StringToStringP("define", "D", make(map[string]string), "define key=value pairs, use comma to separate multiple pairs")
	rootCmd.Flags().StringP("values", "f", "", "a yaml file to load values from")
}

func mkOutput(output string) (io.WriteCloser, error) {
	if output == "" {
		return os.Stdout, nil
	}

	of, err := os.Create(output)
	if err != nil {
		return nil, err
	}

	return of, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func inferTypes(kv map[string]string) map[string]any {
	ret := make(map[string]any)
	for k, v := range kv {
		if _, err := strconv.Atoi(v); err == nil {
			ret[k] = v
			continue
		}

		if _, err := strconv.ParseFloat(v, 64); err == nil {
			ret[k] = v
			continue
		}

		if _, err := time.Parse(time.RFC3339, v); err == nil {
			ret[k] = v
			continue
		}

		if _, err := time.ParseDuration(v); err == nil {
			ret[k] = v
			continue
		}

		ret[k] = v
	}

	return ret
}

func loadValues(path string, vals *map[string]any) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read values file: %w", err)
	}

	err = yaml.Unmarshal(fileBytes, vals)
	if err != nil {
		return fmt.Errorf("failed to unmarshal values file: %w", err)
	}

	return nil
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
