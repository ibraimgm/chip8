package chip8_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ibraimgm/chip8"
)

func runEmulator(rom []byte, videoImage ...byte) (*chip8.Emulator, error) {
	c := &chip8.Emulator{}
	if err := c.LoadROM(bytes.NewReader(rom)); err != nil {
		return nil, err
	}

	// starting video image (used for draw tests)
	copy(c.Memory[chip8.AddrVideo:], videoImage)

	// run as fast as possible, until the end of the
	// memory or a zero instruction is reached
	for {
		if _, err := c.Execute(1); err != nil {
			return nil, err
		}

		if int(c.PC) >= len(c.Memory)-1 {
			break
		}

		if c.Memory[int(c.PC)] == 0 && c.Memory[int(c.PC+1)] == 0 {
			break
		}
	}

	return c, nil
}

func TestOpSys(t *testing.T) {
	t.SkipNow()

	c, err := runEmulator([]byte{0x01, 0x23, 0x00, 0xE1, 0x00, 0xEA})
	if err != nil {
		t.Fatal(err)
	}

	expectectedPC := chip8.AddrStart + 6
	if expectectedPC != int(c.PC) {
		t.Fatalf("expected PC to be at address 0x%X, but was at 0x%X", expectectedPC, c.PC)
	}

	if c.I != 0 {
		t.Fatalf("Register I should be zero, but was 0x%X", c.I)
	}

	if c.DT != 0 {
		t.Fatalf("Timer DT should be zero, but was 0x%02X", c.DT)
	}

	if c.ST != 0 {
		t.Fatalf("Timer ST should be zero, but was 0x%02X", c.ST)
	}

	if c.SP != 0 {
		t.Fatalf("Stack Pointer should be zero, but was 0x%02X", c.SP)
	}

	for i := 0; i < len(c.V); i++ {
		if c.V[i] != 0 {
			t.Fatalf("register V%X should be zero, but was 0x%02X", i, c.V[i])
		}
	}
}

