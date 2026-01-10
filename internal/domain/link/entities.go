package link

type ID int32

type LinkType string

const (
	LinkTypeJSON  LinkType = "json"
	LinkTypeHTML  LinkType = "html"
	LinkTypeRegex LinkType = "regex"
)

type LinkDescription struct {
	ID            ID       `db:"id"`
	URL           string   `db:"url"`
	ProductName   string   `db:"product_name"`
	LinkType      LinkType `db:"link_type"`
	PriceSelector string   `db:"price_selector"`
	CountryCode   string   `db:"country_code"`
}
