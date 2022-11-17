package model_test

import (
	"encoding/json"
	"testing"

	"github.com/phoops/ngsi-gold/model"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalRelationship(t *testing.T) {
	type testCase struct {
		name         string
		json         string
		relationship model.Relationship
	}

	testDataset := "test"
	tests := []testCase{
		{
			name: "shallow",
			json: `{"type":"Relationship","object":"urn:object_id"}`,
			relationship: model.Relationship{
				Object: "urn:object_id",
			},
		},
		{
			name: "nested relationship",
			json: `{"type":"Relationship","object":"urn:object_id", "r1":{"type":"Relationship","object":"urn:nested_object_id"}}`,
			relationship: model.Relationship{
				Object: "urn:object_id",
				Relationships: model.Relationships{
					"r1": model.Relationship{Object: "urn:nested_object_id"},
				},
			},
		},
		{
			name: "nested property",
			json: `{"type":"Relationship","object":"urn:object_id", "p1":{"type":"Property","value":"property value"}}`,
			relationship: model.Relationship{
				Object: "urn:object_id",
				Properties: model.Properties{
					"p1": model.Property{Value: "property value"},
				},
			},
		},
		{
			name: "ignored extra key",
			json: `{"type":"Relationship","object":"urn:object_id", "e1":{"type":"something", "color":"red"}}`,
			relationship: model.Relationship{
				Object: "urn:object_id",
			},
		},
		{
			name: "considered optional attribute",
			json: `{"type":"Relationship","object":"urn:object_id", "datasetId":"test"}`,
			relationship: model.Relationship{
				Object:    "urn:object_id",
				DatasetID: &testDataset,
			},
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// test self-check
			var js any
			isJSON := json.Unmarshal([]byte(y.json), &js) == nil
			assert.True(t, isJSON)

			// test unmarshaling
			r := model.Relationship{}
			err := json.Unmarshal([]byte(y.json), &r)
			assert.NoError(t, err)
			assert.EqualValues(t, y.relationship, r)
		})
	}
}

func TestUnmarshalRelationshipErrors(t *testing.T) {
	type testCase struct {
		name   string
		json   string
		err    error
		errMsg string
	}

	tests := []testCase{
		{
			name: "missing object",
			json: `{"type":"Relationship"}`,
			err:  model.ErrRelationshipMissingObject,
		},
		{
			name:   "malformed object",
			json:   `{"type":"Relationship","object":[2]}`,
			errMsg: "json: cannot unmarshal array into Go struct field readRelationship.object of type string",
		},
		{
			name: "no type",
			json: `{"object":"entity_url"}`,
			err:  model.ErrRelationshipWrongType,
		},
		{
			name: "wrong type",
			json: `{"type":"kind","object":"entity_url"}`,
			err:  model.ErrRelationshipWrongType,
		},
		{
			name:   "invalid nested attribute",
			json:   `{"type":"Relationship","object":"urn:object_id","a1":"definitely am attribute"}`,
			errMsg: "cannot unmarshal attribute a1: json: cannot unmarshal string into Go value of type model.Attribute",
		},
		{
			name:   "invalid nested relationship",
			json:   `{"type":"Relationship","object":"urn:object_id","r1":{"type":"Relationship","object": 1}}`,
			errMsg: "cannot unmarshal relationship r1: json: cannot unmarshal number into Go struct field readRelationship.object of type string",
		},
		{
			name:   "invalid nested property",
			json:   `{"type":"Relationship","object":"urn:object_id","p1":{"type":"Property"}}`,
			errMsg: "cannot unmarshal property p1",
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// unmarshal
			r := model.Relationship{}
			err := json.Unmarshal([]byte(y.json), &r)
			assert.Error(t, err)

			if y.err != nil {
				// Check the precise error instance
				assert.ErrorIs(t, err, y.err)
			} else {
				// Check the correct error type and resulting message
				_, ok := err.(model.ErrInvalidRelationship)
				assert.True(t, ok)
				assert.ErrorContains(t, err, y.errMsg)
			}
		})
	}
}

