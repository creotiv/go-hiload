package ifacevaluecopy

import (
	"runtime"
	"testing"
)

// bigValue is intentionally large so copying it is visible in benchmarks.
type bigValue struct {
	payload [32 * 1024]byte
}

var (
	ptrSink   *bigValue
	ifaceSink any
)

func consumePointer(v *bigValue) {
	ptrSink = v
	runtime.KeepAlive(v)
}

func consumeInterface(v any) {
	ifaceSink = v
	runtime.KeepAlive(v)
}

func BenchmarkPointerNoInterface(b *testing.B) {
	b.ReportAllocs()

	v := &bigValue{}
	for i := 0; i < b.N; i++ {
		v.payload[0] = byte(i) // mutate to keep compiler from eliding the value
		consumePointer(v)
	}
}

func BenchmarkInterfacePointer(b *testing.B) {
	b.ReportAllocs()

	v := &bigValue{}
	for i := 0; i < b.N; i++ {
		v.payload[0] = byte(i)
		consumeInterface(v) // interface{} holds the pointer, no copy of payload
	}
}

func BenchmarkInterfaceValueCopy(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var v bigValue
		v.payload[0] = byte(i)
		consumeInterface(v) // interface{} boxes the value => copies payload, forces escape
	}
}
