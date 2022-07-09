package model

type Page struct {
	SEO      SEO
	Products []Product
}

type SEO struct {
	PageId      int    `bson:"page_id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Robots      string `bson:"robots"`
}

type Product struct {
	Id          int     `bson:"id"`
	PageId      int     `bson:"page_id"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
}
