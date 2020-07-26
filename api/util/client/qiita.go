package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"

	"github.com/mintak21/qiitaWrapper/api/model"
)

const (
	qiitaDomain      = "https://qiita.com"
	itemsEndpoint    = "/api/v2/items"
	stockersEndpoint = "/api/v2/items/%s/stockers"
	timeout          = 3 * time.Second
)

// QiitaClient client to Qiita API
type QiitaClient interface {
	GetItems(param *GetItemsParameter) ([]*model.QiitaItem, error)
}

type qiitaClient struct {
	client *resty.Client
}

// CommonParameter common request parameter
type CommonParameter struct {
	Page    int `validate:"gt=0"`
	PerPage int `validate:"gt=0,lte=100"`
}

// GetItemsParameter parameters for GET /api/v2/items
type GetItemsParameter struct {
	Common *CommonParameter
	Query  string // ここの組み立ては改善ポイント
}

// NewQiitaClient creates qiitaClient
func NewQiitaClient() QiitaClient {
	return &qiitaClient{
		client: resty.New().
			SetTimeout(timeout).
			SetHeader("Content-Type", "application/json").
			SetContentLength(true),
	}
}

// NewGetItemsParameter creates GetItemsParameter
func NewGetItemsParameter(page, perPage int, query string) *GetItemsParameter {
	return &GetItemsParameter{
		Common: &CommonParameter{
			Page:    page,
			PerPage: perPage,
		},
		Query: query,
	}
}

// GetItems send request to /api/v2/items
func (c *qiitaClient) GetItems(param *GetItemsParameter) ([]*model.QiitaItem, error) {
	validate := validator.New()
	if err := validate.Struct(param); err != nil {
		return nil, err
	}
	fmt.Println("query", param.Query)
	resp, err := c.sendGetRequest(
		fmt.Sprintf("%s%s", qiitaDomain, itemsEndpoint),
		map[string]string{
			"page":     strconv.Itoa(param.Common.Page),
			"per_page": strconv.Itoa(param.Common.PerPage),
			"query":    param.Query,
		},
	)
	if err != nil {
		return nil, err
	}

	var results []*model.QiitaItem
	err = json.Unmarshal(resp.Body(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *qiitaClient) sendGetRequest(endpoint string, params map[string]string) (*resty.Response, error) {
	resp, err := c.client.R().SetQueryParams(params).Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("error occurs via send request to Qiita: %s", resp.Status())
	}
	return resp, nil
}
