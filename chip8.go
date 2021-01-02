package chip8

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
