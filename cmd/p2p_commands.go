package cmd

import (
	"errors"
	"fmt"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

func loadP2PCommand() *cli.Command {
	return &cli.Command{
		Name:        "p2p",
		Usage:       "control tcp/udp libp2p tunnels",
		Description: "allows forwarding to/from libp2p, and local tcp/udp services. Think SSH tunnels but for libp2p",
		Subcommands: cli.Commands{
			&cli.Command{
				Name:  "listen",
				Usage: "listen creates a libp2p service that will send connections to localhost",
				Description: `
Used to create a libp2p service, which forwards connected to the target address.
It requires specifying the libp2p protocol handler, and by default must be prefixed with "/x".
The following command creates a libp2p service called dns that forwards connections to 127.0.0.1:53
	tex-cli client node p2p --target.address /ip4/127.0.0.1/tcp/53 --protocol /x/dns				
`,
				Action: p2pAction,
				Flags: P2pFlags(&cli.StringFlag{
					Name:    "command",
					Aliases: []string{"cmd"},
					Value:   "listen",
				}),
			},
			&cli.Command{
				Name:  "forward",
				Usage: "forward connections from a local port to a libp2p service",
				Description: `
Forwards connectinos to a libp2p service, and is the reverse of the litsen command.
It requires specifying the libp2p protocol handler, and by default must be prefixed with "/x".
The following commands forwards connections from 127.0.0.1:53 to a libp2p service called libdns to peer /p2p/temporalxisbest
    tex-cli client node p2p --listen.address /ip4/127.0.0.1:53 --protocol /x/libdns --target.address /p2p/temporalxisbest
`,
				Action: p2pAction,
				Flags: P2pFlags(&cli.StringFlag{
					Name:    "command",
					Aliases: []string{"cmd"},
					Value:   "forward",
				}),
			},
			&cli.Command{
				Name:   "ls",
				Usage:  "list various p2p streams",
				Action: p2pAction,
				Flags: P2pFlags(&cli.StringFlag{
					Name:    "command",
					Aliases: []string{"cmd"},
					Value:   "ls",
				}),
			},
			&cli.Command{
				Name:   "close",
				Usage:  "close libp2p streams",
				Action: p2pAction,
				Flags: P2pFlags(&cli.StringFlag{
					Name:    "command",
					Aliases: []string{"cmd"},
					Value:   "close",
				}),
			},
		},
	}
}

func p2pAction(c *cli.Context) error {
	var cmd pb.P2PREQTYPE
	switch c.String("command") {
	case "listen":
		cmd = pb.P2PREQTYPE_LISTEN
		break
	case "close":
		cmd = pb.P2PREQTYPE_CLOSE
		break
	case "forward":
		cmd = pb.P2PREQTYPE_FORWARD
		break
	case "ls":
		cmd = pb.P2PREQTYPE_LS
		break
	default:
		return errors.New("invalid command option")
	}
	cl, err := client.NewClient(optsFromFlags(c))
	if err != nil {
		return err
	}
	resp, err := cl.NodeAPIClient.P2P(ctx, &pb.P2PRequest{
		RequestType:          cmd,
		All:                  c.Bool("all"),
		Verbose:              c.Bool("verbose"),
		ProtocolName:         c.String("protocol.name"),
		ListenAddress:        c.String("listen.address"),
		TargetAddress:        c.String("target.address"),
		RemoteAddress:        c.String("remote.address"),
		AllowCustomProtocols: c.Bool("custom.protocols"),
		ReportPeerID:         c.Bool("report.peerid"),
	})
	if err != nil {
		return err
	}
	switch c.String("command") {
	case "listen":
		fmt.Printf("%s\n", au.Bold(au.Green("listening")))
		break
	case "close":
		fmt.Printf("%s\n", au.Bold(au.Green("closed")))
		break
	case "forward":
		fmt.Printf("%s\n", au.Bold(au.Green("forwarded")))
		break
	case "ls":
		fmt.Printf(
			"%s: %s\n",
			au.Bold(au.Green("number of streams")),
			au.Bold(au.White(fmt.Sprint(len(resp.GetStreamInfos())))),
		)
		break
	default:
		return errors.New("invalid command option")
	}
	return nil
}
