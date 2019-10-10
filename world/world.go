package world

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type Type int

const (
	Overworld Type = iota
	Nether
	End
)

type World struct {
	Name   string
	Type   Type
	chunks map[int]map[int]Chunk
}

func (w *World) LoadChunk(x int, z int) error {
	_, ok := w.chunks[x][z]
	if ok {
		return errors.New("Chunk already loaded")
	}

	fn := fmt.Sprintf("worlds/%s/%d-%d.gob", w.Name, x, z)

	file, err := os.Create(fn)
	if err != nil {
		return err
	}

	chunk := Chunk{}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&chunk)
	if err != nil {
		return err
	}

	w.chunks[x][z] = chunk
	return nil
}

func (w *World) SaveChunk(x int, z int) error {
	chunk := w.chunks[x][z]
	fn := fmt.Sprintf("worlds/%s/%d-%d.gob", w.Name, x, z)

	file, err := os.Create(fn)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(chunk)
	if err != nil {
		return err
	}
	return nil
}

func (w *World) UnloadChunk(x int, z int) error {
	err := w.SaveChunk(x, z)
	if err != nil {
		return err
	}
	delete(w.chunks[x], z)
	return nil
}

type Chunk struct {
	X             int
	Z             int
	Sections      [16]ChunkSection
	Entities      []Entity
	BlockEntities []BlockEntity
}

type ChunkSection struct {
	// Y Index of the Section (0 - 15)
	Y       int
	Palette []Block
}

type Block struct {
}

type Position struct {
	X float64
	Y float64
	Z float64
}

type Entity struct {
	ID       string
	Position Position
	Motion   struct {
		DX float64
		DY float64
		DZ float64
	}
	Rotation struct {
		Yaw   float32
		Pitch float32
	}
}

type BlockEntity struct {
	ID       string
	Position Position
}