func TestOpCls(t *testing.T) {
	t.SkipNow()

	demoImage := []byte{
		0b00000000, 0b00000000, 0b00000001, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, // _______________________#________________________________________
		0b00100000, 0b00000100, 0b00000011, 0b10000000, 0b00011111, 0b11000000, 0b00000000, 0b00000011, // __#__________#________###__________#######____________________##
		0b01010000, 0b00001010, 0b00000111, 0b11000000, 0b00010000, 0b01000000, 0b00000000, 0b00000111, // _#_#________#_#______#####_________#_____#___________________###
		0b00100000, 0b00001010, 0b00001111, 0b11100000, 0b00110000, 0b11000001, 0b00000111, 0b11111101, // __#_________#_#_____#######_______##____##_____#_____#########_#
		0b00000000, 0b10000100, 0b00011111, 0b11110000, 0b00100001, 0b10000010, 0b10001100, 0b00011001, // ________#____#_____#########______#____##_____#_#___##_____##__#
		0b00000001, 0b11000000, 0b00111111, 0b11111000, 0b00100011, 0b00000001, 0b00011000, 0b00000001, // _______###________###########_____#___##_______#___##__________#
		0b00000011, 0b11100000, 0b01111111, 0b11111100, 0b00100011, 0b11111000, 0b00011000, 0b00000001, // ______#####______#############____#___#######______##__________#
		0b00000111, 0b11110000, 0b00111111, 0b11111000, 0b00110000, 0b00001000, 0b00001100, 0b00000001, // _____#######______###########_____##________#_______##_________#
		0b00001111, 0b11111000, 0b00011111, 0b11110000, 0b00010001, 0b00011000, 0b00000011, 0b10000001, // ____#########______#########_______#___#___##_________###______#
		0b00011111, 0b11111100, 0b00001111, 0b11100000, 0b00010011, 0b00110000, 0b01000000, 0b01100001, // ___###########______#######________#__##__##_____#_______##____#
		0b00111111, 0b11111110, 0b00000111, 0b11000000, 0b00010010, 0b00100000, 0b10100000, 0b01100001, // __#############______#####_________#__#___#_____#_#______##____#
		0b00011111, 0b11111100, 0b00000011, 0b10000000, 0b00010000, 0b00100000, 0b01000001, 0b11000001, // ___###########________###__________#______#______#_____###_____#
		0b00001111, 0b11111000, 0b00000001, 0b00000000, 0b00010000, 0b00111111, 0b00000011, 0b00000011, // ____#########__________#___________#______######______##______##
		0b00000111, 0b11110000, 0b00000000, 0b00001110, 0b00011000, 0b00000001, 0b10000011, 0b00000011, // _____#######________________###____##__________##_____##______##
		0b00000011, 0b11100000, 0b01111100, 0b00011010, 0b00001000, 0b00000011, 0b00010001, 0b10000011, // ______#####______#####_____##_#_____#_________##___#___##_____##
		0b00000001, 0b11000000, 0b01100100, 0b00010011, 0b00001000, 0b00001110, 0b00101000, 0b11110011, // _______###_______##__#_____#__##____#_______###___#_#___####__##
		0b00000000, 0b10000000, 0b00001100, 0b00110001, 0b00001000, 0b00001000, 0b00010000, 0b00011111, // ________#___________##____##___#____#_______#______#_______#####
		0b00000000, 0b00000000, 0b00011000, 0b00110001, 0b10011000, 0b00001000, 0b00000000, 0b00000111, // ___________________##_____##___##__##_______#________________###
		0b00110000, 0b00000001, 0b00010000, 0b00011000, 0b10010000, 0b00001111, 0b11100001, 0b00000011, // __##___________#___#_______##___#__#________#######____#______##
		0b00110000, 0b00000011, 0b10011111, 0b10001000, 0b10010001, 0b00100000, 0b00100010, 0b10000011, // __##__________###__######___#___#__#___#__#_______#___#_#_____##
		0b00000000, 0b00000111, 0b11000000, 0b10001001, 0b10011000, 0b11000000, 0b01100001, 0b00010011, // _____________#####______#___#__##__##___##_______##____#___#__##
		0b00010000, 0b00001111, 0b11100001, 0b10001001, 0b00001000, 0b11000011, 0b11000000, 0b00101011, // ___#________#######____##___#__#____#___##____####________#_#_##
		0b00101000, 0b00011111, 0b11110001, 0b00011001, 0b00001000, 0b11000010, 0b00000000, 0b00010011, // __#_#______#########___#___##__#____#___##____#____________#__##
		0b00010000, 0b00111111, 0b11111001, 0b00010001, 0b00011000, 0b00000010, 0b00001100, 0b00000111, // ___#______###########__#___#___#___##_________#_____##_______###
		0b00000000, 0b01111111, 0b11111101, 0b00010001, 0b10010000, 0b00000011, 0b11111100, 0b10001101, // _________#############_#___#___##__#__________########__#___##_#
		0b10000000, 0b00111111, 0b11111001, 0b10010000, 0b10010000, 0b00110000, 0b00000101, 0b01011001, // #_________###########__##__#____#__#______##_________#_#_#_##__#
		0b10100110, 0b00011111, 0b11110001, 0b10011001, 0b10010000, 0b00110000, 0b00001100, 0b10010001, // #_#__##____#########___##__##__##__#______##________##__#__#___#
		0b01000110, 0b00001111, 0b11100011, 0b00001001, 0b00010000, 0b00000000, 0b00001000, 0b00110001, // _#___##_____#######___##____#__#___#________________#_____##___#
		0b00000000, 0b00000111, 0b11000010, 0b00001001, 0b00010000, 0b00000000, 0b00011000, 0b01100001, // _____________#####____#_____#__#___#_______________##____##____#
		0b00000000, 0b00000011, 0b10000110, 0b00001001, 0b11110000, 0b01011100, 0b00011111, 0b11000001, // ______________###____##_____#__#####_____#_###_____#######_____#
		0b00011110, 0b00000001, 0b00000111, 0b11111000, 0b01100000, 0b00000000, 0b00000000, 0b00000001, // ___####________#_____########____##____________________________#
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000011, // ______________________________________________________________##
	}

	c, err := runEmulator([]byte{0x00, 0xE0}, demoImage...)
	if err != nil {
		t.Fatal(err)
	}

	// video should be cleared
	for i := chip8.AddrVideo; i < 256; i++ {
		if c.Memory[i] != 0 {
			t.Fatalf("video memory at address 0x%X should be zero, but was 0x%02X", i, c.Memory[i])
		}
	}
}