func TestUnmarshalProperty(t *testing.T) {
	type testCase struct {
		name     string
		json     string
		property model.Property
	}

	testDataset := "UV index"
	tests := []testCase{
		{
			name: "shallow",
			json: `{"type":"Property","value":"words or phrases"}`,
			property: model.Property{
				Value: "words or phrases",
			},
		},
		{
			name: "nested relationship",
			json: `{"type":"Property","value":"words or phrases", "r1":{"type":"Relationship","object":"urn:nested_object_id"}}`,
			property: model.Property{
				Value: "words or phrases",
				Relationships: model.Relationships{
					"r1": model.Relationship{Object: "urn:nested_object_id"},
				},
			},
		},
		{
			name: "nested property",
			json: `{"type":"Property","value":"words or phrases", "p1":{"type":"Property","value":"property value"}}`,
			property: model.Property{
				Value: "words or phrases",
				Properties: model.Properties{
					"p1": model.Property{Value: "property value"},
				},
			},
		},
		{
			name: "ignored extra key",
			json: `{"type":"Property","value":"words or phrases", "e1":{"type":"something", "color":"red"}}`,
			property: model.Property{
				Value: "words or phrases",
			},
		},
		{
			name: "integer numeric value",
			json: `{"type":"Property","value":7}`,
			property: model.Property{
				Value: 7.0,
			},
		},
		{
			name: "floating numeric value",
			json: `{"type":"Property","value":7.2}`,
			property: model.Property{
				Value: 7.2,
			},
		},
		{
			name: "list multi-value",
			json: `{"type":"Property","value":["one","two"]}`,
			property: model.Property{
				Value: []any{"one", "two"},
			},
		},
		{
			name: "set multi-value",
			json: `{"type":"Property","value":{"one": 1,"two": 2}}`,
			property: model.Property{
				Value: map[string]any{
					"one": float64(1),
					"two": float64(2),
				},
			},
		},
		{
			name: "set multi-value",
			json: `{"type":"Property","value":9.4,"datasetId":"UV index"}`,
			property: model.Property{
				Value:     float64(9.4),
				DatasetID: &testDataset,
			},
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// test self-check
			var js any
			isJSON := json.Unmarshal([]byte(y.json), &js) == nil
			assert.True(t, isJSON)

			// test unmarshaling
			p := model.Property{}
			err := json.Unmarshal([]byte(y.json), &p)
			assert.NoError(t, err)
			assert.EqualValues(t, y.property, p)
		})
	}
}

func TestUnmarshalPropertyErrors(t *testing.T) {
	type testCase struct {
		name   string
		json   string
		err    error
		errMsg string
	}

	tests := []testCase{
		// TODO: need fields that are annotated and are not the empty interface
		// {
		// 	name:   "malformed json",
		// 	json:   `{"type":"Property","value":abc}`,
		// 	errMsg: "cannot unmarshal property",
		// },
		{
			name: "missing value",
			json: `{"type":"Property"}`,
			err:  model.ErrPropertyMissingValue,
		},
		{
			name: "no type",
			json: `{"value":"35°C"}`,
			err:  model.ErrPropertyWrongType,
		},
		{
			name: "wrong type",
			json: `{"type":"kind","value":"35°C"}`,
			err:  model.ErrPropertyWrongType,
		},
		{
			name:   "invalid nested attribute",
			json:   `{"type":"Property","value":"35°C","a1":"definitely am attribute"}`,
			errMsg: "cannot unmarshal attribute a1: json: cannot unmarshal string into Go value of type model.Attribute",
		},
		{
			name:   "invalid nested relationship",
			json:   `{"type":"Property","value":"35°C","r1":{"type":"Relationship","object": 1}}`,
			errMsg: "cannot unmarshal relationship r1",
		},
		{
			name:   "invalid nested property",
			json:   `{"type":"Property","value":"35°C","p1":{"type":"Property"}}`,
			errMsg: "cannot unmarshal property p1",
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// unmarshal
			p := model.Property{}
			err := json.Unmarshal([]byte(y.json), &p)
			assert.Error(t, err)

			if y.err != nil {
				// Check the precise error instance
				assert.ErrorIs(t, err, y.err)
			} else {
				// Check the correct error type and resulting message
				_, ok := err.(model.ErrInvalidProperty)
				assert.True(t, ok)
				assert.ErrorContains(t, err, y.errMsg)
			}
		})
	}
}

