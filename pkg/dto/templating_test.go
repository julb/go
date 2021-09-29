package dto

import (
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func Test_WhenMarshallingTemplatingMode_ShouldReturnAppropriateResponse(t *testing.T) {
	result, err := json.MarshalToString(TextTemplatingMode)
	assert.Nil(t, err)
	assert.Equal(t, "\"TEXT\"", result)

	result, err = json.MarshalToString(HtmlTemplatingMode)
	assert.Nil(t, err)
	assert.Equal(t, "\"HTML\"", result)

	result, err = json.MarshalToString(NoneTemplatingMode)
	assert.Nil(t, err)
	assert.Equal(t, "\"NONE\"", result)
}

func Test_WhenUnmarshallingTemplatingMode_ShouldReturnAppropriateResponse(t *testing.T) {
	var result TemplatingMode
	err := json.UnmarshalFromString("\"TEXT\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, TextTemplatingMode, &result)

	err = json.UnmarshalFromString("\"HTML\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, HtmlTemplatingMode, &result)

	err = json.UnmarshalFromString("\"NONE\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, NoneTemplatingMode, &result)
}
