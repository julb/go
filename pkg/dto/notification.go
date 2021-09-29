package dto

type NotificationDispatchType struct {
	Alias          string
	TemplatingMode TemplatingMode
	HasSubject     bool
}

func (obj *NotificationDispatchType) GetAlias() string {
	return obj.Alias
}

var (
	MailNotificationDispatchType = &NotificationDispatchType{
		Alias:          "MAIL",
		TemplatingMode: *HtmlTemplatingMode,
		HasSubject:     true,
	}
	SmsNotificationDispatchType = &NotificationDispatchType{
		Alias:          "SMS",
		TemplatingMode: *TextTemplatingMode,
		HasSubject:     false,
	}
	GoogleChatNotificationDispatchType = &NotificationDispatchType{
		Alias:          "GOOGLE_CHAT",
		TemplatingMode: *TextTemplatingMode,
		HasSubject:     false,
	}
	WebNotificationDispatchType = &NotificationDispatchType{
		Alias:          "WEB",
		TemplatingMode: *NoneTemplatingMode,
		HasSubject:     false,
	}
	PushNotificationDispatchType = &NotificationDispatchType{
		Alias:          "PUSH",
		TemplatingMode: *TextTemplatingMode,
		HasSubject:     false,
	}
	NotificationDispatchTypes = []interface{}{
		MailNotificationDispatchType,
		SmsNotificationDispatchType,
		GoogleChatNotificationDispatchType,
		WebNotificationDispatchType,
		PushNotificationDispatchType,
	}
)

func (obj *NotificationDispatchType) MarshalJSON() ([]byte, error) {
	return MarshalEnumToJSON(obj)
}

func (obj *NotificationDispatchType) UnmarshalJSON(data []byte) error {
	value, err := UnmarshalJSONToEnum(data, NotificationDispatchTypes)
	if err != nil {
		return err
	}
	*obj = *value.(*NotificationDispatchType)
	return nil
}

// DTO to generate notification content
type GenerateNotificationContentDTO struct {
	Trademark  string                    `yaml:"tm" json:"tm"`
	Locale     *LocaleDTO                `yaml:"locale" json:"locale"`
	Type       *NotificationDispatchType `yaml:"type" json:"type"`
	Name       string                    `yaml:"name" json:"name"`
	Parameters map[string]interface{}    `yaml:"parameters" json:"parameters"`
}
