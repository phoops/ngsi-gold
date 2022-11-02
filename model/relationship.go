package model

import "encoding/json"

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
