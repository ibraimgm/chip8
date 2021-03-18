package chip8

import "errors"

// Keys representing CHIP-8 keyboard
const (
	Key0 = iota
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
)

// ErrStackOverflow is returned when the memory reserved for the call
// stack is exceeded due to a programming error.
var ErrStackOverflow = errors.New("stack overflow")

// ErrInputHalt is returned when the emulator is stopped due to waiting for
// a key press from the user. This happens when the instruction Fx0A is requested
// to execute, since it halts the emulation until a user input is received.
var ErrInputHalt = errors.New("awaiting for key input")

// Emulator is the main CHIP-8 CPU emulator.
// It holds the memory, register, stack and all ohter
// components of the CHIP-8 spec.
//
// All the CHIP-8 internals are exported, so you can easily check
// or even modify memory, registers, etc.
//
// Be wary that only the 'logic' of CHIP-8 is emulated; the
// IO (ex: graphics and keyboard) must be implemented separately.
type Emulator struct {
	Memory [4096]byte // main memory
	V      [16]byte   // Vx registers
	I      uint16     // register to store memory address
	DT     byte       // delay timer
	ST     byte       // sound timer
	PC     uint16     // program counter
	SP     int8       // stack pointer
	Stack  [16]uint16 // the stack itself
	keys   [16]bool   // pressed keys
}

// PressKey signal to the emulator that a given key is pressed.
// the key will keep being counted as pressed until a call to
// ReleaseKey. Pressing an already pressed key is a noop.
func (e *Emulator) PressKey(key int) {
	if key >= Key0 && key <= KeyF {
		e.keys[key] = true
	}
}

// ReleaseKey signal to the emulator that a given key is released.
// Releasing an unpressed key is a noop.
func (e *Emulator) ReleaseKey(key int) {
	if key >= Key0 && key <= KeyF {
		e.keys[key] = false
	}
}
