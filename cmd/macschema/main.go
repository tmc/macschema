package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/progrium/macschema/pkg/schema"
	"github.com/progrium/macschema/pkg/topic"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&fetch{}, "")
	subcommands.Register(&crawl{}, "")
	subcommands.Register(&stats{}, "")
	subcommands.Register(&parse{}, "")
	subcommands.Register(&types{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

type types struct{}

func (*types) Name() string             { return "types" }
func (*types) Synopsis() string         { return "types topic as json" }
func (*types) Usage() string            { return "types <topic>" }
func (*types) SetFlags(f *flag.FlagSet) {}
func (p *types) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	schema.Types(f.Arg(0))
	return subcommands.ExitSuccess
}

type parse struct{}

func (*parse) Name() string             { return "parse" }
func (*parse) Synopsis() string         { return "parse topic as json" }
func (*parse) Usage() string            { return "parse <topic>" }
func (*parse) SetFlags(f *flag.FlagSet) {}
func (p *parse) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	schema.Parse(f.Arg(0))
	return subcommands.ExitSuccess
}

type fetch struct{}

func (*fetch) Name() string             { return "fetch" }
func (*fetch) Synopsis() string         { return "fetch topic as json" }
func (*fetch) Usage() string            { return "fetch <topic>" }
func (*fetch) SetFlags(f *flag.FlagSet) {}
func (p *fetch) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	topic.FetchURL(f.Arg(0))
	return subcommands.ExitSuccess
}

type crawl struct{}

func (*crawl) Name() string             { return "crawl" }
func (*crawl) Synopsis() string         { return "crawl local topic" }
func (*crawl) Usage() string            { return "crawl <topic>" }
func (*crawl) SetFlags(f *flag.FlagSet) {}
func (p *crawl) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	topic.CrawlTopic(f.Arg(0))
	return subcommands.ExitSuccess
}

type stats struct{}

func (*stats) Name() string             { return "stats" }
func (*stats) Synopsis() string         { return "stats of local docs" }
func (*stats) Usage() string            { return "stats" }
func (*stats) SetFlags(f *flag.FlagSet) {}
func (p *stats) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	topic.Stats()
	return subcommands.ExitSuccess
}
