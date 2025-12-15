package goroutinestack

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	goroutineCount = 256
	bufferSize     = 32 * 1024 // 32KB scratch buffer per goroutine
)

var sink uint64

type memSnapshot struct {
	heapInuse  uint64
	stackInuse uint64
	sysTotal   uint64
}

func BenchmarkGoroutineHeapBuffersRetained(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		runtime.GC()
		before := snapshotMem()
		retained := make([][]byte, goroutineCount)
		b.StartTimer()

		runHeapBuffers(retained)

		b.StopTimer()
		runtime.GC()
		after := snapshotMem()
		b.ReportMetric(float64(delta(after.stackInuse, before.stackInuse)), "stack_inuse_bytes")
		b.ReportMetric(float64(delta(after.heapInuse, before.heapInuse)), "heap_inuse_bytes")
		b.ReportMetric(float64(after.sysTotal), "process_sys_bytes")

		for i := range retained {
			retained[i] = nil
		}
		runtime.KeepAlive(retained)
		b.StartTimer()
	}
}

func BenchmarkGoroutineStackBuffers(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		runtime.GC()
		before := snapshotMem()
		b.StartTimer()

		runStackBuffers()

		b.StopTimer()
		runtime.GC()
		after := snapshotMem()
		b.ReportMetric(float64(delta(after.stackInuse, before.stackInuse)), "stack_inuse_bytes")
		b.ReportMetric(float64(delta(after.heapInuse, before.heapInuse)), "heap_inuse_bytes")
		b.ReportMetric(float64(after.sysTotal), "process_sys_bytes")
		b.StartTimer()
	}
}

func runStackBuffers() {
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			var buf [bufferSize]byte
			local := uint64(0)
			for j := 0; j < len(buf); j += 512 {
				buf[j]++
				local += uint64(buf[j])
			}
			atomic.AddUint64(&sink, local)
		}()
	}

	wg.Wait()
}

func runHeapBuffers(retained [][]byte) {
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		idx := i
		go func() {
			defer wg.Done()

			buf := make([]byte, bufferSize)
			local := uint64(0)
			for j := 0; j < len(buf); j += 512 {
				buf[j] = byte(j)
				local += uint64(buf[j])
			}

			retained[idx] = buf
			atomic.AddUint64(&sink, local)
		}()
	}

	wg.Wait()
}

func snapshotMem() memSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return memSnapshot{
		heapInuse:  m.HeapInuse,
		stackInuse: m.StackInuse,
		sysTotal:   m.Sys,
	}
}

func delta(after, before uint64) uint64 {
	if after > before {
		return after - before
	}
	return 0
}
