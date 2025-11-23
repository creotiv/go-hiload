package lockfreeringbuffer

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"golang.org/x/sys/cpu"
)

// ====================================
//   Lock-free SPSC padded ring buffer
// ====================================

type Ring struct {
	_    cpu.CacheLinePad
	head atomic.Uint64
	_    cpu.CacheLinePad
	tail atomic.Uint64
	_    cpu.CacheLinePad
	buf  []uint64
	mask uint64
}

func NewRing(size int) *Ring {
	// must be pow2
	buf := make([]uint64, size)
	return &Ring{
		buf:  buf,
		mask: uint64(size - 1),
	}
}

func (r *Ring) Enqueue(v uint64) {
	h := r.head.Load()
	r.buf[h&r.mask] = v
	r.head.Store(h + 1)
}

func (r *Ring) Dequeue() uint64 {
	t := r.tail.Load()
	v := r.buf[t&r.mask]
	r.tail.Store(t + 1)
	return v
}

// ====================================
//            CHANNEL BENCH
// ====================================

func benchChannel(b *testing.B, ch chan uint64) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < b.N; i++ {
			ch <- uint64(i)
		}
	}()

	// Consumer
	go func() {
		defer wg.Done()
		<-start
		var sum uint64
		for i := 0; i < b.N; i++ {
			sum += <-ch
		}
		runtime.KeepAlive(sum)
	}()

	b.ResetTimer()
	close(start)
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	// Bounded buffered channel
	ch := make(chan uint64, 1024)
	benchChannel(b, ch)
}

// ====================================
//          RING BUFFER BENCH
// ====================================

func benchRing(b *testing.B, r *Ring) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < b.N; i++ {
			r.Enqueue(uint64(i))
		}
	}()

	// Consumer
	go func() {
		defer wg.Done()
		<-start
		var sum uint64
		for i := 0; i < b.N; i++ {
			sum += r.Dequeue()
		}
		runtime.KeepAlive(sum)
	}()

	b.ResetTimer()
	close(start)
	wg.Wait()
}

func BenchmarkRingBuffer(b *testing.B) {
	r := NewRing(1024)
	benchRing(b, r)
}