func TestOpLdValue(t *testing.T) {
	t.SkipNow()

	tests := []byte{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32}
	rom := make([]byte, 0, len(tests)*2)

	for i := range tests {
		rom = append(rom, 0x60+byte(i), tests[i])
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, expected := range tests {
		if c.V[i] != expected {
			t.Fatalf("expected register V%X to have value of 0x%02X but found 0x%02X", i, expected, c.V[i])
		}
	}
}

func TestOpLdVx(t *testing.T) {
	t.SkipNow()

	const registers = 16

	for y := byte(0); y < registers; y++ {
		y := y

		t.Run(fmt.Sprintf("LDV%X", y), func(t *testing.T) {
			value := y*2 + 1

			rom := make([]byte, 0, registers*2)
			rom = append(rom, 0x60+y, value)

			for x := byte(0); x < registers; x++ {
				if x == y {
					continue
				}

				rom = append(rom, 0x80+x, y<<4)
			}

			c, err := runEmulator(rom)
			if err != nil {
				t.Fatal(err)
			}

			for i := range c.V {
				if c.V[i] != value {
					t.Fatalf("expected register V%X to be 0x%02X, but it was 0x%02X", i, value, c.V[i])
				}
			}
		})
	}
}

func TestOpJP(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name  string
		rom   []byte
		vx    []byte
		isErr bool
	}{
		{name: "JP", rom: []byte{0x60, 0x01, 0x12, 0x06, 0x61, 0x01, 0x62, 0x01}, vx: []byte{1, 0, 1}},
		{name: "JP Error", rom: []byte{0x1F, 0xFF}, isErr: true},
		{name: "JP+", rom: []byte{0x60, 0x04, 0xB2, 0x02, 0x61, 0x01, 0x61, 0x01}, vx: []byte{1, 0, 1}},
		{name: "JP+ Error", rom: []byte{0x60, 0xAA, 0xBF, 0xF0, 0x61, 0x01, 0x61, 0x01}, isErr: true},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			c, err := runEmulator(test.rom)

			if test.isErr {
				if err == nil {
					t.Fatal("expected error, but received nil")
				}

				if !errors.Is(err, chip8.ErrInvalidAddress) {
					t.Fatalf("expected invalid address error, but got: %v", err)
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			for i, value := range test.vx {
				if c.V[i] != value {
					t.Fatalf("expected register V%X to be 0x%02X but was 0x%02X", i, value, c.V[i])
				}
			}
		})
	}
}

func TestOpCallRet(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x12, 0x08, // JP 0x208 (skip the next 3 lines)
		0x61, 0x01, // V1 = 1
		0x62, 0x02, // V2 = 2
		0x00, 0xEE, // RET
		0x60, 0x01, // V0 = 1
		0x63, 0x03, // V3 = 3
		0x22, 0x02, // CALL 0x202
		0x64, 0x04, // V4 = 4
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range []byte{1, 1, 2, 3, 4} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to have value 0x%02X but found 0x%02X", i, v, c.V[i])
		}
	}
}

func TestStackOverflow(t *testing.T) {
	t.SkipNow()

	_, err := runEmulator([]byte{0x22, 0x00}) // call self (inf. loop)
	if !errors.Is(err, chip8.ErrStackOverflow) {
		t.Fatalf("expected stack overflow error, but got %v", err)
	}
}

