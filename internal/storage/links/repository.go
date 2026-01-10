package links

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/turbak/bigmacindex/internal/domain/link"
)

const tableName = "links"

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
func (r *repository) DeleteLink(ctx context.Context, ID link.ID) error {
	_, err := r.db.Delete(tableName).Where(squirrel.Eq{"id": ID}).ExecContext(ctx)
	return err
}

func (r *repository) AddLink(ctx context.Context, linkDesc link.LinkDescription) (link.LinkDescription, error) {
	res, err := r.db.Insert(tableName).
		Columns("url", "link_type", "price_selector", "country_code", "product_name").
		Values(linkDesc.URL, linkDesc.LinkType, linkDesc.PriceSelector, linkDesc.CountryCode, linkDesc.ProductName).
		ExecContext(ctx)
	if err != nil {
		return link.LinkDescription{}, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return link.LinkDescription{}, err
	}
	linkDesc.ID = link.ID(ID)
	return linkDesc, nil
}

func (r *repository) ListLinks(ctx context.Context) ([]link.LinkDescription, error) {
	rows, err := r.db.Select("id", "url", "link_type", "price_selector", "country_code", "product_name").
		From(tableName).
		QueryContext(ctx)
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
