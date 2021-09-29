package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenSystemStatusIsUp_ShouldReturnAppropriateMetricValue(t *testing.T) {
	status := Up
	expectedMetricValue := float64(1)

	assert.Equal(t, expectedMetricValue, mapSystemStatusToMetricValue(status))
}

func Test_WhenSystemStatusIsDown_ShouldReturnAppropriateMetricValue(t *testing.T) {
	status := Down
	expectedMetricValue := float64(0)

	assert.Equal(t, expectedMetricValue, mapSystemStatusToMetricValue(status))
}

func Test_WhenSystemStatusIsOutOfService_ShouldReturnAppropriateMetricValue(t *testing.T) {
	status := OutOfService
	expectedMetricValue := float64(-1)

	assert.Equal(t, expectedMetricValue, mapSystemStatusToMetricValue(status))
}

func Test_WhenSystemStatusIsPartial_ShouldReturnAppropriateMetricValue(t *testing.T) {
	status := Partial
	expectedMetricValue := float64(-2)

	assert.Equal(t, expectedMetricValue, mapSystemStatusToMetricValue(status))
}

func Test_WhenSystemStatusIsUnknown_ShouldReturnAppropriateMetricValue(t *testing.T) {
	status := Unknown
	expectedMetricValue := float64(-3)

	assert.Equal(t, expectedMetricValue, mapSystemStatusToMetricValue(status))
}
