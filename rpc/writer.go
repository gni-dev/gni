package rpc

import (
	"io"
	"math"
)

type Writer struct {
	scratch [8]byte
	out     io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{out: w}
}

func (w *Writer) WriteByte(v byte) error {
	w.scratch[0] = v
	_, err := w.out.Write(w.scratch[:1])
	return err
}

func (w *Writer) WriteBool(v bool) error {
	if v {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

func (w *Writer) WriteShort(v uint16) error {
	if nativeOrder == littleEndian {
		w.scratch[0] = byte(v)
		w.scratch[1] = byte(v >> 8)
	} else {
		w.scratch[1] = byte(v)
		w.scratch[0] = byte(v >> 8)
	}
	_, err := w.out.Write(w.scratch[:2])
	return err
}

func (w *Writer) WriteInt(v uint32) error {
	if nativeOrder == littleEndian {
		w.scratch[0] = byte(v)
		w.scratch[1] = byte(v >> 8)
		w.scratch[2] = byte(v >> 16)
		w.scratch[3] = byte(v >> 24)
	} else {
		w.scratch[3] = byte(v)
		w.scratch[2] = byte(v >> 8)
		w.scratch[1] = byte(v >> 16)
		w.scratch[0] = byte(v >> 24)
	}
	_, err := w.out.Write(w.scratch[:4])
	return err
}

func (w *Writer) WriteLong(v uint64) error {
	if nativeOrder == littleEndian {
		w.scratch[0] = byte(v)
		w.scratch[1] = byte(v >> 8)
		w.scratch[2] = byte(v >> 16)
		w.scratch[3] = byte(v >> 24)
		w.scratch[4] = byte(v >> 32)
		w.scratch[5] = byte(v >> 40)
		w.scratch[6] = byte(v >> 48)
		w.scratch[7] = byte(v >> 56)
	} else {
		w.scratch[7] = byte(v)
		w.scratch[6] = byte(v >> 8)
		w.scratch[5] = byte(v >> 16)
		w.scratch[4] = byte(v >> 24)
		w.scratch[3] = byte(v >> 32)
		w.scratch[2] = byte(v >> 40)
		w.scratch[1] = byte(v >> 48)
		w.scratch[0] = byte(v >> 56)
	}
	_, err := w.out.Write(w.scratch[:8])
	return err
}

func (w *Writer) WriteFloat(v float32) error {
	return w.WriteInt(math.Float32bits(v))
}

func (w *Writer) WriteDouble(v float64) error {
	return w.WriteLong(math.Float64bits(v))
}

func (w *Writer) WriteSize(v uint32) error {
	if v > 254 {
		if err := w.WriteByte(255); err != nil {
			return err
		}
		return w.WriteInt(v)
	} else {
		return w.WriteByte(byte(v))
	}
}

func (w *Writer) WriteByteArr(v []byte) error {
	l := len(v)
	if l > math.MaxUint32 {
		return errOverflow
	}
	if err := w.WriteSize(uint32(l)); err != nil {
		return err
	}
	_, err := w.out.Write(v)
	return err
}

func (w *Writer) WriteString(v string) error {
	return w.WriteByteArr([]byte(v))
}
