package rpc

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkWriteSeq(t assert.TestingT, w *Writer) {
	assert.NoError(t, w.WriteBool(false))
	assert.NoError(t, w.WriteBool(true))

	assert.NoError(t, w.WriteByte(0x01))
	assert.NoError(t, w.WriteShort(0x0203))
	assert.NoError(t, w.WriteInt(0x04050607))
	assert.NoError(t, w.WriteLong(0x08090A0B0C0D0E0F))

	assert.NoError(t, w.WriteFloat(math.Float32frombits(0x10111213)))
	assert.NoError(t, w.WriteDouble(math.Float64frombits(0x1415161718191A1B)))

	assert.NoError(t, w.WriteSize(0x01))
	assert.NoError(t, w.WriteSize(0x0102))

	assert.NoError(t, w.WriteByteArr([]byte{1, 2, 3}))
	assert.NoError(t, w.WriteString("abc"))
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	nativeOrder = littleEndian
	checkWriteSeq(t, w)
	assert.Equal(t, littleEndianSeq, buf.Bytes())

	buf.Reset()

	nativeOrder = bigEndian
	checkWriteSeq(t, w)
	assert.Equal(t, bigEndianSeq, buf.Bytes())
}

func TestWriteOverflow(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	b := make([]byte, math.MaxUint32+1)
	assert.ErrorIs(t, w.WriteByteArr(b), errOverflow)
}

func BenchmarkWriteInts(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	b.SetBytes(1 + 2 + 4 + 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.WriteByte(1)
		w.WriteShort(1)
		w.WriteInt(1)
		w.WriteLong(1)
		buf.Reset()
	}
}

func BenchmarkWriteFloats(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	b.SetBytes(4 + 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.WriteFloat(1.1)
		w.WriteDouble(1.1)
		buf.Reset()
	}
}

func BenchmarkWriteStrings(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	s := string(make([]byte, 1000))

	b.SetBytes(1 + 4 + 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.WriteString(s)
		buf.Reset()
	}
}
