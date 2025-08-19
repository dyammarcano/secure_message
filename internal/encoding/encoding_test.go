package encoding

import (
	"fmt"
	"os"
	"testing"

	"github.com/dyammarcano/secure_message/internal/algorithm/cryptography"
	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	serialized, err := Serialize("hello world")
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, "hello world")
}

func TestCompression(t *testing.T) {
	i := 0
	for {
		i++
		data, _ := cryptography.GenerateKeys(i)

		serialized, err := Serialize(fmt.Sprintf("%x", data))
		assert.Nil(t, err)

		if len(serialized) < len(fmt.Sprintf("%x", data)) {
			_, _ = fmt.Fprintf(os.Stdout, "serialized data len: %d, original data len: %d\n\n", len(serialized), len(fmt.Sprintf("%x", data)))
			break
		}
	}
}
