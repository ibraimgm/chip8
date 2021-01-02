package chip8

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
}