func TestOpSkips(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x60, 0x01, // V0 = 1
		0x30, 0x01, // SE V0, 1
		0x61, 0x01, // V1 = 1 (skipped)
		0x30, 0x00, // SE V0,0
		0x65, 0x05, // V5 = 5
		0x40, 0x00, // SNE V0,0
		0x62, 0x01, // V2 = 1 (skipped)
		0x40, 0x00, // SNE V0,0
		0x66, 0x05, // V6 = 5
		0x50, 0x10, // SE V0, V1
		0x67, 0x05, // V7 = 7
		0x55, 0x60, // SE V5, V6
		0x63, 0x01, // V3 = 1
		0x95, 0x60, // SNE V5,V6
		0x60, 0x03, // V0 = 3 (skipped)
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range []byte{1, 0, 0, 0, 0, 5, 5, 7} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to have value 0x%02X but found 0x%02X", i, v, c.V[i])
		}
	}
}

func TestOpMath(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x60, 0x03, // V0 = 3
		0x61, 0x04, // V1 = 4
		0x6A, 0x01, // VA = 1
		0x62, 0x00, // V2 = 0
		0x63, 0x00, // V3 = 0
		0x72, 0x05, // V2 = V2 + 5
		0x83, 0x04, // V3 = V3 + V0
		0x81, 0xA5, // V1 = V1 - VA
		0x41, 0x00, // SNE V1,0
		0x12, 0x0A, // JP 0x20A
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range []byte{3, 0, 20, 12, 0, 0, 0, 0, 0, 0, 1} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to have value 0x%02X but found 0x%02X", i, v, c.V[i])
		}
	}
}

func TestOpMathCarry(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x60, 0x05, // V0 = 5
		0x61, 0x05, // V1 = 6
		0x62, 0x0A, // V2 = 10
		0x80, 0x25, // V0 = V0 - V2 (SUB)
		0x4F, 0x00, // SNE VF,0
		0x63, 0x01, // V3 = 1 (skipped)
		0x81, 0x27, // V1 = V2 - V1 (SUBN)
		0x4F, 0x00, // SNE VF,1
		0x63, 0x01, // V4 = 1 (skipped)
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range []byte{250, 4, 10, 0, 0} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to have value 0x%02X but found 0x%02X", i, v, c.V[i])
		}
	}
}

func TestOpLogic(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x60, 0x6E, // V0 = 110
		0x61, 0x6E, // V1 = 110
		0x62, 0x38, // V2 = 56
		0x63, 0xAA, // V3 = 170
		0x80, 0x22, // V0 = V0 &  V2
		0x81, 0x21, // V1 = V1 |  V2
		0x82, 0x23, // V2 = V2 ^  V2
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range []byte{40, 126, 146, 170} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to have value 0x%02X but found 0x%02X", i, v, c.V[i])
		}
	}
}

func TestOpShift(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name string
		rom  []byte
		v0   byte
		vf   byte
	}{
		{name: "SHR0", rom: []byte{0x61, 0xFE, 0x80, 0x16}, v0: 0x7F, vf: 0},
		{name: "SHR1", rom: []byte{0x61, 0xFF, 0x80, 0x16}, v0: 0x7F, vf: 1},
		{name: "SHL0", rom: []byte{0x61, 0x7F, 0x80, 0x1E}, v0: 0xFE, vf: 0},
		{name: "SHL1", rom: []byte{0x61, 0xFF, 0x80, 0x1E}, v0: 0xFE, vf: 1},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			c, err := runEmulator(test.rom)
			if err != nil {
				t.Fatal(err)
			}

			if c.V[0] != test.v0 {
				t.Fatalf("expected register V0 to be 0x%X but it was 0x%02X", test.v0, c.V[0])
			}

			if c.V[0xF] != test.vf {
				t.Fatalf("expected register VF to be 0x%X but it was 0x%02X", test.vf, c.V[0xF])
			}
		})
	}
}

