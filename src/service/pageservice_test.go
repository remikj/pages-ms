package service

import (
	"context"
	"fmt"
	"github.com/remikj/pages-ms/src/model"
	"github.com/remikj/pages-ms/src/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	sleepTime50ms = 50 * time.Millisecond
)

var (
	sampleSeoError      = fmt.Errorf("seo error")
	sampleProductsError = fmt.Errorf("products error")
)

func TestPageServiceImpl_GetPage(t *testing.T) {
	tests := []struct {
		name         string
		repository   repository.PageRepositoryAsync
		pageId       int
		expectedPage *model.Page
		epectedErr   error
	}{
		{
			name: "should succeed when seo and products in repository",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(&sampleModelPage.SEO, nil, 0),
				GetProductsForPageFunc: createGetProductsForPageFunc(sampleModelPage.Products, nil, 0),
			},
			pageId:       0,
			expectedPage: &sampleModelPage,
			epectedErr:   nil,
		},
		{
			name: "should return nil when seo not in repository",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(nil, nil, 0),
				GetProductsForPageFunc: createGetProductsForPageFunc(sampleModelPage.Products, nil, 0),
			},
			pageId:       0,
			expectedPage: nil,
			epectedErr:   nil,
		},
		{
			name: "should return error when error while getting seo",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(nil, sampleSeoError, 0),
				GetProductsForPageFunc: createGetProductsForPageFunc(sampleModelPage.Products, nil, 0),
			},
			pageId:     0,
			epectedErr: sampleSeoError,
		},
		{
			name: "should return error when error while getting products",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(&sampleModelPage.SEO, nil, 0),
				GetProductsForPageFunc: createGetProductsForPageFunc(nil, sampleProductsError, 0),
			},
			pageId:     0,
			epectedErr: sampleProductsError,
		},
		{
			name: "should return seo error when error from seo before error from products",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(nil, sampleSeoError, 0),
				GetProductsForPageFunc: createGetProductsForPageFunc(nil, sampleProductsError, sleepTime50ms),
			},
			pageId:     0,
			epectedErr: sampleSeoError,
		},
		{
			name: "should return seo error when error from seo before error from products",
			repository: pageRepositoryAsyncMock{
				GetSeoForPageFunc:      createGetSeoForPageFunc(nil, sampleSeoError, sleepTime50ms),
				GetProductsForPageFunc: createGetProductsForPageFunc(nil, sampleProductsError, 0),
			},
			pageId:     0,
			epectedErr: sampleProductsError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PageServiceImpl{
				PageRepositoryAsync: tt.repository,
			}

			resultPage, err := ps.GetPage(tt.pageId)

			if tt.epectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.epectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPage, resultPage)
			}

		})
	}
}

func TestPageServiceImpl_GetPage_shouldGetSeoAndProductsAsynchronously(t *testing.T) {
	ps := &PageServiceImpl{
		PageRepositoryAsync: pageRepositoryAsyncMock{
			GetSeoForPageFunc:      createGetSeoForPageFunc(&sampleModelPage.SEO, nil, 100*time.Millisecond),
			GetProductsForPageFunc: createGetProductsForPageFunc(sampleModelPage.Products, nil, 100*time.Millisecond),
		},
	}

	testStart := time.Now()
	page, err := ps.GetPage(0)
	testTime := time.Since(testStart)

	assert.NoError(t, err)
	assert.Equal(t, &sampleModelPage, page)
	assert.Less(t, testTime, 150*time.Millisecond)
	assert.Greater(t, testTime, 99*time.Millisecond)
}

func TestPageServiceImpl_GetPage_shouldCancelGetSeoImmediately_whenProductsFailFirst(t *testing.T) {
	seoCancelTimer := &timeStruct{
		startTime: time.Now(),
	}
	ps := &PageServiceImpl{
		PageRepositoryAsync: pageRepositoryAsyncMock{
			GetSeoForPageFunc:      createGetSeoForPageFuncWithCancelTime(&sampleModelPage.SEO, nil, 100*time.Millisecond, seoCancelTimer),
			GetProductsForPageFunc: createGetProductsForPageFunc(nil, sampleProductsError, 10*time.Millisecond),
		},
	}

	page, err := ps.GetPage(0)

	assert.Error(t, err)
	assert.Nil(t, page)
	assert.Less(t, seoCancelTimer.timeToCancel, 15*time.Millisecond)
	assert.Greater(t, seoCancelTimer.timeToCancel, 9*time.Millisecond)
}

