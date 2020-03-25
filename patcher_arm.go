// +build arm

package mpatch

// Gets the jump function rewrite bytes
func getJumpFuncBytes(to uintptr) []byte {
	//ldr r3, =0x[XXXXXXXX]
	//blx r3
	return []byte {
		0x00, 0x30,
		0x9F, 0xE5,
		0x33, 0xFF,
		0x2F, 0xE1,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
	}
}