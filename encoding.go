package secure_message

import "github.com/dyammarcano/secure_message/internal/encoding"

func Serialize(message string) (string, error) {
	return encoding.Serialize(message)
}

func SerializeWithPassword(message, password []byte) ([]byte, error) {
	return encoding.SerializeWithPassword(message, password)
}

func Deserialize(message string) (string, error) {
	return encoding.Deserialize(message)
}

func DeserializeWithPassword(message, password []byte) ([]byte, error) {
	return encoding.DeserializeWithPassword(message, password)
}
