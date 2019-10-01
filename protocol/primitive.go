package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Tnze/go-mc/nbt"
	"io"
	"math"
	"strings"
	"unicode/utf8"
)

func readBoolean(r io.Reader) (bool, error) {
	temp := make([]byte, 1)
	_, err := r.Read(temp)

	if err != nil {
		return false, err
	}

	if temp[0] == 0x00 {
		return false, nil
	} else if temp[0] == 0x01 {
		return true, nil
	}

	return false, errors.New("Neither true nor false")
}

func readByte(r io.Reader) (byte, error) {
	temp := make([]byte, 1)
	_, err := r.Read(temp)
	return temp[0], err
}

func readUnsignedByte(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)

	_, err := r.Read(buf)

	return buf[0], err
}

func readShort(r io.Reader) (int16, error) {
	temp := make([]byte, 2)
	_, err := r.Read(temp)
	if err != nil {
		return 0, nil
	}

	return int16(binary.BigEndian.Uint16(temp)), nil
}

func readUnsignedShort(r io.Reader) (uint16, error) {
	num := make([]byte, 2)

	_, err := r.Read(num)
	if err != nil {
		return 0, err
	}

	result := uint16(num[0]) << 8
	result |= uint16(num[1])

	return result, nil
}

func readInt(r io.Reader) (int32, error) {
	temp := make([]byte, 4)

	_, err := r.Read(temp)
	if err != nil {
		return 0, err
	}

	result := int32(binary.BigEndian.Uint32(temp))
	return result, nil
}

func readLong(r io.Reader) (int64, error) {
	temp := make([]byte, 8)

	_, err := r.Read(temp)
	if err != nil {
		return 0, err
	}

	result := int64(binary.BigEndian.Uint64(temp))
	return result, nil
}

func readFloat(r io.Reader) (float32, error) {
	u, err := readInt(r)
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(uint32(u)), nil
}

func readDouble(r io.Reader) (float64, error) {
	u, err := readInt(r)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(uint64(u)), nil
}

func readString(r interface {
	io.Reader
	io.RuneReader
}) (string, error) {
	runes, err := readVarInt(r)
	if err != nil {
		return "", err
	}

	var result string
	for ; runes > 0; runes-- {
		r, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		result += string(r)
	}
	return result, nil
}

func readChat(r io.Reader) (string, error) {
	return "", errors.New("Not implemented")
}

func readIdentifier(r interface {
	io.Reader
	io.RuneReader
}) (string, error) {
	temp, err := readString(r)

	if err != nil {
		return "", err
	}

	if !strings.Contains(temp, ":") {
		return "minecraft:" + temp, nil
	}

	return temp, nil
}

func readVarInt(r io.Reader) (int32, error) {
	var result uint32
	buf := make([]byte, 1)
	for i := 0; i < 5; i++ {
		_, err := r.Read(buf)
		if err != nil {
			return 0, errors.New("fail")
		}
		b := buf[0]

		result |= (uint32(b&0x7F) << uint(7*i))
		if b&0x80 == 0 {
			break
		}
	}

	return int32(result), nil
}

func readVarLong(r io.Reader) (int64, error) {
	return 0, errors.New("not implemented")
}

func readEntityMetadata(r io.Reader) (string, error) {
	return "", errors.New("not implemented")
}

func readSlot(r io.Reader) (string, error) {
	return "", errors.New("not implemented")
}

func readNBTTag(r io.Reader) (interface{}, error) {
	var x interface{}
	err := nbt.NewDecoder(r).Decode(&x)
	return x, err
}

func readPosition(r io.Reader) (Position, error) {
	raw, err := readLong(r)
	if err != nil {
		return Position{}, err
	}

	p := Position{
		X: int(raw >> 38),
		Y: int(raw & 0xFFF),
		Z: int(raw << 26 >> 38),
	}
	return p, nil
}

func readAngle(r io.Reader) (byte, error) {
	return 0, errors.New("not implemented")
}

func readUUID(r io.Reader) (string, error) {
	return "", errors.New("not implemented")
}

func readByteArray(r io.Reader, length int) ([]byte, error) {
	temp := make([]byte, length)
	_, err := r.Read(temp)

	return temp, err
}

func writeBoolean(w io.Writer, b bool) error {
	var val byte
	if b {
		val = 0x01
	} else {
		val = 0x00
	}

	_, err := w.Write([]byte{val})
	return err
}

func writeByte(w io.Writer, b byte) error {
	_, err := w.Write([]byte{b})
	return err
}

func writeUnsignedByte(w io.Writer, u uint8) error {
	_, err := w.Write([]byte{u})
	return err
}

func writeShort(w io.Writer, i int16) error {
	temp := make([]byte, 2)
	binary.BigEndian.PutUint16(temp, uint16(i))
	_, err := w.Write(temp)
	if err != nil {
		return err
	}
	return nil
}

func writeUnsignedShort(w io.Writer, u uint16) error {
	var a, b byte = byte(u >> 8), byte(u & 0xFF)
	_, err := w.Write([]byte{a, b})
	return err
}

func writeInt(w io.Writer, i int32) error {
	temp := make([]byte, 4)
	binary.BigEndian.PutUint32(temp, uint32(i))
	_, err := w.Write(temp)
	if err != nil {
		return err
	}
	return nil
}

func writeLong(w io.Writer, i int64) error {
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, uint64(i))
	_, err := w.Write(temp)
	if err != nil {
		return err
	}
	return nil
}

func writeFloat(w io.Writer, f float32) error {
	return writeInt(w, int32(math.Float32bits(f)))
}

func writeDouble(w io.Writer, f float64) error {
	return writeLong(w, int64(math.Float64bits(f)))
}

func writeString(w *bytes.Buffer, s string) error {
	runes := len(s)

	err := writeVarInt(w, int32(runes))

	if err != nil {
		return err
	}

	for s != "" {
		r, l := utf8.DecodeRuneInString(s)
		s = s[l:]
		_, err := w.WriteRune(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeChat(w *bytes.Buffer, s string) error {
	return errors.New("not implemented")
}

func writeIdentifier(w *bytes.Buffer, s string) error {
	return writeString(w, s)
}

func writeVarInt(w io.Writer, i int32) error {
	value := uint32(i) // We have to shift the sign bit too
	for ok := true; ok; ok = (value != 0) {
		temp := (byte)(value & 0x7F) // 0x7F == 0b01111111
		value >>= 7
		if value != 0 {
			temp |= 0x80 // 0x80 == 0b10000000
		}

		_, err := w.Write([]byte{temp})
		if err != nil {
			return err
		}
	}
	return nil
}

func writeVarLong(w io.Writer, i int64) error {
	return errors.New("not implemented")
}

func writeEntityMetadata(w io.Writer, s string) error {
	return errors.New("not implemented")
}

func writeSlot(w io.Writer, s string) error {
	return errors.New("not implemented")
}

func writeNBTTag(w io.Writer, i interface{}) error {
	return nbt.Marshal(w, i)
}

func writePosition(w io.Writer, p Position) error {
	raw := uint64(((p.X & 0x3FFFFFF) << 38) | ((p.Z & 0x3FFFFFF) << 12) | (p.Y & 0xFF))
	return writeLong(w, int64(raw))
}

func writeAngle(w io.Writer, b byte) error {
	return errors.New("not implemented")
}

func writeUUID(w io.Writer, s string) error {
	return errors.New("not implemented")
}

func writeByteArray(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}
