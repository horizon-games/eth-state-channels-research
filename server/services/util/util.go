package util

import "encoding/json"

func Jsonify(o interface{}) string {
	byteArray, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(byteArray)
}
