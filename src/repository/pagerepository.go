package repository

import (
	"context"
	"github.com/remikj/pages-ms/src/model"
	"github.com/remikj/pages-ms/src/repository/mongoimpl"
)

type PageRepository interface {
	GetSeoForPage(ctx context.Context, pageId int) (*model.SEO, error)
	GetProductsForPage(ctx context.Context, pageId int) ([]model.Product, error)
	CloseRepository() error
}

func InitPageRepositoryFromEnv() (PageRepository, error) {
	return mongoimpl.InitPageRepositoryMongoFromEnv()
}
