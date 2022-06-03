package main

import (
	"context"
	err2 "db-server/err"
	cmd "db-server/subcommands"
	"flag"
	"github.com/google/subcommands"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	err2.PanicErr(err)

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.Backup{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

}
