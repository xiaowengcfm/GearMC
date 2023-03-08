package packet

import "io"

type Packet struct {
	ID   int32
	Data []byte
}

func (p *Packet) Pack(w io.Writer) error {
	return nil
}

func (p *Packet) UnPack(r io.Reader) error {
	return nil
}
