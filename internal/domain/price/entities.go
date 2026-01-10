package price

type ID int32

type PriceRecord struct {
	ID          ID     `db:"id"`
	ProductName string `db:"product_name"`
	Price       int32  `db:"price"`
	PriceCents  int32  `db:"price_cents"`
	Currency    string `db:"currency"`
	CountryCode string `db:"country_code"`
	CreatedDate string `db:"created_date"`
}
