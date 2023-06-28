package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/progrium/macschema/schema"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Generate a schema in api dir fetching topics if needed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		ctx, cancel := schema.WithBrowserContext(context.Background())
		defer cancel()
		l := schema.NewLookup(args[0], flagLang)
		if !l.DocExists() {
			fmt.Fprintln(os.Stderr, "=> Fetching topic...")
			t := schema.FetchTopic(ctx, l, fetchOptions(cmd))
			fatal(writeTopic(l, t))
		}
		t, err := schema.ReadTopic(l)
		fatal(err)

		fmt.Fprintln(os.Stderr, "=> Fetching sub-topics...")
		sem := semaphore.NewWeighted(int64(flagPullConcurrency))

		for _, link := range t.Topics {
			ll := schema.LookupFromPath(link.Path)
			if ll.DocExists() {
				// TODO: check last fetch, version
				continue
			}
			sem.Acquire(ctx, 1)
			go func() {
				defer sem.Release(1)
				fmt.Fprintln(os.Stderr, "  ", ll.DocPath)
				tt := schema.FetchTopic(ctx, ll, fetchOptions(cmd))
				fatal(writeTopic(ll, tt))
			}()
		}
		fmt.Fprintln(os.Stderr, "=> Waiting for workers to finish...")
		sem.Acquire(ctx, int64(flagPullConcurrency))

		fmt.Fprintln(os.Stderr, "=> Generating schema...")
		s := schema.PullSchema(l)
		fatal(writeSchema(l, s))
		fmt.Fprintf(os.Stderr, "=> %s [%s]\n", l.APIPath, time.Since(start))
	},
}

func writeSchema(l schema.Lookup, s schema.Schema) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(l.APIPath), 0755)
	if err := ioutil.WriteFile(l.APIPath, b, 0644); err != nil {
		return err
	}

	if flagShow {
		os.Stdout.Write(append(b, '\n'))
	}

	return nil
}
