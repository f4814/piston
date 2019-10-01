package protocol

type Direction int

const (
	Clientbound Direction = iota
	Serverbound
)

type Packet interface {
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

type Position struct {
	X int
	Y int
	Z int
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

type Pong struct {
	//generator:idmap Direction=Clientbound State=Status ID=0x01
	Payload int64 `minecraft:"Long"`
}

// Login
type Disconnect struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x00
	//generator:idmap Direction=Clientbound State=Play ID=0x1A
	Reason string `minecraft:"String"`
}

type EncryptionRequest struct {
	//generator:idmap Direction=Clientbound State=Login ID=0x01
	ServerID          string `minecraft:"String"`
	PublicKeyLength   int32  `minecraft:"VarInt"`
	PublicKey         []byte `minecraft:"Array Byte"`
	VerifyTokenLength int32  `minecraft:"VarInte"`
	VerifyToken       []byte `minecraft:"Array Byte"`
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
	SharedSecret       []byte `minecraft:"Array Byte"`
	VerifyTokenLength  int32  `minecraft:"VarInt"`
	VerifyToken        []byte `minecraft:"Array Byte"`
}

type AnimationClientbound struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x06
	EID       int32 `minecraft:"VarInt"`
	Animation uint8 `minecraft:"UnsignedByte"`
}

// Play
type SpawnObject struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x00
	EID   int32 `minecraft:"VarInt"`
	UUID  interface{}
	Type  int32   `minecraft:"VarInt"`
	X     float64 `minecraft:"Double"`
	Y     float64 `minecraft:"Double"`
	Z     float64 `minecraft:"Double"`
	Pitch byte    `minecraft:"Angle"`
	Yaw   byte    `minecraft:"Angle"`
	Data  int32   `minecraft:"Int"`
	VX    int16   `minecraft:"Short"`
	VY    int16   `minecraft:"Short"`
	VZ    int16   `minecraft:"Short"`
}

type ServerDifficulty struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x0D
	Difficulty       uint8 `minecraft:"UnsignedByte"`
	DifficultyLocked bool  `minecraft:"Boolean"`
}

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

type ChunkData struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x21
	X                   int32         `minecraft:"Int"`
	Z                   int32         `minecraft:"Int"`
	FullChunk           bool          `minecraft:"Boolean"`
	PrimaryBitMask      int32         `minecraft:"VarInt"`
	Heightmaps          []int         `minecraft:"ChunkHeightMap"`
	Size                int32         `minecraft:"VarInt"`
	Data                []byte        `minecraft:"Array Byte"`
	NumberBlockEntities int32         `minecraft:"VarInt"`
	BlockEntities       []interface{} `minecraft:"Array NBT"`
}

type ChunkSection struct {
	BlockCount   int16
	BitsPerBlock uint8
	Palette      []int32
	DataArray    []uint
}

type PlayerAbilities struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x31
	Flags               byte    `minecraft:"Byte"`
	FlyingSpeed         float32 `minecraft:"Float"`
	FieldOfViewModifier float32 `minecraft:"Float"`
}

type PlayerPositionAndLookClientbound struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x35
	X          float64 `minecraft:"Double"`
	Y          float64 `minecraft:"Double"`
	Z          float64 `minecraft:"Double"`
	Yaw        float32 `minecraft:"Float"`
	Pitch      float32 `minecraft:"Float"`
	Flags      byte    `minecraft:"Byte"`
	TeleportID int32   `minecraft:"VarInt"`
}

type HeldItemChangeClientbound struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x3f
	Slot byte `minecraft:"Byte"`
}

type SpawnPosition struct {
	//generator:idmap Direction=Clientbound State=Play ID=0x4d
	Location Position `minecraft:"Position"`
}

type TeleportConfirm struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x00
	TeleportID int32 `minecraft:"VarInt"`
}

type ClientStatus struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x04
	ActionID int32 `minecraft:"VarInt"`
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
	//generator:idmap Direction=Clientbound State=Play ID=0x18
	Channel string `minecraft:"Identifier"`
	Data    []byte `minecraft:"Array Byte"`
}

type PlayerPositionAndLookServerbound struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x12
	X        float64 `minecraft:"Double"`
	Y        float64 `minecraft:"Double"`
	Z        float64 `minecraft:"Double"`
	Yaw      float32 `minecraft:"Float"`
	Pitch    float32 `minecraft:"Float"`
	OnGround bool    `minecraft:"Boolean"`
}

type AnimationServerbound struct {
	//generator:idmap Direction=Serverbound State=Play ID=0x2a
	Hand int32 `minecraft:"VarInt"`
}
