
# Core Ideas Behind Memory Alignment

Before looking at how Go handles struct layout, it helps to understand a few basic ideas about how memory works.

# Memory Address

Every byte in RAM has its own number (address).
When we say “a value is stored at address X,” we mean the first byte of that value starts there.

# Word Size

A CPU reads data in fixed-size chunks called words.

* On a 64-bit CPU, a word is usually **8 bytes**
* On a 32-bit CPU, a word is usually **4 bytes**

If a value crosses a word boundary, the CPU may need **extra operations**, making it slower.

# Alignment Requirement

Different data types prefer to start at addresses divisible by their size.

On a 64-bit system this typically means:

* `byte (1B)`: can be anywhere
* `int16 / short (2B)`: start at address divisible by 2
* `int32 (4B)`: start at address divisible by 4
* `int64 / float64 (8B)`: start at address divisible by 8

This lets the CPU load the value quickly in one aligned read.

# Padding

If fields in a struct don’t naturally fall on aligned boundaries,
the compiler inserts padding bytes between them.

Padding:
* wastes memory
* doesn’t store real data
* exists only to make the next field start at the right alignment

This is why struct field order changes the total size.

# Struct aligment

## Field Alignment

Each field in a struct must start at an address that matches its own alignment requirement, unless the struct’s alignment is smaller — in that case, the struct’s alignment rules apply.

## Struct Alignment

A struct’s alignment is determined by the field inside it that has the largest alignment requirement.

## Struct Size

The final size of a struct is always rounded up to a multiple of its alignment.
If the natural size doesn’t fit that rule, Go adds padding bytes at the end.

# Example

```go
type S1 struct {
	A bool  // 1 byte
	B int32 // 4 bytes, requires 4-byte alignment
	C bool  // 1 byte
}
/*
Offset 00: A (1 byte)
Offset 01–03: ... padding (to align int32 at offset 4)
Offset 04–07: B (4 bytes)
Offset 08:    C (1 byte)
Offset 09–11: ... padding (struct must end on 4-byte boundary)

FINAL SIZE = 12 bytes
*/

type S2 struct {
	A int32
	B bool
	C bool
}


/*
    Offset 00: A (4 bytes)
    Offset 04: C (1 byte)
    Offset 05: B (1 byte)
    Offset 05-7: B (2 bytes)

    FINAL SIZE = 8 bytes  (33% smaller than S1)
*/
```
