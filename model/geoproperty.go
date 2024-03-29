package model

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/philiphil/geojson"
	"github.com/pkg/errors"
)

// Property is an Attribute that holds a value
type GeoProperty struct {
	Value         *geojson.Geometry `json:"value"`
	Properties    Properties        `json:"-"`
	Relationships Relationships     `json:"-"`
	ObservedAt    *time.Time        `json:"observedAt,omitempty"`
	DatasetID     *string           `json:"datasetId,omitempty"`
}

func (p *GeoProperty) Type() string {
	return "GeoProperty"
}

func (p GeoProperty) MarshalJSON() ([]byte, error) {
	data := map[string]any{}

	data["type"] = p.Type()
	data["value"] = p.Value
	if p.ObservedAt != nil {
		data["observedAt"] = p.ObservedAt.UTC().Format(timeRFC3339Micro)
	}
	if p.DatasetID != nil {
		data["datasetId"] = p.DatasetID
	}

	for k, v := range p.Properties {
		data[k] = v
	}

	for k, v := range p.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}

func (p *GeoProperty) UnmarshalJSON(b []byte) error {

	// Use an alias to avoid recursion into this function
	type readGeoProperty GeoProperty

	d := readGeoProperty{}

	// First pass - extract annotated fields
	if err := json.Unmarshal(b, &d); err != nil {
		return ErrInvalidGeoProperty(err)
	}

	// Check for missing or malformed mandatory value
	if d.Value == nil {
		return ErrGeoPropertyMissingValue
	}
	if !validGeoPropertyValue(d.Value) {
		return ErrGeoPropertyInvalidValue
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
	if !ok || string(rtype) != `"GeoProperty"` {
		return ErrGeoPropertyWrongType
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
			return ErrInvalidGeoProperty(errors.Wrapf(err, "cannot unmarshal attribute %s", k))
		}

		// Fourth pass - decode according to type
		switch a.Type {
		case "Relationship":
			r := Relationship{}
			err = json.Unmarshal(v, &r)
			if err != nil {
				return ErrInvalidGeoProperty(errors.Wrapf(err, "cannot unmarshal relationship %s", k))
			}
			d.Relationships[k] = r
		case "Property":
			p := Property{}
			err = json.Unmarshal(v, &p)
			if err != nil {
				return ErrInvalidGeoProperty(errors.Wrapf(err, "cannot unmarshal property %s", k))
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
	*p = GeoProperty(d)

	return nil
}

func (p *GeoProperty) Validate(strictness bool) ValidationResult {
	if p.Value == nil {
		return ErrGeoPropertyMissingValue
	}
	if !validGeoPropertyValue(p.Value) {
		return ErrGeoPropertyInvalidValue
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

func validGeoPropertyValue(g *geojson.Geometry) bool {
	return g.Type != geojson.GeometryCollection
}
