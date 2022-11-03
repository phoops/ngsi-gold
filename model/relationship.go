package model

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type Relationship struct {
	Object        string        `json:"object"`
	Properties    Properties    `json:"-"`
	Relationships Relationships `json:"-"`
}

func (r *Relationship) Type() string {
	return "Relationship"
}

func (r Relationship) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	data["type"] = r.Type()
	data["object"] = r.Object

	for k, v := range r.Properties {
		data[k] = v
	}

	for k, v := range r.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}

func (r *Relationship) UnmarshalJSON(b []byte) error {

	// Use an alias to avoid recursion into this function
	type readRelationship Relationship

	d := readRelationship{}

	// First pass - extract annotated fields
	if err := json.Unmarshal(b, &d); err != nil {
		return ErrInvalidRelationship(err)
	}

	// Check for missing mandatory value
	if d.Object == "" {
		return ErrRelationshipMissingObject
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

	// Check type or bail out
	rtype, ok := jsonValues["type"]
	if !ok || string(rtype) != `"Relationship"` {
		return ErrRelationshipWrongType
	}
	delete(jsonValues, "type")

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
			return ErrInvalidRelationship(errors.Wrapf(err, "cannot unmarshal attribute %s", k))
		}

		// Fourth pass - decode according to type
		switch a.Type {
		case "Relationship":
			r := Relationship{}
			err = json.Unmarshal(v, &r)
			if err != nil {
				return ErrInvalidRelationship(errors.Wrapf(err, "cannot unmarshal relationship %s", k))
			}
			d.Relationships[k] = r
		case "Property":
			p := Property{}
			err = json.Unmarshal(v, &p)
			if err != nil {
				return ErrInvalidRelationship(errors.Wrapf(err, "cannot unmarshal property %s", k))
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
	*r = Relationship(d)

	return nil
}
