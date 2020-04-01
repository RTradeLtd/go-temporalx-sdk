package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	nemonify "github.com/bonedaddy/nemonify"
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
			Subcommands: cli.Commands{createKey()},
		},
	}
}

func createKey() *cli.Command {
	return &cli.Command{
		Name:        "new-ipfs-key",
		Aliases:     []string{"nik"},
		Usage:       "creates n ipfs key, returning the peerID and hex encoded private key",
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
			hexPK := hex.EncodeToString(pkBytes)
			fmt.Printf(
				"%s\n%s\n%s\n%s\n",
				au.Bold(au.Green("hex encoded key: ")),
				au.Bold(au.White(hexPK)),
				au.Bold(au.Green("peerID: ")),
				au.Bold(au.White(pid.String())),
			)
			if c.String("save.mnemonic") != "" {
				mnemonic, err := nemonify.ToMnemonic(hexPK)
				if err != nil {
					return err
				}
				return ioutil.WriteFile(c.String("save.mnemonic"), []byte(mnemonic), os.FileMode(0640))
			}
			return nil
		},
		Flags: []cli.Flag{keyName(), keyType(), keySize(), mnemonicFlag()},
	}
}
