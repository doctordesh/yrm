package yrm

import (
	"testing"

	check "gitlab.com/MaxIV/lib-maxiv-go-check"
)

var input = `
// Some comment
foo: 5
value: 0.0
bar:
	baz: 5
	bool:
		ways: "lorem ipsum"
	port: 5.5
host: "localhost"
`

func TestYrm(t *testing.T) {
	m, err := Parse(input)
	check.OK(t, err)

	exp := map[string]interface{}{
		"foo":   5,
		"value": 0.0,
		"bar": map[string]interface{}{
			"baz": 5,
			"bool": map[string]interface{}{
				"ways": "lorem ipsum",
			},
			"port": 5.5,
		},
		"host": "localhost",
	}
	check.Equals(t, exp, m)
}
