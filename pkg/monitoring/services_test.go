package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenGetInitialHealthStatus_ShouldBeUp(t *testing.T) {
	setSystemStatus(Up)
	assert.Equal(t, Up, GetSystemStatus())
}

func Test_WhenSettingSystemStatus_ShouldReturnNewStatus(t *testing.T) {
	setSystemStatus(Down)
	assert.Equal(t, Down, GetSystemStatus())
}
