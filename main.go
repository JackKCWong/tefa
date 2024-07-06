package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tefa [options] <template file path>",
	Short: "TEmplating with FAke data.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		count := must(cmd.Flags().GetInt("count"))
		if count < 1 {
			return fmt.Errorf("count must be greater than 0")
		}

		out, err := mkOutput(must(cmd.Flags().GetString("output")))
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		// Read the template file
		mainTemplatePath := args[0]
		mainTemplateBytes, err := os.ReadFile(mainTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to read main template file: %w", err)
		}

		preTemplatePath := must(cmd.Flags().GetString("pre"))
		preTemplateBytes := []byte{}
		if preTemplatePath != "" {
			preTemplateBytes, err = os.ReadFile(preTemplatePath)
			if err != nil {
				return fmt.Errorf("failed to read pre template file: %w", err)
			}
		}

		te, err := newTefa(string(preTemplateBytes), string(mainTemplateBytes))
		if err != nil {
			return fmt.Errorf("failed to parse template file: %w", err)
		}

		if err := te.Execute(out, count); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		return nil
	},
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
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

func init() {
	rootCmd.Flags().StringP("output", "o", "", "output file")
	rootCmd.Flags().IntP("count", "c", 1, "number of times to repeat the template")
	rootCmd.Flags().StringP("pre", "p", "", "a snippet of template code to run before the main template")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
