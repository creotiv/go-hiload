package main

import (
	"fmt"
	"unsafe"
)

/*
   ────────────────────────────────────────────────────────────────
   HOW TO READ THE OUTPUT
   ────────────────────────────────────────────────────────────────

   unsafe.Sizeof(T)  → total size of struct in bytes
   unsafe.Alignof(T) → required alignment (largest field alignment)
   unsafe.Offsetof   → where each field begins inside the struct

   Go inserts *padding bytes* between fields to ensure each field
   begins at an address divisible by its alignment.

   The examples below show how *field order* affects padding,
   and therefore struct total size and cache efficiency.
*/

//
// ──────────────────────────────────────────────────────────────
// Example 1: S1 vs S2 — same fields, different order
// ──────────────────────────────────────────────────────────────
//

type S1 struct {
	A bool  // 1 byte
	B int32 // 4 bytes, requires 4-byte alignment
	C bool  // 1 byte
}

type S2 struct {
	A int32
	B bool
	C bool
}

func main() {
	//
	// ──────────────────────────────────────────────────────────
	// S1 — inefficient ordering
	// ──────────────────────────────────────────────────────────
	//
	fmt.Println("=== S1 (A bool, B int32, C bool) ===")
	fmt.Printf("Sizeof(S1): %d bytes (because of padding)\n", unsafe.Sizeof(S1{}))

	/*
	   Memory layout (with padding shown as dots):

	   Offset 00: A (1 byte)
	   Offset 01–03: ... padding (to align int32 at offset 4)
	   Offset 04–07: B (4 bytes)
	   Offset 08:    C (1 byte)
	   Offset 09–11: ... padding (struct must end on 4-byte boundary)

	   FINAL SIZE = 12 bytes
	*/

	printS1()

	//
	// ──────────────────────────────────────────────────────────
	// S2 — optimal ordering for these fields
	// ──────────────────────────────────────────────────────────
	//
	fmt.Println("=== S2 (A bool, C bool, B int32) ===")
	fmt.Printf("Sizeof(S2): %d bytes (no wasted padding)\n", unsafe.Sizeof(S2{}))

	/*
		   Memory layout:

			Offset 00: A (4 bytes)
			Offset 04: C (1 byte)
			Offset 05: B (1 byte)
			Offset 06-07: ... padding (struct must end on 4-byte boundary)

		   FINAL SIZE = 8 bytes  (33% smaller than S1)
	*/

	printS2()
}

//
// --- Helper print functions f

func printS1() {
	fmt.Printf("A offset=%d size=%d align=%d\n", unsafe.Offsetof(S1{}.A), unsafe.Sizeof(true), unsafe.Alignof(true))
	fmt.Printf("B offset=%d size=%d align=%d\n", unsafe.Offsetof(S1{}.B), unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	fmt.Printf("C offset=%d size=%d align=%d\n", unsafe.Offsetof(S1{}.C), unsafe.Sizeof(true), unsafe.Alignof(true))
	fmt.Println()
}

func printS2() {
	fmt.Printf("A offset=%d size=%d align=%d\n", unsafe.Offsetof(S2{}.A), unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	fmt.Printf("B offset=%d size=%d align=%d\n", unsafe.Offsetof(S2{}.B), unsafe.Sizeof(true), unsafe.Alignof(true))
	fmt.Printf("C offset=%d size=%d align=%d\n", unsafe.Offsetof(S2{}.C), unsafe.Sizeof(true), unsafe.Alignof(true))
	fmt.Println()
}
