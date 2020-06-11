package cmd

import (
	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/ipfs/go-cid"
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
					resp, err := cl.NodeAPIClient.Dag(c.Context, &pb.DagRequest{
						RequestType:         pb.DAGREQTYPE_DAG_PUT,
						Data:                []byte(c.String("data")),
						ObjectEncoding:      c.String("object.encoding"),
						SerializationFormat: c.String("serialization.format"),
						Hash:                c.String("multihash"),
					})
					if err != nil {
						return err
					}
					print(
						"%s: %s\n",
						getArgs("hash(es)", resp.GetHashes()),
					)
					return nil
				},
				Flags: []cli.Flag{
					DataFlag("data to store inside the dag"),
					SerializationFormatFlag(),
					ObjectEncodingFlag(),
					MultiHashFlag(),
				},
			},
			&cli.Command{
				Name:        "get",
				Usage:       "get the contents of the IPLD object",
				Description: "this is essentially like the unix `cat` command, except for IPLD",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.NodeAPIClient.Dag(c.Context, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_GET,
						Hash:        c.String("cid"),
					})
					if err != nil {
						return err
					}
					print(
						"%s: %s\n",
						getArgs(
							"data", string(resp.GetRawData()),
						),
					)
					return nil
				},
				Flags: []cli.Flag{
					CidFlag("the ipld node to retrieve"),
				},
			},
			&cli.Command{
				Name:  "add-links",
				Usage: "add links an IPLD object",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.NodeAPIClient.Dag(c.Context, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_ADD_LINKS,
						Hash:        c.String("cid"),
						Links: map[string]string{
							c.String("link.name"): c.String("link.cid"),
						},
					})
					if err != nil {
						return err
					}
					print(
						"%s: %s\n",
						getArgs("hash(es)", resp.GetHashes()),
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
					resp, err := cl.NodeAPIClient.Dag(c.Context, &pb.DagRequest{
						RequestType: pb.DAGREQTYPE_DAG_GET_LINKS,
						Hash:        c.String("cid"),
					})
					if err != nil {
						return err
					}
					for _, link := range resp.GetLinks() {
						gocid, err := cid.Parse(link.GetHash())
						if err != nil {
							return err
						}
						print(
							"%s: %s\n\t- %s: %s\n\t- %s: %v\n",
							getArgs(
								"name", link.GetName(), "cid", gocid.String(), "size", link.GetSize_(),
							),
						)
					}
					return nil
				},
				Flags: []cli.Flag{CidFlag("ipld cid to get links from")},
			},
		},
	}
}
