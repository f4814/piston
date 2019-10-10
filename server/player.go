package server

import (
	"github.com/f4814n/piston/protocol"
	"github.com/google/uuid"
	"math/rand"
)

type Player struct {
	Username string
	UUID     uuid.UUID

	X int32
	Y int32
	Z int32

	Yaw   uint8
	Pitch uint8

	conn *protocol.Conn

	rx chan protocol.Packet
	tx chan protocol.Packet
}

func NewPlayer(username string, uuid uuid.UUID, conn *protocol.Conn) Player {
	p := Player{Username: username, UUID: uuid, conn: conn}

	p.rx = make(chan protocol.Packet)
	p.tx = make(chan protocol.Packet)

	return p
}

func (p *Player) Start() {
	defer p.conn.Close()

	go p.readQueue()
	go p.writeQueue()

	p.initialize()

	for {
		packet := <-p.rx
		switch v := packet.(type) {
		case *protocol.AnimationServerbound:
			var animation uint8
			if v.Hand == 0 {
				animation = 0
			} else {
				animation = 3
			}
			p.tx <- &protocol.AnimationClientbound{EID: 6, Animation: animation}
		}
	}
}

func (p *Player) initialize() {
	p.tx <- protocol.JoinGame{
		EID:              6,
		Gamemode:         0,
		Dimension:        0,
		MaxPlayers:       0,
		LevelType:        "default",
		ViewDistance:     10,
		ReducedDebugInfo: false,
	}

	p.tx <- protocol.ServerDifficulty{
		Difficulty:       0,
		DifficultyLocked: true,
	}

	p.tx <- protocol.PlayerAbilities{
		Flags:               0,
		FlyingSpeed:         0.05,
		FieldOfViewModifier: 0.1,
	}

	// for x := 0; x < 7; x++ {
	// 	for z := 0; z < 7; z++ {
	// 		chunkData := protocol.ChunkData{
	// 			X:                   int32(x),
	// 			Z:                   int32(z),
	// 			FullChunk:           true,
	// 			PrimaryBitMask:      1,
	// 			NumberBlockEntities: 0,
	// 		}
	//
	// 		section := protocol.ChunkSection{
	// 			BlockCount:   16 * 16 * 16,
	// 			BitsPerBlock: 4,
	// 			Palette:      []int32{2},
	// 		}
	//
	// 		section.DataArray = make([]uint, 16*16*16)
	//
	// 		biomes := make([]int32, 256)
	//
	// 		for x := 0; x < 16*16; x++ {
	// 			section.DataArray[x] = 0
	// 		}
	//
	// 		err := chunkData.LoadData([]protocol.ChunkSection{section}, biomes)
	// 		if err != nil {
	// 			p.conn.GetLogger().Error(err)
	// 		}
	//
	// 		p.tx <- chunkData
	// 	}
	// }

	p.tx <- protocol.SpawnPosition{
		Location: protocol.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}

	p.tx <- protocol.PlayerPositionAndLookClientbound{
		X:          0,
		Y:          0,
		Z:          0,
		Yaw:        180,
		Pitch:      90,
		Flags:      0,
		TeleportID: rand.Int31(),
	}
}

func (p *Player) readQueue() {
	for {
		packet, err := p.conn.ReadPacket()
		if err != nil {
			p.conn.GetLogger().Error(err)
			return
		}

		p.rx <- packet
	}
}

func (p *Player) writeQueue() {
	for {
		packet := <-p.tx

		err := p.conn.WritePacket(packet)
		if err != nil {
			p.conn.GetLogger().Error(err)
			return
		}
	}
}
