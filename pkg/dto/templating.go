package dto

type TemplatingMode struct {
	Alias         string
	MimeType      string
	FileExtension string
}

func (obj *TemplatingMode) GetAlias() string {
	return obj.Alias
}

var (
	TextTemplatingMode = &TemplatingMode{
		Alias:         "TEXT",
		MimeType:      "text/plain",
		FileExtension: "txt",
	}
	HtmlTemplatingMode = &TemplatingMode{
		Alias:         "HTML",
		MimeType:      "text/html",
		FileExtension: "html",
	}
	NoneTemplatingMode = &TemplatingMode{
		Alias:         "NONE",
		MimeType:      "application/octet-stream",
		FileExtension: "none",
	}

	TemplatingModes = []interface{}{TextTemplatingMode, HtmlTemplatingMode, NoneTemplatingMode}
)

func (obj *TemplatingMode) MarshalJSON() ([]byte, error) {
	return MarshalEnumToJSON(obj)
}

func (obj *TemplatingMode) UnmarshalJSON(data []byte) error {
	value, err := UnmarshalJSONToEnum(data, TemplatingModes)
	if err != nil {
		return err
	}
	*obj = *value.(*TemplatingMode)
	return nil
}
