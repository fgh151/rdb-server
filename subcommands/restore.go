package subcommands

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

type Restore struct {
	dbPath    string
	mongoPath string
}

func (*Restore) Name() string     { return "restore" }
func (*Restore) Synopsis() string { return "Restore database" }
func (*Restore) Usage() string {
	return `restore:
  Restore database.
`
}

func (p *Restore) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.dbPath, "dbPath", "", "Path to db backup file")
	f.StringVar(&p.mongoPath, "mongoPath", "", "Path to mongo backup")
}

func (p *Restore) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	switch os.Getenv("META_DB_TYPE") {
	case "sqlite":
	case "mysql":
		p.restoreMysql()
	case "postgres":
	default:
		println("Unsupported db type")
	}

	p.restoreMongo()

	return subcommands.ExitSuccess
}

func (p Restore) restoreMysql() {
	conn := drivers.NewMysqlConnectionFromEnv()
	options := []string{
		fmt.Sprintf(`-h%v`, conn.Host),
		fmt.Sprintf(`-P%v`, conn.Port),
		fmt.Sprintf(`-u%v`, conn.User),
		fmt.Sprintf(`-p%v`, conn.Password),
		fmt.Sprintf(`%v > %v`, conn.DbName, p.dbPath),
	}

	out, err := exec.Command("mysql", options...).Output()
	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)
}

func (p Restore) restorePostgres() {
	conn := drivers.NewPostgresConnectionFromEnv()
	options := []string{
		fmt.Sprintf(`-d%v`, conn.DbName),
		fmt.Sprintf(`-h%v`, conn.Host),
		fmt.Sprintf(`-p%v`, conn.Port),
		fmt.Sprintf(`-U%v`, conn.User),
		fmt.Sprintf(`-W%v`, conn.Password),
		fmt.Sprintf(`-f%v`, p.dbPath),
	}

	out, err := exec.Command("pg_restore", options...).Output()
	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)
}

func (p Restore) restoreMongo() {
	conn := drivers.NewMongoConnectionFromEnv()
	options := []string{
		fmt.Sprintf(`--db%v`, conn.DbName),
		conn.GetDsn(),
		p.mongoPath,
	}

	out, err := exec.Command("mongorestore", options...).Output()
	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)
}
