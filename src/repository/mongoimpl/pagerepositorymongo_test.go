package mongoimpl

import (
	"context"
	"fmt"
	"github.com/remikj/pages-ms/src/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var (
	sampleSeo = model.SEO{
		PageId:      0,
		Title:       "title",
		Description: "description",
		Robots:      "robots",
	}
	sampleProduct1 = model.Product{
		Id:          0,
		PageId:      0,
		Name:        "name0",
		Description: "description0",
		Price:       1.0,
	}
	sampleProduct2 = model.Product{
		Id:          1,
		PageId:      0,
		Name:        "name1",
		Description: "description1",
		Price:       11.99,
	}
	sampleProducts = []model.Product{
		{
			Id:          0,
			PageId:      0,
			Name:        "name0",
			Description: "description0",
			Price:       1.0,
		},
		{
			Id:          1,
			PageId:      0,
			Name:        "name1",
			Description: "description1",
			Price:       11.99,
		},
	}
)

func TestPageRepositoryMongo_GetSeoForPage(t *testing.T) {
	tests := []struct {
		name        string
		mongoClient Client
		pageId      int
		expectedSeo *model.SEO
		expectedErr error
	}{
		{
			name: "should successfully find seo, when one seo in cursor",
			mongoClient: mongoClientMock{
				findSeosFunc: createFindFunc(mockMongoCursor([][]byte{marshal(sampleSeo)}), nil),
			},
			pageId:      0,
			expectedSeo: &sampleSeo,
			expectedErr: nil,
		},
		{
			name: "should return nil, when no seos in cursor",
			mongoClient: mongoClientMock{
				findSeosFunc: createFindFunc(mockMongoCursor([][]byte{}), nil),
			},
			pageId:      0,
			expectedSeo: nil,
			expectedErr: nil,
		},
		{
			name: "should return err, when findSeos fails",
			mongoClient: mongoClientMock{
				findSeosFunc: createFindFunc(nil, fmt.Errorf("findSeos error")),
			},
			pageId:      0,
			expectedSeo: nil,
			expectedErr: fmt.Errorf("error happened when using db: findSeos error"),
		},
		{
			name: "should fail, when multiple seos in cursor",
			mongoClient: mongoClientMock{
				findSeosFunc: createFindFunc(mockMongoCursor([][]byte{marshal(sampleSeo), marshal(sampleSeo)}), nil),
			},
			pageId:      0,
			expectedSeo: nil,
			expectedErr: fmt.Errorf("too many results"),
		},
		{
			name: "should return err, when decode fails",
			mongoClient: mongoClientMock{
				findSeosFunc: createFindFunc(mockMongoCursor([][]byte{[]byte("incorrectBytes")}), nil),
			},
			pageId:      0,
			expectedSeo: nil,
			expectedErr: fmt.Errorf("error happened when decoding results: invalid document length"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageRepositoryMongo{
				mongoClient: tt.mongoClient,
			}
			resultSeo, err := p.GetSeoForPage(context.Background(), tt.pageId)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedSeo, resultSeo)

		})
	}
}

func TestPageRepositoryMongo_GetProductsForPage(t *testing.T) {
	tests := []struct {
		name             string
		mongoClient      Client
		pageId           int
		expectedProducts []model.Product
		expectedErr      error
	}{
		{
			name: "should return empty products array, when no products in cursor",
			mongoClient: mongoClientMock{
				findProductsFunc: createFindFunc(mockMongoCursor([][]byte{marshal(sampleProduct1), marshal(sampleProduct2)}), nil),
			},
			pageId:           0,
			expectedProducts: sampleProducts,
			expectedErr:      nil,
		},
		{
			name: "should return err products array, when no products in cursor",
			mongoClient: mongoClientMock{
				findProductsFunc: createFindFunc(nil, fmt.Errorf("findProducts error")),
			},
			pageId:           0,
			expectedProducts: nil,
			expectedErr:      fmt.Errorf("error happened when using db: findProducts error"),
		},
		{
			name: "should return empty products array, when no products in cursor",
			mongoClient: mongoClientMock{
				findProductsFunc: createFindFunc(mockMongoCursor([][]byte{[]byte("incorrectBytes")}), nil),
			},
			pageId:           0,
			expectedProducts: nil,
			expectedErr:      fmt.Errorf("error happened when decoding results: invalid document length"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageRepositoryMongo{
				mongoClient: tt.mongoClient,
			}
			resultSeo, err := p.GetProductsForPage(context.Background(), tt.pageId)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedProducts, resultSeo)
		})
	}
}

func createFindFunc(cursor MongoCursor, err error) func(ctx context.Context, pageId int) (MongoCursor, error) {
	return func(ctx context.Context, pageId int) (MongoCursor, error) {
		if pageId == 0 {
			return cursor, err
		} else {
			return nil, fmt.Errorf("unexpected pageId")
		}
	}
}

func mockMongoCursor(results [][]byte) MongoCursor {
	return &mongoCursosMock{
		idx:     -1,
		results: results,
	}
}

func marshal(val interface{}) []byte {
	bytes, _ := bson.Marshal(val)
	return bytes
}

type mongoClientMock struct {
	findSeosFunc     func(ctx context.Context, pageId int) (MongoCursor, error)
	findProductsFunc func(ctx context.Context, pageId int) (MongoCursor, error)
}

func (m mongoClientMock) FindSeos(ctx context.Context, pageId int) (MongoCursor, error) {
	return m.findSeosFunc(ctx, pageId)

}

func (m mongoClientMock) FindProducts(ctx context.Context, pageId int) (MongoCursor, error) {
	return m.findProductsFunc(ctx, pageId)
}

func (m mongoClientMock) CloseMongoClient() error {
	return nil
}

type mongoCursosMock struct {
	idx     int
	results [][]byte
}

func (m *mongoCursosMock) Next(_ context.Context) bool {
	m.idx++
	return m.idx < len(m.results)
}

func (m *mongoCursosMock) All(_ context.Context, vals interface{}) error {
	valsArrPointer := vals.(*[]model.Product)
	for _, result := range m.results {
		resultUnmarshal := model.Product{}
		err := bson.Unmarshal(result, &resultUnmarshal)
		if err != nil {
			return err
		}
		*valsArrPointer = append(*valsArrPointer, resultUnmarshal)
	}
	return nil
}

func (m *mongoCursosMock) Decode(val interface{}) error {
	return bson.Unmarshal(m.results[m.idx], val)
}

func (m *mongoCursosMock) Close(ctx context.Context) error {
	return nil
}

func (m *mongoCursosMock) Err() error {
	return nil
}
