package model_test

import (
	"encoding/json"
	"testing"

	"github.com/phoops/ngsild/model"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalRelationship(t *testing.T) {
	type testCase struct {
		name         string
		json         string
		relationship model.Relationship
	}

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
					"r1": {Object: "urn:nested_object_id"},
				},
			},
		},
		{
			name: "nested property",
			json: `{"type":"Relationship","object":"urn:object_id", "p1":{"type":"Property","value":"property value"}}`,
			relationship: model.Relationship{
				Object: "urn:object_id",
				Properties: model.Properties{
					"p1": {Value: "property value"},
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
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// test self-check
			var js interface{}
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
			name:   "malformed nested relationship",
			json:   `{"type":"Relationship","object":"urn:object_id","r1":"definitely a relationship"}`,
			errMsg: "cannot unmarshal attribute r1: json: cannot unmarshal string into Go value of type model.Attribute",
		},
	}

	for _, y := range tests {
		t.Run(y.name, func(t *testing.T) {
			// test self-check
			var js interface{}
			isJSON := json.Unmarshal([]byte(y.json), &js) == nil
			assert.True(t, isJSON)

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
