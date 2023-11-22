package models

type Immovable struct {
	Title          string `json:"title" bson:"title"`
	Link           string `json:"link" bson:"link"`
	Data           string `json:"data" bson:"data"`
	Price          int    `json:"price" bson:"price"`
	PriceInitially int    `json:"priceprice_initially" bson:"priceInitially"`
}
