package repository

import (
	"context"
	"fmt"
	"github.com/remikj/pages-ms/src/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	sampleSeo = model.SEO{
		PageId:      0,
		Title:       "Sample page title",
		Description: "Sample page description",
		Robots:      "Sample robots",
	}
)

func TestPageRepositoryAsyncImpl_GetSeoForPage(t *testing.T) {
	tests := []struct {
		name           string
		pageRepo       PageRepository
		pageId         int
		expectedResult ResultSEO
	}{
		{
			name: "should return result from channel",
			pageRepo: pageRepositoryMock{
				getSeoForPageFunc: createGetSeoForPageFunc(&sampleSeo, nil),
			},
			pageId: 0,
			expectedResult: ResultSEO{
				SEO: &sampleSeo,
				Err: nil,
			},
		},
		{
			name: "should return result with error from channel, when getSeo returns err",
			pageRepo: pageRepositoryMock{
				getSeoForPageFunc: createGetSeoForPageFunc(nil, fmt.Errorf("error when getting seo")),
			},
			pageId: 0,
			expectedResult: ResultSEO{
				SEO: nil,
				Err: fmt.Errorf("error when getting seo"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageRepositoryAsyncImpl{
				pageRepo: tt.pageRepo,
			}

			resultChan, cancelFunc := p.GetSeoForPage(tt.pageId)
			result := <-resultChan

			assert.NotNil(t, cancelFunc)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestPageRepositoryAsyncImpl_GetProductsForPage(t *testing.T) {
	tests := []struct {
		name           string
		pageRepo       PageRepository
		pageId         int
		expectedResult ResultProducts
	}{
		{
			name: "should return result from channel",
			pageRepo: pageRepositoryMock{
				getProductsForPageFunc: createGetProductsForPageFunc([]model.Product{}, nil),
			},
			pageId: 0,
			expectedResult: ResultProducts{
				Products: []model.Product{},
				Err:      nil,
			},
		},
		{
			name: "should return result with error from channel, when getSeo returns err",
			pageRepo: pageRepositoryMock{
				getProductsForPageFunc: createGetProductsForPageFunc(nil, fmt.Errorf("error when getting seo")),
			},
			pageId: 0,
			expectedResult: ResultProducts{
				Products: nil,
				Err:      fmt.Errorf("error when getting seo"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageRepositoryAsyncImpl{
				pageRepo: tt.pageRepo,
			}

			resultChan, cancelFunc := p.GetProductsForPage(tt.pageId)
			result := <-resultChan

			assert.NotNil(t, cancelFunc)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func createGetSeoForPageFunc(seo *model.SEO, err error) func(ctx context.Context, pageId int) (*model.SEO, error) {
	return func(ctx context.Context, pageId int) (*model.SEO, error) {
		if pageId == 0 {
			return seo, err
		} else {
			return nil, fmt.Errorf("unexpected pageId")
		}
	}
}

func createGetProductsForPageFunc(products []model.Product, err error) func(ctx context.Context, pageId int) ([]model.Product, error) {
	return func(ctx context.Context, pageId int) ([]model.Product, error) {
		if pageId == 0 {
			return products, err
		} else {
			return nil, fmt.Errorf("unexpected pageId")
		}
	}
}

type pageRepositoryMock struct {
	getSeoForPageFunc      func(ctx context.Context, pageId int) (*model.SEO, error)
	getProductsForPageFunc func(ctx context.Context, pageId int) ([]model.Product, error)
}

func (p pageRepositoryMock) GetSeoForPage(ctx context.Context, pageId int) (*model.SEO, error) {
	return p.getSeoForPageFunc(ctx, pageId)
}

func (p pageRepositoryMock) GetProductsForPage(ctx context.Context, pageId int) ([]model.Product, error) {
	return p.getProductsForPageFunc(ctx, pageId)
}

func (p pageRepositoryMock) CloseRepository() error {
	return nil
}
