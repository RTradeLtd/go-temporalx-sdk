package cmd

import "github.com/urfave/cli/v2"

// command-line flags used by multiple different tex-cli commands

// MultiAddrFlag indicates the multi address to use
func MultiAddrFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "multi.address",
		Aliases: []string{"ma"},
		Usage:   usage,
		EnvVars: []string{"MULTI_ADDRESS", "MA"},
	}
}

// PeerIDFlag indicates the peerID
func PeerIDFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "peer.id",
		Aliases: []string{"pid"},
		Usage:   usage,
		EnvVars: []string{"PEER_ID", "PID"},
	}
}

// CidFlag is used to indicate the cid to proccess
func CidFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "cid",
		Usage:   usage,
		EnvVars: []string{"CID"},
	}
}

// PrintProgressFlag enables printing the progress of uploads
func PrintProgressFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "print.progress",
		Aliases: []string{"pp"},
		Usage:   "print progress of uploads/downloads",
		EnvVars: []string{"PRINT_PROGRESS", "PP"},
	}
}

// OutputFlag helps to control the style of output
func OutputFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "output",
		Usage:   "control output, accepts 'print' or 'monitor'",
		Value:   "print",
		EnvVars: []string{"OUTPUT"},
	}
}

// KeyName indicates the name of the key
func KeyName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.name",
		Aliases: []string{"kn"},
		Usage:   "name of the key used in the keystore",
		EnvVars: []string{"KEY_NAME", "KN"},
	}
}

// KeyType indicates the type of key
func KeyType() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.type",
		Aliases: []string{"kt"},
		Usage:   "type of key: ed25519, ecdsa, rsa, secp256k1",
		EnvVars: []string{"KEY_TYPE", "KT"},
	}
}

// KeySize indicates the size of a key, default size suitable for all but RSA
func KeySize() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "key.size",
		Value:   256,
		Aliases: []string{"ks"},
		Usage:   "size of key in bytes",
		EnvVars: []string{"KEY_SIZE", "KS"},
	}
}

// MnemonicFlag allows saving data as a mnemonic
func MnemonicFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "save.mnemonic",
		Aliases: []string{"sm"},
		Usage:   "save mnemonic to `PATH` if not empty",
		Value:   "",
		EnvVars: []string{"SAVE_MNEMONIC", "SM"},
	}
}

// InputFileFlag allows reading data from a file
func InputFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "input.file",
		Aliases: []string{"in.fi", "if"},
		Usage:   "load data contained in file at `PATH`",
		EnvVars: []string{"INPUT_FILE", "IF"},
	}
}

// IsHexEncodedFlag indicates if the data is hex encoded
func IsHexEncodedFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "hex.encoded",
		Value:   false,
		Usage:   "whether or not the key has been hex encoded",
		EnvVars: []string{"HEX_ENCODED"},
	}
}

// IsMnemonicEncodedFlag indicates if the data is encoded as mnemonic
func IsMnemonicEncodedFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "mnemonic.encoded",
		Value:   true,
		Usage:   "whether or not the key has been converted into a mnemonic",
		EnvVars: []string{"MNEMONIC_ENCODED"},
	}
}

// P2pFlags are used to control p2p stream
// takes in an argument which is a command that should be
// loaded with a default value. This is appended to the default
// p2p command flag list
func P2pFlags(cmdFlag *cli.StringFlag) []cli.Flag {
	if cmdFlag.Value == "" {
		panic("flag value is nil")
	}
	return append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Usage:   "close all listeners. used by: close",
			EnvVars: []string{"ALL"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "print protocol, listen and target information. used by ls",
			EnvVars: []string{"VERBOSE"},
		},
		&cli.BoolFlag{
			Name:    "custom.protocols",
			Usage:   "disables requiring /x/ prefix. used by: listen, forward",
			EnvVars: []string{"CUSTOM_PROTOCOLS"},
		},
		&cli.BoolFlag{
			Name:    "report.peerid",
			Usage:   "send base58 peerID to target. used by: listen",
			EnvVars: []string{"REPORT_PEERID"},
		},
		&cli.StringFlag{
			Name:    "protocol.name",
			Usage:   "match/set protocol name. used by: close, forward, listen",
			EnvVars: []string{"PROTOCOL_NAME"},
		},
		&cli.StringFlag{
			Name:    "listen.address",
			Usage:   "match/set against listen address. used by: close, forward",
			EnvVars: []string{"LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "target.address",
			Usage:   "match/set against target address. used by: close, forward, listen",
			EnvVars: []string{"TARGET_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "remote.address",
			Usage:   "note currently used but here for compatability",
			EnvVars: []string{"REMOTE_ADDRESS"},
		},
	}, cmdFlag)
}
