package util

import (
	"encoding/json"
	"encoding/hex"
)

func Jsonify(o interface{}) string {
	byteArray, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(byteArray)
}

func DecodeHexString(s string) ([]byte, error) {
	if s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}
	r, err := hex.DecodeString(string(s))
	if err != nil {
		return nil, err
	}
	return r, nil
}
