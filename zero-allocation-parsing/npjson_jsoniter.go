package zeroallocationparsing

import jsoniter "github.com/json-iterator/go"

// json-iter os low allocation library
func ParseBatchJsonIter(b []byte) ([]LogRecordStd, error) {
	var out []LogRecordStd
	err := jsoniter.Unmarshal(b, &out)
	return out, err
}
