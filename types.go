package deb822

const (
	FieldNamePattern = "(?:[!\"$-,.-9;-~][!-9;-~]*)"
)

type Field struct {
	Name  string
	Value string
}
