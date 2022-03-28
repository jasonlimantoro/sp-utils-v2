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
	"github.com/google/go-querystring/query"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/config"
)

var ErrHTTPStatusNon2xx = errors.New("err_http_status_non_2xx")

type Accessor interface {
	CreateCard(ctx context.Context, listID string, name string, description string) (*Card, error)
	CreateList(ctx context.Context, boardID string, name string, pos float64) (*List, error)
	GetList(ctx context.Context, boardID string) ([]*List, error)
	GetCards(ctx context.Context, listID string) ([]*Card, error)
	GetCardActions(ctx context.Context, request *GetCardActionsRequest) ([]*CardAction, error)
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

func (a accessor) CreateList(ctx context.Context, boardID string, name string, pos float64) (*List, error) {
	res := &List{}
	req := &CreateListRequest{
		Name: name,
		Pos:  fmt.Sprintf("%.2f", pos),
	}
	q, _ := query.Values(req)

	if err := a.postJSON(ctx, fmt.Sprintf(RouteCreateListOnBoard, boardID, q.Encode()), req, res); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) GetList(ctx context.Context, boardID string) ([]*List, error) {
	res := []*List{}

	if err := a.getJSON(ctx, fmt.Sprintf(RouteGetListOnBoard, boardID), &res); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) GetCards(ctx context.Context, listID string) ([]*Card, error) {
	res := []*Card{}
	if err := a.getJSON(ctx, fmt.Sprintf(RouteGetCardsOnList, listID), &res); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) GetCardActions(ctx context.Context, req *GetCardActionsRequest) ([]*CardAction, error) {
	res := []*CardAction{}
	q, _ := query.Values(req)

	if err := a.getJSON(ctx, fmt.Sprintf(RouteGetCardActions, req.CardID, q.Encode()), &res); err != nil {
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
