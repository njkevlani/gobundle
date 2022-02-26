package stdpkgdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBuiltInFunc(t *testing.T) {
	assert.True(t, IsStdPkg("fmt"))
	assert.True(t, IsStdPkg("math/rand"))
	assert.False(t, IsStdPkg("algo"))
}
