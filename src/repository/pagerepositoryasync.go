package repository

import (
	"context"
	"github.com/remikj/pages-ms/src/model"
)

type ResultSEO struct {
	SEO *model.SEO
	Err error
}

type ResultProducts struct {
	Products []model.Product
	Err      error
}

type PageRepositoryAsync interface {
	GetSeoForPage(pageId int) (<-chan ResultSEO, context.CancelFunc)
	GetProductsForPage(pageId int) (<-chan ResultProducts, context.CancelFunc)
}

type PageRepositoryAsyncImpl struct {
	pageRepo PageRepository
}

func InitPageRepositoryAsyncFromEnv() (PageRepositoryAsync, error) {
	pageRepo, err := InitPageRepositoryFromEnv()
	if err != nil {
		return nil, err
	}
	return &PageRepositoryAsyncImpl{pageRepo: pageRepo}, nil
}

func (p PageRepositoryAsyncImpl) GetSeoForPage(pageId int) (<-chan ResultSEO, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	seoChan := make(chan ResultSEO, 1)
	go func() {
		page, err := p.pageRepo.GetSeoForPage(ctx, pageId)
		seoChan <- ResultSEO{
			SEO: page,
			Err: err,
		}
	}()
	return seoChan, cancelFunc
}

func (p PageRepositoryAsyncImpl) GetProductsForPage(pageId int) (<-chan ResultProducts, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	productsChan := make(chan ResultProducts, 1)
	go func() {
		products, err := p.pageRepo.GetProductsForPage(ctx, pageId)
		productsChan <- ResultProducts{
			Products: products,
			Err:      err,
		}
	}()
	return productsChan, cancelFunc
}
