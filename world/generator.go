package world

import (
	"fmt"
)

type Generator interface {
	Chunk(x int, z int) Chunk
}

func NewGenerator(name string, options interface{}) (Generator, error) {
	switch name {
	case "superflat":
		return &Superflat{opts: options.(SuperflatOptions)}, nil
	default:
		return nil, fmt.Errorf("Unknown generator type: %s", name)
	}
}

type SuperflatOptions struct {
}

type Superflat struct {
	opts SuperflatOptions
}

func (s *Superflat) Chunk(x int, z int) Chunk {
	return Chunk{}
}
