package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
)

type EcdsaSignature struct {
	R, S *big.Int
}

type Signature struct {
	V uint8  `json:"v"`
	R []byte `json:"r,string"` // unmarshalled base64
	S []byte `json:"s,string"` // unmarshalled base64
}

func (s *Signature) ECDSA() ([]byte, error) {
	return asn1.Marshal(EcdsaSignature{
		new(big.Int).SetBytes(s.R[:]),
		new(big.Int).SetBytes(s.S[:])})
}

func Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}
