package contoller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/remikj/pages-ms/src/model"
	"github.com/remikj/pages-ms/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	sampleModelPage = model.Page{
		SEO: model.SEO{
			PageId:      0,
			Title:       "Sample page title",
			Description: "Sample page description",
			Robots:      "Sample robots",
		},
		Products: []model.Product{
			{
				Id:          0,
				PageId:      0,
				Name:        "Sample product 0 name",
				Description: "Sample product 0 description",
				Price:       2.50,
			},
			{
				Id:          1,
				PageId:      0,
				Name:        "Sample product 1 name",
				Description: "Sample product 1 description",
				Price:       19.99,
			},
		},
	}
	sampleModelPageString = getSampleModelPageString()
)

func TestPageControllerImpl_HandlePageGet(t *testing.T) {
	type expectedWrite struct {
		code       int
		bodyString string
	}
	tests := []struct {
		name        string
		pageService service.PageService
		request     *http.Request
		expected    expectedWrite
	}{
		{
			name: "should return correct json when all data valid",
			pageService: &pageServiceMock{
				getPageFn: func(pageId int) (*model.Page, error) {
					if pageId == 1 {
						return &sampleModelPage, nil
					} else {
						return nil, errors.New("incorrect argument passed to PageService")
					}
				},
			},
			request:  requestWithParam("1"),
			expected: expectedWrite{code: http.StatusOK, bodyString: sampleModelPageString},
		},
		{
			name:     "should return correct json when all data valid",
			request:  requestWithParam("not-a-number"),
			expected: expectedWrite{code: http.StatusBadRequest, bodyString: "Expected pageId to be number"},
		},
		{
			name:    "should return internal server error when PageService fails",
			request: requestWithParam("1"),
			pageService: &pageServiceMock{
				getPageFn: func(pageId int) (*model.Page, error) {
					return nil, errors.New("PageService failed")
				},
			},
			expected: expectedWrite{code: http.StatusInternalServerError, bodyString: "Unexpected error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := PageControllerImpl{
				PageService: tt.pageService,
			}
			responseRecorder := httptest.NewRecorder()

			pc.HandlePageGet(responseRecorder, tt.request)

			assert.Equal(t, tt.expected.code, responseRecorder.Code)
			responseBodyBytes, err := ioutil.ReadAll(responseRecorder.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.bodyString, string(responseBodyBytes))
		})
	}
}

func requestWithParam(s string) *http.Request {
	request := httptest.NewRequest("GET", "/pages/"+s, nil)

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", s)

	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
}

func getSampleModelPageString() string {
	marshal, _ := json.Marshal(sampleModelPage)
	return string(marshal)
}

type pageServiceMock struct {
	getPageFn func(pageId int) (*model.Page, error)
}

func (p pageServiceMock) GetPage(pageId int) (*model.Page, error) {
	return p.getPageFn(pageId)
}
