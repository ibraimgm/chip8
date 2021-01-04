package chip8

import "fmt"

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
	nil,
	nil,
	nil,
	nil,
	nil,
	handleOp6,
	nil,
	handleOp8,
	nil,
	nil,
	nil,
	nil,
	nil,
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
	if a == 0x00 && b == 0xE0 {
		for i := 0; i < 256; i++ {
			c.Memory[AddrVideo+i] = 0
		}
	}

	return nil
}

func handleOp6(c *Emulator, a byte, b byte) error {
	c.V[a&lsnMask] = b
	return nil
}

func handleOp8(c *Emulator, a byte, b byte) error {
	x := a & lsnMask
	y := b & msnMask >> 4
	n := b & lsnMask

	if n == 0 {
		c.V[x] = c.V[y]
		return nil
	}

	return NoOpError{A: a, B: b}
}
