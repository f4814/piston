package protocol

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"net"
	"io"
	"bytes"
	"reflect"
)

type State int

const (
	Handshaking State = iota
	Play
	Status
	Login
)

func (s State) String() string {
	if s == Handshaking {
		return "Handshaking"
	} else if s == Play {
		return "Play"
	} else if s == Status {
		return "Status"
	} else if s == Login {
		return "Login"
	}
	return ""
}

type Connection struct {
	State State
	Logger log.Entry
	compression int
	conn	net.TCPConn
}

func NewConnection(conn *net.TCPConn) (*Connection, error) {
	return &Connection{
		State: Handshaking,
		conn:  *conn,
		compression: 0,
	}, nil
}

func (c *Connection) GetLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"remote": c.conn.RemoteAddr(),
		"state": c.State,
		"compression": c.compression,
	})
}

func (c *Connection) Close() {
	c.GetLogger().Info("Closing Connection")
	c.conn.Close()
}

func (c *Connection) ReadPacket() (Packet, error) {
	length, err := readVarInt(&c.conn)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	io.ReadFull(&c.conn, buf)
	reader := bytes.NewReader(buf)

	id, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}

	packetType := idToPacketType(Serverbound, c.State, id)
	if packetType == nil {
		return nil, &PacketError{id}
	}

	n := reflect.New(packetType).Elem()

	for i := 0; i < n.NumField(); i++ {
		tag := n.Type().Field(i).Tag.Get("minecraft")
		
		if tag == "" || tag == "-" {
			continue
		}

		var x interface{}

		switch tag {
		case "VarInt":
			x, err = readVarInt(reader)
		case "String":
			x, err = readString(reader)
		case "UnsignedShort":
			x, err = readUnsignedShort(reader)
		case "Int":
			x, err = readInt(reader)
		case "UnsignedByte":
			x, err = readUnsignedByte(reader)
		case "Boolean":
			x, err = readBoolean(reader)
		case "Byte":
			x, err = readByte(reader)
		case "Identifier":
			x, err = readIdentifier(reader)
		case "ByteArray":
			if packetType == reflect.TypeOf(PluginMessage{}) {
				x, err = readByteArray(reader, reader.Len()) // XXX
			} else {
				log.WithFields(log.Fields{
					"packet": packetType,
				}).Panic("Cannot handle ByteArray in this packet")
			}
		default:
			return nil, &TypeError{tag}
		}

		if err != nil {
			return nil, err
		}

		n.Field(i).Set(reflect.ValueOf(x))
	}

	p := n.Interface().(Packet)

	c.GetLogger().WithFields(log.Fields{
		"content": fmt.Sprintf("%+v", p),
		"type": fmt.Sprintf("%T", p),
	}).Trace("Recieved Package")

	return p, nil
}

func (c *Connection) WritePacket(p Packet) error {
	n := reflect.ValueOf(p)

	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)

	err := writeVarInt(buffer, packetToID(Clientbound, c.State, p))
	if err != nil {
		c.GetLogger().Trace(err)
		return err
	}

	for i := 0; i < n.NumField(); i++ {
		tag := n.Type().Field(i).Tag.Get("minecraft")

		if tag == "" || tag == "-" {
			continue
		}

		var err error
		val := n.Field(i).Interface()

		switch tag {
		case "VarInt":
			err = writeVarInt(buffer, val.(int32))
		case "String":
			err = writeString(buffer, val.(string))
		case "UnsignedShort":
			err = writeUnsignedShort(buffer, val.(uint16))
		case "Int":
			err = writeInt(buffer, val.(int32))
		case "UnsignedByte":
			err = writeUnsignedByte(buffer, val.(uint8))
		case "Boolean":
			err = writeBoolean(buffer, val.(bool))
		case "Byte":
			err = writeByte(buffer, val.(byte))
		case "Identifier":
			err = writeIdentifier(buffer, val.(string))
		case "ByteArray":
			err = writeByteArray(buffer, val.([]byte))
		default:
			return &TypeError{tag}
		}

		if err != nil {
			c.GetLogger().Trace(err)
			return err
		}
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
		"type": fmt.Sprintf("%T", p),
		"content": fmt.Sprintf("%+v", p),
	}).Trace("Sent packet")
	return nil
}

type PacketError struct {
	ID	int32
}

type TypeError struct {
	Type string
}

func (e *PacketError) Error() string {
	return fmt.Sprintf("Unknown packet ID: %d", e.ID)
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("Cannot read/write minecraft type: %s", e.Type)
}
