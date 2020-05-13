package cmd

import (
	"fmt"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/ipfs/go-cid"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

func loadDagCommand() *cli.Command {
	return &cli.Command{
		Name:        "dag",
		Usage:       "manual dag manipulation",
		Description: "a low-level API allowing manipulation of IPLD objects",
		Subcommands: cli.Commands{
			&cli.Command{
				Name:  "put",
				Usage: "put data into an IPLD object",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.NodeAPIClient.Dag(ctx, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_PUT,
						Data:        []byte(c.String("data")),
					})
					if err != nil {
						return err
					}
					fmt.Printf(
						"%s: %s\n",
						au.Bold(au.Green("hash(es)")),
						au.Bold(au.White(resp.GetHashes())),
					)
					return nil
				},
				Flags: []cli.Flag{DataFlag("data to store inside the dag")},
			},
			&cli.Command{
				Name:  "add-links",
				Usage: "add links an IPLD object",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.NodeAPIClient.Dag(ctx, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_ADD_LINKS,
						Hash:        c.String("root.cid"),
						Links: map[string]string{
							c.String("link.name"): c.String("link.cid"),
						},
					})
					if err != nil {
						return err
					}
					fmt.Printf(
						"%s: %s\n",
						au.Bold(au.Green("hash(es)")),
						au.Bold(au.White(resp.GetHashes())),
					)
					return nil
				},
				Flags: []cli.Flag{
					CidFlag("the ipld node to add links to"),
					LinkCidFlag("the cid to link to"),
					LinkNameFlag("the name to give the link"),
				},
			},
			&cli.Command{
				Name:  "get-links",
				Usage: "get links in an IPLD object",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.NodeAPIClient.Dag(ctx, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_GET_LINKS,
						Hash:        c.String("cid"),
					})
					if err != nil {
						return err
					}
					var cids []string
					for _, link := range resp.GetLinks() {
						gocid, err := cid.Parse(link.GetHash())
						if err != nil {
							return err
						}
						cids = append(cids, gocid.String())
					}
					fmt.Printf(
						"%s: %s\n",
						au.Bold(au.Green("link(s)")),
						au.Bold(au.White(cids)),
					)
					return nil
				},
				Flags: []cli.Flag{CidFlag("ipld cid to get links from")},
			},
		},
	}
}
