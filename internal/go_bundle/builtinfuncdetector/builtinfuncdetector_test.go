package builtinfuncdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBuiltInFunc(t *testing.T) {
	assert.True(t, IsBuiltinFunc("append"))
	assert.True(t, IsBuiltinFunc("panic"))
	assert.False(t, IsBuiltinFunc("lorem"))
}
