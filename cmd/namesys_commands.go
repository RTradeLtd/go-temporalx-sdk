package cmd

import (
	"errors"
	"fmt"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

// LoadNameSysCommands returns namesys api commands
func LoadNameSysCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "namesys",
			Usage:       "namesys commands",
			Description: "Enables access to NameSysAPI",
			Subcommands: cli.Commands{
				namesysPublish(),
				namesysResolve(),
			},
		},
	}
}

func namesysPublish() *cli.Command {
	return &cli.Command{
		Name:  "publish",
		Usage: "publish a new (or update) an ipns record",
		Action: func(c *cli.Context) error {
			if c.String("key.name") == "" {
				return errors.New("key.name flag is empty")
			}
			if c.String("cid") == "" {
				return errors.New("cid flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			resp, err := cl.Keystore(c.Context, &pb.KeystoreRequest{
				RequestType: pb.KSREQTYPE_KS_GET,
				Name:        c.String("key.name"),
			})
			if err != nil {
				return err
			}
			_, err = cl.NameSysPublish(c.Context, &pb.NameSysPublishRequest{
				PrivateKey: resp.GetPrivateKey(),
				Value:      c.String("cid"),
			})
			return err
		},
		Flags: []cli.Flag{KeyName(), CidFlag("the cid to publish")},
	}
}

func namesysResolve() *cli.Command {
	return &cli.Command{
		Name:  "resolve",
		Usage: "resolve an ipns record",
		Action: func(c *cli.Context) error {
			if c.String("cid") == "" {
				return errors.New("cid flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			resp, err := cl.NameSysResolve(c.Context, &pb.NameSysResolveRequest{
				Name: c.String("cid"),
			})
			if err != nil {
				return err
			}
			if resp.GetError() != "" {
				return errors.New(resp.GetError())
			}
			fmt.Printf(
				"%s %s\n",
				au.Bold(au.Green("path:")),
				au.Bold(au.White(resp.GetPath())),
			)
			return nil
		},
		Flags: []cli.Flag{CidFlag("the cid to resolve")},
	}
}
