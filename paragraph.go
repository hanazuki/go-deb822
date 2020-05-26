package deb822

import (
	"fmt"
	"strings"
)

type Paragraph struct {
	Fields []Field
}

func NewParagraph(fields ...Field) (*Paragraph, error) {
	var paragraph Paragraph
	for _, field := range fields {
		err := paragraph.Add(field)
		if err != nil {
			return nil, err
		}
	}
	return &paragraph, nil
}

func MustNewParagraph(fields ...Field) *Paragraph {
	paragraph, err := NewParagraph(fields...)
	if err != nil {
		panic(err)
	}

	return paragraph
}

func (p *Paragraph) Find(name string) *Field {
	for i, field := range p.Fields {
		if strings.EqualFold(field.Name, name) {
			return &p.Fields[i]
		}
	}
	return nil
}

func (p *Paragraph) Add(field Field) error {
	if p.Find(field.Name) != nil {
		return fmt.Errorf("Duplicate field `%s'", field.Name)
	}
	p.Fields = append(p.Fields, field)
	return nil
}
