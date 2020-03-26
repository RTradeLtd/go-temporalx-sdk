package cmd

import "github.com/urfave/cli/v2"

// command-line flags used by multiple different tex-cli commands

func multiAddrFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "multi.address",
		Aliases: []string{"ma"},
		Usage:   usage,
	}
}

func peerIDFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "peer.id",
		Aliases: []string{"pid"},
		Usage:   usage,
	}
}

func cidFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "cid",
		Usage: usage,
	}
}

func printProgressFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "print.progress",
		Aliases: []string{"pp"},
		Usage:   "print progress of uploads/downloads",
	}
}

func outputFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "output",
		Usage: "control output, accepts 'print' or 'monitor'",
		Value: "print",
	}
}

func keyName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.name",
		Aliases: []string{"kn"},
		Usage:   "name of the key used in the keystore",
	}
}

func keyType() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.type",
		Aliases: []string{"kt"},
		Usage:   "type of key: ed25519, ecdsa, rsa",
	}
}

// takes in an argument which is a command that should be
// loaded with a default value. This is appended to the default
// p2p command flag list
func p2pFlags(cmdFlag *cli.StringFlag) []cli.Flag {
	if cmdFlag.Value == "" {
		panic("flag value is nil")
	}
	return append([]cli.Flag{
		&cli.BoolFlag{
			Name: "all",
		},
		&cli.BoolFlag{
			Name: "verbose",
		},
		&cli.BoolFlag{
			Name: "custom.protocols",
		},
		&cli.BoolFlag{
			Name: "report.peerid",
		},
		&cli.StringFlag{
			Name: "protocol.name",
		},
		&cli.StringFlag{
			Name: "listen.address",
		},
		&cli.StringFlag{
			Name: "target.address",
		},
		&cli.StringFlag{
			Name: "remote.address",
		},
	}, cmdFlag)
}
