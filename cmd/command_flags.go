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
		Usage:   "type of key: ed25519, ecdsa, rsa, secp256k1",
	}
}

func keySize() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "key.size",
		Aliases: []string{"ks"},
		Usage:   "size of key in bytes",
	}
}

func mnemonicFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "save.mnemonic",
		Aliases: []string{"sm"},
		Usage:   "save mnemonic to `PATH` if not empty",
		Value:   "",
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
			Name:  "all",
			Usage: "close all listeners. used by: close",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "print protocol, listen and target information. used by ls",
		},
		&cli.BoolFlag{
			Name:  "custom.protocols",
			Usage: "disables requiring /x/ prefix. used by: listen, forward",
		},
		&cli.BoolFlag{
			Name:  "report.peerid",
			Usage: "send base58 peerID to target. used by: listen",
		},
		&cli.StringFlag{
			Name:  "protocol.name",
			Usage: "match/set protocol name. used by: close, forward, listen",
		},
		&cli.StringFlag{
			Name:  "listen.address",
			Usage: "match/set against listen address. used by: close, forward",
		},
		&cli.StringFlag{
			Name:  "target.address",
			Usage: "match/set against target address. used by: close, forward, listen",
		},
		&cli.StringFlag{
			Name:  "remote.address",
			Usage: "note currently used but here for compatability",
		},
	}, cmdFlag)
}
