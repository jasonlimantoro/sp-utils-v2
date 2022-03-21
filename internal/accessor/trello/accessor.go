package trello

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/config"
)

var ErrHTTPStatusNon2xx = errors.New("err_http_status_non_2xx")

type Accessor interface {
	CreateCard(ctx context.Context, listID string, name string, description string) (*Card, error)
	CreateList(ctx context.Context, boardID string, name string, pos int) (*List, error)
	GetList(ctx context.Context, boardID string) ([]*List, error)
}

type accessor struct {
	httpClient *http.Client
}

func NewAccessor(httpClient *http.Client) *accessor {
	return &accessor{httpClient: httpClient}
}

func (a accessor) CreateCard(ctx context.Context, listID string, name string, description string) (*Card, error) {
	req := &CreateCardRequest{
		Name: name,
		Desc: description,
	}
	res := &Card{}
	if err := a.postJSON(ctx, fmt.Sprintf(RouteCreateCardOnList, listID), req, res); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) CreateList(ctx context.Context, boardID string, name string, pos int) (*List, error) {
	//TODO implement me
	panic("implement me")
}

func (a accessor) GetList(ctx context.Context, boardID string) ([]*List, error) {
	res := []*List{}

	if err := a.getJSON(ctx, fmt.Sprintf(RouteGetListOnBoard, boardID), &res); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) postJSON(ctx context.Context, path string, req, res interface{}) error {
	reqBuf := &bytes.Buffer{}
	err := json.NewEncoder(reqBuf).Encode(req)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	fullURL := fmt.Sprintf("https://%s/%s", TrelloHost, path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(reqBuf.Bytes()))
	if err != nil {
		return errlib.WrapFunc(err)
	}

	httpReq.Header.Add("Authorization", fmt.Sprintf(
		`OAuth oauth_consumer_key="%s", oauth_token="%s"`,
		config.GetTrelloAPIKey(),
		config.GetTrelloAPIToken(),
	))
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	resBuf := &bytes.Buffer{}
	_, err = io.Copy(resBuf, resp.Body)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	if resp.StatusCode >= 400 {
		return errlib.WrapFunc(errlib.WithFields(ErrHTTPStatusNon2xx, errlib.Fields{
			"endpoint": fullURL,
			"request":  reqBuf.String(),
			"status":   resp.StatusCode,
			"response": resBuf.String(),
		}))
	}

	err = json.Unmarshal(resBuf.Bytes(), res)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func (a accessor) getJSON(ctx context.Context, path string, res interface{}) error {
	fullURL := fmt.Sprintf("https://%s/%s", TrelloHost, path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	httpReq.Header.Add("Authorization", fmt.Sprintf(
		`OAuth oauth_consumer_key="%s", oauth_token="%s"`,
		config.GetTrelloAPIKey(),
		config.GetTrelloAPIToken(),
	))

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	resBuf := &bytes.Buffer{}
	_, err = io.Copy(resBuf, resp.Body)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	if resp.StatusCode >= 400 {
		return errlib.WrapFunc(errlib.WithFields(ErrHTTPStatusNon2xx, errlib.Fields{
			"response": resBuf.String(),
			"status":   resp.StatusCode,
			"endpoint": fullURL,
		}))
	}

	err = json.Unmarshal(resBuf.Bytes(), res)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}
