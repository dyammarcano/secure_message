package management

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagementExportKeys(t *testing.T) {
	data, err := ExportKeys()
	assert.Nil(t, err)

	assert.NotNil(t, data)
}
