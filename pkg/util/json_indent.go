package util

import "encoding/json"

func JsonIndent(v any) string {
	byteVal, _ := json.MarshalIndent(v, " ", " ")
	return string(byteVal)
}
