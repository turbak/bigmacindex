package poller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/turbak/bigmacindex/internal/domain/link"
	"github.com/turbak/bigmacindex/internal/domain/price"
	"github.com/turbak/bigmacindex/internal/poller/parsers"
)

type LinksLister interface {
	ListLinks(ctx context.Context) ([]link.LinkDescription, error)
}

type PricesUpserter interface {
	UpsertPrice(ctx context.Context, priceRec price.PriceRecord) error
}

type poller struct {
	linksLister    LinksLister
	pricesUpserter PricesUpserter
	httpClient     *http.Client
	parsersByType  map[link.LinkType]Parser
}

func NewPoller(linksLister LinksLister, pricesUpserter PricesUpserter) *poller {
	return &poller{
		linksLister:    linksLister,
		pricesUpserter: pricesUpserter,
		httpClient:     &http.Client{},
		parsersByType: map[link.LinkType]Parser{
			link.LinkTypeHTML:  parsers.HTMLParser{},
			link.LinkTypeJSON:  parsers.JSONParser{},
			link.LinkTypeRegex: parsers.RegexParser{},
		},
	}
}

func (p *poller) Poll(ctx context.Context) error {
	links, err := p.linksLister.ListLinks(ctx)
	if err != nil {
		return err
	}

	for _, linkDesc := range links {
		priceRec, err := p.fetchPriceData(ctx, linkDesc)
		if err != nil {
			return err
		}

		log.Printf("Fetched price for %s: %d.%02d %s", priceRec.ProductName, priceRec.Price, priceRec.PriceCents, priceRec.Currency)

		err = p.pricesUpserter.UpsertPrice(ctx, priceRec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *poller) fetchPriceData(ctx context.Context, linkDesc link.LinkDescription) (price.PriceRecord, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkDesc.URL, nil)
	if err != nil {
		return price.PriceRecord{}, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return price.PriceRecord{}, err
	}
	defer resp.Body.Close()

	parser, ok := p.parsersByType[linkDesc.LinkType]
	if !ok {
		return price.PriceRecord{}, fmt.Errorf("no parser for link type %s", linkDesc.LinkType)
	}

	priceValueStr, err := parser.ParsePriceStringFromReader(resp.Body, linkDesc.PriceSelector)
	if err != nil {
		return price.PriceRecord{}, err
	}

	priceValueStr = sanitizePriceString(priceValueStr)
	priceComp := parsePriceString(priceValueStr)

	var priceInt, priceCents int
	if priceComp.IntegerPart != "" {
		priceInt, err = strconv.Atoi(priceComp.IntegerPart)
		if err != nil {
			return price.PriceRecord{}, fmt.Errorf("invalid integer part %s: %w", priceComp.IntegerPart, err)
		}
	}

	if priceComp.DecimalPart != "" {
		priceCents, err = strconv.Atoi(priceComp.DecimalPart)
		if err != nil {
			return price.PriceRecord{}, fmt.Errorf("invalid decimal part %s: %w", priceComp.DecimalPart, err)
		}
	}

	return price.PriceRecord{
		ProductName: linkDesc.ProductName,
		Price:       int32(priceInt),
		PriceCents:  int32(priceCents),
		Currency:    priceComp.Currency,
		CountryCode: linkDesc.CountryCode,
		CreatedDate: time.Now().Format(time.DateOnly),
	}, nil
}

func sanitizePriceString(priceStr string) string {
	var result strings.Builder
	for _, r := range priceStr {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func parsePriceString(priceStr string) priceComponent {
	var integerPart, decimalPart string

	var priceParts []string
	if strings.Contains(priceStr, ".") {
		priceParts = strings.SplitN(priceStr, ".", 2)
	} else if strings.Contains(priceStr, ",") {
		priceParts = strings.SplitN(priceStr, ",", 2)
	} else {
		priceParts = []string{priceStr}
	}

	integerPart = priceParts[0]
	if len(priceParts) == 2 {
		decimalPart = priceParts[1]
	}

	return priceComponent{
		IntegerPart: integerPart,
		DecimalPart: decimalPart,
	}
}

type priceComponent struct {
	IntegerPart string
	DecimalPart string
	Currency    string
}
