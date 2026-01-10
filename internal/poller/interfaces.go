package poller

import (
	"io"
)

type Parser interface {
	ParsePriceStringFromReader(reader io.Reader, priceSelector string) (string, error)
}
