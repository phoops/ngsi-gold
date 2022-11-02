package model

import "encoding/json"

// Properties is a helper type, defines a set of Properties, identified by a string
type Properties map[string]Property

// Relationships is a helper type, defines a set of Relationships, identified by a string
type Relationships map[string]Relationship

// Entity is the top-level abstraction of the user domain and thus it can represent anything.
// Entities have mandatory type and id and a number of Attributes
// https://github.com/FIWARE/context.Orion-LD/blob/develop/doc/manuals-ld/entities-and-attributes.md
type Entity struct {
	ID            string        `json:"id"`   // ID of the entity used to identify the single entity
	Type          string        `json:"type"` // Type of the entity used for categorization
	Properties    Properties    `json:"-"`    // Values that define the entity
	Relationships Relationships `json:"-"`    // Links to other entities
}

func (e Entity) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	data["type"] = e.Type
	data["id"] = e.ID

	for k, v := range e.Properties {
		data[k] = v
	}

	for k, v := range e.Relationships {
		data[k] = v
	}

	return json.Marshal(data)
}
