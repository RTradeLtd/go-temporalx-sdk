package cmd

import (
	"context"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/urfave/cli/v2"
)

var (
	// Version is git commit information at build time
	Version string
)

// SetupCommands MUST be called to properly setup the commands repository
// DEPRECATED: is now a no-op
func SetupCommands(cctx context.Context, ccancel context.CancelFunc) {}

func optsFromFlags(c *cli.Context) client.Opts {
	return client.Opts{
		ListenAddress: c.String("endpoint.address"),
		Insecure:      c.Bool("insecure"),
	}
}
