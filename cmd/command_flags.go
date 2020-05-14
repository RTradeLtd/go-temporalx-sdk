package cmd

import "github.com/urfave/cli/v2"

// command-line flags used by multiple different tex-cli commands

// MultiAddrFlag indicates the multi address to use
func MultiAddrFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "multi.address",
		Aliases: []string{"ma"},
		Usage:   usage,
		EnvVars: []string{"TEX_MULTI_ADDRESS", "TEX_MA"},
	}
}

// PeerIDFlag indicates the peerID
func PeerIDFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "peer.id",
		Aliases: []string{"pid"},
		Usage:   usage,
		EnvVars: []string{"TEX_PEER_ID", "TEX_PID"},
	}
}

// CidFlag is used to indicate the cid to proccess
func CidFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "cid",
		Usage:   usage,
		EnvVars: []string{"TEX_CID"},
	}
}

// PrintProgressFlag enables printing the progress of uploads
func PrintProgressFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "print.progress",
		Aliases: []string{"pp"},
		Usage:   "print progress of uploads/downloads",
		EnvVars: []string{"TEX_PRINT_PROGRESS", "TEX_PP"},
	}
}

// OutputFlag helps to control the style of output
func OutputFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "output",
		Usage:   "control output, accepts 'print' or 'monitor'",
		Value:   "print",
		EnvVars: []string{"TEX_OUTPUT"},
	}
}

// KeyName indicates the name of the key
func KeyName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.name",
		Aliases: []string{"kn"},
		Usage:   "name of the key used in the keystore",
		EnvVars: []string{"TEX_KEY_NAME", "TEX_KN"},
	}
}

// KeyType indicates the type of key
func KeyType() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "key.type",
		Aliases: []string{"kt"},
		Usage:   "type of key: ed25519, ecdsa, rsa, secp256k1",
		EnvVars: []string{"TEX_KEY_TYPE", "TEX_KT"},
	}
}

// KeySize indicates the size of a key, default size suitable for all but RSA
func KeySize() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "key.size",
		Value:   256,
		Aliases: []string{"ks"},
		Usage:   "size of key in bytes",
		EnvVars: []string{"TEX_KEY_SIZE", "TEX_KS"},
	}
}

// MnemonicFlag allows saving data as a mnemonic
func MnemonicFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "save.mnemonic",
		Aliases: []string{"sm"},
		Usage:   "save mnemonic to `PATH` if not empty",
		Value:   "",
		EnvVars: []string{"TEX_SAVE_MNEMONIC", "TEX_SM"},
	}
}

// InputFileFlag allows reading data from a file
func InputFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "input.file",
		Aliases: []string{"in.fi", "if"},
		Usage:   "load data contained in file at `PATH`",
		EnvVars: []string{"TEX_INPUT_FILE", "TEX_IF"},
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
		EnvVars: []string{"TEX_MNEMONIC_ENCODED"},
	}
}

// LinkNameFlag indicates the name of a link in IPLD object management
func LinkNameFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "link.name",
		Aliases: []string{"ln"},
		Usage:   usage,
		EnvVars: []string{"TEX_LINK_NAME"},
	}
}

// LinkCidFlag indicates the cid of a link
func LinkCidFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "link.cid",
		Aliases: []string{"lc"},
		Usage:   usage,
		EnvVars: []string{"TEX_LINK_CID"},
	}
}

// DataFlag indicates some arbitrary string data
func DataFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "data",
		Aliases: []string{"d"},
		Usage:   usage,
		EnvVars: []string{"TEX_DATA"},
	}
}

// ObjectEncodingFlag enables specifying alternate encoding formats for creating objects
func ObjectEncodingFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "object.encoding",
		Aliases: []string{"obj.enc", "oe"},
		Usage:   "the type of encoding used to create an object",
		Value:   "raw",
	}
}

// SerializationFormatFlag enables specifying alternate serialization formats
func SerializationFormatFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "serialization.format",
		Aliases: []string{"serial.form", "sfmt", "st"},
		Usage:   "the type of serialization to use when creating the object",
		Value:   "raw",
		EnvVars: []string{"TEX_SERIALIZATION_FORMAT", "TEX_SERIAL_FORM"},
	}
}

// MultiHashFlag enables specifying the multihash to use
func MultiHashFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "multihash",
		Aliases: []string{"mhash", "mh"},
		Usage:   "specify the multihash to use in object  creation",
		Value:   "sha2-256",
		EnvVars: []string{"TEX_MULTI_HASH"},
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
			EnvVars: []string{"TEX_ALL"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "print protocol, listen and target information. used by ls",
			EnvVars: []string{"TEX_VERBOSE"},
		},
		&cli.BoolFlag{
			Name:    "custom.protocols",
			Usage:   "disables requiring /x/ prefix. used by: listen, forward",
			EnvVars: []string{"TEX_CUSTOM_PROTOCOLS"},
		},
		&cli.BoolFlag{
			Name:    "report.peerid",
			Usage:   "send base58 peerID to target. used by: listen",
			EnvVars: []string{"TEX_REPORT_PEERID"},
		},
		&cli.StringFlag{
			Name:    "protocol.name",
			Usage:   "match/set protocol name. used by: close, forward, listen",
			EnvVars: []string{"TEX_PROTOCOL_NAME"},
		},
		&cli.StringFlag{
			Name:    "listen.address",
			Usage:   "match/set against listen address. used by: close, forward",
			EnvVars: []string{"TEX_LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "target.address",
			Usage:   "match/set against target address. used by: close, forward, listen",
			EnvVars: []string{"TEX_TARGET_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "remote.address",
			Usage:   "note currently used but here for compatability",
			EnvVars: []string{"TEX_REMOTE_ADDRESS"},
		},
	}, cmdFlag)
}
