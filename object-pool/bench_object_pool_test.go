package bufferpool

import (
	"runtime"
	"sync"
	"testing"
)

const bufSize = 16 * 1024 // 16 KB realistic WAL/log page

// Sink is global so compiler cannot eliminate work
var sink []byte

func consume(b *[]byte) {
	// Store reference so slice escapes
	sink = *b
	// Prevent compiler from optimizing away usage
	runtime.KeepAlive(b)
}

// --- Section: No pool ---

func BenchmarkNoPoolReal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := make([]byte, bufSize)
		consume(&buf)
	}
}

// --- Section: sync.Pool ---

var pool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, bufSize)
		return &b // <- only pointers!!!
	},
}

func BenchmarkPoolReal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*[]byte)
		consume(buf)
		pool.Put(buf)
	}
}

// --- Section: Parallel benchmarks ---

func BenchmarkNoPoolParallelReal(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := make([]byte, bufSize)
			consume(&buf)
		}
	})
}

func BenchmarkPoolParallelReal(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get().(*[]byte)
			consume(buf)
			pool.Put(buf)
		}
	})
}
