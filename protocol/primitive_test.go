package protocol

import (
	"bytes"
	"testing"
)

type varIntTestpair struct {
	val int32
	raw []byte
}

var varIntTests = []varIntTestpair{
	varIntTestpair{0, []byte{0x00}},
	varIntTestpair{1, []byte{0x01}},
	varIntTestpair{2, []byte{0x02}},
	varIntTestpair{127, []byte{0x7f}},
	varIntTestpair{128, []byte{0x80, 0x01}},
	varIntTestpair{255, []byte{0xff, 0x01}},
	varIntTestpair{2147483647, []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
	varIntTestpair{-1, []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
	varIntTestpair{-2147483648, []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
}

func TestReadVarInt(t *testing.T) {
	for _, d := range varIntTests {
		reader := bytes.NewReader(d.raw)
		ret, err := readVarInt(reader)

		if err != nil {
			t.Error(err)
		}

		if ret != d.val {
			t.Fatalf("Parsed %d (should be %d)", ret, d.val)
		}
	}
}

func TestWriteVarInt(t *testing.T) {
	for _, d := range varIntTests {
		buf := make([]byte, 5)
		writer := bytes.NewBuffer(buf)

		err := writeVarInt(writer, d.val)

		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(writer.Bytes(), d.raw) {
			t.Fatalf("Wrote %v (should by %v)", writer.Bytes(), d.raw)
		}
	}
}

func TestReadWriteString(t *testing.T) {
	tests := []string{"Heyho", "", "With\nNewline"}

	for _, d := range tests {
		buf := make([]byte, 0)
		writer := bytes.NewBuffer(buf)

		err := writeVarInt(writer, int32(len(d)))
		if err != nil {
			t.Error(err)
		}

		err = writeString(writer, d)
		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(writer.Bytes())

		str, err := readString(reader)

		if err != nil {
			t.Error(err)
		}

		if str != d {
			t.Fatalf("Read %s (should be %s)", str, d)
		}
	}
}

func TestReadWriteInt(t *testing.T) {
	tests := []int32{1, -32, 16000}

	for _, d := range tests {
		buf := make([]byte, 0)
		writer := bytes.NewBuffer(buf)

		err := writeInt(writer, d)

		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(writer.Bytes())

		str, err := readInt(reader)

		if err != nil {
			t.Error(err)
		}

		if str != d {
			t.Fatalf("Read %d (should be %d)", str, d)
		}
	}
}

func TestReadWritePosition(t *testing.T) {
	tests := []Position{Position{100, 100, 100}}

	for _, d := range tests {
		buf := make([]byte, 0)
		writer := bytes.NewBuffer(buf)

		err := writePosition(writer, d)

		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(writer.Bytes())

		pos, err := readPosition(reader)

		if err != nil {
			t.Error(err)
		}

		if pos != d {
			t.Fatalf("Read %d (should be %d)", pos, d)
		}
	}
}

// func TestReadWriteNBT(t *testing.T) {
// 	tests := []interface{}{
// 		struct {
// 			Name string `nbt:"name"`
// 		}{"test"},
// 		struct {
// 			ID int `nbt:"id"`
// 		}{22},
// 	}
//
// 	for _, d := range tests {
// 		buf := make([]byte, 0)
// 		writer := bytes.NewBuffer(buf)
//
// 		err := writeNBTTag(writer, d)
//
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		reader := bytes.NewReader(writer.Bytes())
//
// 		nbt, err := readNBTTag(reader)
//
// 		if err != nil {
// 			t.Error(err)
// 		}
//
// 		if nbt != d {
// 			t.Fatalf("Read %+v (should be %+v)", nbt, d)
// 		}
// 	}
// }

func TestReadWriteUnsignedShort(t *testing.T) {
	tests := []uint16{0, 255, 16000}

	for _, d := range tests {
		buf := make([]byte, 0)
		writer := bytes.NewBuffer(buf)

		err := writeUnsignedShort(writer, d)

		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(writer.Bytes())

		short, err := readUnsignedShort(reader)

		if err != nil {
			t.Error(err)
		}

		if short != d {
			t.Fatalf("Read %d (should be %d)", short, d)
		}
	}
}
