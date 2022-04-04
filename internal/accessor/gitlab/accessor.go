package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/google/go-querystring/query"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/config"
)

var ErrHTTPStatusNon2xx = errors.New("err_http_status_non_2xx")

type Accessor interface {
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	CreateMergeRequest(ctx context.Context, req *CreateMergeRequestRequest) (*MergeRequest, error)
	ListMergeRequests(ctx context.Context, req *ListMergeRequestRequest) ([]*MergeRequest, error)
}

type accessor struct {
	httpClient *http.Client
}

func NewAccessor(httpClient *http.Client) *accessor {
	return &accessor{httpClient: httpClient}
}

func (a accessor) postJSON(ctx context.Context, path string, req, res interface{}) error {

	reqBuf := &bytes.Buffer{}
	err := json.NewEncoder(reqBuf).Encode(req)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	fullURL := fmt.Sprintf("https://%s/%s", GitlabHost, path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(reqBuf.Bytes()))
	if err != nil {
		return errlib.WrapFunc(err)
	}

	httpReq.Header.Add("PRIVATE-TOKEN", config.GetGitlabPrivateToken())
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

func (a accessor) getJSON(ctx context.Context, path string, res interface{}) (*http.Response, error) {
	fullURL := a.getEndpoint(path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	httpReq.Header.Add("PRIVATE-TOKEN", config.GetGitlabPrivateToken())

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}
	defer resp.Body.Close()

	resBuf := &bytes.Buffer{}
	_, err = io.Copy(resBuf, resp.Body)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	if resp.StatusCode >= 400 {
		return nil, errlib.WrapFunc(errlib.WithFields(ErrHTTPStatusNon2xx, errlib.Fields{
			"response": resBuf.String(),
			"status":   resp.StatusCode,
			"endpoint": fullURL,
		}))
	}

	err = json.Unmarshal(resBuf.Bytes(), res)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return resp, nil
}

func (a accessor) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	res := &Project{}

	_, err := a.getJSON(ctx, fmt.Sprintf(RouteGetProjectsByName, url.QueryEscape(name)), res)

	if err != nil {
		return nil, errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
			"name": name,
		}))
	}

	return res, nil
}

func (a accessor) CreateMergeRequest(ctx context.Context, req *CreateMergeRequestRequest) (*MergeRequest, error) {
	res := &MergeRequest{}

	err := a.postJSON(ctx, fmt.Sprintf(RouteCreateMergeRequest, req.ID), req, res)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) ListMergeRequests(ctx context.Context, req *ListMergeRequestRequest) ([]*MergeRequest, error) {
	res := []*MergeRequest{}

	q, _ := query.Values(req)
	endpoint := fmt.Sprintf(RouteListMergeRequests, req.ID, q.Encode())

	if err := paginate(
		endpoint,
		func(nextEndpoint string) (interface{}, error) {
			currentBatch := []*MergeRequest{}
			resp, err := a.getJSON(ctx, nextEndpoint, &currentBatch)
			if err != nil {
				return nil, errlib.WrapFunc(err)
			}
			res = append(res, currentBatch...)
			return resp, nil
		}, func(resp interface{}) string {
			return defaultGetLinkHeader(resp.(*http.Response))
		},
	); err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (a accessor) getEndpoint(path string) string {
	if isValidURL(path) {
		return path
	}

	return fmt.Sprintf("https://%s/%s", GitlabHost, path)
}
