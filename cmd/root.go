package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version string

	flagShow bool
	flagLang string

	flagPullConcurrency int

	rootCmd = &cobra.Command{
		Version: Version,
		Use:     "macschema",
		Short:   "Generates JSON definitions for Apple APIs",
	}
)

func init() {
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().IntVar(&flagPullConcurrency, "concurrency", 1, "number of concurrent workers")

	rootCmd.PersistentFlags().BoolVar(&flagShow, "show", false, "show resulting JSON to stdout")
	rootCmd.PersistentFlags().StringVar(&flagLang, "lang", "objc", "use language")
}

func Execute() {
	fatal(rootCmd.Execute())
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
