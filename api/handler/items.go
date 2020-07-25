package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	apiModel "github.com/mintak21/qiitaWrapper/api/model"
	genModel "github.com/mintak21/qiitaWrapper/gen/models"
	"github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper/items"
	log "github.com/sirupsen/logrus"
)

const (
	qiitaItemsURL = "https://qiita.com/api/v2/items"
	perPage       = 5
)

var client = &http.Client{
	Timeout: 3 * time.Second,
}

// NewGetTagItemsHandler handles a request for getting tag items
func NewGetTagItemsHandler() items.GetTagItemsHandler {
	return &tagItemsHandler{}
}

type tagItemsHandler struct{}

// Handle the get entry request
func (h *tagItemsHandler) Handle(params items.GetTagItemsParams) middleware.Responder {
	response, err := h.sendRequest2Qiita(params)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error(("failed to send request to Qiita API"))
		return items.NewGetTagItemsInternalServerError().WithPayload(&genModel.Error{Message: err.Error()})
	}
	return items.NewGetTagItemsOK().WithPayload(h.res2Model(response))
}

func (h *tagItemsHandler) res2Model(resItems []*apiModel.QiitaItem) *genModel.Items {
	var items []*genModel.Item
	for _, resItem := range resItems {
		item := genModel.Item{
			Title: resItem.Title,
			Link:  resItem.URL,
			User: &genModel.User{
				Name:          resItem.User.Name,
				ThumbnailLink: resItem.User.ProfileImageURL,
			},
			Statistics: &genModel.Statistics{
				Lgtms: int64(resItem.LikesCount),
			},
			CreatedAt: strfmt.DateTime(resItem.CreatedAt),
		}
		items = append(items, &item)
	}
	return &genModel.Items{
		HasNext: false,
		Page:    1,
		Items:   items,
	}
}

func (h *tagItemsHandler) sendRequest2Qiita(params items.GetTagItemsParams) ([]*apiModel.QiitaItem, error) {
	request, err := http.NewRequest("GET", qiitaItemsURL, nil)
	if err != nil {
		return nil, err
	}
	// build query
	requestParams := request.URL.Query()
	requestParams.Add("page", fmt.Sprintf("%d", 1))
	requestParams.Add("per_page", fmt.Sprintf("%d", perPage))
	requestParams.Add("query", fmt.Sprintf("tag:%s", params.Tag))
	request.URL.RawQuery = requestParams.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request does not succeeded: %s", resp.Status)
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var qiitaItems []*apiModel.QiitaItem
	err = json.Unmarshal(byteArray, &qiitaItems)
	if err != nil {
		return nil, err
	}
	return qiitaItems, nil
}
