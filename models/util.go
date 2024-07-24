package models

import (
	"github.com/LetsFocus/goLF/errors"
	"github.com/google/uuid"
	"strings"
)

func ValidateUUID(id string, field string) (uuid.UUID, error) {
	if strings.TrimSpace(id) == "" {
		return uuid.UUID{}, errors.MissingParam([]string{field})
	}

	UUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, errors.InvalidParam([]string{field})
	}

	return UUID, nil
}
