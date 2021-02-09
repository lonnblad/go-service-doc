package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stroem/go-service-doc/utils"
)

var camelCaseTCs = []tc{
	{input: "a a", expected: "aA"},
	{input: "a-a", expected: "aA"},
	{input: "a_a", expected: "aA"},
	{input: "aA", expected: "aA"},
}

var kebabCaseTCs = []tc{
	{input: "a a", expected: "a-a"},
	{input: "a-a", expected: "a-a"},
	{input: "a_a", expected: "a-a"},
	{input: "aA", expected: "a-a"},
}

type tc struct {
	name     string
	input    string
	expected string
}

func Test_ConvertToCamelCase(t *testing.T) {
	for _, tc := range camelCaseTCs {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.ConvertToCamelCase(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ConvertToKebabCase(t *testing.T) {
	for _, tc := range kebabCaseTCs {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.ConvertToKebabCase(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
