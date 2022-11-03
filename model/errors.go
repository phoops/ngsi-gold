package model

import "github.com/pkg/errors"

type ErrInvalidRelationship error

var (
	ErrRelationshipWrongType     ErrInvalidRelationship = errors.New(`relationships must have "Relationship" type`)
	ErrRelationshipMissingObject ErrInvalidRelationship = errors.New(`relationships must have an "object" field`)
)

type ErrInvalidProperty error

var (
	ErrPropertyWrongType    ErrInvalidProperty = errors.New(`Property must have "Property" type`)
	ErrPropertyMissingValue ErrInvalidProperty = errors.New(`Property must have an "value" field`)
)
