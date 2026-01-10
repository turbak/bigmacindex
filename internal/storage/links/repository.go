package links

import (
	"context"
	"database/sql"

	"github.com/turbak/bigmacindex/internal/domain/link"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) DeleteLink(ctx context.Context, ID link.ID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM links WHERE id = ?", ID)
	return err
}

func (r *repository) ListLinks(ctx context.Context) ([]link.LinkDescription, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, url, link_type, price_selector, country_code, product_name FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var linkDescs []link.LinkDescription
	for rows.Next() {
		var linkDesc link.LinkDescription
		err := rows.Scan(&linkDesc.ID, &linkDesc.URL, &linkDesc.LinkType, &linkDesc.PriceSelector, &linkDesc.CountryCode, &linkDesc.ProductName)
		if err != nil {
			return nil, err
		}
		linkDescs = append(linkDescs, linkDesc)
	}
	return linkDescs, nil
}
