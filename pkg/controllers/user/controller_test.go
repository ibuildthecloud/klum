package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, 25, sanitizedVersion("25"))
	assert.Equal(t, 25, sanitizedVersion("25+"))
	assert.Equal(t, 25, sanitizedVersion("+25"))
}
