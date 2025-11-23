package zeroallocationparsing

import (
	"github.com/valyala/fastjson"
)

// Zero-allocation fastjson library
func ParseBatchFast(b []byte) ([]LogRecord, error) {
	var out []LogRecord
	var p fastjson.Parser
	v, err := p.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	arr := v.GetArray()
	if cap(out) < len(arr) {
		out = make([]LogRecord, len(arr))
	} else {
		out = out[:len(arr)]
	}

	for i, item := range arr {
		out[i].TS = item.GetInt64("ts")
		out[i].Msg = string(item.GetStringBytes("msg"))
		out[i].Lev = string(item.GetStringBytes("lev"))
		out[i].App = string(item.GetStringBytes("app"))
	}

	return out, nil
}
