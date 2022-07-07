package subcommands

import (
	"context"
	"db-server/models"
	"db-server/security"
	"db-server/server"
	"flag"
	"github.com/google/subcommands"
	"github.com/google/uuid"
)

type CreateAdmin struct {
	UserName string
	Password string
}

func (*CreateAdmin) Name() string     { return "admin" }
func (*CreateAdmin) Synopsis() string { return "Create admin user" }
func (*CreateAdmin) Usage() string {
	return `admin -e=admin -p=passwd:
  Create user admin with password passwd.
`
}

func (p *CreateAdmin) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.UserName, "e", "", "Set user email")
	f.StringVar(&p.UserName, "p", "", "Set user password")
}

func (p *CreateAdmin) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	var u = models.User{
		Email:        p.UserName,
		PasswordHash: security.HashPassword(p.Password),
		Token:        security.GenerateRandomString(10),
		Active:       true,
		Admin:        true,
	}

	u.Id, _ = uuid.NewUUID()

	server.MetaDb.GetConnection().Create(&u)

	return subcommands.ExitSuccess
}