func TestOpDraw(t *testing.T) {
	t.SkipNow()

	// a square in the first 3 screen rows, on top left
	image := []byte{
		0b11111111, 0b11111111, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b11111111, 0b11111111, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b11111111, 0b11111111, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
	}

	tests := []struct {
		name   string
		x      byte
		y      byte
		offset byte
		size   byte
		video  map[uint16]byte
		vf     byte
	}{
		{name: "TopLeft", x: 0, y: 0, size: 2, video: map[uint16]byte{0: 0, 1: 0xFF, 8: 0, 9: 0xFF, 16: 0xFF, 17: 0xFF}, vf: 1},
		{name: "TopRight", x: 56, y: 0, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 7: 0xFF, 8: 0xFF, 9: 0xFF, 15: 0xFF, 16: 0xFF, 17: 0xFF}, vf: 0},
		{name: "BottomCenter", x: 24, y: 30, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xFF, 9: 0xFF, 16: 0xFF, 17: 0xFF, 243: 0xFF, 251: 0xFF}, vf: 0},
		{name: "Misaligned1", x: 44, y: 9, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xFF, 9: 0xFF, 16: 0xFF, 17: 0xFF, 77: 0x1F, 78: 0xE0, 85: 0x1F, 86: 0xE0}, vf: 0},
		{name: "Misaligned2", x: 4, y: 1, offset: 2, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xF3, 9: 0xCF, 16: 0xF3, 17: 0xCF}, vf: 1},
		{name: "OutOfBounds1", x: 60, y: 8, offset: 2, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xFF, 9: 0xFF, 16: 0xFF, 17: 0xFF, 71: 0x0C, 79: 0x0C}, vf: 0},
		{name: "OutOfBounds2", x: 250, y: 0, size: 2, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xFF, 9: 0xFF, 16: 0xFF, 17: 0xFF}, vf: 0},
		{name: "OutOfBounds3", x: 4, y: 29, size: 4, video: map[uint16]byte{0: 0xFF, 1: 0xFF, 8: 0xFF, 9: 0xFF, 16: 0xFF, 17: 0xFF, 232: 0x0F, 233: 0xF0, 240: 0x0F, 241: 0xF0, 248: 0x0C, 249: 0x30}, vf: 0},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			rom := []byte{
				0x12, 0x06, // Jump to address 206 (i. e. skip sprite data)
				0xFF,                     // ########
				0xFF,                     // ########
				0xC3,                     // ##____##
				0xC3,                     // ##____##
				0xA2, 0x02 + test.offset, // LD I, addr (considering the offset on test)
				0x60, test.x, // V0 = X
				0x61, test.y, // V1 = Y
				0xD0, 0x10 + test.size, // DRW V0, V1, size bytes
			}

			c, err := runEmulator(rom, image...)
			if err != nil {
				t.Fatal(err)
			}

			for position, expected := range test.video {
				addr := chip8.AddrVideo + position
				actual := c.Memory[addr]

				if actual != expected {
					t.Fatalf("video memory at 0x%04X should be 0x%02X, but was 0x%02X", addr, expected, actual)
				}
			}
		})
	}
}

func TestOpRand(t *testing.T) {
	t.SkipNow()

	const registers = 16

	tests := []struct {
		name       string
		mask       byte
		expectZero bool
	}{
		{name: "RandTRUE", mask: 0xFF},
		{name: "RandFALSE", mask: 0x00, expectZero: true},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			rom := make([]byte, 0, registers*2)

			for i := byte(0); i < registers; i++ {
				rom = append(rom, 0xC0+i, test.mask)
			}

			c, err := runEmulator(rom)
			if err != nil {
				t.Fatal(err)
			}

			zeros := 0
			for i := range c.V {
				if c.V[i] == 0 {
					zeros++
				}
			}

			if test.expectZero && zeros != registers {
				t.Fatalf("expected all registers to be zero, but found %v", c.V)
			}

			if !test.expectZero && zeros == registers {
				t.Fatalf("at least one register should have a nonzero value")
			}
		})
	}
}

