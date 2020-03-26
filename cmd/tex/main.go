package main

import (
	"context"
	"fmt"
	"os"

	clientCmd "github.com/RTradeLtd/go-temporalx-sdk/cmd"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"

	// auto register pprof handlers
	_ "net/http/pprof"
)

// This is the main command line file for TemporalX
// the way this directory is broken up, is such that big subcommand groups
// such as node management, file commands, etc... are broken out into their own files.
// Smaller commands, such as config management commands reside in this file.
// This is done to make it easier to manage a complex cli library

var (
	bootstrapEnabled, insecure               bool
	endpointAddress, tlsFilePath, tlsKeyPath string
	// Version is git commit information at build time
	Version string
	ctx     context.Context
	cancel  context.CancelFunc
)

func main() {
	// initialize context
	ctx, cancel = context.WithCancel(context.Background())
	// defer cancel to make sure that we don't get a context leak
	// calling cancel() multiple times wont cause issues
	defer cancel()
	clientCmd.SetupCommands(ctx, cancel)
	// generate the actual cli app
	app := newApp()
	// run the cli app
	if err := app.Run(os.Args); err != nil {
		fmt.Printf(
			"%s %s\n",
			au.Bold(au.Red("error encountered:")),
			au.Red(err.Error()),
		)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "tex-cli"
	app.Usage = "TemporalX command-line management tool"
	app.Version = Version
	app.Authors = loadAuthors()
	app.Flags = loadFlags()
	app.Commands = LoadCommands()
	return app
}

func loadAuthors() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "Alex Trottier",
			Email: "postables@rtradetechnologies.com",
		},
		{
			Name:  "George Xie",
			Email: "georgex@rtradetechnologies.com",
		},
	}
}

func loadFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "bootstrap",
			Aliases:     []string{"bp"},
			Usage:       "bootstrap against public ipfs",
			Destination: &bootstrapEnabled,
		},
	}
}

// LoadCommands returns the root commands object containing
// access to all cli functionality
func LoadCommands() cli.Commands {
	commands := cli.Commands{}
	// load grpc client commands object
	commands = append(commands, clientCmd.LoadClientCommands()...)
	return commands
}
