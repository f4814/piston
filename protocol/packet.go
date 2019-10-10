package protocol

import (
	"io"
)

type Packet interface {
	// Write the packet into w
	// This function can be created automatically for most packets, but sometimes
	// writing it by hand is easier
	encode(w io.Writer) (err error)

	// Return the packets ID
	id() int32
}

// We generate most of the implementations of Packet
//go:generate go run internal/generate/main.go

// Convenience
type ResponseJSON struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon string `json:"favicon,omitempty"`
}

type Position struct {
	X int
	Y int
	Z int
}

type EncryptionResponse struct {
	//minecraft:id Direction=Serverbound State=Login ID=0x01
	SharedSecretLength int32  `minecraft:"VarInt"`
	SharedSecret       []byte `minecraft:"Array Byte"`
	VerifyTokenLength  int32  `minecraft:"VarInt"`
	VerifyToken        []byte `minecraft:"Array Byte"`
}