func TestUnmarshalEntity(t *testing.T) {
	type testCase struct {
		name   string
		json   string
		entity model.Entity
	}

	tests := []testCase{
		{
			name: "shallow",
			json: `{"type":"Room", "id":"urn:room:1"}`,
			entity: model.Entity{
				ID:   "urn:room:1",
				Type: "Room",
			},
		},
		{
			name: "nested relationship",
			json: `{"type":"Room", "id":"urn:room:1", "r1":{"type":"Relationship","object":"urn:nested_object_id"}}`,
			entity: model.Entity{
				ID:   "urn:room:1",
				Type: "Room",
				Relationships: model.Relationships{
					"r1": model.Relationship{Object: "urn:nested_object_id"},
				},
			},
		},
		{
			name: "nested property",
			json: `{"type":"Room", "id":"urn:room:1", "p1":{"type":"Property","value":"property value"}}`,
			entity: model.Entity{
				ID:   "urn:room:1",
				Type: "Room",
				Properties: model.Properties{
					"p1": model.Property{Value: "property value"},
				},
			},
		},
		{
			name: "ignored extra key",
			json: `{"type":"Room", "id":"urn:room:1", "e1":{"type":"something", "color":"red"}}`,
			entity: model.Entity{
				ID:   "urn:room:1",
				Type: "Room",
			},
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// test self-check
			var js any
			isJSON := json.Unmarshal([]byte(y.json), &js) == nil
			assert.True(t, isJSON)

			// test unmarshaling
			e := model.Entity{}
			err := json.Unmarshal([]byte(y.json), &e)
			assert.NoError(t, err)
			assert.EqualValues(t, y.entity, e)
		})
	}
}

func TestUnmarshalEntityErrors(t *testing.T) {
	type testCase struct {
		name   string
		json   string
		err    error
		errMsg string
	}

	tests := []testCase{
		{
			name:   "malformed type",
			json:   `{"type":[],"id":"book:1"}`,
			errMsg: "json: cannot unmarshal array into Go struct field readEntity.type of type string",
		},
		{
			name:   "malformed id",
			json:   `{"type":"Book","id":0}`,
			errMsg: "json: cannot unmarshal number into Go struct field readEntity.id of type string",
		},
		{
			name: "missing id",
			json: `{"type":"Property"}`,
			err:  model.ErrEntityMissingID,
		},
		{
			name: "no type",
			json: `{"id":"book:22"}`,
			err:  model.ErrEntityMissingType,
		},
		{
			name:   "invalid nested attribute",
			json:   `{"type":"Room", "id":"urn:room:1","a1":"definitely am attribute"}`,
			errMsg: "cannot unmarshal attribute a1: json: cannot unmarshal string into Go value of type model.Attribute",
		},
		{
			name:   "invalid nested relationship",
			json:   `{"type":"Room", "id":"urn:room:1","r1":{"type":"Relationship","object": 1}}`,
			errMsg: "cannot unmarshal relationship r1",
		},
		{
			name:   "invalid nested property",
			json:   `{"type":"Room", "id":"urn:room:1","p1":{"type":"Property"}}`,
			errMsg: "cannot unmarshal property p1",
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// unmarshal
			p := model.Entity{}
			err := json.Unmarshal([]byte(y.json), &p)
			assert.Error(t, err)

			if y.err != nil {
				// Check the precise error instance
				assert.ErrorIs(t, err, y.err)
			} else {
				// Check the correct error type and resulting message
				_, ok := err.(model.ErrInvalidEntity)
				assert.True(t, ok)
				assert.ErrorContains(t, err, y.errMsg)
			}
		})
	}
}
