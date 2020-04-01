package cmd

import (
	"errors"
	"strings"

	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// createIPFSKey is a helper function to create an IPFS key
func createIPFSKey(keyType string, keySize int) (crypto.PrivKey, peer.ID, error) {
	var (
		pk  crypto.PrivKey
		err error
	)
	switch strings.ToLower(keyType) {
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
	case "secp256k1":
		pk, _, err = crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	default:
		err = errors.New("key.type flag is empty or contains incorrect value")
	}
	if err != nil {
		return nil, "", err
	}
	pid, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, "", err
	}
	return pk, pid, nil
}
