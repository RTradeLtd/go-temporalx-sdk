package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
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
			Subcommands: cli.Commands{keystoreCreate(), keystoreImport()},
		},
	}
}

func keystoreImport() *cli.Command {
	return &cli.Command{
		Name:        "import",
		Usage:       "import a mnemonic phrased private key",
		Description: "enables importing private keys exported as mnemonic phrases",
		Action: func(c *cli.Context) error {
			if c.String("input.file") == "" {
				return errors.New("input.file flag is empty")
			}
			if c.String("key.name") == "" {
				return errors.New("key.name flag is empty")
			}
			contents, err := ioutil.ReadFile(c.String("input.file"))
			if err != nil {
				return err
			}
			var pk crypto.PrivKey
			if c.Bool("hex.encoded") {
				pk, err = hexToKey(string(contents))
			} else if c.Bool("mnemonic.encoded") {
				pk, err = keyFromMnemonic(string(contents))
				if err != nil {
					return err
				}
			} else {
				return errors.New("invalid on-disk format for key")
			}
			pkBytes, err := pk.Bytes()
			if err != nil {
				return err
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			_, err = cl.Keystore(c.Context, &pb.KeystoreRequest{
				RequestType: pb.KSREQTYPE_KS_PUT,
				PrivateKey:  pkBytes,
				Name:        c.String("key.name"),
			})
			return err
		},
		Flags: []cli.Flag{KeyName(), InputFileFlag(), IsHexEncodedFlag(), IsMnemonicEncodedFlag()},
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
			_, err = cl.Keystore(c.Context, &pb.KeystoreRequest{
				RequestType: pb.KSREQTYPE_KS_PUT,
				PrivateKey:  pkBytes,
				Name:        c.String("key.name"),
			})
			return err
		},
		Flags: []cli.Flag{KeyName(), KeyType(), KeySize(), MnemonicFlag()},
	}
}
