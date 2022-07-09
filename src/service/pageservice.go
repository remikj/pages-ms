package service

import (
	"fmt"
	"github.com/remikj/pages-ms/src/model"
	"github.com/remikj/pages-ms/src/repository"
)

type PageService interface {
	GetPage(pageId int) (*model.Page, error)
}

type PageServiceImpl struct {
	PageRepositoryAsync repository.PageRepositoryAsync
}

func NewPageService(pageRepositoryAsync repository.PageRepositoryAsync) *PageServiceImpl {
	return &PageServiceImpl{
		PageRepositoryAsync: pageRepositoryAsync,
	}
}

func (ps *PageServiceImpl) GetPage(pageId int) (*model.Page, error) {
	fmt.Printf("Getting page for id: %v\n", pageId)
	seoChan, getSeoCancelFunc := ps.PageRepositoryAsync.GetSeoForPage(pageId)
	defer getSeoCancelFunc()
	productsChan, getProductsCancelFunc := ps.PageRepositoryAsync.GetProductsForPage(pageId)
	defer getProductsCancelFunc()

	seoReceived, productsReceived := false, false
	page := model.Page{
		Products: []model.Product{},
	}
	for !seoReceived || !productsReceived {
		select {
		case seoResult := <-seoChan:
			if seoResult.Err != nil {
				return nil, seoResult.Err
			}
			if seoResult.SEO == nil {
				return nil, nil
			}
			page.SEO = *seoResult.SEO
			seoReceived = true
		case productsResult := <-productsChan:
			if productsResult.Err != nil {
				return nil, productsResult.Err
			}
			if productsResult.Products != nil {
				page.Products = productsResult.Products
			}
			productsReceived = true
		}
	}
	return &page, nil
}
