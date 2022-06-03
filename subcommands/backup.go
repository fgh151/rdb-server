package subcommands

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/server"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

var (
	TarCmd       = "tar"
	MysqlDumpCmd = "mysqldump"
	PGDumpCmd    = "pg_dump"
	MongoDumpCmd = "mongodump"
)

type Backup struct {
	s3 bool
}

func (*Backup) Name() string     { return "backup" }
func (*Backup) Synopsis() string { return "Backup database" }
func (*Backup) Usage() string {
	return `backup [-s3]:
  Backup database.
`
}

func (p *Backup) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.s3, "s3", false, "Upload to s3")
}

type ExportResult struct {
	// Path to exported file
	Path string
	// MIME type of the exported file (e.g. application/x-tar)
	MIME string
	// Any error that occured during `Export()`
	Error *error
}

func (r ExportResult) UploadToS3() (string, minio.UploadInfo) {
	upload, _ := os.Open(r.Path)
	return server.UploadToS3(upload, "backup", r.Path, r.MIME)
}

func (p *Backup) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	switch os.Getenv("META_DB_TYPE") {
	case "sqlite":
		p.backupSqlite()
	case "mysql":
		p.backupMysql()
	case "postgres":
		p.backupPostgres()
	default:
		println("Unsupported db type")
	}

	p.backupMongo()

	return subcommands.ExitSuccess
}

func (p Backup) backupSqlite() {
	result := &ExportResult{
		MIME: "application/vnd.sqlite3",
		Path: os.Getenv("META_DB_DSN"),
	}

	result.UploadToS3()
}

func (p Backup) backupMysql() {
	conn := drivers.NewMysqlConnectionFromEnv()

	result := &ExportResult{MIME: "application/x-tar"}

	dumpPath := fmt.Sprintf(`me_%v_%v.sql`, conn.DbName, time.Now().Unix())

	options := []string{
		fmt.Sprintf(`-r%v`, dumpPath),
		fmt.Sprintf(`-h%v`, conn.Host),
		fmt.Sprintf(`-P%v`, conn.Port),
		fmt.Sprintf(`-u%v`, conn.User),
		fmt.Sprintf(`-p%v`, conn.Password),
	}
	out, err := exec.Command(MysqlDumpCmd, options...).Output()

	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)

	result.Path = dumpPath + ".tar.gz"

	_, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()

	err2.PanicErr(err)

	if p.s3 {
		result.UploadToS3()
	}
}

func (p Backup) backupPostgres() {

	conn := drivers.NewPostgresConnectionFromEnv()

	result := &ExportResult{MIME: "application/x-tar"}
	result.Path = fmt.Sprintf(`me_%v_%v.sql.tar.gz`, conn.DbName, time.Now().Unix())
	options := []string{
		"-Fc",
		fmt.Sprintf(`-f%v`, result.Path),
		fmt.Sprintf(`-d%v`, conn.DbName),
		fmt.Sprintf(`-h%v`, conn.Host),
		fmt.Sprintf(`-p%v`, conn.Port),
		fmt.Sprintf(`-U%v`, conn.User),
		fmt.Sprintf(`-W%v`, conn.Password),
	}
	out, err := exec.Command(PGDumpCmd, options...).Output()

	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)

	if p.s3 {
		result.UploadToS3()
	}
}

func (p Backup) backupMongo() {

	conn := drivers.NewMongoConnectionFromEnv()
	result := &ExportResult{MIME: "application/x-tar"}
	dumpPath := fmt.Sprintf(`db_%v_%v`, conn.DbName, time.Now().Unix())

	options := []string{
		fmt.Sprintf(`--db%v`, conn.DbName),
		fmt.Sprintf(`--out%v`, dumpPath),
		conn.GetDsn(),
	}

	out, err := exec.Command(MongoDumpCmd, options...).Output()

	if err != nil {
		logrus.Debug(out)
	}
	err2.PanicErr(err)

	result.Path = dumpPath + ".tar.gz"

	_, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()

	err2.PanicErr(err)

	if p.s3 {
		result.UploadToS3()
	}
}
