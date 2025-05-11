package app

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func (a *Abdd) ExtractData(t *Test) error {
	if a.LastResponse == nil {
		return fmt.Errorf("no response to extract data from")
	}

	if t.Extract == nil {
		return nil
	}

	for _, ex := range t.Extract {
		if ex.Path == "" {
			return fmt.Errorf("%w: extraction path cannot be empty", ErrExtractionPathNotFound)
		}
		if ex.As == "" {
			return fmt.Errorf("%w: extraction variable name cannot be empty", ErrExtractionVariableNameEmpty)
		}

		value := gjson.Get(*a.LastResponse.Body, ex.Path)
		if !value.Exists() {
			return fmt.Errorf("%w: expected %s to be present", ErrExtractionPathNotFound, ex.Path)
		}

		switch value.Type {
		case gjson.String:
			a.Store[ex.As] = value.String()
		case gjson.Number:
			a.Store[ex.As] = value.Float()
		case gjson.True, gjson.False:
			a.Store[ex.As] = value.Bool()
		default:
			a.Store[ex.As] = value.Raw
		}

		if value.IsObject() || value.IsArray() || value.IsBool() || value.Type == gjson.Number {
			a.Store[ex.As] = value.Raw
		} else {
			a.Store[ex.As] = value.String()
		}
	}

	return nil
}
