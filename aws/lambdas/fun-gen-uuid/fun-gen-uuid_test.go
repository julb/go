package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenCallingLambda_ShouldReturnAppropriateResponse(t *testing.T) {
	output, err := HandleRequest(context.Background())
	assert.NotNil(t, output)
	assert.Nil(t, err)

	assert.NotNil(t, output.Value)
}
