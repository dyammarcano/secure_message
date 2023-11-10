package cryptography

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

var globalKeys *Keys

func init() {
	globalKeys = NewKeys() // initialize globalKeys
}

type (
	Keys struct {
		length int
	}
)

func NewKeys() *Keys {
	return &Keys{
		length: keysLenght,
	}
}

// ConvertToKeys its return 256 bits for key and 96 bits GCM nonce
func (k *Keys) ConvertToKeys(inputKey []byte) ([]byte, []byte, error) {
	if len(inputKey) < 2 {
		return nil, nil, errors.New("input key too short")
	}

	selectedKey, err := k.generateFromKey(inputKey)
	if err != nil {
		return nil, nil, errors.New("error generating key")
	}

	inputKeySize := len(inputKey)
	matrixKey := make([]byte, k.length)

	for i, value := range selectedKey {
		matrixKey[i] = value ^ inputKey[i%inputKeySize]
	}

	key := make([]byte, 32)         // 256 bits for AES-256
	copy(key, matrixKey[:len(key)]) // get first bytes from matrixKey

	nonce := make([]byte, 12)                            // 96 bits for GCM nonce
	copy(nonce, matrixKey[len(key):len(key)+len(nonce)]) // get last bytes from matrixKey

	return key, nonce, nil
}

// generateFromKey return specific key from map of keys
func (k *Keys) generateFromKey(inputKey []byte) ([]byte, error) {
	index := int(binary.BigEndian.Uint16(inputKey)) % k.length
	key := stringKeys[index]

	if len(key) < 2 {
		return nil, errors.New("AmountKeys is too short")
	}

	data, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	keyData := make([]byte, k.length) // 354 bits matrix key
	copy(keyData, data)               // copy data to keyData

	return keyData, nil
}
