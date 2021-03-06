// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// FileMetadata file metadata
// swagger:model FileMetadata

type FileMetadata struct {

	// digest
	Digest string `json:"digest,omitempty"`

	// digest algorithm
	DigestAlgorithm string `json:"digestAlgorithm,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// size
	Size float64 `json:"size,omitempty"`

	// space
	Space string `json:"space,omitempty"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}

/* polymorph FileMetadata digest false */

/* polymorph FileMetadata digestAlgorithm false */

/* polymorph FileMetadata id false */

/* polymorph FileMetadata name false */

/* polymorph FileMetadata size false */

/* polymorph FileMetadata space false */

// Validate validates this file metadata
func (m *FileMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *FileMetadata) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FileMetadata) UnmarshalBinary(b []byte) error {
	var res FileMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
