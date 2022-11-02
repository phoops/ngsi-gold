package model

import "encoding/json"

// Property is an Attribute that holds a value
type Property struct {
	Value         any           `json:"value"`
	Properties    Properties    `json:"-"`
	Relationships Relationships `json:"-"`
}

func (p *Property) Type() string {
	return "Property"
}

func (p Property) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	data["type"] = p.Type()
	data["value"] = p.Value

	for k, v := range p.Properties {
		data[k] = v
	}

	for k, v := range p.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}
