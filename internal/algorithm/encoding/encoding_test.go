package encoding

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoding1kChars(t *testing.T) {
	serialized, err := Serialize("hello world")
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, "hello world")

	_, _ = fmt.Fprintf(os.Stdout, "serialized: %s, deserialized: %s\n", serialized, deserialized)
}
