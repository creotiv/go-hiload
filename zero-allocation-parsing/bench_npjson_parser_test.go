package zeroallocationparsing

import (
	"encoding/json"
	"testing"
)

var testNP []byte

func init() {
	// build sample batch of 200 records
	tmp := make([]map[string]any, 200)
	for i := range tmp {
		tmp[i] = map[string]any{
			"ts":  123456789,
			"msg": "hello",
			"lev": "info",
			"app": "gateway",
		}
	}
	testNP, _ = json.Marshal(tmp)
}

func BenchmarkStdJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := ParseBatchStd(testNP)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
		var dst []LogRecordStd
		json.Unmarshal(testNP, &dst)
	}
}

func BenchmarkZeroAllocationJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := ParseBatchCustom(testNP)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
	}
}

func BenchmarkFastZeroAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := ParseBatchFast(testNP)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
	}
}

func BenchmarkJsonIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := ParseBatchJsonIter(testNP)
		if err != nil {
			b.Fatal(err)
		}
		_ = res
	}
}
