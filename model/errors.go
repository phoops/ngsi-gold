package model

import "github.com/pkg/errors"

type ErrInvalidRelationship error

var (
	ErrRelationshipWrongType     ErrInvalidRelationship = errors.New(`relationships must have "Relationship" type`)
	ErrRelationshipMissingObject ErrInvalidRelationship = errors.New(`relationships must have an "object" field`)
)
