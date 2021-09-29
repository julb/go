package dto

import (
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func Test_WhenMarshallingNotificationDispatchType_ShouldReturnAppropriateResponse(t *testing.T) {
	result, err := json.MarshalToString(MailNotificationDispatchType)
	assert.Nil(t, err)
	assert.Equal(t, "\"MAIL\"", result)

	result, err = json.MarshalToString(SmsNotificationDispatchType)
	assert.Nil(t, err)
	assert.Equal(t, "\"SMS\"", result)

	result, err = json.MarshalToString(GoogleChatNotificationDispatchType)
	assert.Nil(t, err)
	assert.Equal(t, "\"GOOGLE_CHAT\"", result)

	result, err = json.MarshalToString(WebNotificationDispatchType)
	assert.Nil(t, err)
	assert.Equal(t, "\"WEB\"", result)

	result, err = json.MarshalToString(PushNotificationDispatchType)
	assert.Nil(t, err)
	assert.Equal(t, "\"PUSH\"", result)
}

func Test_WhenUnmarshallingNotificationDispatchType_ShouldReturnAppropriateResponse(t *testing.T) {
	var result NotificationDispatchType
	err := json.UnmarshalFromString("\"MAIL\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, MailNotificationDispatchType, &result)

	err = json.UnmarshalFromString("\"SMS\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, SmsNotificationDispatchType, &result)

	err = json.UnmarshalFromString("\"GOOGLE_CHAT\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, GoogleChatNotificationDispatchType, &result)

	err = json.UnmarshalFromString("\"WEB\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, WebNotificationDispatchType, &result)

	err = json.UnmarshalFromString("\"PUSH\"", &result)
	assert.Nil(t, err)
	assert.Equal(t, PushNotificationDispatchType, &result)
}
