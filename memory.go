package chip8

import (
	"errors"
	"io"
)

// Common addresses used by the CHIP-8 emulator.
const (
	AddrVideo  = 0x000
	AddrSprite = 0x100
	AddrStart  = 0x200
)

// ErrLoadOverflow is the error returned when the emulator tries to
// load a image that is bigger than the maximum available memory.
var ErrLoadOverflow = errors.New("error loading ROM: size exceeds CHIP-8 memory limit")

// ErrInvalidAddress is returned when a jump or similar command tries to
// go to an invalid address.
var ErrInvalidAddress = errors.New("invalid memory address")

// ErrMemWrite is returned by the emulator when an instruction tries to
// write to the reserved memory of the emulator.
var ErrMemWrite = errors.New("cannot write into reserved memory address")

// static sprite data
var sprites = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// Reset resets the emulator state. This clears all memory and
// resets all registers to the initial values.
func (c *Emulator) Reset() {
	// clear memory
	for i := 0; i < len(c.Memory); i++ {
		if i >= AddrSprite && i < AddrSprite+len(sprites) {
			c.Memory[i] = sprites[i-AddrSprite]
		} else {
			c.Memory[i] = 0
		}
	}

	// clear Vx registers
	for i := range c.V {
		c.V[i] = 0
	}

	// clear stack
	for i := range c.Stack {
		c.Stack[i] = 0
	}

	c.I = 0
	c.DT = 0
	c.ST = 0
	c.SP = 0
	c.PC = AddrStart
}

// LoadROM loads a given ROM to the emulator memory. Before
// loading, Reset is called to keep the emulator in a 'clean'
// state.
func (c *Emulator) LoadROM(rom io.Reader) error {
	c.Reset()

	buffer := make([]byte, 4096-AddrStart)
	addr := AddrStart

	for {
		count, err := rom.Read(buffer)
		var i int

		for i = 0; i < count && addr < len(c.Memory); i++ {
			c.Memory[addr] = buffer[i]
			addr++
		}

		if i < count {
			return ErrLoadOverflow
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}
