package identifier

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenGenerate_ShouldReturnIdentifier(t *testing.T) {
	identifier := Generate()

	assert.NotNil(t, identifier)
	rx, _ := regexp.Compile("^[a-z0-9]{32,32}$")
	assert.Regexp(t, rx, identifier)
}
