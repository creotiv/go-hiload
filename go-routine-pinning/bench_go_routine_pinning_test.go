package goroutinepinning

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	workers    = 8
	iterations = 100_000
)

// Hot shared cache line
type hotCounter struct {
	v uint64
	_ [56]byte // pad to full cache line (64B)
}

func Benchmark_UnpinnedWorkers(b *testing.B) {
	var counter hotCounter

	runtime.GOMAXPROCS(runtime.NumCPU())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()

				for j := 0; j < iterations; j++ {
					atomic.AddUint64(&counter.v, 1)
				}
			}()
		}

		wg.Wait()
	}
}

func Benchmark_PinnedWorkers(b *testing.B) {
	var counter hotCounter

	runtime.GOMAXPROCS(runtime.NumCPU())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()

				runtime.LockOSThread()
				defer runtime.UnlockOSThread()

				for j := 0; j < iterations; j++ {
					atomic.AddUint64(&counter.v, 1)
				}
			}()
		}

		wg.Wait()
	}
}
