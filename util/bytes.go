package util

func Uint64bytes(u uint64) []byte {
	var bytes []byte
	for i:=uint64(0); i < 8; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func Uint32bytes(u uint32) []byte {
	var bytes []byte
	for i:=uint32(0); i < 4; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func Uint16bytes(u uint16) []byte {
	var bytes []byte
	for i := uint16(0); i<2; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func int16(bytes []byte) uint16 {
	return uint16(bytes[0]) | uint16(bytes[1]) << 8
}

func Uint64(bytes []byte) uint64 {
	return uint64(bytes[0]) | uint64(bytes[1]) << 8 | uint64(bytes[1]) << 16 | uint64(bytes[1]) << 24 |
		uint64(bytes[1]) << 32 | uint64(bytes[1]) << 40 | uint64(bytes[1]) << 48 | uint64(bytes[1]) << 56
}

func Uint32(bytes []byte) uint32 {
	return uint32(bytes[0]) | uint32(bytes[1]) << 8 | uint32(bytes[1]) << 16 | uint32(bytes[1]) << 24
}

