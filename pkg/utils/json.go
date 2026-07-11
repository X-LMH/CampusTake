package utils

import (
	"encoding/json"
)

// ParseJSONToStringSlice JSON 转 []string
func ParseJSONToStringSlice(data []byte) []string {
	var res []string
	_ = json.Unmarshal(data, &res)
	return res
}

func ParseStringSliceToJSON(data []string) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}
