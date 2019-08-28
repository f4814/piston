package protocol

import (
	"io"
	"strings"
	"errors"
	"unicode/utf8"
	"bytes"
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
	panic("readShort")
	return 0, nil
}

func readUnsignedShort(r io.Reader) (uint16, error) {
	num := make([]byte, 2)

	_, err := r.Read(num)
	if err != nil {
		return 0, err
	}

	var result uint16
	result = uint16(num[0]) << 8
	result |= uint16(num[1])

	return result, nil
}

func readInt(r io.Reader) (int32, error) {
	var result int32
	temp := make([]byte, 4)

	_, err := r.Read(temp)
	if err != nil {
		return 0, nil
	}

	result = int32(temp[0])
	result <<= 8
	result |= int32(temp[1])
	result <<= 8
	result |= int32(temp[2])
	result <<= 8
	result |= int32(temp[3])
	
	return result, nil
}

func readLong(r io.Reader) (int64, error) {
	panic("readLong")
	return 0, nil
}

func readFloat(r io.Reader) (float32, error) {
	panic("readLong")
	return 0, nil
}

func readDouble(r io.Reader) (float64, error) {
	panic("readDouble")
	return 0, nil
}

func readString(r interface{io.Reader; io.RuneReader}) (string, error) {
	runes, err := readVarInt(r)
	if err != nil {
		return "", err
	}

	var result string
	for ;runes > 0; runes-- {
		r, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		result += string(r)
	}
	return result, nil
}

func readChat(r io.Reader) (string, error) {
	panic("readChat")
	return "", nil
}

func readIdentifier(r interface{io.Reader; io.RuneReader}) (string, error) {
	temp, err := readString(r)

	if err != nil {
		return "", err
	}

	if strings.Index(temp, ":") == -1 {
		return "minecraft:" + temp, nil
	} 

	return temp, nil
}

func readVarInt(r io.Reader) (int32, error) {
	var (
		numRead int
		result	int32
		read	byte
		err		error
	)

	buf := make([]byte, 1)

	for ok := true; ok; ok = ((read & 0x80) != 0) { // 0x80 == 0b10000000
		if _, err = r.Read(buf); err != nil {
			return 0, err
		}

		read = buf[0]
		value := read & 0x7F // 0x7F == 0b01111111
		result |= int32(value << uint(7 * numRead))

		numRead++
		if numRead > 5 {
			return 0, errors.New("Unable to decode VarInt")
		}
	}

	return result, nil
}

func readVarLong(r io.Reader) (int64, error) {
	panic("readVarLong")
	return 0, nil
}

func readEntityMetadata(r io.Reader) (string, error) {
	panic("readEntityMetadata")
	return "", nil
}

func readSlot(r io.Reader) (string ,error) {
	panic("readSlot")
	return "", nil
}

func readNBTTag(r io.Reader) (string, error) {
	panic("readNBTTag")
	return "", nil
}

func readPosition(r io.Reader) (string, error) {
	panic("readPosition")
	return "", nil
}

func readAngle(r io.Reader) (byte, error) {
	panic("readAngle")
	return 0, nil
}

func readUUID(r io.Reader) (string, error) {
	panic("readUUID")
	return "", nil
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
	panic("writeShort")
	return nil
}

func writeUnsignedShort(w io.Writer, u uint16) error {
	panic("writeUnsignedShort")
	return nil
}

func writeInt(w io.Writer, i int32) error {
	var result int32
	for i := 3; i >= 0; i-- {
		temp := byte(result >> uint(8 * i))
		_, err := w.Write([]byte{temp})
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLong(w io.Writer, i int64) error {
	panic("writeLong")
	return nil
}

func writeFloat(w io.Writer, f float32) error {
	panic("writeLong")
	return nil
}

func writeDouble(w io.Writer, f float64) error {
	panic("writeDouble")
	return nil
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
	panic("writeChat")
	return nil
}

func writeIdentifier(w *bytes.Buffer, s string) error {
	panic("writeIdentifier")
	return nil
}

func writeVarInt(w io.Writer, i int32) error {
	value := uint32(i) // We have to shift the sign bit too
	for ok := true; ok; ok = (value != 0) {
		temp := (byte)(value & 0x7F) // 0x7F == 0b01111111
		value >>= 7
		if (value != 0) {
			temp |= 0x80 // 0x80 == 0b10000000
		}
		w.Write([]byte{temp})
	}
	return nil
}

func writeVarLong(w io.Writer, i int64) error {
	panic("writeVarLong")
	return nil
}

func writeEntityMetadata(w io.Writer, s string) error {
	panic("writeEntityMetadata")
	return nil
}

func writeSlot(w io.Writer, s string) error {
	panic("writeSlot")
	return nil
}

func writeNBTTag(w io.Writer, s string) error {
	panic("writeNBTTag")
	return nil
}

func writePosition(w io.Writer, s string) error {
	panic("writePosition")
	return nil
}

func writeAngle(w io.Writer, b byte) error {
	panic("writeAngle")
	return nil
}

func writeUUID(w io.Writer, s string) error {
	panic("writeUUID")
	return nil
}

func writeByteArray(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}
