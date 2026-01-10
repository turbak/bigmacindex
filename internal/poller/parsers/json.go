package parsers

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/oliveagle/jsonpath"
)

type JSONParser struct{}

func (p JSONParser) ParsePriceStringFromReader(reader io.Reader, priceSelector string) (string, error) {
	path, err := jsonpath.Compile(priceSelector)
	if err != nil {
		return "", err
	}

	var data any
	err = json.NewDecoder(reader).Decode(&data)
	if err != nil {
		return "", err
	}

	result, err := path.Lookup(data)
	if err != nil {
		return "", err
	}

	if str, ok := result.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("expected a string but got %v of type %T", result, result)
}
