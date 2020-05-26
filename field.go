package deb822

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	FieldNamePattern = `(?:[!"$-,.-9;-~][!-9;-~]*)`
)

var FieldNameRegepx = regexp.MustCompile(`\A` + FieldNamePattern + `\z`)

type Field struct {
	Name  string
	Value string
}

var emptyLineRegexp = regexp.MustCompile(`\n(\n|\z)`)

func NewField(name, value string) (Field, error) {
	if !FieldNameRegepx.MatchString(name) {
		return Field{}, fmt.Errorf("Invalid field name: %s", name)
	}

	if emptyLineRegexp.MatchString(value) {
		return Field{}, fmt.Errorf("Field value contains empty line")
	}

	return Field{
		Name:  name,
		Value: value,
	}, nil
}

func MustNewField(name, value string) Field {
	field, err := NewField(name, value)
	if err != nil {
		panic(err)
	}

	return field
}

func EscapeFieldValue(value string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		if line == "" {
			lines[i] = "."
		} else if line[0] == '.' {
			lines[i] = "." + line
		}
	}
	return strings.Join(lines, "\n")
}

func UnescapeFieldValue(value string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		if line[0] == '.' {
			lines[i] = line[1:]
		}
	}
	return strings.Join(lines, "\n")
}
