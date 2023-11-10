package secure_message

import "github.com/dyammarcano/secure_message/internal/encoding"

func Serialize(message string) (string, error) {
	return encoding.Serialize(message)
}

func Deserialize(message string) (string, error) {
	return encoding.Deserialize(message)
}
