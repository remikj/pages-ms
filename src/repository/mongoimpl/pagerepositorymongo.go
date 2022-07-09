package mongoimpl

import (
	"context"
	"fmt"
	"github.com/remikj/pages-ms/src/model"
)

type PageRepositoryMongo struct {
	mongoClient Client
}

func InitPageRepositoryMongoFromEnv() (*PageRepositoryMongo, error) {
	mongoClient, err := InitMongoFromEnv()
	if err != nil {
		return nil, err
	}
	return &PageRepositoryMongo{
		mongoClient: mongoClient,
	}, nil
}

func (p PageRepositoryMongo) GetSeoForPage(ctx context.Context, pageId int) (*model.SEO, error) {
	fmt.Printf("Getting seo for page_id: %v\n", pageId)
	seosCursor, err := p.mongoClient.FindSeos(ctx, pageId)
	if err != nil {
		return nil, fmt.Errorf("error happened when using db: %w", err)
	}
	defer seosCursor.Close(ctx)

	seo := &model.SEO{}
	if seosCursor.Next(ctx) {
		err := seosCursor.Decode(seo)
		if err != nil {
			return nil, fmt.Errorf("error happened when decoding results: %w", err)
		}
	} else {
		return nil, nil
	}

	if seosCursor.Next(ctx) {
		return nil, fmt.Errorf("too many results")
	}
	return seo, seosCursor.Err()
}

func (p PageRepositoryMongo) GetProductsForPage(ctx context.Context, pageId int) ([]model.Product, error) {
	fmt.Printf("Getting products for page_id: %v\n", pageId)
	productsCursor, err := p.mongoClient.FindProducts(ctx, pageId)
	if err != nil {
		return nil, fmt.Errorf("error happened when using db: %w", err)
	}

	var products []model.Product
	if err = productsCursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("error happened when decoding results: %w", err)
	}
	return products, nil
}

func (p PageRepositoryMongo) CloseRepository() error {
	return p.mongoClient.CloseMongoClient()
}
