package encoding

import (
	"github.com/dyammarcano/base58"
	"github.com/dyammarcano/secure_message/internal/algorithm/compression"
	"github.com/dyammarcano/secure_message/internal/algorithm/crypto"
)

func Serialize(message string) (string, error) {
	comp, err := compression.CompressData([]byte(message))
	if err != nil {
		return "", err
	}

	enc, err := crypto.AutoEncryptBytes(comp)
	if err != nil {
		return "", err
	}

	return base58.StdEncoding.EncodeToString(enc), nil
}

func SerializeWithPassword(message, password []byte) ([]byte, error) {
	comp, err := compression.CompressData(message)
	if err != nil {
		return nil, err
	}

	enc, err := crypto.EncryptPassword(comp, password)
	if err != nil {
		return nil, err
	}

	return []byte(base58.StdEncoding.EncodeToString(enc)), nil
}

func DeserializeWithPassword(message, password []byte) ([]byte, error) {
	dec, err := base58.StdEncoding.DecodeString(string(message))
	if err != nil {
		return nil, err
	}

	dec, err = crypto.DecryptPassword(dec, password)
	if err != nil {
		return nil, err
	}

	dec, err = compression.DecompressData(dec)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

func Deserialize(message string) (string, error) {
	dec, err := base58.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	dec, err = crypto.AutoDecryptBytes(dec)
	if err != nil {
		return "", err
	}

	dec, err = compression.DecompressData(dec)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
