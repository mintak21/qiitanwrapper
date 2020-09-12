package handler

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	apiModel "github.com/mintak21/qiitaWrapper/api/model"
	"github.com/mintak21/qiitaWrapper/api/util/client"
	genModel "github.com/mintak21/qiitaWrapper/gen/models"
	"github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper/items"
	log "github.com/sirupsen/logrus"
)

const (
	perPage          = 100
	rfc3339DateMonth = "2006-01"
)

var contentTags = cascadia.MustCompile("h1, h2")

// NewGetTagItemsHandler handles a request for getting tag items
func NewGetTagItemsHandler() items.GetTagItemsHandler {
	return &tagItemsHandler{
		client: client.NewQiitaClient(),
	}
}

type tagItemsHandler struct {
	client client.QiitaClient
}

// NewSyncTagItemsHandler handles a request for getting target day tag items
func NewSyncTagItemsHandler() items.SyncTagItemsHandler {
	return &syncTagItemsHandler{
		client: client.NewQiitaClient(),
	}
}

type syncTagItemsHandler struct {
	client client.QiitaClient
}

// NewMonthlyTrendItemsHandler handles a request for getting monthly tread tag items
func NewMonthlyTrendItemsHandler() items.GetMonthlyTrendItemsHandler {
	return &monthlyTrendItemsHandler{
		client: client.NewQiitaClient(),
	}
}

type monthlyTrendItemsHandler struct {
	client client.QiitaClient
}

// Handle the get entry request
func (h *tagItemsHandler) Handle(params items.GetTagItemsParams) middleware.Responder {
	query := fmt.Sprintf("tag:%s", params.Tag)
	response, hasNext, err := sendGetItemRequest(h.client, int(*params.Page), query)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewGetTagItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	stocksMap, err := createStocksMap(h.client, response)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewSyncTagItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	return items.NewGetTagItemsOK().WithPayload(toModel(response, stocksMap, *params.Page, hasNext))
}

func (h *syncTagItemsHandler) Handle(params items.SyncTagItemsParams) middleware.Responder {
	var targetDate string
	if params.Date == nil {
		targetDate = time.Now().Format(strfmt.RFC3339FullDate)
	} else {
		targetDate = params.Date.String()
	}
	query := fmt.Sprintf("tag:%s created:<=%s created:>=%s", params.Tag, targetDate, targetDate)

	response, hasNext, err := sendGetItemRequest(h.client, int(*params.Page), query)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewSyncTagItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	stocksMap, err := createStocksMap(h.client, response)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewSyncTagItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}

	return items.NewSyncTagItemsOK().WithPayload(toModel(response, stocksMap, *params.Page, hasNext))
}

func (h *monthlyTrendItemsHandler) Handle(params items.GetMonthlyTrendItemsParams) middleware.Responder {
	var fromTime time.Time
	if params.TargetMonth == nil {
		fromTime = time.Now()
	} else {
		parseTime, err := time.Parse(rfc3339DateMonth, *params.TargetMonth)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error(("failed to parse param"))
			return items.NewGetMonthlyTrendItemsBadRequest().WithPayload(&genModel.Error{Message: err.Error()})
		}
		fromTime = parseTime
	}
	toTime := fromTime.AddDate(0, 1, 0)

	fromMonth := fromTime.Format(rfc3339DateMonth)
	toMonth := toTime.Format(rfc3339DateMonth)
	fixedPage := 1

	query := fmt.Sprintf("created:>=%s created:<%s stocks:>=%v", fromMonth, toMonth, 150)

	response, _, err := sendGetItemRequest(h.client, fixedPage, query)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewGetMonthlyTrendItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	stocksMap, err := createStocksMap(h.client, response)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewGetMonthlyTrendItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	sort.SliceStable(response, func(i, j int) bool { return response[i].LikesCount > response[j].LikesCount })

	return items.NewGetMonthlyTrendItemsOK().WithPayload(toModel(response, stocksMap, int64(fixedPage), false))
}

func createStocksMap(clt client.QiitaClient, items []*apiModel.QiitaItem) (map[string]int, error) {
	result := map[string]int{}
	// FIXME: リクエスト数がそれなりになるので、一旦コメントアウト
	// for _, item := range items {
	// 	stocks, err := sendGetStockerRequest(clt, item.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	result[item.ID] = stocks
	// }
	return result, nil
}

func sendGetItemRequest(clt client.QiitaClient, page int, query string) ([]*apiModel.QiitaItem, bool, error) {
	parameter := client.NewGetItemsParameter(page, perPage, query)
	qiitaItems, err := clt.GetItems(parameter)
	if err != nil {
		return nil, false, err
	}
	return qiitaItems, perPage <= len(qiitaItems), nil
}

// func sendGetStockerRequest(cl client.QiitaClient, itemID string) (int, error) {
// 	parameter := client.NewGetStockersParameter(itemID)
// 	stockers, err := cl.GetStockers(parameter)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return len(stockers), nil
// }

func toModel(resItems []*apiModel.QiitaItem, stocks map[string]int, page int64, hasNext bool) *genModel.Items {
	items := make([]*genModel.Item, 0, len(resItems))
	for _, resItem := range resItems {
		var tags []string
		for _, tag := range resItem.Tags {
			tags = append(tags, tag.Name)
		}
		item := genModel.Item{
			Title:         resItem.Title,
			Link:          resItem.URL,
			Tags:          tags,
			TableContents: contents(resItem),
			User: &genModel.User{
				Name:          resItem.User.ID,
				ThumbnailLink: resItem.User.ProfileImageURL,
			},
			Statistics: &genModel.Statistics{
				Lgtms:  int64(resItem.LikesCount),
				Stocks: int64(stocks[resItem.ID]),
			},
			CreatedAt: strfmt.DateTime(resItem.CreatedAt),
		}
		items = append(items, &item)
	}
	return &genModel.Items{
		HasNext: hasNext,
		Page:    page,
		Items:   items,
	}
}

func contents(apiItem *apiModel.QiitaItem) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(apiItem.RenderedBody))
	if err != nil {
		// ここのエラーはログに出すだけで握りつぶし
		log.Warn("failed to scrape contents", "err", err)
		return ""
	}
	return doc.FindMatcher(contentTags).Text()
}

func unusedFunction2() int {
	return 100
}
