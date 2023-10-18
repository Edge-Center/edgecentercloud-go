package edgecloud

import (
	"fmt"

	"github.com/google/uuid"
)

func isValidUUID(id, name string) error {
	if _, err := uuid.Parse(id); err != nil {
		return NewArgError(name, fmt.Sprintf("should be the correct UUID. current value is: %s", id))
	}

	return nil
}
