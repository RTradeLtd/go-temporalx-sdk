package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	nemonify "github.com/bonedaddy/nemonify"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

// LoadUtilCommands returns generalized utility commands
func LoadUtilCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "util",
			Usage:       "generalized utility functions",
			Description: "provides a series of utility commands not strictly related to TemporalX but useless none the less",
			Subcommands: cli.Commands{createKey(), peerIDFromKey()},
		},
	}
}

func createKey() *cli.Command {
	return &cli.Command{
		Name:        "new-ipfs-key",
		Aliases:     []string{"nik"},
		Usage:       "creates an ipfs key, returning the peerID and hex encoded private key",
		Description: "if the save.mnemonic flag is not empty, a mnemonic phrase of the hex encoded key is save at that path",
		Action: func(c *cli.Context) error {
			pk, pid, err := createIPFSKey(c.String("key.type"), c.Int("key.size"))
			if err != nil {
				return err
			}
			pkBytes, err := pk.Bytes()
			if err != nil {
				return err
			}
			hexPK := hex.EncodeToString(append(pkBytes[0:0:0], pkBytes...))
			fmt.Printf(
				"%s\n%s\n%s\n%s\n",
				au.Bold(au.Green("hex encoded key: ")),
				au.Bold(au.White(hexPK)),
				au.Bold(au.Green("peerID: ")),
				au.Bold(au.White(pid.String())),
			)
			if c.String("save.mnemonic") != "" {
				mnemonic, err := nemonify.ToMnemonic(string(pkBytes))
				if err != nil {
					return err
				}
				return ioutil.WriteFile(c.String("save.mnemonic"), []byte(mnemonic), os.FileMode(0640))
			}
			return nil
		},
		Flags: []cli.Flag{KeyName(), KeyType(), KeySize(), MnemonicFlag()},
	}
}

func peerIDFromKey() *cli.Command {
	return &cli.Command{
		Name:    "peerid-from-key",
		Aliases: []string{"pfk"},
		Usage:   "returns the peerID associated with a private key",
		Action: func(c *cli.Context) error {
			if c.String("input.file") == "" {
				return errors.New("input.file flag is empty")
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
			pid, err := peer.IDFromPrivateKey(pk)
			if err != nil {
				return err
			}
			_, err = fmt.Printf(
				"%s\n%s\n",
				au.Bold(au.Green("peerID: ")),
				au.Bold(au.White(pid.String())),
			)
			return err
		},
		Flags: []cli.Flag{InputFileFlag(), IsHexEncodedFlag(), IsMnemonicEncodedFlag()},
	}
}
