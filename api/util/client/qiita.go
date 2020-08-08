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
	GetStockers(param *GetStockersParameter) ([]*model.QiitaStocker, error)
}

type qiitaClient struct {
	client *resty.Client
}

// CommonParameter common request parameter
type CommonParameter struct {
	Page    int `validate:"required,gt=0"`
	PerPage int `validate:"required,gt=0,lte=100"`
}

// GetItemsParameter parameters for GET /api/v2/items
type GetItemsParameter struct {
	Common *CommonParameter
	Query  string // ここの組み立ては改善ポイント
}

// GetStockersParameter parameters for GET /api/v2/items/:item_id/stockers
type GetStockersParameter struct {
	Common *CommonParameter
	ItemID string `validate:"required"`
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
	resp, err := c.sendGetRequest(
		fmt.Sprintf("%s%s", qiitaDomain, itemsEndpoint),
		map[string]string{
			"page":     strconv.Itoa(param.Common.Page),
			"per_page": strconv.Itoa(param.Common.PerPage),
			"query":    param.Query,
		},
	)
	fmt.Print(param.Common, param.Query)
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

// NewGetStockersParameter creates GetStockersParameter
func NewGetStockersParameter(itemID string) *GetStockersParameter {
	return &GetStockersParameter{
		ItemID: itemID,
	}
}

// GetStockers send request to /api/v2/items/:item_id/stockers
func (c *qiitaClient) GetStockers(param *GetStockersParameter) ([]*model.QiitaStocker, error) {
	var currentPage int
	validate := validator.New()
	if err := validate.Struct(param); err != nil {
		return nil, err
	}
	var results []*model.QiitaStocker
	endPoint := fmt.Sprintf(stockersEndpoint, param.ItemID)
	for {
		var stockers []*model.QiitaStocker
		currentPage++
		resp, err := c.sendGetRequest(
			fmt.Sprintf("%s%s", qiitaDomain, endPoint),
			map[string]string{
				"page":     strconv.Itoa(currentPage),
				"per_page": strconv.Itoa(param.Common.PerPage),
			},
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(resp.Body(), &stockers)
		if err != nil {
			return nil, err
		}
		results = append(results, stockers...)
		if len(stockers) < param.Common.PerPage {
			break
		}
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
