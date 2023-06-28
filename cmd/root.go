package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version string

	flagShow bool
	flagLang string

	flagDebug   bool
	flagTimeout time.Duration

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

	pullCmd.Flags().IntVar(&flagPullConcurrency, "concurrency", runtime.NumCPU(), "number of concurrent workers")

	rootCmd.PersistentFlags().BoolVar(&flagShow, "show", false, "show resulting JSON to stdout")
	rootCmd.PersistentFlags().StringVar(&flagLang, "lang", "objc", "use language")

	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 20*time.Second, "timeout duration")
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
