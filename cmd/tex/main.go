package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	clientCmd "github.com/RTradeLtd/go-temporalx-sdk/cmd"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
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
	// CompileDate is the date at which this binary was compiled
	CompileDate string
)

func main() {
	// initialize context
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel to make sure that we don't get a context leak
	// calling cancel() multiple times wont cause issues
	defer cancel()
	// generate the actual cli app
	app := newApp()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	// run the cli app
	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Printf(
			"%s %s\n",
			au.Bold(au.Red("error encountered:")),
			au.Red(err.Error()),
		)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	cli.VersionPrinter = versionPrinter()
	app := cli.NewApp()
	app.Name = "tex-cli"
	app.Usage = "TemporalX client cli"
	app.Description = `
This is the publicly available version of TemporalX's CLI tool intended for using the gRPC API exposed by TemporalX, stripped of all configuration+service management
`
	app.EnableBashCompletion = true
	app.Copyright = "(c) 2020 RTrade Technologies Ltd"
	app.Version = Version
	app.Authors = loadAuthors()
	app.Commands = LoadCommands()

	return app
}

func versionPrinter() func(c *cli.Context) {
	return func(c *cli.Context) {
		fmt.Fprintf(
			c.App.Writer,
			"version:\t\t%s\nreleased:\t\t%s\n",
			c.App.Version,
			CompileDate,
		)
	}
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

// LoadCommands returns the root commands object containing
// access to all cli functionality
func LoadCommands() cli.Commands {
	commands := cli.Commands{}
	// load grpc client commands object
	commands = append(commands, clientCmd.LoadClientCommands()...)
	commands = append(commands, clientCmd.LoadUtilCommands()...)
	return commands
}
