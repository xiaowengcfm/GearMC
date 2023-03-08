package packet

import (
	"errors"
	"github.com/google/uuid"
	"io"
)

const (
	BooleanTrue   = 0x01
	BooleanFalse  = 0x00
	MaxVarIntLen  = 5
	MaxVarLongLen = 10
)

func (b Boolean) WriteTo(w io.Writer) (n int64, err error) {
	var v byte
	if b {
		v = BooleanTrue
	} else {
		v = BooleanFalse
	}
	nn, err := w.Write([]byte{v})
	return int64(nn), err
}

func (b *Boolean) ReadFrom(r io.Reader) (n int64, err error) {
	n, v, err := readByte(r)
	if err != nil {
		return n, err
	}
	*b = v != BooleanFalse
	return n, err
}

func (b Byte) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write([]byte{byte(b)})
	return int64(nn), err
}

func (b *Byte) ReadFrom(r io.Reader) (n int64, err error) {
	n, v, err := readByte(r)
	if err != nil {
		return n, err
	}
	*b = Byte(v)
	return n, err
}

func (ub UByte) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write([]byte{byte(ub)})
	return int64(nn), err
}

func (ub *UByte) ReadFrom(r io.Reader) (n int64, err error) {
	n, v, err := readByte(r)
	if err != nil {
		return n, err
	}
	*ub = UByte(v)
	return n, err
}

func (s Short) WriteTo(w io.Writer) (n int64, err error) {
	ss := uint16(s)
	nn, err := w.Write([]byte{byte(ss >> 8), byte(ss)})
	return int64(nn), err
}

func (s *Short) ReadFrom(r io.Reader) (n int64, err error) {
	var v [2]byte
	nn, err := io.ReadFull(r, v[:])
	if err != nil {
		return int64(nn), err
	}
	*s = Short(int16(v[0]<<8) | int16(v[1]))
	return int64(nn), err
}

func (us UShort) WriteTo(w io.Writer) (n int64, err error) {
	uss := uint16(us)
	nn, err := w.Write([]byte{byte(uss >> 8), byte(uss)})
	return int64(nn), err
}

func (us *UShort) ReadFrom(r io.Reader) (n int64, err error) {
	var v [2]byte
	nn, err := io.ReadFull(r, v[:])
	if err != nil {
		return int64(nn), err
	}
	*us = UShort(uint16(v[0]<<8) | uint16(v[1]))
	return int64(nn), err
}

func (i Int) WriteTo(w io.Writer) (n int64, err error) {
	ii := uint32(i)
	var v [4]byte
	for j := 0; j < len(v); j++ {
		v[j] = byte(ii >> (j + 1) * 8)
	}
	nn, err := w.Write(v[:])
	return int64(nn), err
}

func (i *Int) ReadFrom(r io.Reader) (n int64, err error) {
	var v [4]byte
	var ii int32
	nn, err := io.ReadFull(r, v[:])
	if err != nil {
		return int64(nn), err
	}
	for j := 0; j < len(v); j++ {
		ii |= int32(v[j] << (j + 1) * 8)
	}
	*i = Int(ii)
	return int64(nn), err
}

func (l Long) WriteTo(w io.Writer) (n int64, err error) {
	var v [8]byte
	ll := uint64(l)
	for j := 0; j < len(v); j++ {
		v[j] = byte(ll >> (j + 1) * 8)
	}
	nn, err := w.Write(v[:])
	return int64(nn), err
}

func (l *Long) ReadFrom(r io.Reader) (n int64, err error) {
	var v [8]byte
	var ll int64
	nn, err := io.ReadFull(r, v[:])
	if err != nil {
		return int64(nn), err
	}
	for j := 0; j < len(v); j++ {
		ll |= int64(v[j] << (j + 1) * 8)
	}
	*l = Long(ll)
	return int64(nn), err
}

func (v VarInt) WriteTo(w io.Writer) (n int64, err error) {
	vi := make([]byte, 0, MaxVarIntLen)
	num := uint32(v)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		vi = append(vi, byte(b))
		if num == 0 {
			break
		}
	}
	nn, err := w.Write(vi)
	return int64(nn), err
}

func (v *VarInt) ReadFrom(r io.Reader) (n int64, err error) {
	var V uint32
	var num int64
	for sec := byte(0x80); sec&0x80 != 0; num++ {
		if num > MaxVarIntLen {
			return n, errors.New("VarInt is too big")
		}
		nn, sec, err := readByte(r)
		n += nn
		if err != nil {
			return n, err
		}
		V |= uint32(sec&0x7F) << uint32(7*num)
	}
	*v = VarInt(V)
	return n, err
}

func (v VarLong) WriteTo(w io.Writer) (n int64, err error) {
	vi := make([]byte, 0, MaxVarLongLen)
	num := uint64(v)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		vi = append(vi, byte(b))
		if num == 0 {
			break
		}
	}
	nn, err := w.Write(vi)
	return int64(nn), err
}

func (v *VarLong) ReadFrom(r io.Reader) (n int64, err error) {
	var V uint64
	var num int64
	for sec := byte(0x80); sec&0x80 != 0; num++ {
		if num > MaxVarLongLen {
			return n, errors.New("VarLong is too big")
		}
		nn, sec, err := readByte(r)
		n += nn
		if err != nil {
			return n, err
		}
		V |= uint64(sec&0x7F) << uint64(7*num)
	}
	*v = VarLong(V)
	return n, err
}

func (s String) WriteTo(w io.Writer) (n int64, err error) {
	byteStr := []byte(s)
	n1, err := VarInt(len(byteStr)).WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(byteStr)
	return n1 + int64(n2), err
}

func (s *String) ReadFrom(r io.Reader) (n int64, err error) {
	//TODO implement me
	panic("implement me")
}

type Field interface {
	FieldEncoder
	FieldDecoder
}

type (
	FieldEncoder io.WriterTo
	FieldDecoder io.ReaderFrom
)

type (
	Boolean    bool
	Byte       int8
	UByte      uint8
	Short      int16
	UShort     uint16
	Int        int32
	Long       int64
	VarInt     int32
	VarLong    int64
	Float      float32
	Double     float64
	String     string
	Chat       = String
	Identifier = String
	Angle      byte
	UUID       uuid.UUID
	ByteArray  []byte
	Position   struct {
		X, Y, Z int
	}
)

func readByte(r io.Reader) (int64, byte, error) {
	rb, ok := r.(io.ByteReader)
	if ok {
		v, err := rb.ReadByte()
		return 1, v, err
	}
	var v [1]byte
	n, err := r.Read(v[:])
	return int64(n), v[0], err
}
