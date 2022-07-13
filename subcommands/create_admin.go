package subcommands

import (
	"context"
	"db-server/modules/user"
	"db-server/server/db"
	"db-server/utils"
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

	var u = user.User{
		Email:        p.UserName,
		PasswordHash: utils.HashPassword(p.Password),
		Token:        utils.GenerateRandomString(10),
		Active:       true,
		Admin:        true,
	}

	u.Id, _ = uuid.NewUUID()

	db.MetaDb.GetConnection().Create(&u)

	return subcommands.ExitSuccess
}
