package rpc

import "unsafe"

type byteOrder int8

const (
	littleEndian byteOrder = iota
	bigEndian
)

var nativeOrder byteOrder

func init() {
	x := uint32(0x01020304)
	switch *(*byte)(unsafe.Pointer(&x)) {
	case 0x01:
		nativeOrder = bigEndian
	case 0x04:
		nativeOrder = littleEndian
	}
}
