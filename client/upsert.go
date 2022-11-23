package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/phoops/ngsi-gold/ldcontext"
	"github.com/phoops/ngsi-gold/model"
	"github.com/pkg/errors"
)

const batchUpsertEndpoint string = "ngsi-ld/v1/entityOperations/upsert"

type EntityWithContext struct {
	LdCtx  *ldcontext.LdContext
	Entity *model.Entity
}

type upsertMode string

const upsertModeReplace upsertMode = "replace"
const upsertModeUpdate upsertMode = "update"

type batchUpsertOptions struct {
	mode upsertMode
}

func newBatchUpsertOptions() *batchUpsertOptions {
	return &batchUpsertOptions{
		mode: upsertModeReplace,
	}
}

type UpsertOptionFunc func(*batchUpsertOptions) error

var UpsertSetUpdateMode UpsertOptionFunc = func(o *batchUpsertOptions) error {
	o.mode = upsertModeUpdate
	return nil
}

var UpsertSetReplaceMode UpsertOptionFunc = func(o *batchUpsertOptions) error {
	o.mode = upsertModeReplace
	return nil
}

func (client *NgsiLdClient) BatchUpsertEntities(ctx context.Context, payload []*EntityWithContext, opts ...UpsertOptionFunc) error {
	batchUpsertURL := strings.Join([]string{client.url, batchUpsertEndpoint}, "/")
	batchRequest := batchRequestBody{}

	for _, x := range payload {
		ldCtx := x.LdCtx
		entity := x.Entity
		// Set default context whenever missing
		if ldCtx == nil {
			ldCtx = &ldcontext.DefaultContext
		}

		// Validate entity to be created before contacting the server
		err := entity.Validate(true)
		if err != nil {
			return errors.Wrap(err, "invalid Entity")
		}
		inner, err := addContext(entity, ldCtx)
		if err != nil {
			return err
		}
		batchRequest = append(batchRequest, &inner)
	}

	upsertRequestBody, err := json.Marshal(&batchRequest)
	if err != nil {
		return err
	}

	req, err := client.newRequest(
		ctx,
		http.MethodPost,
		batchUpsertURL,
		bytes.NewBuffer(upsertRequestBody),
		jsonLdBody,
	)
	if err != nil {
		return err
	}

	requestOptions := newBatchUpsertOptions()
	for _, o := range opts {
		err := o(requestOptions)
		if err != nil {
			return errors.Wrap(ErrInvalidUpsertOptions, err.Error())
		}
	}

	q := req.URL.Query()
	switch requestOptions.mode {
	case upsertModeReplace:
		q.Add("options", string(upsertModeReplace))
	case upsertModeUpdate:
		q.Add("options", string(upsertModeUpdate))
	}

	resp, err := client.c.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't upsert Entities")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	conflict := ProblemDetails{}
	err = json.Unmarshal(bodyBytes, &conflict)

	if err == nil {
		if conflict.ErrType == ngsiLdErrInvalidRequest {
			return ErrNgsiLdInvalidRequest
		}
		if conflict.ErrType == ngsiLdErrBadData {
			return errors.Wrapf(ErrNgsiBadData, "Detail: %s", conflict.Detail)
		}
	}

	multiError := BatchOperationResult{}
	err = json.Unmarshal(bodyBytes, &multiError)
	if err != nil {
		return fmt.Errorf("Unexpected status code: '%d'\nResponse body: %s", resp.StatusCode, string(bodyBytes))
	}
	return errors.Wrapf(ErrNgsiMixedResponse, "Success: %v, Errors: %v", multiError.Success, multiError.Errors)
}
