package protocol

type Direction int

const (
	Clientbound Direction = iota
	Serverbound
)

type Packet interface {
}

// Handshake
type Handshake struct {
	//generator:idmap Direction=Serverbound State=Handshaking ID=0x00
	ProtocolVersion int32  `minecraft:"VarInt"`
	ServerAddress   string `minecraft:"String"`
	ServerPort      uint16 `minecraft:"UnsignedShort"`
	NextState       int32  `minecraft:"VarInt"`
}

// Status
type Request struct {
	//generator:idmap Direction=Serverbound State=Status ID=0x00
}

type Ping struct {
	//generator:idmap Direction=Serverbound State=Status ID=0x01
	Payload int64 `minecraft:"Long"`
}

type Response struct {
	//generator:idmap Direction=Clientbound State=Status ID=0x00
	JSONResponse string `minecraft:"String"`
}

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

type Pong struct {
	//generator:idmap Direction=Clientbound State=Status ID=0x01
	Payload int64 `minecraft:"Long"`
}

// Login
type Disconnect struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x00
	//generator:idmap Direction=Clientbound State=Play ID=0x1A
	Reason string `minecraft:"Chat"`
}

type EncryptionRequest struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x01
	ServerID          string `minecraft:"String"`
	PublicKeyLength   int32  `minecraft:"VarInt"`
	PublicKey         []byte `minecraft:"ByteArray"`
	VerifyTokenLength int32  `minecraft:"VarInte"`
	VerifyToken       []byte `minecraft:"ByteArray"`
}

type LoginSuccess struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x02
	UUID     string `minecraft:"String"`
	Username string `minecraft:"String"`
}

type SetCompression struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x03
	Treshold int32 `minecraft:"VarInt"`
}

type LoginStart struct {
	//generator:idmap Direction=Serverbound State=Login ID=0x00
	Name string `minecraft:"String"`
}

type EncryptionResponse struct {
	//generator:idmap Direction=Serverbound State=Login ID=0x01
	SharedSecretLength int32  `minecraft:"VarInt"`
	SharedSecret       []byte `minecraft:"ByteArray"`
	VerifyTokenLength  int32  `minecraft:"VarInt"`
	VerifyToken        []byte `minecraft:"ByteArray"`
}

// Play
type JoinGame struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x25
	EID              int32  `minecraft:"Int"`
	Gamemode         uint8  `minecraft:"UnsignedByte"`
	Dimension        int32  `minecraft:"Int"`
	MaxPlayers       uint8  `minecraft:"UnsignedByte"`
	LevelType        string `minecraft:"String"`
	ViewDistance     int32  `minecraft:"VarInt"`
	ReducedDebugInfo bool   `minecraft:"Boolean"`
}

type ClientSettings struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x05
	Locale           string `minecraft:"String"`
	ViewDistance     byte   `minecraft:"Byte"`
	ChatMode         int32  `minecraft:"VarInt"`
	ChatColors       bool   `minecraft:"Boolean"`
	DisplaySkinParts uint8  `minecraft:"UnsignedByte"`
	MainHand         int32  `minecraft:"VarInt"`
}

type PluginMessage struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x0b
	Channel string `minecraft:"Identifier"`
	Date    []byte `minecraft:"ByteArray"`
}
