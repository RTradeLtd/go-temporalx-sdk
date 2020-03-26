package cmd

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	adminURL string
)

// LoadAdminCommands returns an admin api commands object
func LoadAdminCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "admin",
			Usage:       "admin commands",
			Description: "Enables access to AdminAPI",
			Subcommands: cli.Commands{refCount(), gcCommands()},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "admin.url",
					Usage:       "the admin api endpoint",
					Value:       "localhost:9999",
					Destination: &adminURL,
				},
			},
		},
	}
}

func refCount() *cli.Command {
	return &cli.Command{
		Name:  "ref-count",
		Usage: "get reference count information",
		Action: func(c *cli.Context) error {
			if c.String("cid") == "" {
				return errors.New("cid flag empty")
			}
			conn, err := grpc.Dial(adminURL, grpc.WithInsecure())
			if err != nil {
				return err
			}
			admin := pb.NewAdminAPIClient(conn)
			resp, err := admin.RefCount(ctx, &pb.RefCountRequest{Cids: []string{c.String("cid")}})
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s %s\n",
				au.Bold(au.Green("count:")),
				au.Bold(au.White(resp.GetCids()[c.String("cid")])),
			)
			return nil
		},
		Flags: []cli.Flag{cidFlag("cid to lookup")},
	}
}

func gcCommands() *cli.Command {
	return &cli.Command{
		Name:        "gc",
		Usage:       "manage garbage collection",
		Subcommands: cli.Commands{gcControl()},
	}
}

func gcControl() *cli.Command {
	return &cli.Command{
		Name:  "control",
		Usage: "control garbage collection",
		Action: func(c *cli.Context) error {
			if strings.ToLower(c.String("operation")) == "" {
				return errors.New("operation flag is empty")
			}
			conn, err := grpc.Dial(adminURL, grpc.WithInsecure())
			if err != nil {
				return err
			}
			operation := strings.ToUpper(c.String("operation"))
			admin := pb.NewAdminAPIClient(conn)
			resp, err := admin.ManageGC(ctx, &pb.ManageGCRequest{
				Type: pb.GCREQTYPE(
					pb.GCREQTYPE_value[operation],
				),
			})
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s %s\n",
				au.Bold(au.Green("status:")),
				au.Bold(au.White(resp.GetStatus())),
			)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "operation",
				Usage: "the operation to perform",
			},
		},
	}
}
