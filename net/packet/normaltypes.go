package packet

import (
	"errors"
	"io"
	"math"

	"github.com/google/uuid"
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
	var v VarInt
	nn, err := v.ReadFrom(r)
	if err != nil {
		return nn, err
	}
	n += nn
	bs := make([]byte, v)
	_, err = io.ReadFull(r, bs)
	if err != nil {
		return n, err
	}
	n += int64(v)
	*s = String(bs)
	return n, err
}

func (f Float) WriteTo(w io.Writer) (n int64, err error) {
	return Int(math.Float32bits(float32(f))).WriteTo(w)
}

func (f *Float) ReadFrom(r io.Reader) (n int64, err error) {
	var v Int
	n, err = v.ReadFrom(r)
	if err != nil {
		return n, err
	}
	*f = Float(math.Float32frombits(uint32(v)))
	return n, err
}

func (d Double) WriteTo(w io.Writer) (n int64, err error) {
	return Long(math.Float64bits(float64(d))).WriteTo(w)
}

func (d *Double) ReadFrom(r io.Reader) (n int64, err error) {
	var v Long
	n, err = v.ReadFrom(r)
	if err != nil {
		return n, err
	}
	*d = Double(math.Float64frombits(uint64(v)))
	return n, err
}

func (a Angle) WriteTo(w io.Writer) (n int64, err error) {
	return Byte(a).WriteTo(w)
}

func (a *Angle) ReadFrom(r io.Reader) (n int64, err error) {
	return (*Byte)(a).ReadFrom(r)
}

func (u UUID) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(u[:])
	return int64(nn), err
}

func (u *UUID) ReadFrom(r io.Reader) (n int64, err error) {
	nn, err := io.ReadFull(r, (*u)[:])
	return int64(nn), err
}

func (b ByteArray) WriteTo(w io.Writer) (n int64, err error) {
	n1, err := VarInt(len(b)).WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(b)
	return n1 + int64(n2), err
}

func (b *ByteArray) ReadFrom(r io.Reader) (n int64, err error) {
	var Len VarInt
	n1, err := Len.ReadFrom(r)
	if err != nil {
		return n1, err
	}
	if cap(*b) < int(Len) {
		*b = make(ByteArray, Len)
	} else {
		*b = (*b)[:Len]
	}
	n2, err := io.ReadFull(r, *b)
	return n1 + int64(n2), err
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
	Angle      Byte
	UUID       uuid.UUID
	ByteArray  []byte
	Position   struct {
		X, Y, Z int
	}
)

func (a Angle) ToDeg() float64 {
	return 360 * float64(a) / 256
}

func (a Angle) ToRad() float64 {
	return 2 * math.Pi * float64(a) / 256
}

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
