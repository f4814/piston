package protocol

import (
	"reflect"
)

type Direction int

const (
	Clientbound Direction = iota
	Serverbound
)

type Packet interface {
}

func packetToID(direction Direction, state State, packet Packet) int32 {
	switch packet.(type) {
	// Handshaking
	case Handshake:
		return 0x00
	// Status
	case Request:
		if direction == Serverbound && state == Status {
			return 0x00
		}
	case Ping:
		if direction == Serverbound && state == Status {
			return 0x01
		}
	case Response:
		if direction == Clientbound && state == Status {
			return 0x00
		}
	case Pong:
		if direction == Clientbound && state == Status {
			return 0x01
		}
	// Login
	case DisconnectLogin:
		if direction == Clientbound && state == Login {
			return 0x00
		}
	case EncryptionRequest:
		if direction == Clientbound && state == Login {
			return 0x01
		}
	case LoginSuccess:
		if direction == Clientbound && state == Login {
			return 0x02
		}
	case SetCompression:
		if direction == Clientbound && state == Login {
			return 0x03
		}
	case LoginStart:
		if direction == Serverbound && state == Login {
			return 0x00
		}
	case EncryptionResponse:
		if direction == Serverbound && state == Login {
			return 0x01
		}
	// Play
	case JoinGame:
		if direction == Clientbound && state == Play {
			return 0x25
		}
	case ClientSettings:
		if direction == Serverbound && state == Play {
			return 0x05
		}
	}

	return -1
}

func idToPacketType(direction Direction, state State, id int32) reflect.Type {
	if direction == Serverbound {
		if state == Handshaking {
			if id == 0x00 {
				return reflect.TypeOf(Handshake{})
			}
		} else if state == Status {
			if id == 0x00 {
				return reflect.TypeOf(Request{})
			} else if id == 0x01 {
				return reflect.TypeOf(Ping{})
			}
		} else if state == Login {
			if id == 0x00 {
				return reflect.TypeOf(LoginStart{})
			} else if id == 0x01 {
				return reflect.TypeOf(EncryptionResponse{})
			}
		} else if state == Play {
			if id == 0x05 {
				return reflect.TypeOf(ClientSettings{})
			} else if id == 0x0b {
				return reflect.TypeOf(PluginMessage{})
			}
		}
	} else if direction == Clientbound {
		if state == Handshaking {
			return nil
		} else if state == Status {
		} else if state == Login {
			if id == 0x00 {
				return reflect.TypeOf(DisconnectLogin{})
			} else if id == 0x01 {
				return reflect.TypeOf(EncryptionRequest{})
			} else if id == 0x02 {
				return reflect.TypeOf(LoginSuccess{})
			} else if id == 0x03 {
				return reflect.TypeOf(SetCompression{})
			}
		} else if state == Play {
			if id == 0x25 {
				return reflect.TypeOf(JoinGame{})
			}
		}
	}

	return nil
}

// Handshake
type Handshake struct {
	ProtocolVersion int32  `minecraft:"VarInt"`
	ServerAddress   string `minecraft:"String"`
	ServerPort      uint16 `minecraft:"UnsignedShort"`
	NextState       int32  `minecraft:"VarInt"`
}

// Status
type Request struct {
}

type Ping struct {
	Payload int64 `minecraft:"Long"`
}

type Response struct {
	JSONResponse string `minecraft:"String"`
}

type ResponseJSON struct {
	// Convenience
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

type Pong struct {
	Payload int64 `minecraft:"Long"`
}

// Login
type DisconnectLogin struct {
	Reason string `minecraft:"Chat"`
}

type EncryptionRequest struct {
	ServerID          string `minecraft:"String"`
	PublicKeyLength   int32  `minecraft:"VarInt"`
	PublicKey         []byte `minecraft:"ByteArray"`
	VerifyTokenLength int32  `minecraft:"VarInte"`
	VerifyToken       []byte `minecraft:"ByteArray"`
}

type LoginSuccess struct {
	UUID     string `minecraft:"String"`
	Username string `minecraft:"String"`
}

type SetCompression struct {
	Treshold int32 `minecraft:"VarInt"`
}

type LoginStart struct {
	Name string `minecraft:"String"`
}

type EncryptionResponse struct {
	SharedSecretLength int32  `minecraft:"VarInt"`
	SharedSecret       []byte `minecraft:"ByteArray"`
	VerifyTokenLength  int32  `minecraft:"VarInt"`
	VerifyToken        []byte `minecraft:"ByteArray"`
}

// Play
type JoinGame struct {
	EID	int32 `minecraft:"Int"`
	Gamemode uint8 `minecraft:"UnsignedByte"`
	Dimension	int32 `minecraft:"Int"`
	MaxPlayers	uint8 `minecraft:"UnsignedByte"`
	LevelType	string `minecraft:"String"`
	ViewDistance int32 `minecraft:"VarInt"`
	ReducedDebugInfo bool `minecraft:"Boolean"`
}

type ClientSettings struct {
	Locale string `minecraft:"String"`
	ViewDistance byte `minecraft:"Byte"`
	ChatMode int32 `minecraft:"VarInt"`
	ChatColors bool `minecraft:"Boolean"`
	DisplaySkinParts uint8 `minecraft:"UnsignedByte"`
	MainHand int32 `minecraft:"VarInt"`
}

type PluginMessage struct {
	Channel	string `minecraft:"Identifier"`
	Date []byte `minecraft:"ByteArray"`
}
