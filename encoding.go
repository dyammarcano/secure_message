package secure_message

import (
	"github.com/dyammarcano/secure_message/internal/encoding"
	"github.com/dyammarcano/secure_message/internal/version"
)

func init() {
	version.AddFeature("secure_message")
}

func Serialize(message string) (string, error) {
	return encoding.Serialize(message)
}

func Deserialize(message string) (string, error) {
	return encoding.Deserialize(message)
}
