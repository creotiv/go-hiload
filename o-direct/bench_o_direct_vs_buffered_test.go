//go:build linux

package directio

import (
	"os"
	"sync"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	blockSize   = 4 * 1024 // 4KB
	commitEvery = 64       // group commit frequency
)

// --- Section: Alignment helper ---

func aligned(size int) []byte {
	mem := make([]byte, size+4095)
	ptr := uintptr(unsafe.Pointer(&mem[0]))
	offset := int(ptr & 4095)
	start := (4096 - offset) & 4095
	return mem[start : start+size]
}

// --- Section: Buffered I/O ---

var (
	bufOnce sync.Once
	bufFile *os.File
	bufData []byte
)

func setupBuffered(b *testing.B) {
	bufOnce.Do(func() {
		var err error
		bufData = make([]byte, blockSize)

		bufFile, err = os.OpenFile("/tmp/buffered_wal.dat",
			os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			b.Fatalf("open buffered: %v", err)
		}
	})
}

func BenchmarkBufferedWrite(b *testing.B) {
	setupBuffered(b)
	b.SetBytes(blockSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := bufFile.WriteAt(bufData, 0); err != nil {
			b.Fatalf("WriteAt: %v", err)
		}
	}
	b.StopTimer()
}

// --- Section: Buffered + sync (group commit) ---

func BenchmarkBufferedWriteSync(b *testing.B) {
	setupBuffered(b)
	b.SetBytes(blockSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := bufFile.WriteAt(bufData, 0); err != nil {
			b.Fatalf("WriteAt: %v", err)
		}

		if i%commitEvery == 0 {
			if err := bufFile.Sync(); err != nil {
				b.Fatalf("Sync: %v", err)
			}
		}
	}

	b.StopTimer()
}

// --- Section: O_DIRECT I/O ---

var (
	dirOnce sync.Once
	dirFD   int
	dirData []byte
)

func setupDirect(b *testing.B) {
	dirOnce.Do(func() {
		var err error
		dirData = aligned(blockSize)

		dirFD, err = unix.Open(
			"/tmp/direct_wal.dat",
			unix.O_CREAT|unix.O_WRONLY|unix.O_DIRECT,
			0o644,
		)
		if err != nil {
			b.Fatalf("open direct: %v", err)
		}

		// Optional: preallocate to avoid metadata work
		// _ = unix.Fallocate(dirFD, 0, 0, blockSize)
	})
}

func BenchmarkDirectWrite(b *testing.B) {
	setupDirect(b)
	b.SetBytes(blockSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := unix.Pwrite(dirFD, dirData, 0)
		if err != nil {
			b.Fatalf("Pwrite: %v", err)
		}
		if n != blockSize {
			b.Fatalf("short write: %d", n)
		}
	}

	b.StopTimer()
}

// --- Section: Direct + sync (fdatasync group commit) ---

func BenchmarkDirectWriteSync(b *testing.B) {
	setupDirect(b)
	b.SetBytes(blockSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := unix.Pwrite(dirFD, dirData, 0)
		if err != nil {
			b.Fatalf("Pwrite: %v", err)
		}
		if n != blockSize {
			b.Fatalf("short write: %d", n)
		}

		if i%commitEvery == 0 {
			if err := unix.Fdatasync(dirFD); err != nil {
				b.Fatalf("fdatasync: %v", err)
			}
		}
	}

	b.StopTimer()
}

// --- Section: Cleanup ---

func TestMain(m *testing.M) {
	code := m.Run()

	if bufFile != nil {
		_ = bufFile.Close()
		_ = os.Remove("/tmp/buffered_wal.dat")
	}
	if dirFD != 0 {
		_ = unix.Close(dirFD)
		_ = os.Remove("/tmp/direct_wal.dat")
	}

	os.Exit(code)
}
