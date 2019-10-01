package protocol

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
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
	State       State
	Logger      log.Entry
	compression int
	conn        net.TCPConn

	read  Direction
	write Direction
}

func NewConnection(t string, conn *net.TCPConn) (*Connection, error) {
	var read, write Direction

	switch t {
	case "server":
		read = Serverbound
		write = Clientbound
	case "client":
		read = Clientbound
		write = Serverbound
	default:
		return nil, errors.New("Fail")
	}

	return &Connection{
		State:       Handshaking,
		conn:        *conn,
		compression: 0,
		read:        read,
		write:       write,
	}, nil
}

func (c *Connection) GetLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"remote":      c.conn.RemoteAddr(),
		"state":       c.State,
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
	_, err = io.ReadFull(&c.conn, buf)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf)

	id, err := readVarInt(reader)
	c.GetLogger().Tracef("%d %x", length, id)
	if err != nil {
		return nil, err
	}

	packetType := idToPacketType(c.read, c.State, id)
	if packetType == nil {
		return nil, &PacketError{id}
	}

	packet := reflect.New(packetType).Elem()

	for i := 0; i < packet.NumField(); i++ {
		tag := packet.Type().Field(i).Tag.Get("minecraft")

		if tag == "" || tag == "-" {
			continue
		}

		x, err := readValue(tag, reader, packetType, packet)

		if err != nil {
			return nil, err
		}

		packet.Field(i).Set(reflect.ValueOf(x))
	}

	p := packet.Interface().(Packet)

	c.GetLogger().WithFields(log.Fields{
		"content": fmt.Sprintf("%+v", p),
		"type":    fmt.Sprintf("%T", p),
	}).Trace("Recieved Package")

	return p, nil
}

func readValue(mcType string, reader interface { // nolint: funlen
	io.Reader
	io.RuneReader
	Len() int
},
	packetType reflect.Type, value reflect.Value) (interface{}, error) {
	var (
		x   interface{}
		err error
	)

	switch mcType {
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
	case "Array Byte":
		if packetType == reflect.TypeOf(PluginMessage{}) {
			x, err = readByteArray(reader, reader.Len()) // XXX
		} else if packetType == reflect.TypeOf(ChunkData{}) {
			x, err = readByteArray(reader, int(value.FieldByName("Size").Interface().(int32)))
		} else {
			return nil, &TypeError{mcType}
		}
	case "Array NBT":
		if packetType == reflect.TypeOf(ChunkData{}) {
			num := value.FieldByName("NumberBlockEntities").Interface().(int32)
			array := make([]interface{}, num)
			for i := int32(0); i < num; i++ {
				nbt, err := readNBTTag(reader)
				if err != nil {
					return nil, err
				}

				array[i] = nbt
			}
			x = array
		} else {
			return nil, &TypeError{mcType}
		}
	case "Long":
		x, err = readLong(reader)
	case "Float":
		x, err = readFloat(reader)
	case "Position":
		x, err = readPosition(reader)
	case "Double":
		x, err = readDouble(reader)
	case "NBT":
		x, err = readNBTTag(reader)
	case "Chat":
		x, err = readChat(reader)
	case "VarLong":
		x, err = readVarLong(reader)
	case "EntityMetadata":
		x, err = readEntityMetadata(reader)
	case "Slot":
		x, err = readSlot(reader)
	case "Angle":
		x, err = readAngle(reader)
	case "UUID":
		x, err = readUUID(reader)
	case "ChunkHeightMap":
		nbt, err := readNBTTag(reader)
		if err != nil {
			return nil, err
		}

		compacted := nbt.(map[string]interface{})["MOTION_BLOCKING"].([]int64)

		r := make([]int, 256)

		for i := 0; i < 256; i++ {
			start := int(i * 9 / 64)
			startOffs := uint(i * 9 % 64)

			end := int(((i+1)*9 - 1) / 64)
			var val int
			if start == end {
				val = int(compacted[start] >> startOffs)
			} else {
				endOffs := 64 - startOffs
				val = int(compacted[start]>>startOffs | compacted[end]<<endOffs)
			}

			r[i] = val
		}

		x = r
	default:
		return nil, &TypeError{mcType}
	}

	return x, err
}

func (c *Connection) WritePacket(p Packet) error {
	n := reflect.ValueOf(p)

	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)

	err := writeVarInt(buffer, packetToID(c.write, c.State, p))
	if err != nil {
		c.GetLogger().Trace(err)
		return err
	}

	for i := 0; i < n.NumField(); i++ {
		tag := n.Type().Field(i).Tag.Get("minecraft")

		if tag == "" || tag == "-" {
			continue
		}

		val := n.Field(i).Interface()

		err := writeValue(val, tag, buffer)

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
		"type":    fmt.Sprintf("%T", p),
		"content": fmt.Sprintf("%+v", p),
	}).Trace("Sent packet")
	return nil
}

func writeValue(val interface{}, tag string, buffer *bytes.Buffer) error {
	var err error

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
	case "Array Byte":
		err = writeByteArray(buffer, val.([]byte))
	case "Array Int":
		for _, v := range val.([]int32) {
			err = writeInt(buffer, v)
			if err != nil {
				return err
			}
		}
	case "Array VarInt":
		for _, v := range val.([]int32) {
			err = writeVarInt(buffer, v)
			if err != nil {
				return err
			}
		}
	case "Array NBT":
		for _, v := range val.([]interface{}) {
			err = writeNBTTag(buffer, v)
			if err != nil {
				return err
			}
		}
	case "Long":
		err = writeLong(buffer, val.(int64))
	case "Float":
		err = writeFloat(buffer, val.(float32))
	case "Position":
		err = writePosition(buffer, val.(Position))
	case "Double":
		err = writeDouble(buffer, val.(float64))
	case "NBT":
		err = writeNBTTag(buffer, val)
	case "Chat":
		err = writeChat(buffer, val.(string))
	case "VarLong":
		err = writeVarLong(buffer, val.(int64))
	case "EntityMetadata":
		err = writeEntityMetadata(buffer, val.(string))
	case "Slot":
		err = writeSlot(buffer, val.(string))
	case "Angle":
		err = writeAngle(buffer, val.(byte))
	case "UUID":
		err = writeUUID(buffer, val.(string))
	case "ChunkHeightMap":
	default:
		return &TypeError{tag}
	}

	return err
}

type PacketError struct {
	ID int32
}

type TypeError struct {
	Type string
}

func (e *PacketError) Error() string {
	return fmt.Sprintf("Unknown packet ID: %x", e.ID)
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("Cannot read/write minecraft type: %s", e.Type)
}
