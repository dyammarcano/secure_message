package cryptography

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptString(t *testing.T) {
	encrypted, err := AutoEncryptString("hello world")
	assert.Nil(t, err)

	decrypted, err := AutoDecryptString(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, "hello world")
}

func TestEncryptBytes(t *testing.T) {
	encrypted, err := AutoEncryptBytes([]byte("hello world"))
	assert.Nil(t, err)

	decrypted, err := AutoDecryptBytes(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, []byte("hello world"))
}

func TestEncrypt(t *testing.T) {
	hello := []byte("hello")

	enc, err := encrypt(hello, []byte("supersecret"))
	assert.Nil(t, err)

	dec, err := decrypt(enc)
	assert.Nil(t, err)

	assert.Equal(t, hello, dec)
}

func TestEncryptRandomKeys(t *testing.T) {
	hello := []byte("hello")
	value := []int{16, 24, 32, 64}

	for i := range value {
		key, err := GenerateKeys(value[i])
		assert.Nil(t, err)

		enc, err := encrypt(hello, key)
		assert.Nil(t, err)

		dec, err := decrypt(enc)
		assert.Nil(t, err)

		assert.Equal(t, hello, dec)
	}
}
