package protocol

// import (
// 	"encoding/gob"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"testing"
// )
//
// func TestChunkData(t *testing.T) {
// 	file, err := os.Open("testdata/ChunkData.gob")
// 	if err != nil {
// 		t.Error(fmt.Errorf("Unable to open example: %w", err))
// 	}
//
// 	var original ChunkData
// 	err = gob.NewDecoder(file).Decode(&original)
// 	if err != nil {
// 		t.Error(fmt.Errorf("Unable to decode example: %w", err))
// 	}
//
// 	sections, biomes, err := original.ParseData()
// 	if err != nil {
// 		t.Error(fmt.Errorf("Unable to parse data: %w", err))
// 	}
//
// 	test := original
// 	test.Data = nil
// 	test.Size = 0
// 	err = test.LoadData(sections, biomes)
// 	if err != nil {
// 		t.Error(fmt.Errorf("Unabel to load data: %w", err))
// 	}
//
// 	if !reflect.DeepEqual(original, test) {
// 		sections_, biomes_, err := test.ParseData()
// 		if err != nil {
// 			t.Error(fmt.Errorf("Unable to parse data a second time: %w", err))
// 		}
//
// 		if !reflect.DeepEqual(sections, sections_) {
// 			t.Fatalf("Sections Differ (old, new): %+v // %+v", sections, sections_)
// 		}
//
// 		if !reflect.DeepEqual(biomes, biomes_) {
// 			t.Fatalf("biomes Differ (old, new): %+v // %+v", biomes, biomes_)
// 		}
//
// 		t.Logf("(old, new): %+v, %+v", original, test)
// 		t.Fatal("Parsing and loading modified Chunk data")
// 	}
// }
