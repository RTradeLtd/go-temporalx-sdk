package cmd

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
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
			var pk crypto.PrivKey
			switch strings.ToLower(c.String("key.type")) {
			case "rsa":
				pk, _, err = crypto.GenerateKeyPair(
					crypto.RSA,
					4096,
				)
			case "ed25519":
				pk, _, err = crypto.GenerateKeyPair(
					crypto.Ed25519,
					256,
				)
			case "ecdsa":
				pk, _, err = crypto.GenerateKeyPair(
					crypto.ECDSA,
					256,
				)
			default:
				return errors.New("key.type flag is empty or contains incorrect value")
			}
			if err != nil {
				return err
			}
			pid, err := peer.IDFromPrivateKey(pk)
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
		Flags: []cli.Flag{keyName(), keyType()},
	}
}
