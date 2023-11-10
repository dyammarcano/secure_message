package cryptography

import (
	"github.com/dyammarcano/secure_message/internal/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptString(t *testing.T) {
	encrypted, err := AutoEncryptString(mock.Message1kChars)
	assert.Nil(t, err)

	decrypted, err := AutoDecryptString(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, mock.Message1kChars)
}

func TestEncryptBytes(t *testing.T) {
	encrypted, err := AutoEncryptBytes([]byte(mock.Message1kChars))
	assert.Nil(t, err)

	decrypted, err := AutoDecryptBytes(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, []byte(mock.Message1kChars))
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
