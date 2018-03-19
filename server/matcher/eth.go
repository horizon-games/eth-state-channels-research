package matcher

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func Hash(message string) []byte {
	m := fmt.Sprintf("\x19Ethereum Signed Message:\n%v%v", len(message), message)
	return crypto.Keccak256([]byte(m))
}

// Return the hex address (pub key) of a signed message given the message and signature
func Address(message string, signature string) (string, error) {
	h := Hash(message)

	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}
	s, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	// https://github.com/ethereum/go-ethereum/issues/2053
	// https://github.com/ethereum/go-ethereum/blob/v1.7.3/core/types/transaction_signing.go#L201
	// This will form the V value (recovery ID) at s[64]
	s[64] -= 27

	pubkey, err := crypto.SigToPub(h, s)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*pubkey).Hex(), nil
}
