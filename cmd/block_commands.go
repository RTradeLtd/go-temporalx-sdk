package cmd

import (
	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/urfave/cli/v2"
)

func loadBlockCommand() *cli.Command {
	return &cli.Command{
		Name:        "blockstore",
		Aliases:     []string{"block"},
		Usage:       "manual block manipulation",
		Description: "a low-level API allowing manipulation of blocks",
		Subcommands: cli.Commands{
			&cli.Command{
				Name:  "put",
				Usage: "put data into a block",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.Blockstore(c.Context, &pb.BlockstoreRequest{
						RequestType: pb.BSREQTYPE_BS_PUT,
						Data:        [][]byte{[]byte(c.String("data"))},
						HashFunc:    c.String("multihash"),
					})
					if err != nil {
						return err
					}
					for _, block := range resp.GetBlocks() {
						print(
							"%s: %s\n",
							getArgs("hash", block.GetCid()),
						)
					}
					return nil
				},
				Flags: []cli.Flag{
					DataFlag("data to store inside the block"),
					MultiHashFlag(),
				},
			},
			&cli.Command{
				Name:  "get",
				Usage: "get the contents of a block",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.Blockstore(c.Context, &pb.BlockstoreRequest{
						RequestType: pb.BSREQTYPE_BS_GET,
						Cids:        []string{c.String("cid")},
					})
					if err != nil {
						return err
					}
					for _, block := range resp.GetBlocks() {
						print(
							"%s: %s\n",
							getArgs("data", string(block.GetData())),
						)
					}
					return nil
				},
				Flags: []cli.Flag{
					CidFlag("hash of the block to lookup"),
				},
			},
			&cli.Command{
				Name:        "size",
				Usage:       "get the size of a block",
				Description: "returns the size of the block in bytes",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.Blockstore(c.Context, &pb.BlockstoreRequest{
						RequestType: pb.BSREQTYPE_BS_GET_STATS,
						Cids:        []string{c.String("cid")},
					})
					if err != nil {
						return err
					}
					for _, block := range resp.GetBlocks() {
						print(
							"%s: %v(%s)\n",
							getArgs(c.String("cid"), "bytes", block.GetSize_()),
						)
					}
					return nil
				},
				Flags: []cli.Flag{
					CidFlag("hash of the block to lookup"),
				},
			},
			&cli.Command{
				Name:  "has",
				Usage: "check if are storing the block",
				Action: func(c *cli.Context) error {
					cl, err := client.NewClient(optsFromFlags(c))
					if err != nil {
						return err
					}
					resp, err := cl.Blockstore(c.Context, &pb.BlockstoreRequest{
						RequestType: pb.BSREQTYPE_BS_HAS,
						Cids:        []string{c.String("cid")},
					})
					if err != nil {
						return err
					}
					for _, block := range resp.GetBlocks() {
						if block.GetCid() == c.String("cid") {
							print(
								"%s: %s\n",
								getArgs(c.String("cid"), "true"),
							)
							return nil
						}
					}
					print(
						"%s: %s\n",
						getArgs(c.String("cid"), "false"),
					)
					return nil
				},
				Flags: []cli.Flag{
					CidFlag("hash of the block to lookup"),
				},
			},
		},
	}
}
