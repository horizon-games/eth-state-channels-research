package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
)

type EcdsaSignature struct {
	R, S *big.Int
}

type Signature struct {
	V uint8    `json:"v"`
	R [32]byte `json:"r"`
	S [32]byte `json:"s"`
}

func (s *Signature) MarshalJSON() ([]byte, error) {
	object := make(map[string]interface{})
	object["v"] = s.V
	object["r"] = base64.StdEncoding.EncodeToString(s.R[:])
	object["s"] = base64.StdEncoding.EncodeToString(s.S[:])
	return json.Marshal(object)
}

func (s *Signature) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	s.V = uint8(object["v"].(float64))

	rBytes, err := base64.StdEncoding.DecodeString(object["r"].(string))
	if err != nil {
		return err
	}
	copy(s.R[:], rBytes)

	sBytes, err := base64.StdEncoding.DecodeString(object["s"].(string))
	if err != nil {
		return err
	}
	copy(s.S[:], sBytes)

	return nil
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