func TestPageServiceImpl_GetPage_shouldCancelGetProductsImmediately_wheSeoFailsFirst(t *testing.T) {
	seoCancelTimer := &timeStruct{
		startTime: time.Now(),
	}
	ps := &PageServiceImpl{
		PageRepositoryAsync: pageRepositoryAsyncMock{
			GetSeoForPageFunc:      createGetSeoForPageFunc(nil, sampleSeoError, 10*time.Millisecond),
			GetProductsForPageFunc: createGetProductsForPageFuncWithCancelTime(sampleModelPage.Products, nil, 100*time.Millisecond, seoCancelTimer),
		},
	}

	page, err := ps.GetPage(0)

	assert.Error(t, err)
	assert.Nil(t, page)
	assert.Less(t, seoCancelTimer.timeToCancel, 15*time.Millisecond)
	assert.Greater(t, seoCancelTimer.timeToCancel, 9*time.Millisecond)
}

func createGetSeoForPageFunc(seo *model.SEO, err error, sleepTime time.Duration) func(pageId int) (<-chan repository.ResultSEO, context.CancelFunc) {
	return createGetSeoForPageFuncWithCancelFunc(seo, err, sleepTime, func() {})
}

func createGetSeoForPageFuncWithCancelTime(seo *model.SEO, err error, sleepTime time.Duration, cancelTimer *timeStruct) func(pageId int) (<-chan repository.ResultSEO, context.CancelFunc) {
	return createGetSeoForPageFuncWithCancelFunc(seo, err, sleepTime, func() {
		cancelTimer.timeToCancel = time.Since(cancelTimer.startTime)
	})
}

func createGetSeoForPageFuncWithCancelFunc(seo *model.SEO, err error, sleepTime time.Duration, cancelFunc context.CancelFunc) func(pageId int) (<-chan repository.ResultSEO, context.CancelFunc) {
	return func(pageId int) (<-chan repository.ResultSEO, context.CancelFunc) {
		seoChan := make(chan repository.ResultSEO, 1)
		go func() {
			time.Sleep(sleepTime)
			seoChan <- repository.ResultSEO{
				SEO: seo,
				Err: err,
			}
		}()
		return seoChan, cancelFunc
	}
}

func createGetProductsForPageFunc(products []model.Product, err error, sleepTime time.Duration) func(pageId int) (<-chan repository.ResultProducts, context.CancelFunc) {
	return createGetProductsForPageFuncWithCancelFunc(products, err, sleepTime, func() {})
}

func createGetProductsForPageFuncWithCancelTime(products []model.Product, err error, sleepTime time.Duration, cancelTimer *timeStruct) func(pageId int) (<-chan repository.ResultProducts, context.CancelFunc) {
	return createGetProductsForPageFuncWithCancelFunc(products, err, sleepTime, func() {
		cancelTimer.timeToCancel = time.Since(cancelTimer.startTime)
	})
}

func createGetProductsForPageFuncWithCancelFunc(products []model.Product, err error, sleepTime time.Duration, cancelFunc context.CancelFunc) func(pageId int) (<-chan repository.ResultProducts, context.CancelFunc) {
	return func(pageId int) (<-chan repository.ResultProducts, context.CancelFunc) {
		productsChan := make(chan repository.ResultProducts, 1)
		go func() {
			time.Sleep(sleepTime)
			productsChan <- repository.ResultProducts{
				Products: products,
				Err:      err,
			}
		}()
		return productsChan, cancelFunc
	}
}

type timeStruct struct {
	startTime    time.Time
	timeToCancel time.Duration
}

type pageRepositoryAsyncMock struct {
	GetSeoForPageFunc      func(pageId int) (<-chan repository.ResultSEO, context.CancelFunc)
	GetProductsForPageFunc func(pageId int) (<-chan repository.ResultProducts, context.CancelFunc)
}

func (p pageRepositoryAsyncMock) GetSeoForPage(pageId int) (<-chan repository.ResultSEO, context.CancelFunc) {
	return p.GetSeoForPageFunc(pageId)
}

func (p pageRepositoryAsyncMock) GetProductsForPage(pageId int) (<-chan repository.ResultProducts, context.CancelFunc) {
	return p.GetProductsForPageFunc(pageId)
}
