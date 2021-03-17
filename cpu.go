package chip8

import (
	"fmt"
)

// masks to extract the most/least nibbles
const (
	lsnMask = 0b00001111
	msnMask = 0b11110000
)

// NoOpError is the error returned when an unknown instruction
// is found during the execution. Receiving this error most likely
// means a programming error occurred or the CHIP-8 program is
// using an unsupported instruction set.
//
// Unlike other errors, receiving a NoOpError will still move the
// program counter and still counts the CPU cycle as executed. This
// is done this way to make ignoring NOOP easier to handle for users.
type NoOpError struct {
	A byte
	B byte
}

func (e NoOpError) Error() string {
	return fmt.Sprintf("unknown instruction 0x%02X%02X (noop)", e.A, e.B)
}

// handlers for each 'category' of instruction
var handlers = [16]func(*Emulator, byte, byte) error{
	handleOp0,
	handleOp1,
	handleOp2,
	handleOp34,
	handleOp34,
	handleOp59,
	handleOp6,
	handleOp7,
	handleOp8,
	handleOp59,
	handleOpA,
	handleOpB,
	nil,
	handleOpD,
	nil,
	nil,
}

// Execute runs at most 'cycles' CPU cycles, returning the number of cycles
// executed and possibly an error value.
//
// Not all errors are fatal; particularly, the instruction Fx0A will keep
// returning ErrInputHalt until a key is pressed.
func (c *Emulator) Execute(cycles int) (int, error) {
	executed := 0

	for executed < cycles {
		a := c.Memory[int(c.PC)]
		b := c.Memory[int(c.PC+1)]
		c.PC += 2

		if err := handlers[a&msnMask>>4](c, a, b); err != nil {
			return executed, err
		}

		executed++
	}

	return executed, nil
}

func handleOp0(c *Emulator, a byte, b byte) error {
	switch {
	// clear screen
	case a == 0x00 && b == 0xE0:
		for i := 0; i < 256; i++ {
			c.Memory[AddrVideo+i] = 0
		}

	// RET (ignores if there is nothing on the stack)
	case a == 0x00 && b == 0xEE && c.SP > 0:
		c.SP--
		c.PC = c.Stack[int(c.SP)]
	}

	return nil
}

func handleOp1(c *Emulator, a byte, b byte) error {
	// no need to error check; max value is 4095 which is still in range
	c.PC = (uint16(a&lsnMask) << 8) + uint16(b)
	return nil
}

func handleOp2(c *Emulator, a byte, b byte) error {
	if int(c.SP) >= len(c.Stack) {
		return ErrStackOverflow
	}

	c.Stack[c.SP] = c.PC
	c.SP++

	c.PC = (uint16(a&lsnMask) << 8) + uint16(b)
	return nil
}

func handleOp34(c *Emulator, a byte, b byte) error {
	x := int(a & lsnMask)

	if (c.V[x] == b) == (a&msnMask>>4 == 3) {
		c.PC += 2
	}

	return nil
}

func handleOp59(c *Emulator, a byte, b byte) error {
	if (b & lsnMask) != 0 {
		return NoOpError{A: a, B: b}
	}

	x := int(a & lsnMask)
	y := ((b & msnMask) >> 4)

	if (c.V[x] == c.V[y]) == (a&msnMask>>4 == 5) {
		c.PC += 2
	}

	return nil
}

func handleOp6(c *Emulator, a byte, b byte) error {
	c.V[a&lsnMask] = b
	return nil
}

func handleOp7(c *Emulator, a byte, b byte) error {
	c.V[a&lsnMask] += b
	return nil
}

func handleOp8(c *Emulator, a byte, b byte) error {
	x := a & lsnMask
	y := b & msnMask >> 4
	n := b & lsnMask

	switch n {
	case 0:
		c.V[x] = c.V[y]
	case 1:
		c.V[x] |= c.V[y]
	case 2:
		c.V[x] &= c.V[y]
	case 3:
		c.V[x] ^= c.V[y]
	case 4:
		v := uint16(c.V[x]) + uint16(c.V[y])
		c.V[x] = byte(v & 0x00FF)
		c.V[0xF] = byte(v >> 8)
	case 5:
		if c.V[x] > c.V[y] {
			c.V[0xF] = 1
		} else {
			c.V[0xF] = 0
		}
		c.V[x] -= c.V[y]
	case 6:
		c.V[0xF] = c.V[y] & 0b00000001
		c.V[x] = c.V[y] >> 1
	case 7:
		if c.V[y] > c.V[x] {
			c.V[0xF] = 1
		} else {
			c.V[0xF] = 0
		}
		c.V[x] = c.V[y] - c.V[x]
	case 0xE:
		c.V[0xF] = c.V[y] >> 7
		c.V[x] = c.V[y] << 1
	default:
		return NoOpError{A: a, B: b}
	}

	return nil
}

func handleOpA(c *Emulator, a byte, b byte) error {
	c.I = (uint16(a&lsnMask) << 8) + uint16(b)
	return nil
}

func handleOpB(c *Emulator, a byte, b byte) error {
	addr := (uint16(a&lsnMask) << 8) + uint16(b) + uint16(c.V[0])
	if int(addr) >= len(c.Memory) {
		return ErrInvalidAddress
	}

	c.PC = addr
	return nil
}

func handleOpD(c *Emulator, a byte, b byte) error {
	const maxX = 63
	const maxY = 31
	const maxVideo = (maxX+1)*(maxY+1)/8 - 1

	x := c.V[a&lsnMask]
	y := c.V[b&msnMask>>4]
	n := b & lsnMask

	// out of bounds or zero, just do nothing
	if x > maxX || y > maxY || n == 0 {
		return nil
	}

	// subslice the sprite data
	c.V[0xF] = 0
	sprites := c.Memory[int(c.I) : int(c.I)+int(n)]

	// draw every byte, accounting for misaligned sprites
	for _, sprite := range sprites {
		// out of bounds: get out
		if y > maxY {
			break
		}

		// compute memory position and byte offset
		pos := AddrVideo + x/8 + y*8
		offset := x % 8 // 0 == aligned byte
		y++

		var before byte

		// extract the whole byte
		if offset == 0 {
			before = c.Memory[pos]
		} else {
			s1, s2 := c.Memory[pos]<<offset, c.Memory[pos+1]>>(8-offset)
			before = s1 | s2
		}

		// apply the sprite and update VF
		after := before ^ sprite
		if before&after != before {
			c.V[0xF] = 1
		}

		// update memory
		if offset == 0 {
			c.Memory[pos] = after
		} else {
			s1, s2 := after>>offset, after<<(8-offset)
			c.Memory[pos] = ((c.Memory[pos] >> (8 - offset)) << (8 - offset)) | s1

			if pos <= AddrVideo+maxVideo {
				c.Memory[pos+1] = ((c.Memory[pos+1] << offset) >> offset) | s2
			}
		}
	}

	return nil
}