func TestOpInputSkip(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x65, byte(chip8.KeyA), // V5 = 'A'
		0x66, byte(chip8.KeyB), // V6 = 'B'
		0xE5, 0x9E, // SKP V5 (skips if 'A' is pressed)
		0x60, 0x01, // V0 = 1 (skipped)
		0xE5, 0xA1, // SKP V5 (skips if 'A' is not pressed)
		0x61, 0x01, // V1 = 1
		0xE5, 0x9E, // SKP V6 (skips if 'B' is pressed)
		0x62, 0x01, // V2 = 1
		0xE5, 0xA1, // SKP V6 (skips if 'B' is not pressed)
		0x61, 0x01, // V3 = 1 (skipped)
	}

	var c chip8.Emulator
	if err := c.LoadROM(bytes.NewReader(rom)); err != nil {
		t.Fatal(err)
	}

	c.PressKey(chip8.KeyA)
	cycles := len(rom)/2 - 2 // two skips

	if n, err := c.Execute(cycles); err != nil {
		t.Fatal(err)
	} else if n != cycles {
		t.Fatalf("expected to run %d cycles, but ran %d", cycles, n)
	}

	for i, v := range []byte{0, 1, 1, 0} {
		if c.V[i] != v {
			t.Fatalf("expected register V%X to be 0x%02X, but was 0x%02X", i, v, c.V[i])
		}
	}
}

func TestOpInputHalt(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0xF0, 0x0A, // LD V0, K (waiting)
	}

	var c chip8.Emulator
	if err := c.LoadROM(bytes.NewReader(rom)); err != nil {
		t.Fatal(err)
	}

	// should be halted
	n, err := c.Execute(1)
	if !errors.Is(err, chip8.ErrInputHalt) {
		t.Fatalf("expected input halt error, but received %v", err)
	}

	if n != 0 {
		t.Fatal("should not have executed any instructions")
	}

	// press any key and try again
	c.PressKey(chip8.KeyA)
	if n, err := c.Execute(1); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal("should have executed one instruction")
	}

	if c.V[0] != chip8.KeyA {
		t.Fatalf("expected key 0x%02X on V0, but received %02X", chip8.KeyA, c.V[0])
	}
}

func TestOpLdTimers(t *testing.T) {
	t.SkipNow()

	const expected = 0xAA

	rom := []byte{
		0x60, expected, // V0 = {expected}
		0xF0, 0x15, // DT = V0
		0xF1, 0x07, // V1 = DT
		0xF1, 0x18, // ST = V1
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	if c.DT != expected {
		t.Fatalf("expected DT to be 0x%02X, but was 0x%02X", expected, c.DT)
	}

	if c.V[1] != expected {
		t.Fatalf("expected V1 to be 0x%02X, but was 0x%02X", expected, c.V[1])
	}

	if c.ST != expected {
		t.Fatalf("expected ST to be 0x%02X, but was 0x%02X", expected, c.ST)
	}
}

func TestOpLdSprites(t *testing.T) {
	t.SkipNow()

	const numSprites = 16
	const spriteSize = 5

	for i := 0; i < numSprites; i++ {
		i := i

		t.Run(fmt.Sprintf("LDSprite-%X", i), func(t *testing.T) {
			c, err := runEmulator([]byte{0x60, byte(i), 0xF0, 0x29})
			if err != nil {
				t.Fatal(err)
			}

			expected := chip8.AddrSprite + (i * spriteSize)
			if c.I != uint16(expected) {
				t.Fatalf("expected I to be 0x%04X, but was 0x%04X", expected, c.I)
			}
		})
	}
}

func TestOpBCD(t *testing.T) {
	t.SkipNow()

	const baseAddr = 0x202 // location of the first BCD digit

	tests := []struct {
		name  string
		value byte
		bcd   [3]byte
	}{
		{name: "BCD-1s", value: 8, bcd: [3]byte{0x00, 0x00, 0x08}},
		{name: "BCD-10s", value: 54, bcd: [3]byte{0x00, 0x05, 0x04}},
		{name: "BCD-100s", value: 153, bcd: [3]byte{0x01, 0x05, 0x03}},
		{name: "BCD-00", value: 0, bcd: [3]byte{0x00, 0x00, 0x00}},
		{name: "BCD-FF", value: 255, bcd: [3]byte{0x02, 0x05, 0x05}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			rom := []byte{
				0x12, 0x05, // JP 0x205
				0xFF, 0xFF, 0xFF, // BCD digits (not the misaligned execution)
				0x60, test.value, // V0 = {value}
				0xA2, 0x02, // LD I, 0x202
				0xF0, 0x33, // LD Vx (BCD)
			}

			c, err := runEmulator(rom)
			if err != nil {
				t.Fatal(err)
			}

			for i, digit := range test.bcd {
				addr := baseAddr + i

				if c.Memory[addr] != digit {
					t.Fatalf("expected memory position at 0x%04X to be 0x%02X, but was 0x%02X", addr, digit, c.Memory[addr])
				}
			}
		})
	}
}

