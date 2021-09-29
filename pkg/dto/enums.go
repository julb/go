package dto

import (
	"fmt"
	"strconv"
)

type Enum interface {
	GetAlias() string
}

func MarshalEnumToJSON(obj Enum) ([]byte, error) {
	return []byte(strconv.Quote(obj.GetAlias())), nil
}

func UnmarshalJSONToEnum(data []byte, values []interface{}) (Enum, error) {
	key, err := strconv.Unquote(string(data))
	if err != nil {
		return nil, err
	}

	for _, value := range values {
		switch typedValue := value.(type) {
		case Enum:
			if typedValue.GetAlias() == key {
				return typedValue, nil
			}
		default:
			return nil, fmt.Errorf("not a enum type")
		}

	}

	return nil, fmt.Errorf("invalid value for enum: %s", key)
}
