package util

import (
	"encoding/hex"
	"encoding/json"
)

func Jsonify(o interface{}) (string, error) {
	byteArray, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(byteArray), nil
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