func TestOpLdVxI(t *testing.T) {
	t.SkipNow()

	rom := []byte{
		0x60, 0x00, // V0 = 0
		0x61, 0x00, // V1 = 1
		0x62, 0x00, // V2 = 2
		0x63, 0x00, // V3 = 3
		0x64, 0x00, // V4 = 4
		0x65, 0x00, // V5 = 5
		0x66, 0x00, // V6 = 6
		0x67, 0x00, // V7 = 7
		0x68, 0x00, // V8 = 8
		0x69, 0x00, // V9 = 9
		0x12, 0x20, // JP 0x220
		0xFF, 0xFF, // will be read/written
		0xFF, 0xFF, // will be read/written
		0xFF, 0xFF, // will be read/written
		0xFF, 0xFF, // will be read/written
		0xFF, 0xFF, // will be read/written
		0xA2, 0x16, // LD I, 0x216 (534)
		0xF9, 0x55, // LD [I], Vx (store V0-V9)
		0x85, 0x00, // V5 = V0
		0x86, 0x01, // V6 = V1
		0x87, 0x02, // V7 = V2
		0x88, 0x03, // V8 = V3
		0x89, 0x04, // V9 = V4
		0x6A, 0x05, // VA = 5
		0xFA, 0x1E, // I = I + VA = 0x21B (539)
		0xF4, 0x65, // LD Vx [I] (read V0-V4)
	}

	c, err := runEmulator(rom)
	if err != nil {
		t.Fatal(err)
	}

	for i, value := range []byte{5, 6, 7, 8, 9, 0, 1, 2, 3, 4} {
		if c.V[i] != value {
			t.Fatalf("expected register V%X to be 0x%02X, but was 0x%02X", i, value, c.V[i])
		}
	}
}

func TestWriteViolation(t *testing.T) {
	t.SkipNow()

	tests := []struct {
		name string
		addr uint16
		ok   bool
	}{
		{name: "AddrVideo-1", addr: chip8.AddrVideo},
		{name: "AddrVideo-2", addr: chip8.AddrVideo + 0xA0},
		{name: "AddrSprite-1", addr: chip8.AddrSprite},
		{name: "AddrSprite-2", addr: chip8.AddrSprite + 0x28},
		{name: "AddrStart-1", addr: chip8.AddrStart, ok: true},
		{name: "AddrStart-2", addr: chip8.AddrStart - 1},
		{name: "End-1", addr: 4095, ok: true},
		{name: "End-2", addr: 4096},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			var c chip8.Emulator
			if err := c.LoadROM(bytes.NewReader([]byte{0xF0, 0x55})); err != nil {
				t.Fatal(err)
			}

			c.I = test.addr
			_, err := c.Execute(1)

			if test.ok && err != nil {
				t.Fatal(err)
			}

			if !test.ok && !errors.Is(err, chip8.ErrMemWrite) {
				t.Fatalf("expected memory write error, but got %v", err)
			}
		})
	}
}
