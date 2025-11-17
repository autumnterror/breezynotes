package uid

import "github.com/google/uuid"

func New() string {
	return uuid.NewString()
}

func Validate(id string) bool {
	return uuid.Validate(id) == nil
}
