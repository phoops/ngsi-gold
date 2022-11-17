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
	ErrPropertyMissingValue ErrInvalidProperty = errors.New(`Property must have a "value" field`)
)

type ErrInvalidEntity error

var (
	ErrEntityMissingType ErrInvalidEntity = errors.New(`Entity must have a type`)
	ErrEntityMissingID   ErrInvalidEntity = errors.New(`Entity must have an ID`)
)

type ErrInvalidGeoProperty error

var (
	ErrGeoPropertyWrongType    ErrInvalidGeoProperty = errors.New(`GeoProperty must have "GeoProperty" type`)
	ErrGeoPropertyMissingValue ErrInvalidGeoProperty = errors.New(`GeoProperty must have a "value" field`)
	ErrGeoPropertyInvalidValue ErrInvalidGeoProperty = errors.New(`GeoProperty value must be a valid GeoJson geometry except GeometryCollection`)
)
