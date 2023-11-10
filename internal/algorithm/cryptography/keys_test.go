package cryptography

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeys(t *testing.T) {
	key, nonce, err := globalKeys.ConvertToKeys([]byte("SuperSecret"))
	assert.Nil(t, err)

	assert.NotNil(t, key)
	assert.NotNil(t, nonce)

	fmt.Printf("key: %x, nonce: %x\n", key, nonce)
}
