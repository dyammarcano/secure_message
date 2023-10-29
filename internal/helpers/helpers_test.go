package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	originalData := DataWrapper{
		Data: struct {
			Field1 int
			Field2 string
			Field3 []byte
		}{
			Field1: 42,
			Field2: "Hello, World!",
			Field3: []byte{1, 2, 3, 4, 5},
		},
	}

	// Test Serialize
	data := Binarization(originalData)
	assert.NotNil(t, data)

	//// Test Compress
	//compressedData, err := CompressData(data)
	//assert.Nil(t, err)
	//
	//// Test Deserialize
	//var decodedData map[string]any
	//decompressedData, err := DecompressData(compressedData)
	//assert.Nil(t, err)
	//
	//err = Deserialize(data, decodedData)
	//assert.Nil(t, err)
	//
	//assert.Equal(t, originalData, decodedData)
}
