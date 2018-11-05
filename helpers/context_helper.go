package helpers

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gobuffalo/buffalo"
)

// Param - Gets a query parameter inside a buffalo Context
func Param(c buffalo.Context, param string) (string, error) {
	if m, ok := c.Params().(url.Values); ok {
		for k, v := range m {
			fmt.Println(k, v)
			if k == param && v != nil && len(v) > 0 {
				return v[0], nil
			}
		}
	}

	return "", errors.New("Unable to find parameter")
}
