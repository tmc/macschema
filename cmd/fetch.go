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
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Download a topic to doc dir",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		l := schema.NewLookup(args[0], flagLang)

		if l.DocExists() && flagShow {
			b, err := ioutil.ReadFile(l.DocPath)
			fatal(err)
			os.Stdout.Write(append(b, '\n'))
			return
		}

		ctx := context.Background()
		t := schema.FetchTopic(ctx, l)
		fatal(writeTopic(l, t))
		fmt.Fprintf(os.Stderr, "=> %s [%s]\n", l.DocPath, time.Since(t.LastFetch))
	},
}

func writeTopic(l schema.Lookup, t schema.Topic) error {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(l.DocPath), 0755)
	if err := ioutil.WriteFile(l.DocPath, b, 0644); err != nil {
		return err
	}

	if flagShow {
		os.Stdout.Write(append(b, '\n'))
	}

	return nil
}
