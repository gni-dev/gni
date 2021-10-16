package rpc

import (
	"io"
	"math"
)

type Reader struct {
	scratch [8]byte
	in      io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{in: r}
}

func (r *Reader) ReadByte() (byte, error) {
	b := r.scratch[:1]
	if _, err := r.in.Read(b); err != nil {
		// explicitly return default value. since scratch may contain trash
		return 0, err
	}
	return b[0], nil
}

func (r *Reader) ReadBool() (bool, error) {
	v, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return v != 0, nil
}

func (r *Reader) ReadShort() (uint16, error) {
	b := r.scratch[:2]
	if _, err := r.in.Read(b); err != nil {
		return 0, err
	}
	if nativeOrder == littleEndian {
		return uint16(b[0]) | uint16(b[1])<<8, nil
	} else {
		return uint16(b[1]) | uint16(b[0])<<8, nil
	}
}

func (r *Reader) ReadInt() (uint32, error) {
	b := r.scratch[:4]
	if _, err := r.in.Read(b); err != nil {
		return 0, err
	}
	if nativeOrder == littleEndian {
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
	} else {
		return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24, nil
	}
}

func (r *Reader) ReadLong() (uint64, error) {
	b := r.scratch[:8]
	if _, err := r.in.Read(b); err != nil {
		return 0, err
	}
	if nativeOrder == littleEndian {
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56, nil
	} else {
		return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
			uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56, nil
	}
}

func (r *Reader) ReadFloat() (float32, error) {
	v, err := r.ReadInt()
	return math.Float32frombits(v), err
}

func (r *Reader) ReadDouble() (float64, error) {
	v, err := r.ReadLong()
	return math.Float64frombits(v), err
}

func (r *Reader) ReadSize() (uint32, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if b == 255 {
		return r.ReadInt()
	} else {
		return uint32(b), nil
	}
}

func (r *Reader) ReadByteArr() ([]byte, error) {
	l, err := r.ReadSize()
	if err != nil {
		return nil, err
	}
	b := make([]byte, l)
	if _, err = r.in.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (r *Reader) ReadString() (string, error) {
	v, err := r.ReadByteArr()
	if err != nil {
		return "", err
	}
	return string(v), nil
}
