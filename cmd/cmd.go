package cmd

import (
	"context"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/urfave/cli/v2"
)

var (
	bootstrapEnabled bool
	// Version is git commit information at build time
	Version string
	ctx     context.Context
	cancel  context.CancelFunc
)

// SetupCommands MUST be called to properly setup the commands repository
func SetupCommands(cctx context.Context, ccancel context.CancelFunc) {
	ctx = cctx
	cancel = ccancel
}

func optsFromFlags(c *cli.Context) client.Opts {
	return client.Opts{
		ListenAddress: c.String("endpoint.address"),
		Insecure:      c.Bool("insecure"),
	}
}
