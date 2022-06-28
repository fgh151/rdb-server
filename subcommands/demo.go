package subcommands

import (
	"context"
	"db-server/modules/user"
	"db-server/security"
	"db-server/server/db"
	"flag"
	"github.com/google/subcommands"
	"github.com/google/uuid"
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
	CreateDemo()

	return subcommands.ExitSuccess
}

func CreateDemo() {
	var u = user.User{
		Email:        "test@example.com",
		PasswordHash: security.HashPassword("test"),
		Token:        "123",
		Active:       true,
		Admin:        true,
	}

	u.Id, _ = uuid.NewUUID()

	db.MetaDb.GetConnection().Create(&u)
}
