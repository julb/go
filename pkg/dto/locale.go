package dto

import (
	"strconv"
	"strings"
)

// The locale DTO.
type LocaleDTO struct {
	Language string
	Country  string
}

func (obj *LocaleDTO) String() string {
	if obj.Country != "" {
		return strings.Join([]string{strings.ToLower(obj.Language), strings.ToUpper(obj.Country)}, "-")
	} else {
		return strings.ToLower(obj.Language)
	}
}

func ParseLocale(localeString string) (*LocaleDTO, error) {
	localeParts := strings.Split(localeString, "-")

	return &LocaleDTO{
		Language: strings.ToLower(localeParts[0]),
		Country:  strings.ToUpper(localeParts[1]),
	}, nil
}

func (obj *LocaleDTO) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(obj.String())), nil
}

func (obj *LocaleDTO) UnmarshalJSON(data []byte) error {
	localeString, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	parsedLocale, err := ParseLocale(localeString)
	*obj = *parsedLocale
	return err
}
