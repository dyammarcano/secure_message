package cryptography

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var rkeys map[int][]byte

func setup() {
	rkeys = make(map[int][]byte)
	for i := 0; i < 10; i++ {
		k, err := GenerateKeys(4)
		if err != nil {
			panic(fmt.Sprintf("Failed to generate keys: %v", err))
		}
		rkeys[i] = k
	}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestNewKeys(t *testing.T) {
	for i := 0; i < 4; i++ {
		index := i % len(rkeys)
		key, nonce, err := g.convertToKeys(rkeys[index])
		assert.NoError(t, err)

		assert.NotNil(t, key)
		assert.NotNil(t, nonce)

		fmt.Printf("original: %x, matrix: %x, nonce: %x\n", rkeys[index], key, nonce)
	}
}

func BenchmarkNewKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := i % len(rkeys)
		_, _, err := g.convertToKeys(rkeys[index])
		assert.NoError(b, err)
	}
}
