package cmd

import (
	"errors"
	"fmt"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

// LoadKeystoreCommands returns keystore management commands
func LoadKeystoreCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "keystore",
			Usage:       "keystore commands",
			Description: "Enables access to KeystoreAPI",
			Subcommands: cli.Commands{keystoreCreate()},
		},
	}
}

func keystoreCreate() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create a new key",
		Action: func(c *cli.Context) error {
			if c.String("key.name") == "" {
				return errors.New("key.name flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			pk, pid, err := createIPFSKey(c.String("key.type"), c.Int("key.size"))
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s %s\n",
				au.Bold(au.Green("key peer id:")),
				au.Bold(au.White(pid)),
			)
			pkBytes, err := pk.Bytes()
			if err != nil {
				return err
			}
			_, err = cl.Keystore(ctx, &pb.KeystoreRequest{
				RequestType: pb.KSREQTYPE_KS_PUT,
				PrivateKey:  pkBytes,
				Name:        c.String("key.name"),
			})
			return err
		},
		Flags: []cli.Flag{keyName(), keyType(), keySize(), mnemonicFlag()},
	}
}
