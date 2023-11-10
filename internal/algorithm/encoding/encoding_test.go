package encoding

import (
	"github.com/dyammarcano/secure_message/internal/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding1kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message1kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message1kChars)
}

func TestEncoding2kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message2kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message2kChars)
}

func TestEncoding3kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message3kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message3kChars)
}

func TestEncoding4kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message4kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message4kChars)
}

func TestEncoding5kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message5kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message5kChars)
}
