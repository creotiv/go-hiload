package zeroallocationparsing

import (
	"encoding/json"
)

type LogRecordStd struct {
	TS  int64  `json:"ts"`
	Msg string `json:"msg"`
	Lev string `json:"lev"`
	App string `json:"app"`
}

// Standard library parser (allocs, reflection)
func ParseBatchStd(b []byte) ([]LogRecordStd, error) {
	var out []LogRecordStd
	err := json.Unmarshal(b, &out)
	return out, err
}
