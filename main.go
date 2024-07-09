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

		te, err := newTefa(args...)
		if err != nil {
			return fmt.Errorf("failed to parse template file: %w", err)
		}

		if err := te.Execute(out, round); err != nil {
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
	rootCmd.Flags().IntP("repeat", "r", 1, "number of times to repeat the template")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
