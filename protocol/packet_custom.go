package protocol

import (
	"io"
)

func (c ChunkData) encode(w io.Writer) (err error) {
	return
}

func (c *ChunkData) decode(r io.Reader) (err error) {
	return
}

func (p PluginMessageClientbound) encode(w io.Writer) (err error) {
	return nil
}

func (p *PluginMessageClientbound) decode(r io.Reader) (err error) {
	return nil
}

func (p PluginMessageServerbound) encode(w io.Writer) (err error) {
	return nil
}

func (p *PluginMessageServerbound) decode(r io.Reader) (err error) {
	return nil
}

func (p Particle) encode(w io.Writer) (err error) {
	return
}

func (p *Particle) decode(r io.Reader) (err error) {
	return
}

// import (
// 	"bytes"
// 	"fmt"
// 	"math/bits"
// )
//
// type ChunkSection struct {
// 	BlockCount   int16
// 	BitsPerBlock uint8
// 	Palette      []int32
// 	DataArray    []uint
// }
//
// func (c *ChunkData) ParseData() ([]ChunkSection, []int32, error) {
// 	chunkCount := bits.OnesCount32(uint32(c.PrimaryBitMask))
// 	chunkSections := make([]ChunkSection, chunkCount)
// 	buffer := bytes.NewBuffer(c.Data)
//
// 	for chunk := 0; chunk < chunkCount; chunk++ {
// 		var err error
// 		section := ChunkSection{}
//
// 		section.BlockCount, err = readShort(buffer)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("unable to read block count: %w", err)
// 		}
//
// 		section.BitsPerBlock, err = readUnsignedByte(buffer)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("unable to read bits per block: %w", err)
// 		}
//
// 		if section.BitsPerBlock < 9 {
// 			length, err := readVarInt(buffer)
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("unable to read palette length: %w", err)
// 			}
//
// 			section.Palette = make([]int32, length)
// 			for i := int32(0); i < length; i++ {
// 				section.Palette[i], err = readVarInt(buffer)
// 				if err != nil {
// 					return nil, nil, fmt.Errorf("unable to read palette: %w", err)
// 				}
// 			}
// 		}
//
// 		dataArrayLength, err := readVarInt(buffer)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("unable to read data array length: %w", err)
// 		}
//
// 		dataArray := make([]byte, dataArrayLength*8)
// 		_, err = buffer.Read(dataArray)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("unable to read data array: %w", err)
// 		}
//
// 		section.DataArray = make([]uint, 16*16*16)
//
// 		for y := 0; y < 16; y++ {
// 			for z := 0; z < 16; z++ {
// 				for x := 0; x < 16; x++ {
// 					blockNumber := (((y * 16) + z) * 16) + x
// 					startLong := (blockNumber * int(section.BitsPerBlock)) / 64
// 					startOffset := uint((blockNumber * int(section.BitsPerBlock)) % 64)
// 					endLong := ((blockNumber+1)*int(section.BitsPerBlock) - 1) / 64
//
// 					var data uint
// 					if startLong == endLong {
// 						data = uint(dataArray[startLong] >> startOffset)
// 					} else {
// 						endOffset := 64 - startOffset
// 						data = uint(dataArray[startLong]>>startOffset | dataArray[endLong]<<endOffset)
// 					}
//
// 					section.DataArray[blockNumber] = data
// 				}
// 			}
// 		}
//
// 		chunkSections[chunk] = section
// 	}
//
// 	var (
// 		biomes []int32
// 		err    error
// 	)
//
// 	if c.FullChunk {
// 		biomes = make([]int32, 256)
// 		for i := 0; i < 256; i++ {
// 			biomes[i], err = readInt(buffer)
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("unable to read biomes: (step %d) %w", i, err)
// 			}
// 		}
// 	}
//
// 	return chunkSections, biomes, nil
// }
//
// func (c *ChunkData) LoadData(sections []ChunkSection, biomes []int32) error {
// 	buffer := bytes.NewBuffer(make([]byte, 0))
//
// 	var err error
// 	for _, section := range sections {
// 		err = writeShort(buffer, section.BlockCount)
// 		if err != nil {
// 			return err
// 		}
//
// 		err = writeUnsignedByte(buffer, section.BitsPerBlock)
// 		if err != nil {
// 			return err
// 		}
//
// 		if section.BitsPerBlock < 9 {
// 			err = writeVarInt(buffer, int32(len(section.Palette)))
// 			if err != nil {
// 				return err
// 			}
//
// 			for i := 0; i < len(section.Palette); i++ {
// 				err = writeVarInt(buffer, section.Palette[i])
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
//
// 		dataArrayLength := int32(section.BitsPerBlock) * 16 * 16 * 16 / 64
// 		dataArray := make([]byte, dataArrayLength*8)
// 		err = writeVarInt(buffer, dataArrayLength)
// 		if err != nil {
// 			return err
// 		}
//
// 		for y := 0; y < 16; y++ {
// 			for z := 0; z < 16; z++ {
// 				for x := 0; x < 16; x++ {
// 					blockNumber := (((y * 16) + z) * 16) + x
// 					startLong := (blockNumber * int(section.BitsPerBlock)) / 64
// 					startOffset := uint((blockNumber * int(section.BitsPerBlock)) % 64)
// 					endLong := ((blockNumber+1)*int(section.BitsPerBlock) - 1) / 64
//
// 					dataArray[startLong] |= byte(section.DataArray[blockNumber] << startOffset)
//
// 					if startLong != endLong {
// 						dataArray[endLong] = byte(section.DataArray[blockNumber] >> (64 - startOffset))
// 					}
// 				}
// 			}
// 		}
//
// 		_, err = buffer.Write(dataArray)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	if c.FullChunk {
// 		if len(biomes) != 256 {
// 			return fmt.Errorf("biomes must be exactly 256 long (got %d)", len(biomes))
// 		}
//
// 		for _, v := range biomes {
// 			err = writeInt(buffer, v)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	c.Data = buffer.Bytes()
// 	c.Size = int32(len(c.Data))
//
// 	return nil
// }
