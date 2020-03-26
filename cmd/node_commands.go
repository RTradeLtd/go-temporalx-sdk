package cmd

import (
	"errors"
	"fmt"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

func loadNodeCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "node",
			Usage:       "node management commands",
			Description: "enables management of a node via its API",
			Subcommands: cli.Commands{
				// peer management commands
				&cli.Command{
					Name:        "peer",
					Usage:       "peer management commands",
					Description: "grants access to peer management calls",
					Subcommands: cli.Commands{
						peerCount(),
						peerConnect(),
						peerDisconnect(),
						peerIsConnected(),
					},
				}, loadP2PCommand()},
		},
	}
}

func peerCount() *cli.Command {
	return &cli.Command{
		Name:  "count",
		Usage: "returns the number of connected peers",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			peerCount, err := cl.GetPeerCount(ctx)
			if err != nil {
				return err
			}
			if c.String("output") == "monitor" {
				fmt.Printf(
					"%s\n",
					au.Bold(au.Green(fmt.Sprint(peerCount))),
				)
			} else {
				fmt.Printf(
					"%s %s\n",
					au.Bold(au.Green("connected peer count:")),
					au.Bold(au.White(fmt.Sprint(peerCount))),
				)
			}
			return nil
		},
		Flags: []cli.Flag{outputFlag()},
	}
}

func peerConnect() *cli.Command {
	return &cli.Command{
		Name:        "connect",
		Usage:       "connect to a peer",
		Description: "connect to a peer by its specified multiaddress",
		Action: func(c *cli.Context) error {
			if c.String("multi.address") == "" {
				return errors.New("multi.address flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			return cl.ConnectToPeer(ctx, c.String("multi.address"))
		},
		Flags: []cli.Flag{multiAddrFlag("the multiaddress to connect to")},
	}
}

func peerDisconnect() *cli.Command {
	return &cli.Command{
		Name:        "disconnect",
		Usage:       "disconnect from a peer",
		Description: "disconnect from a peer identified by peer id",
		Action: func(c *cli.Context) error {
			if c.String("peer.id") == "" {
				return errors.New("peer.id flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			return cl.DisconnectFromPeer(ctx, c.String("peer.id"))
		},
		Flags: []cli.Flag{peerIDFlag("the remote libp2p peer id")},
	}
}

func peerIsConnected() *cli.Command {
	return &cli.Command{
		Name:        "is-connected",
		Usage:       "check if we are connected to a peer",
		Description: "check if we are connected to a peer by it's peer id",
		Action: func(c *cli.Context) error {
			if c.String("peer.id") == "" {
				return errors.New("peer.id flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			connected, err := cl.Connected(ctx, c.String("peer.id"))
			if err != nil {
				return err
			}
			if c.String("output") == "monitor" {
				fmt.Printf(
					"%s\n",
					au.Bold(au.Green(fmt.Sprint(connected))),
				)
			} else {
				fmt.Printf(
					"%s %s\n",
					au.Bold(au.Green("connected:")),
					au.Bold(au.White(fmt.Sprint(connected))),
				)
			}
			return nil
		},
		Flags: []cli.Flag{peerIDFlag("the remote libp2p peer id"), outputFlag()},
	}
}