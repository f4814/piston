package protocol

import (
	"testing"
	"bytes"
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

		err := writeString(writer, d)

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
