package cmd

import (
	"fmt"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

// LoadAPICommands returns an api commands object
func LoadAPICommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "api",
			Usage:       "low level api maintenance commands",
			Subcommands: cli.Commands{apiStatus(), apiVersion()},
		},
	}
}

func apiStatus() *cli.Command {
	return &cli.Command{
		Name:        "status",
		Usage:       "check api status",
		Description: "returns the current status of the api",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			resp, err := cl.Status(c.Context, &pb.Empty{})
			if err != nil {
				return err
			}
			if c.String("output") == "monitor" {
				fmt.Printf("%s\n", au.Bold(au.Green(resp.GetStatus().String())))
			} else {
				fmt.Printf(
					"%s %s\n%s %s\n",
					au.Bold(au.Green("host:")),
					au.Bold(au.White(resp.GetHost())),
					au.Bold(au.Green("status:")),
					au.Bold(au.White(resp.GetStatus())),
				)
			}
			return nil
		},
		Flags: []cli.Flag{OutputFlag()},
	}
}

func apiVersion() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Usage:       "check api version",
		Description: "returns the current version of the api",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			resp, err := cl.Version(c.Context, &pb.Empty{})
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", au.Bold(au.Green(resp.GetVersion())))
			return nil
		},
		Flags: []cli.Flag{OutputFlag()},
	}
}
