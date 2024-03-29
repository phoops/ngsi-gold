package model

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Properties is a helper type, defines a set of Properties, identified by a string
type Properties map[string]Property

// Relationships is a helper type, defines a set of Relationships, identified by a string
type Relationships map[string]Relationship

// Entity is the top-level abstraction of the user domain and thus it can represent anything.
// Entities have mandatory type and id and a number of Attributes
// https://github.com/FIWARE/context.Orion-LD/blob/develop/doc/manuals-ld/entities-and-attributes.md
type Entity struct {
	ID               string        `json:"id"`                         // ID of the entity used to identify the single entity
	Type             string        `json:"type"`                       // Type of the entity used for categorization
	Properties       Properties    `json:"-"`                          // Values that define the entity
	Relationships    Relationships `json:"-"`                          // Links to other entities
	Location         *GeoProperty  `json:"location,omitempty"`         // Position of the Entity
	ObservationSpace *GeoProperty  `json:"observationSpace,omitempty"` // Area observable by the Entity (e.g. a camera)
	OperationSpace   *GeoProperty  `json:"operationSpace,omitempty"`   // Area operable by the Entity (e.g. a sprinkler)
}

const timeRFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

func (e Entity) MarshalJSON() ([]byte, error) {
	data := map[string]any{}

	data["type"] = e.Type
	data["id"] = e.ID
	if e.Location != nil {
		data["location"] = e.Location
	}
	if e.ObservationSpace != nil {
		data["observationSpace"] = e.ObservationSpace
	}
	if e.OperationSpace != nil {
		data["operationSpace"] = e.OperationSpace
	}

	for k, v := range e.Properties {
		data[k] = v
	}

	for k, v := range e.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}

func (e *Entity) UnmarshalJSON(b []byte) error {
	// Use an alias to avoid recursion into this function
	type readEntity Entity

	d := readEntity{}

	// First pass - extract annotated fields
	if err := json.Unmarshal(b, &d); err != nil {
		return ErrInvalidEntity(err)
	}

	// Check for missing mandatory values
	if d.ID == "" {
		return ErrEntityMissingID
	}
	if d.Type == "" {
		return ErrEntityMissingType
	}

	// Second pass - extract rest of the fields present in the JSON
	var jsonValues map[string]json.RawMessage
	_ = json.Unmarshal(b, &jsonValues)

	typ := reflect.TypeOf(d)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			delete(jsonValues, jsonTag)
		}
	}

	// Third pass - partial decode to discover the type of the attribute
	type Attribute struct {
		Type string `json:"type,omitempty"`
	}

	d.Relationships = Relationships{}
	d.Properties = Properties{}

	for k, v := range jsonValues {
		a := Attribute{}
		err := json.Unmarshal(v, &a)
		if err != nil {
			return ErrInvalidEntity(errors.Wrapf(err, "cannot unmarshal attribute %s", k))
		}

		// Fourth pass - decode according to type
		switch a.Type {
		case "Relationship":
			r := Relationship{}
			err = json.Unmarshal(v, &r)
			if err != nil {
				return ErrInvalidEntity(errors.Wrapf(err, "cannot unmarshal relationship %s", k))
			}
			d.Relationships[k] = r
		case "Property":
			p := Property{}
			err = json.Unmarshal(v, &p)
			if err != nil {
				return ErrInvalidEntity(errors.Wrapf(err, "cannot unmarshal property %s", k))
			}
			d.Properties[k] = p
		}
	}

	// Do not return empty maps
	if len(d.Properties) == 0 {
		d.Properties = nil
	}
	if len(d.Relationships) == 0 {
		d.Relationships = nil
	}

	// Assign fields to pointed structure
	*e = Entity(d)

	return nil
}

func (e *Entity) Validate(strictness bool) ValidationResult {
	if len(e.ID) == 0 {
		return ErrEntityMissingID
	}
	if len(e.Type) == 0 {
		return ErrEntityMissingType
	}

	for _, x := range e.Properties {
		err := x.Validate(strictness)
		if err != nil {
			return err
		}
	}
	for _, x := range e.Relationships {
		err := x.Validate(strictness)
		if err != nil {
			return err
		}
	}
	return nil
}
