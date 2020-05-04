package cmd

import "github.com/urfave/cli/v2"

// LoadClientCommands returns a gRPC client commands object
func LoadClientCommands() cli.Commands {
	// load node management commands object
	commands := loadNodeCommands()
	// load file commands object
	commands = append(commands, LoadFileCommands()...)
	// load extras commands object
	commands = append(commands, LoadExtrasCommands()...)
	// load pubsub commands object
	commands = append(commands, LoadPubSubCommands()...)
	// load namesys commands object
	commands = append(commands, LoadNameSysCommands()...)
	// load keystore commands object
	commands = append(commands, LoadKeystoreCommands()...)
	// load api (version, status, etc...) commands
	commands = append(commands, LoadAPICommands()...)
	commands = append(commands, LoadAdminCommands()...)
	// return client commands object
	return cli.Commands{
		&cli.Command{
			Name:        "client",
			Usage:       "gRPC client subcommands",
			Description: "Enables access to a rich gRPC client library",
			Subcommands: commands,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint.address",
					Aliases: []string{"ea"},
					Usage:   "temporalx endpoint address",
					Value:   "127.0.0.1:9090",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "enable insecure connections to temporalx",
					Value: true,
				},
			},
		},
	}
}
