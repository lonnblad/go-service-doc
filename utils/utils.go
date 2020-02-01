package utils

import "strings"

import "regexp"

func ConvertToCamelCase(str string) string {
	str = convertSpaceCaseToCamel(str)
	str = convertKebabCaseToCamel(str)
	str = convertSnakeCaseToCamel(str)
	return str
}

func convertSpaceCaseToCamel(str string) string {
	return convertToCamelCase(str, " ")
}

func convertKebabCaseToCamel(str string) string {
	return convertToCamelCase(str, "-")
}

func convertSnakeCaseToCamel(str string) string {
	return convertToCamelCase(str, "_")
}

func convertToCamelCase(str, fromCase string) string {
	parts := strings.Split(str, fromCase)
	for idx, part := range parts {
		if idx == 0 {
			continue
		}
		parts[idx] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func ConvertToKebabCase(str string) string {
	str = convertSpaceCaseToKebab(str)
	str = convertCamelCaseToKebab(str)
	str = convertSnakeCaseToKebab(str)
	return str
}

func convertSpaceCaseToKebab(str string) string {
	return strings.ReplaceAll(str, " ", "-")
}

func convertCamelCaseToKebab(str string) string {
	reg := regexp.MustCompile(`([A-Z])`)
	str = reg.ReplaceAllString(str, "-$1")
	return strings.ToLower(str)
}

func convertSnakeCaseToKebab(str string) string {
	return strings.ReplaceAll(str, "_", "-")
}
