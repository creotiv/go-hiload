package genericsvsinterface

import "testing"

const dataLen = 1 << 20 // 1,048,576 ints to make per-iteration overhead visible

var (
	intData = makeIntData()
	anyData = makeAnyData(intData)
	sumSink int
)

type intLike interface {
	~int
}

func makeIntData() []int {
	data := make([]int, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = i % 1024
	}
	return data
}

func makeAnyData(src []int) []any {
	data := make([]any, len(src))
	for i, v := range src {
		data[i] = v
	}
	return data
}

func sumGeneric[T intLike](vals []T) int {
	var s int
	for _, v := range vals {
		s += int(v)
	}
	return s
}

func sumInterfaceAssert(vals []any) int {
	var s int
	for _, v := range vals {
		s += v.(int)
	}
	return s
}

func sumInterfaceTypeSwitch(vals []any) int {
	var s int
	for _, v := range vals {
		switch n := v.(type) {
		case int:
			s += n
		case int64:
			s += int(n)
		}
	}
	return s
}
func BenchmarkGenericSum(b *testing.B) {
	data := intData
	b.ReportAllocs()
	b.ResetTimer()

	var s int
	for i := 0; i < b.N; i++ {
		s = sumGeneric(data)
	}
	sumSink = s
}

func BenchmarkInterfaceTypeAssertion(b *testing.B) {
	data := anyData
	b.ReportAllocs()
	b.ResetTimer()

	var s int
	for i := 0; i < b.N; i++ {
		s = sumInterfaceAssert(data)
	}
	sumSink = s
}

func BenchmarkInterfaceTypeSwitch(b *testing.B) {
	data := anyData
	b.ReportAllocs()
	b.ResetTimer()

	var s int
	for i := 0; i < b.N; i++ {
		s = sumInterfaceTypeSwitch(data)
	}
	sumSink = s
}
