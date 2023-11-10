package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/dyammarcano/base58"
)

var globalKeys map[int][]byte

func init() {
	globalKeys = Keys()
}

// GenerateKeys generates a N bytes master key
func GenerateKeys(size int) ([]byte, error) {
	masterKey := make([]byte, size)
	if _, err := rand.Read(masterKey); err != nil {
		return nil, err
	}
	return masterKey, nil
}

// convertToKeys its return 256 bits for key and 96 bits GCM nonce
func convertToKeys(inputKey []byte, keys map[int][]byte) ([]byte, []byte, error) {
	if len(keys) == 0 {
		return nil, nil, errors.New("keys map is empty")
	}

	if len(inputKey) < 2 {
		return nil, nil, errors.New("input key too short")
	}

	sl := int(binary.BigEndian.Uint16(inputKey)) % len(keys)
	selectedKey := keys[sl]

	matrixKey := make([]byte, len(selectedKey))
	for i, value := range selectedKey {
		sk := i % len(inputKey)
		matrixKey[i] = value ^ inputKey[sk]
	}

	key := make([]byte, 32) // 256 bits for AES-256
	copy(key, matrixKey[:len(key)])

	nonce := make([]byte, 12) // 96 bits for GCM nonce
	copy(nonce, matrixKey[:len(nonce)])

	return key, nonce, nil
}

// wrapValues wraps the key and cypherText into a single byte array
func wrapValues(cypherText, key []byte) []byte {
	var response bytes.Buffer
	response.Write(key)                // write 32 bytes to store the key
	response.Write(cypherText)         // write the rest of the bytes to store the cypherText
	response.WriteByte(byte(len(key))) // write 1 byte to store the key size

	return response.Bytes()
}

// unwrapValues unwraps the key, nonce and cypherText from a single byte array
func unwrapValues(cypherText []byte) ([]byte, []byte, []byte, error) {
	if len(cypherText) < 2 {
		return nil, nil, nil, errors.New("cypherText too short")
	}

	keyLen := int(cypherText[len(cypherText)-1])
	key, nonce, err := convertToKeys(cypherText[:keyLen], globalKeys)
	if err != nil {
		return nil, nil, nil, err
	}
	return key, nonce, cypherText[keyLen : len(cypherText)-1], nil
}

// encrypt encrypts a message using AES-256-GCM
func encrypt(message, password []byte) ([]byte, error) {
	key, nonce, err := convertToKeys(password, globalKeys)
	if err != nil {
		return nil, err
	}

	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	sealed := gcm.Seal(nil, nonce, message, nil)
	if sealed == nil {
		return nil, errors.New("encryption error")
	}

	return wrapValues(sealed, password), nil
}

// decrypt decrypts a message using AES-256-GCM
func decrypt(cypherText []byte) ([]byte, error) {
	key, nonce, message, err := unwrapValues(cypherText)
	if err != nil {
		return nil, err
	}

	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	decrypted, err := gcm.Open(nil, nonce, message, nil)
	if err != nil {
		return nil, err
	}

	if decrypted == nil {
		return nil, errors.New("decryption error")
	}

	return decrypted, nil
}

// AutoEncryptString decrypts a message using AES-256-GCM
func AutoEncryptString(message string) (string, error) {
	generateKey, err := GenerateKeys(12)
	if err != nil {
		return "", err
	}

	encrypted, err := encrypt([]byte(message), generateKey)
	if err != nil {
		return "", err
	}

	return base58.StdEncoding.Encode(encrypted), err
}

// AutoEncryptBytes decrypts a message using AES-256-GCM
func AutoEncryptBytes(message []byte) ([]byte, error) {
	generateKey, err := GenerateKeys(12)
	if err != nil {
		return nil, err
	}

	encrypted, err := encrypt(message, generateKey)
	if err != nil {
		return nil, err
	}

	return []byte(base58.StdEncoding.Encode(encrypted)), nil
}

// AutoDecryptString decrypts a message using AES-256-GCM
func AutoDecryptString(message string) (string, error) {
	decoded, err := base58.StdEncoding.Decode(message)
	if err != nil {
		return "", err
	}

	decrypted, err := decrypt(decoded)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// AutoDecryptBytes decrypts a message using AES-256-GCM
func AutoDecryptBytes(message []byte) ([]byte, error) {
	decoded, err := base58.StdEncoding.Decode(string(message))
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(decoded)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
