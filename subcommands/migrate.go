package subcommands

import (
	"context"
	"db-server/migrations"
	"db-server/server"
	"flag"
	"github.com/google/subcommands"
)

type Migrate struct {
}

func (*Migrate) Name() string     { return "migrate" }
func (*Migrate) Synopsis() string { return "Run database migrations" }
func (*Migrate) Usage() string {
	return `migrate:
  Run database migrations.
`
}

func (p *Migrate) SetFlags(f *flag.FlagSet) {
}

func (p *Migrate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	db := server.MetaDb.GetConnection()
	migrations.Migrate(db)

	return subcommands.ExitSuccess
}
