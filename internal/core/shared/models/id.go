package models

import (
	"villainrsty-ecommerce-server/internal/core/shared/errors"

	"github.com/google/uuid"
)

// ID adalah unique identifier untuk semua entities
type ID string

// NewID membuat ID baru menggunakan UUID v4
func NewID() ID {
	return ID(uuid.New().String())
}

// String mengkonversi ID ke string
func (id ID) String() string {
	return string(id)
}

func (id ID) IsEmpty() bool {
	return id == ""
}

// Validate memastikan ID valid
func (id ID) Validate() error {
	if id.IsEmpty() {
		return errors.New(errors.ErrValidation, "id cannot be empty")
	}

	return nil
}
