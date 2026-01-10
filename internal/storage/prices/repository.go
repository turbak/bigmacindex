package prices

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/turbak/bigmacindex/internal/domain/price"
)

const tableName = "prices"

type repository struct {
	db squirrel.StatementBuilderType
}

func NewRepository(db *sql.DB) *repository {
	dbCache := squirrel.NewStmtCache(db)
	sqDB := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).RunWith(dbCache)
	return &repository{
		db: sqDB,
	}
}

func (r *repository) ListPrices(ctx context.Context) ([]price.PriceRecord, error) {
	rows, err := r.db.Select("id", "product_name", "price", "price_cents", "currency", "country_code", "created_date").
		From(tableName).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priceRecs []price.PriceRecord
	for rows.Next() {
		var priceRec price.PriceRecord
		err := rows.Scan(&priceRec.ID, &priceRec.ProductName, &priceRec.Price, &priceRec.PriceCents, &priceRec.Currency, &priceRec.CountryCode, &priceRec.CreatedDate)
		if err != nil {
			return nil, err
		}
		priceRecs = append(priceRecs, priceRec)
	}
	return priceRecs, nil
}

func (r *repository) UpsertPrice(ctx context.Context, priceRec price.PriceRecord) (price.PriceRecord, error) {
	// Squirrel doesn't have built-in UPSERT support, so we'll use raw SQL for the ON CONFLICT part
	res, err := r.db.Insert(tableName).
		Columns("product_name", "price", "price_cents", "currency", "country_code", "created_date").
		Values(priceRec.ProductName, priceRec.Price, priceRec.PriceCents, priceRec.Currency, priceRec.CountryCode, priceRec.CreatedDate).
		SuffixExpr(
			squirrel.Expr(` ON CONFLICT(product_name, country_code, created_date) DO UPDATE SET
									product_name = excluded.product_name,
									price = excluded.price,
									price_cents = excluded.price_cents,
									currency = excluded.currency`),
		).
		ExecContext(ctx)
	if err != nil {
		return price.PriceRecord{}, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return price.PriceRecord{}, err
	}

	priceRec.ID = price.ID(ID)

	return priceRec, nil
}
