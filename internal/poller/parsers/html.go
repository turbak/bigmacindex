package parsers

import (
	"fmt"
	"io"

	"github.com/antchfx/htmlquery"
)

type HTMLParser struct{}

func (p HTMLParser) ParsePriceStringFromReader(reader io.Reader, priceSelector string) (string, error) {
	doc, err := htmlquery.Parse(reader)
	if err != nil {
		return "", err
	}

	node := htmlquery.FindOne(doc, priceSelector)
	if node == nil {
		return "", fmt.Errorf("no price found for XPath '%s'", priceSelector)
	}
	priceStr := htmlquery.InnerText(node)
	return priceStr, nil
}
