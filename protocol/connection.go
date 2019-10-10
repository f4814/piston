package protocol

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type State int

//go:generate stringer -type=State .

const (
	Uninitialized State = iota
	Handshaking
	Play
	Status
	Login
)

type Direction int

//go:generate stringer -type=Direction .

const (
	Clientbound Direction = iota
	Serverbound
)

// Like net.Listener but for minecraft connections
type Listener struct {
	net.TCPListener
}

// Create a Listener
func Listen(address string) (*Listener, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Listener{
		TCPListener: *ln.(*net.TCPListener),
	}, nil
}

// Listen for minecraft connections
func (l *Listener) Accept() (*Conn, error) {
	tcpConn, err := l.TCPListener.Accept()
	if err != nil {
		return nil, err
	}

	return &Conn{
		State: Handshaking,
		conn:  *tcpConn.(*net.TCPConn),
		read:  Serverbound,
		write: Clientbound,
	}, nil
}

// A minecraft connection with an API similar to net.Conn
type Conn struct {
	State       State
	compression int
	conn        net.TCPConn

	read  Direction
	write Direction
}

// Connect to a minecraft server
func Dial(address string) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Conn{
		State: Uninitialized,
		conn:  *c.(*net.TCPConn),
		read:  Clientbound,
		write: Serverbound,
	}, nil
}

// Get a suitable Logging Context for the connection
// TODO Move this
func (c *Conn) GetLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"remote":      c.conn.RemoteAddr(),
		"state":       c.State,
		"compression": c.compression,
	})
}

// Close the connection
func (c *Conn) Close() {
	c.GetLogger().Info("Closing Connection")
	c.conn.Close()
}

// Read a Packet from the connection
func (c *Conn) ReadPacket() (Packet, error) {
	length, err := readVarInt(&c.conn)
	if err != nil {
		return nil, fmt.Errorf("unable to read packet length: %w", err)
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(&c.conn, buf)
	if err != nil {
		return nil, fmt.Errorf("unable to read packet data: %w", err)
	}

	reader := bytes.NewBuffer(buf)

	id, err := readVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read packet id: %w", err)
	}

	packet, err := packetDecoder(c.State, c.read, id, reader)
	if err != nil {
		return nil, err
	}

	c.GetLogger().WithFields(log.Fields{
		"content": fmt.Sprintf("%+v", packet),
		"type":    fmt.Sprintf("%T", packet),
	}).Trace("Recieved Packet")

	return packet, nil
}

// Write a Packet to the connection
func (c *Conn) WritePacket(p Packet) error {
	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)

	err := writeVarInt(buffer, p.id())
	if err != nil {
		c.GetLogger().Trace(err)
		return err
	}

	err = p.encode(buffer)
	if err != nil {
		return err
	}

	err = writeVarInt(&c.conn, int32(buffer.Len()))
	if err != nil {
		c.GetLogger().Trace(err)
		return err
	}

	_, err = c.conn.ReadFrom(buffer)

	if err != nil {
		c.GetLogger().Trace(err)
		return err
	}

	c.GetLogger().WithFields(log.Fields{
		"type":    fmt.Sprintf("%T", p),
		"content": fmt.Sprintf("%+v", p),
	}).Trace("Sent packet")
	return nil
}

type UnknownPacketError struct {
	ID    int32
	State State
}

func (e UnknownPacketError) Error() string {
	return fmt.Sprintf("unknown packet with id %s in state %s",
		strconv.FormatInt(int64(e.ID), 16), e.State)
}
