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
			name: "fully nested properties and relationships",
			entity: model.Entity{
				ID:   "thing:5",
				Type: "thing",
				Properties: model.Properties{
					"p1": {
						Value: 100,
						Relationships: model.Relationships{
							"p1r": {Object: "urn:thing:6"},
						},
						Properties: model.Properties{
							"p1p": {Value: true},
						},
					},
				},
				Relationships: model.Relationships{
					"r1": {
						Object: "urn:thing:7",
						Properties: model.Properties{
							"r1p1": {Value: "hello"},
							"r1p2": {Value: "hello, too"},
						},
						Relationships: model.Relationships{
							"r1r": {Object: "urn:thing:8"},
						},
					},
				},
			},
			json: `{"id":"thing:5","p1":{"p1p":{"type":"Property","value":true},"p1r":{"object":"urn:thing:6","type":"Relationship"},"type":"Property","value":100},"r1":{"object":"urn:thing:7","r1p1":{"type":"Property","value":"hello"},"r1p2":{"type":"Property","value":"hello, too"},"r1r":{"object":"urn:thing:8","type":"Relationship"},"type":"Relationship"},"type":"thing"}`,
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
