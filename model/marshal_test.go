package model_test

import (
	"encoding/json"
	"testing"

	"github.com/phoops/ngsild/model"
	"github.com/stretchr/testify/assert"
)

func TestMarshalRelationship(t *testing.T) {
	r := model.Relationship{
		Object: "sensor:1",
	}
	ej := `{"object":"sensor:1","type":"Relationship"}`
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.EqualValues(t, ej, string(j))
}

func TestMarshalProperty(t *testing.T) {
	// Values can be anything
	r := model.Property{
		Value: struct {
			Field        float32
			privateField string
			MappedField  []int `json:"renamed"`
		}{
			Field:        75.3,
			privateField: "yeah",
			MappedField:  []int{3, 2, 1},
		},
	}

	ej := `{"type":"Property","value":{"Field":75.3,"renamed":[3,2,1]}}`
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.EqualValues(t, ej, string(j))
}

func TestMarshalEntity(t *testing.T) {
	type entityTest struct {
		name   string
		entity model.Entity
		json   string
	}

	tests := []entityTest{
		{
			name: "only ID and type",
			entity: model.Entity{
				ID:   "entity:1",
				Type: "thing",
			},
			json: `{"id":"entity:1","type":"thing"}`,
		},
		{
			name: "with property",
			entity: model.Entity{
				ID:   "room:2",
				Type: "Room",
				Properties: model.Properties{
					"key": {Value: true},
				},
			},
			json: `{"id":"room:2","key":{"type":"Property","value":true},"type":"Room"}`,
		},
		{
			name: "with relationship",
			entity: model.Entity{
				ID:   "bulb:3",
				Type: "Light",
				Relationships: model.Relationships{
					"wall": {Object: "wall:north-east"},
				},
			},
			json: `{"id":"bulb:3","type":"Light","wall":{"object":"wall:north-east","type":"Relationship"}}`,
		},
		{
			name: "property with nested relationship",
			entity: model.Entity{
				ID:   "cabin:4",
				Type: "thing",
				Properties: model.Properties{
					"light": {
						Value: 100,
						Relationships: model.Relationships{
							"wall": {Object: "wall:right"},
						},
					},
				},
			},
			json: `{"id":"cabin:4","light":{"type":"Property","value":100,"wall":{"object":"wall:right","type":"Relationship"}},"type":"thing"}`,
		},
		{
			name: "property and relationship both with nested property and relationship",
			entity: model.Entity{
				ID:   "cabin:4",
				Type: "thing",
				Properties: model.Properties{
					"light": {
						Value: 100,
						Relationships: model.Relationships{
							"wall": {Object: "wall:right"},
						},
						Properties: model.Properties{
							"color": {Value: "white"},
						},
					},
				},
				Relationships: model.Relationships{
					"town": {
						Object: "town:rome",
						Relationships: model.Relationships{
							"neighborhood": {Object: "town:rome:neighborhood:eur"},
						},
						Properties: model.Properties{
							"transient": {Value: false},
						},
					},
				},
			},
			json: `{"id":"cabin:4","light":{"color":{"type":"Property","value":"white"},"type":"Property","value":100,"wall":{"object":"wall:right","type":"Relationship"}},"town":{"neighborhood":{"object":"town:rome:neighborhood:eur","type":"Relationship"},"object":"town:rome","transient":{"type":"Property","value":false},"type":"Relationship"},"type":"thing"}`,
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			j, err := json.Marshal(y.entity)
			assert.NoError(t, err)
			assert.EqualValues(t, y.json, string(j))
		})
	}

}
