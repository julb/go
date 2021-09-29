package identifier

import (
	"strings"

	"github.com/google/uuid"
)

func Generate() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(uuid.String(), "-", "")
}
