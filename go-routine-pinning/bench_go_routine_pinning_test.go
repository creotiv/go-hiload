package goroutinepinning

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// Create scheduler pressure to force unwanted migrations.
// We need this as on new systems without this it hard to see
// difference in synthetic benchmark. That happen because system
// optimize work to low routine migration between CPU.
// So we add additiona preasure to increase migration.
func hammerScheduler() {
	for i := 0; i < runtime.NumCPU()*4; i++ {
		go func() {
			for {
				runtime.Gosched() // constantly yielding
			}
		}()
	}
}

func cpuHotWork() uint64 {
	var x uint64 = 1
	for i := 0; i < 500; i++ {
		x = (x*6364136223846793005 + 1) // strong dependency chain
	}
	runtime.KeepAlive(x)
	return x
}

// ============================
//     NO PINNING (migrates)
// ============================

func BenchmarkNoPinning(b *testing.B) {
	hammerScheduler() // FORCE MIGRATIONS
	time.Sleep(50 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < b.N; i++ {
			cpuHotWork()
			runtime.Gosched() // encourage migration
		}
	}()

	b.ResetTimer()
	wg.Wait()
	b.StopTimer()
}

// ============================
//     PINNED (stable core)
// ============================

func BenchmarkPinned(b *testing.B) {
	hammerScheduler() // same load
	time.Sleep(50 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		for i := 0; i < b.N; i++ {
			cpuHotWork()
			runtime.Gosched() // but no migration occurs
		}
	}()

	b.ResetTimer()
	wg.Wait()
	b.StopTimer()
}
