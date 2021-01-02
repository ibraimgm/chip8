package chip8_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ibraimgm/chip8"
)

func corruptedEmulator() *chip8.Emulator {
	c := chip8.Emulator{}

	// Intentionally corrupt the memory and registers. The values don't matter,
	// the point is testing if 'reset' goes to a valid state
	for i := 0; i < len(c.Memory); i++ {
		c.Memory[i] = byte(rand.Intn(256))
	}

	for i := 0; i < len(c.V); i++ {
		c.V[i] = byte(rand.Intn(256))
	}

	for i := 0; i < len(c.Stack); i++ {
		c.Stack[i] = uint16(rand.Intn(256))
	}

	c.I = 1
	c.DT = 1
	c.ST = 1
	c.PC = 1
	c.SP = 1

	return &c
}

func TestReset(t *testing.T) {
	c := corruptedEmulator()
	c.Reset()

	// the program area should be all zeros
	for i := chip8.AddrStart; i < len(c.Memory); i++ {
		if c.Memory[i] != 0 {
			t.Fatalf("memory address 0x%X should be zero, but was 0x%X", i, c.Memory[i])
		}
	}

	// video should be cleared too
	for i := chip8.AddrVideo; i < 256; i++ {
		if c.Memory[i] != 0 {
			t.Fatalf("video memory at address 0x%X should be zero, but was 0x%X", i, c.Memory[i])
		}
	}

	// we should be able to find some sprite data
	sprites := []struct {
		startAddr int
		expected  []byte
	}{
		{startAddr: chip8.AddrSprite, expected: []byte{0xF0, 0x90, 0x90, 0x90, 0xF0}},      // 0
		{startAddr: chip8.AddrSprite + 25, expected: []byte{0xF0, 0x80, 0xF0, 0x10, 0xF0}}, // 5
		{startAddr: chip8.AddrSprite + 55, expected: []byte{0xE0, 0x90, 0xE0, 0x90, 0xE0}}, // B
		{startAddr: chip8.AddrSprite + 70, expected: []byte{0xF0, 0x80, 0xF0, 0x80, 0xF0}}, // E
		{startAddr: chip8.AddrSprite + 75, expected: []byte{0xF0, 0x80, 0xF0, 0x80, 0x80}}, // F
	}

	for i, sprite := range sprites {
		for offset := 0; offset < len(sprite.expected); offset++ {
			actual := c.Memory[sprite.startAddr+offset]
			if actual != sprite.expected[offset] {
				t.Fatalf("expected byte '0x%X' at position %d of sprite %d, but found '0x%X' (address: 0x%X)", sprite.expected[offset], offset, i, actual, sprite.startAddr)
			}
		}
	}

	// all registers should be zero
	for i := 0; i < len(c.V); i++ {
		if c.V[i] != 0 {
			t.Fatalf("register V%X should be zero, but was 0x%X", i, c.V[i])
		}
	}

	// the entire stack should be zero
	for i := 0; i < len(c.Stack); i++ {
		if c.Stack[i] != 0 {
			t.Fatalf("stack at position %d should be zero, but was 0x%X", i, c.Stack[i])
		}
	}

	if c.I != 0 {
		t.Fatalf("Register I should be zero, but was 0x%X", c.I)
	}

	if c.DT != 0 {
		t.Fatalf("Timer DT should be zero, but was 0x%X", c.DT)
	}

	if c.ST != 0 {
		t.Fatalf("Timer ST should be zero, but was 0x%X", c.ST)
	}

	if c.SP != 0 {
		t.Fatalf("Stack Pointer should be zero, but was 0x%X", c.SP)
	}

	if c.PC != chip8.AddrStart {
		t.Fatalf("PC shoud be pointing to address 0x%X, but was 0x%X", chip8.AddrStart, c.PC)
	}
}

func TestLoadROM(t *testing.T) {
	rom := []byte{0x00, 0x1A, 0x1B, 0x1C, 0x1D, 0xF0}
	c := corruptedEmulator()

	if err := c.LoadROM(bytes.NewReader(rom)); err != nil {
		t.Fatal(err)
	}

	if c.PC != chip8.AddrStart {
		t.Fatalf("expected PC to be 0x%X, but was 0x%X", chip8.AddrStart, c.PC)
	}

	for offset := 0; offset < len(rom); offset++ {
		addr := chip8.AddrStart + offset
		if c.Memory[addr] != rom[offset] {
			t.Fatalf("address 0x%X should be 0x%X, but was 0x%X", addr, rom[offset], c.Memory[addr])
		}
	}
}

func TestLoadBigROM(t *testing.T) {
	rom := make([]byte, 4096-chip8.AddrStart+1)
	for i := 0; i < len(rom)-2; i++ {
		rom[i] = 1
	}
	rom[len(rom)-1] = 9
	rom[len(rom)-2] = 8

	var c chip8.Emulator
	err := c.LoadROM(bytes.NewReader(rom))

	if err == nil {
		t.Fatal("should have returned error, but got none")
	}

	if err != chip8.ErrLoadOverflow {
		t.Fatalf("expected ROM overflow error, but got '%v'", chip8.ErrLoadOverflow)
	}

	actual := c.Memory[4095]
	expected := rom[len(rom)-2]
	if actual != expected {
		t.Fatalf("expected last memory byte to be 0x%X but was 0x%X", expected, actual)
	}
}
