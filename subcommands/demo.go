package subcommands

import (
	"context"
	"db-server/models"
	"flag"
	"github.com/google/subcommands"
)

type Demo struct {
}

func (*Demo) Name() string     { return "demo" }
func (*Demo) Synopsis() string { return "Fill db demo data" }
func (*Demo) Usage() string {
	return `demo:
  Fill db demo data.
`
}

func (p *Demo) SetFlags(f *flag.FlagSet) {
}

func (p *Demo) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	models.CreateDemo()

	return subcommands.ExitSuccess
}
