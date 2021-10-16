package rpc

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var littleEndianSeq = []byte{
	0, 1,

	1,
	3, 2,
	7, 6, 5, 4,
	15, 14, 13, 12, 11, 10, 9, 8,

	19, 18, 17, 16,
	27, 26, 25, 24, 23, 22, 21, 20,

	1,
	255, 2, 1, 0, 0,

	3, 1, 2, 3,
	3, 'a', 'b', 'c',
}

var bigEndianSeq = []byte{
	0, 1,

	1,
	2, 3,
	4, 5, 6, 7,
	8, 9, 10, 11, 12, 13, 14, 15,

	16, 17, 18, 19,
	20, 21, 22, 23, 24, 25, 26, 27,

	1,
	255, 0, 0, 1, 2,

	3, 1, 2, 3,
	3, 'a', 'b', 'c',
}

type mockReader struct {
	remain []byte
	pos    int
}

func (m *mockReader) Read(p []byte) (int, error) {
	n := copy(p, m.remain[m.pos:])
	m.pos += n
	if n == 0 {
		panic("eof")
	}
	return n, nil
}

func checkReadResult(t assert.TestingT, expected interface{}, actual interface{}, err error) {
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func checkReadSeq(t assert.TestingT, r *Reader) {
	var err error
	var v interface{}

	v, err = r.ReadBool()
	checkReadResult(t, false, v, err)
	v, err = r.ReadBool()
	checkReadResult(t, true, v, err)

	v, err = r.ReadByte()
	checkReadResult(t, uint8(0x01), v, err)
	v, err = r.ReadShort()
	checkReadResult(t, uint16(0x0203), v, err)
	v, err = r.ReadInt()
	checkReadResult(t, uint32(0x04050607), v, err)
	v, err = r.ReadLong()
	checkReadResult(t, uint64(0x08090A0B0C0D0E0F), v, err)

	v, err = r.ReadFloat()
	checkReadResult(t, math.Float32frombits(0x10111213), v, err)
	v, err = r.ReadDouble()
	checkReadResult(t, math.Float64frombits(0x1415161718191A1B), v, err)

	v, err = r.ReadSize()
	checkReadResult(t, uint32(0x01), v, err)
	v, err = r.ReadSize()
	checkReadResult(t, uint32(0x0102), v, err)

	v, err = r.ReadByteArr()
	checkReadResult(t, []byte{1, 2, 3}, v, err)
	v, err = r.ReadString()
	checkReadResult(t, "abc", v, err)
}

func TestRead(t *testing.T) {
	nativeOrder = littleEndian
	buf := bytes.NewBuffer(littleEndianSeq)
	r := NewReader(buf)
	checkReadSeq(t, r)
	assert.Equal(t, 0, buf.Len())

	nativeOrder = bigEndian
	buf = bytes.NewBuffer(bigEndianSeq)
	r = NewReader(buf)
	checkReadSeq(t, r)
	assert.Equal(t, 0, buf.Len())
}

func BenchmarkReadInts(b *testing.B) {
	bs := make([]byte, 1+2+4+8)
	buf := &mockReader{remain: bs}
	r := NewReader(buf)

	b.SetBytes(int64(len(bs)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ReadByte()
		r.ReadShort()
		r.ReadInt()
		r.ReadLong()
		buf.pos = 0
	}
}

func BenchmarkReadFloats(b *testing.B) {
	bs := make([]byte, 4+8)
	buf := &mockReader{remain: bs}
	r := NewReader(buf)

	b.SetBytes(int64(len(bs)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ReadFloat()
		r.ReadDouble()
		buf.pos = 0
	}
}

func BenchmarkReadStrings(b *testing.B) {
	bs := make([]byte, 1005)
	bs[0] = 255
	bs[1] = 0xE8
	bs[2] = 0x03
	buf := &mockReader{remain: bs}
	r := NewReader(buf)

	b.SetBytes(int64(len(bs)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ReadString()
		buf.pos = 0
	}
}
