package chip8

import "errors"

// Execute runs at most 'cycles' CPU cycles, returning the number of cycles
// executed and possibly an error value.
//
// Not all errors are fatal; particularly, the instruction Fx0A will keep
// returning ErrInputHalt until a key is pressed.
func (c *Emulator) Execute(cycles int) (int, error) {
	return 0, errors.New("not implemented") //nolint:goerr113
}
