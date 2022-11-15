package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/phoops/ngsild/ldcontext"
	"github.com/phoops/ngsild/model"
	"github.com/pkg/errors"
)

const createEntityEndpoint string = "ngsi-ld/v1/entities"

func (client *NgsiLdClient) CreateEntity(ctx context.Context, ldCtx *ldcontext.LdContext, entity *model.Entity) error {
	// Set default context whenever missing
	if ldCtx == nil {
		ldCtx = &ldcontext.DefaultContext
	}

	// Validate entity to be created before contacting the server
	err := entity.Validate(true)
	if err != nil {
		return errors.Wrap(err, "invalid Entity")
	}

	createURL := strings.Join([]string{client.url, createEntityEndpoint}, "/")
	createRequest, err := addContext(entity, ldCtx)
	if err != nil {
		return err
	}
	createRequestBody, err := json.Marshal(&createRequest)
	if err != nil {
		return err
	}

	req, err := client.newRequest(
		ctx,
		http.MethodPost,
		createURL,
		bytes.NewBuffer(createRequestBody),
		jsonLdBody,
	)
	if err != nil {
		return err
	}

	resp, err := client.c.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't create Entity")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	conflict := requestError{}
	err = json.Unmarshal(bodyBytes, &conflict)
	if err != nil || conflict.ErrType == ngsiLdErrInvalidRequest {
		return errors.Wrapf(ErrNgsiLdInvalidRequest, "ID: %s", entity.ID)
	}
	if err != nil || conflict.ErrType == ngsiLdErrBadData {
		return errors.Wrapf(ErrNgsiLdInvalidID, "ID: %s", entity.ID)
	}
	if err != nil || conflict.ErrType == ngsiLdErrAlreadyExist {
		return errors.Wrapf(ErrNgsiLdEntityExists, "ID: %s", entity.ID)
	}

	return fmt.Errorf("Unexpected status code: '%d'\nResponse body: %s", resp.StatusCode, string(bodyBytes))

}
