package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagShow bool
	flagLang string

	rootCmd = &cobra.Command{
		Version: "0.1.0",
		Use:     "macschema",
		Short:   "Generates JSON definitions for Apple APIs",
	}
)

func init() {
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(pullCmd)

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
