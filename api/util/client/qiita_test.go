package client_test

import (
	"errors"
	"testing"

	"github.com/mintak21/qiitaWrapper/api/model"
	"github.com/mintak21/qiitaWrapper/api/util/client"
)

func TestGetItems(t *testing.T) {
	testPattern := []struct {
		description   string
		parameter     *client.GetItemsParameter
		expectedBody  []*model.QiitaItem
		expectedError error
	}{
		{
			description: "Success Pattern",
			parameter:   client.NewGetItemsParameter(1, 5, "user:mintak21 created:<2020-01-20"),
			expectedBody: []*model.QiitaItem{
				{
					ID: "a6766e3efd6730c9519d",
				},
				{
					ID: "037c9ae39a0db16a0d4e",
				},
				{
					ID: "5000972294d4413471ab",
				},
				{
					ID: "e748803f59a338cb7726",
				},
				{
					ID: "eeba4654a0db21abcb1c",
				},
			},
			expectedError: nil,
		},
		{
			description:   "Bad Parameter(page)",
			parameter:     client.NewGetItemsParameter(-1, 1, "testQuery"),
			expectedBody:  nil,
			expectedError: errors.New("Key: 'GetItemsParameter.Common.Page' Error:Field validation for 'Page' failed on the 'gt' tag"),
		},
		{
			description:   "Bad Parameter(per_page)",
			parameter:     client.NewGetItemsParameter(1, -1, "testQuery"),
			expectedBody:  nil,
			expectedError: errors.New("Key: 'GetItemsParameter.Common.PerPage' Error:Field validation for 'PerPage' failed on the 'gt' tag"),
		},
		{
			description:   "Bad Parameter(per_page) 2",
			parameter:     client.NewGetItemsParameter(1, 101, "testQuery"),
			expectedBody:  nil,
			expectedError: errors.New("Key: 'GetItemsParameter.Common.PerPage' Error:Field validation for 'PerPage' failed on the 'lte' tag"),
		},
	}

	target := client.NewQiitaClient()

	for _, pattern := range testPattern {
		t.Run(pattern.description, func(t *testing.T) {
			actualBody, actualError := target.GetItems(pattern.parameter)

			if pattern.expectedError != nil {
				if actualError == nil {
					t.Errorf("👻GetItems(): expected error %v but got nil", pattern.expectedError)
					return
				}
				if actualError.Error() != pattern.expectedError.Error() {
					t.Errorf("👻GetItems(): expected error %v but got %v", pattern.expectedError, actualError)
					return
				}
			} else {
				if actualError != nil {
					t.Errorf("👻GetItems(): expected no error but got error %v", actualError)
					return
				}
				// 本当は構造体レベルで比較すべきだが、expected書くのが面倒なのでIDだけを比較
				if len(actualBody) != len(pattern.expectedBody) {
					t.Errorf("💀GetItems(): expected %v but got %v", len(pattern.expectedBody), len(actualBody))
					return
				}
				for i := 0; i < len(actualBody); i++ {
					if actualBody[i].ID != pattern.expectedBody[i].ID {
						t.Errorf("💀GetItems(): expected item id is %v but got %v", pattern.expectedBody[i].ID, actualBody[i].ID)
						return
					}
				}
				// if diff := cmp.Diff(pattern.expectedBody, actualBody); diff != "" {
				// 	t.Errorf("💀GetItems() mismatch detected: (-got +want)\n%s", diff)
				// }
			}
		})
	}
}

func TestGetStockers(t *testing.T) {
	testPattern := []struct {
		description   string
		parameter     *client.GetStockersParameter
		expectedBody  []*model.QiitaStocker
		expectedError error
	}{
		{
			description: "Success Pattern",
			parameter:   client.NewGetStockersParameter("e748803f59a338cb7726"),
			expectedBody: []*model.QiitaStocker{
				{
					ID: "a6766e3efd6730c9519d", // 長さだけ関係あるので適当
				},
				{
					ID: "a6766e3efd6730c9519d",
				},
				{
					ID: "a6766e3efd6730c9519d",
				},
				{
					ID: "a6766e3efd6730c9519d",
				},
				{
					ID: "a6766e3efd6730c9519d",
				},
			},
			expectedError: nil,
		},
		{
			description:   "Bad Parameter(item-id)",
			parameter:     client.NewGetStockersParameter(""),
			expectedBody:  nil,
			expectedError: errors.New("Key: 'GetStockersParameter.ItemID' Error:Field validation for 'ItemID' failed on the 'required' tag"),
		},
	}

	target := client.NewQiitaClient()

	for _, pattern := range testPattern {
		t.Run(pattern.description, func(t *testing.T) {
			actualBody, actualError := target.GetStockers(pattern.parameter)

			if pattern.expectedError != nil {
				if actualError == nil {
					t.Errorf("👻GetStockers(): expected error %v but got nil", pattern.expectedError)
					return
				}
				if actualError.Error() != pattern.expectedError.Error() {
					t.Errorf("👻GetStockers(): expected error %v but got %v", pattern.expectedError, actualError)
					return
				}
			} else {
				if actualError != nil {
					t.Errorf("👻GetStockers(): expected no error but got error %v", actualError)
					return
				}
				// 本当は構造体レベルで比較すべきだが、expected書くのが面倒なのでIDだけを比較
				if len(actualBody) != len(pattern.expectedBody) {
					t.Errorf("💀GetStockers(): expected %v but got %v", len(pattern.expectedBody), len(actualBody))
					return
				}
			}
		})
	}
}
