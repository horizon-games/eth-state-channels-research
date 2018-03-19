package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestPubFromHashAndSignature(t *testing.T) {
	message := "This message has been signed!"
	signature := "0xb3825921747a8d3c024c4567d353c25b704e6da53c8bfffa99770ab66aecef0068758554798165e24a7dec85c5e209cc90a19bed424d64b16fdb2767cd771c2c1b"
	address := "0xa5B06b0FF4FBF5D8C5e56F4a6783d28AF72a9a0d"
	addr, err := Address(message, signature)
	if err != nil {
		t.Error(err)
		return
	}
	if addr != address {
		t.Error(fmt.Sprintf("wrong address: %v", addr))
		return
	}
}

type Test struct {
	Address common.Address `json:"address,string"`
	Arr     []byte         `json:"arr,string"`
}

func TestSerial(t *testing.T) {
	test := &Test{}
	s := `
		{
			"address":"0x9700f0f7440179737bb000fba85eba0e5674267b", 
			"arr": "VGhpcyB3b3JrcyE="
		}
	`
	err := json.Unmarshal([]byte(s), test)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if strings.ToLower(test.Address.Hex()) != strings.ToLower("0x9700f0f7440179737bb000fba85eba0e5674267b") {
		t.Errorf("Error hex value %s", test.Address.Hex())
	}
	if string(test.Arr) != "This works!" {
		t.Errorf("Error array value %s", string(test.Arr))
	}

}
