package model

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Property is an Attribute that holds a value
type Property struct {
	Value         any           `json:"value"`
	Properties    Properties    `json:"-"`
	Relationships Relationships `json:"-"`
	ObservedAt    *time.Time    `json:"observedAt,omitempty"`
	UnitCode      *string       `json:"unitCode,omitempty"`
	DatasetID     *string       `json:"datasetId,omitempty"`
}

func (p *Property) Type() string {
	return "Property"
}

func (p Property) MarshalJSON() ([]byte, error) {
	data := map[string]any{}

	data["type"] = p.Type()
	data["value"] = p.Value
	if p.ObservedAt != nil {
		data["observedAt"] = p.ObservedAt
	}
	if p.DatasetID != nil {
		data["datasetId"] = p.DatasetID
	}
	if p.UnitCode != nil {
		data["unitCode"] = p.UnitCode
	}

	for k, v := range p.Properties {
		data[k] = v
	}

	for k, v := range p.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}

func (p *Property) UnmarshalJSON(b []byte) error {

	// Use an alias to avoid recursion into this function
	type readProperty Property

	d := readProperty{}

	// First pass - extract annotated fields
	if err := json.Unmarshal(b, &d); err != nil {
		return ErrInvalidProperty(err)
	}

	// Check for missing mandatory value
	if d.Value == nil {
		return ErrPropertyMissingValue
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
	if !ok || string(rtype) != `"Property"` {
		return ErrPropertyWrongType
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
			return ErrInvalidProperty(errors.Wrapf(err, "cannot unmarshal attribute %s", k))
		}

		// Fourth pass - decode according to type
		switch a.Type {
		case "Relationship":
			r := Relationship{}
			err = json.Unmarshal(v, &r)
			if err != nil {
				return ErrInvalidProperty(errors.Wrapf(err, "cannot unmarshal relationship %s", k))
			}
			d.Relationships[k] = r
		case "Property":
			p := Property{}
			err = json.Unmarshal(v, &p)
			if err != nil {
				return ErrInvalidProperty(errors.Wrapf(err, "cannot unmarshal property %s", k))
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
	*p = Property(d)

	return nil
}

func (p *Property) Validate(strictness bool) ValidationResult {
	if p.Value == nil {
		return ErrPropertyMissingValue
	}

	for _, x := range p.Properties {
		err := x.Validate(strictness)
		if err != nil {
			return err
		}
	}

	for _, x := range p.Relationships {
		err := x.Validate(strictness)
		if err != nil {
			return err
		}
	}
	return nil
}
