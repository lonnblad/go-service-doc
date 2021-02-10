package utils_test

import (
	"testing"

	"github.com/lonnblad/go-service-doc/utils"
	"github.com/stretchr/testify/assert"
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
		testcase := tc

		t.Run(testcase.name, func(t *testing.T) {
			actual := utils.ConvertToCamelCase(testcase.input)
			assert.Equal(t, testcase.expected, actual)
		})
	}
}

func Test_ConvertToKebabCase(t *testing.T) {
	for _, tc := range kebabCaseTCs {
		testcase := tc

		t.Run(testcase.name, func(t *testing.T) {
			actual := utils.ConvertToKebabCase(testcase.input)
			assert.Equal(t, testcase.expected, actual)
		})
	}
}
