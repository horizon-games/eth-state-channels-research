package matcher_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/matcher"
)

func buildService() *matcher.Service {
	gopath := os.Getenv("GOPATH")
	service := matcher.NewService(
		&config.ENVConfig{
			WorkingDir: fmt.Sprintf("%s/%s", gopath, "src/github.com/horizon-games/arcadeum/server/"),
		},
		&config.MatcherConfig{
			PrivKeyFile: "etc/matcher/ec-secp256k1-priv.key",
		},
		&config.ETHConfig{
			NodeURL: "http://localhost:8545",
		},
		&config.ArcadeumConfig{
			ContractAddress: "0x0230CfC895646d34538aE5b684d76Bf40a8B8B88",
		},
		&config.RedisConfig{})
	return service
}
func TestVerifySignature(t *testing.T) {
	service := buildService()
	compact, _ := matcher.Compact("bits", "and", "bytes")
	hash := crypto.Keccak256(compact)
	path := fmt.Sprintf("%s/%s", service.ENV.WorkingDir, service.Config.PrivKeyFile)
	privkey, _ := crypto.LoadECDSA(path)
	r, s, _ := service.SignElliptic("bits", "and", "bytes")
	if !ecdsa.Verify(&privkey.PublicKey, hash, r, s) {
		t.Error("Failed to verify sig")
	}

}
func TestIToBString(t *testing.T) {
	s := "This is a sentence."
	asbytes, err := matcher.IToB(s)
	if err != nil {
		t.Error(err)
		return
	}
	c := string(asbytes)
	if s != c {
		t.Error("invalid string: ", c)
	}
}
func TestIToBUInt32(t *testing.T) {
	s := uint32(3)
	asbytes, err := matcher.IToB(s)
	if err != nil {
		t.Error(err)
		return
	}
	c := binary.BigEndian.Uint32(asbytes)
	if s != c {
		t.Error("invalid value: ", c)
	}
}
func TestIToBInt32(t *testing.T) {
	s := int32(-3)
	asbytes, err := matcher.IToB(s)
	if err != nil {
		t.Error(err)
		return
	}
	c := read_int32(asbytes)
	if s != c {
		t.Error("invalid value: ", c)
	}
}
func TestIToBInt(t *testing.T) {
	s := int(-3)
	asbytes, err := matcher.IToB(s)
	if err != nil {
		t.Error(err)
		return
	}
	c := read_int32(asbytes)
	if int32(s) != c {
		t.Error("invalid value: ", c)
	}
}
func TestIToBBytes(t *testing.T) {
	s := []byte("these are bytes")
	asbytes, err := matcher.IToB(s)
	if err != nil {
		t.Error(err)
		return
	}
	if binary.BigEndian.Uint64(s) != binary.BigEndian.Uint64(asbytes) {
		t.Error("invalid value: ", binary.BigEndian.Uint32(asbytes))
	}
}
func TestCompactOne(t *testing.T) {
	s := []byte("these are bytes")
	res, err := matcher.Compact(s)
	if err != nil {
		t.Error(err)
		return
	}
	if string(res) != string(s) {
		t.Error("Wrong compact value ", string(res))
	}
}
func TestCompactMany(t *testing.T) {
	s := []byte("these are bytes")
	res, err := matcher.Compact("these", " are", " bytes")
	if err != nil {
		t.Error(err)
		return
	}
	if string(res) != string(s) {
		t.Error("Wrong compact value ", string(res))
	}
}

type Sig struct {
	V uint8  `json:"v"`
	R []byte `json:"r"`
	S []byte `json:"s"`
}

type Test struct {
	Signature Sig `json:"signature"`
}

func read_int32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}
