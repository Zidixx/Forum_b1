package utils

import "github.com/gofrs/uuid"

func GenerateUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
