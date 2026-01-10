package parsers

import (
	"fmt"
	"io"
	"regexp"
)

type RegexParser struct{}

func (p RegexParser) ParsePriceStringFromReader(reader io.Reader, priceSelector string) (string, error) {
	re, err := regexp.Compile(priceSelector)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	matches := re.FindSubmatch(data)
	if len(matches) >= 1 {
		return string(matches[1]), nil
	}

	return "", fmt.Errorf("no price found matching regex '%s'", priceSelector)
}
