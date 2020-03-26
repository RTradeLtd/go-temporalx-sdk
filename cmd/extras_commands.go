package cmd

import (
	"errors"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/urfave/cli/v2"
)

var extrasName string

// LoadExtrasCommands returns an extras commands object
func LoadExtrasCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "extras",
			Usage:       "node extras management",
			Description: "enables access to the node extras management system",
			Subcommands: cli.Commands{enableExtras(), disableExtras()},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "extras.name",
					Aliases:     []string{"en"},
					Usage:       "the extras feature to enable",
					Destination: &extrasName,
				},
			},
		},
	}
}

func enableExtras() *cli.Command {
	return &cli.Command{
		Name:  "enable",
		Usage: "enable an extras feature",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			if extrasName == "" {
				return errors.New("extras.name flag is empty")
			}
			_, err = cl.Extras(ctx, &pb.ExtrasRequest{
				RequestType:   pb.EXTRASREQTYPE_EX_ENABLE,
				ExtrasFeature: pb.EXTRASTYPE(pb.EXTRASTYPE_value[extrasName]),
			})
			return err
		},
	}
}

func disableExtras() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "disable an extras feature",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			if extrasName == "" {
				return errors.New("extras.name flag is empty")
			}
			_, err = cl.Extras(ctx, &pb.ExtrasRequest{
				RequestType:   pb.EXTRASREQTYPE_EX_DISABLE,
				ExtrasFeature: pb.EXTRASTYPE(pb.EXTRASTYPE_value[extrasName]),
			})
			return err
		},
	}
}
