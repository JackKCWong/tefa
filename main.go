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

		templatePath := args[0]
		output := must(cmd.Flags().GetString("output"))
		var out io.Writer
		if output == "" {
			out = os.Stdout
		} else {
			of, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}

			defer of.Close()
			out = of
		}

		// Read the template file
		templateData, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}

		templateString := string(templateData)
		te, err := newTefa(templateString)
		if err != nil {
			return fmt.Errorf("failed to parse template file: %w", err)
		}

		for i := 0; i < count; i++ {
			if err := te.Execute(out); err != nil {
				return fmt.Errorf("failed to execute template: %w", err)
			}
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

func init() {
	rootCmd.Flags().StringP("output", "o", "", "output file")
	rootCmd.Flags().IntP("count", "c", 1, "number of times to repeat the template")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
