package net

import (
	"github.com/xiaowengcfm/gearmc/net/packet"
	"io"
	"net"
)

const (
	NetworkProtocol = "tcp"
)

type Listener struct {
	net.Listener
}

type Conn struct {
	io.Writer
	io.Reader
	Socket net.Conn
}

func Listen(addr string) (*Listener, error) {
	l, err := net.Listen(NetworkProtocol, addr)
	return &Listener{l}, err
}

func (l *Listener) Accept() (*Conn, error) {
	conn, err := l.Listener.Accept()
	return &Conn{
		Writer: conn,
		Reader: conn,
		Socket: conn,
	}, err
}

func (c *Conn) WritePacket(p packet.Packet) error {
	return nil
}

func (c *Conn) ReadPacket(p *packet.Packet) error {
	return nil
}
