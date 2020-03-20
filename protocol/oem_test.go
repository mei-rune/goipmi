package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOEM(t *testing.T) {
	assert.Equal(t, "Dell Inc", OemDell.String())
	assert.Equal(t, "Hewlett-Packard", OemHP.String())
}
