package prices

import (
	"context"
	"database/sql"

	"github.com/turbak/bigmacindex/internal/domain/price"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetPriceByCountryCode(ctx context.Context, countryCode string) (price.PriceRecord, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, product_name, priceRec, price_cents, currency, country_code, created_date FROM prices WHERE country_code = ?", countryCode)
	var priceRec price.PriceRecord
	err := row.Scan(&priceRec.ID, &priceRec.ProductName, &priceRec.Price, &priceRec.PriceCents, &priceRec.Currency, &priceRec.CountryCode, &priceRec.CreatedDate)
	if err != nil {
		return price.PriceRecord{}, err
	}
	return priceRec, nil
}
func (r *repository) UpsertPrice(ctx context.Context, price price.PriceRecord) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO prices (product_name, price, price_cents, currency, country_code, created_date) 
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT(product_name, country_code, created_date) DO UPDATE SET 
				product_name = excluded.product_name,
				price = excluded.price,
				price_cents = excluded.price_cents,
				currency = excluded.currency`,
		price.ProductName, price.Price, price.PriceCents, price.Currency, price.CountryCode, price.CreatedDate)
	return err
}
