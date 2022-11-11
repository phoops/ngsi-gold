package client_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phoops/ngsild/client"
	"github.com/phoops/ngsild/ldcontext"
	"github.com/phoops/ngsild/model"

	"github.com/stretchr/testify/assert"
)

func TestCreateSuccess(t *testing.T) {
	testEntity := model.Entity{
		ID:   "entity:1",
		Type: "thing",
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Assertion on the received request
				assert.Equal(t, "application/ld+json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				// client respects global headers and per-request
				assert.Equal(t, "x-custom-value", r.Header.Get("X-Custom-Header"))

				// client must add context
				b, err := ioutil.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(b), `"@context":`)

				// Basic empty response
				w.Header().Set("location", fmt.Sprintf("/ngsi-ld/v1/entities/%s", testEntity.ID))
				w.WriteHeader(http.StatusCreated)
			}))
	defer ts.Close()

	cli, err := client.New(
		client.SetURL(ts.URL),
		client.SetGlobalHeader("X-Custom-Header", "x-custom-value"),
	)

	assert.NoError(t, err)
	err = cli.CreateEntity(
		context.Background(),
		nil,
		&testEntity,
	)
	assert.NoError(t, err)
}

func TestCreateGenericFailure(t *testing.T) {
	testEntity := model.Entity{
		ID:   "entity:1",
		Type: "thing",
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Assertion on the received request
				assert.Equal(t, "application/ld+json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				// Basic empty response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				errInternal := []byte("Internal Error")
				_, err := w.Write(errInternal)
				assert.NoError(t, err)

			}))
	defer ts.Close()

	cli, err := client.New(
		client.SetURL(ts.URL),
	)

	assert.NoError(t, err)
	err = cli.CreateEntity(
		context.Background(),
		&ldcontext.DefaultContext,
		&testEntity,
	)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, client.ErrNgsiLdEntityExists)
}

func TestCreateDuplicate(t *testing.T) {
	testEntity := model.Entity{
		ID:   "entity:1",
		Type: "thing",
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Assertion on the received request
				assert.Equal(t, "application/ld+json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				// Basic empty response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				errExisting := []byte(`
        {
          "type": "https://uri.etsi.org/ngsi-ld/errors/AlreadyExists",
          "title": "Already exists.",
          "detail": "entity:1"
        }
        `)
				_, err := w.Write(errExisting)
				assert.NoError(t, err)

			}))
	defer ts.Close()

	cli, err := client.New(
		client.SetURL(ts.URL),
	)

	assert.NoError(t, err)
	err = cli.CreateEntity(
		context.Background(),
		&ldcontext.DefaultContext,
		&testEntity,
	)
	assert.Error(t, err)
	assert.ErrorIs(t, err, client.ErrNgsiLdEntityExists)
}
