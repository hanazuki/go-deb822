package parser

import (
	"io"
	"strings"

	"github.com/hanazuki/go-deb822"
	"github.com/hanazuki/go-deb822/scanner"
)

type Parser struct {
	scanner *scanner.Scanner
}

func New(source io.Reader) Parser {
	return Parser{
		scanner: scanner.New(source),
	}
}

func (p *Parser) NextParagraph() (*deb822.Paragraph, error) {
	var (
		paragraph  *deb822.Paragraph
		fieldName  string
		fieldValue strings.Builder
	)

	appendField := func() error {
		if fieldName != "" {
			field, err := deb822.NewField(fieldName, fieldValue.String())
			if err != nil {
				return err
			}
			fieldName = ""

			if paragraph == nil {
				paragraph = deb822.MustNewParagraph()
			}
			paragraph.Add(field)
		}
		return nil
	}

Loop:
	for {
		line, err := p.scanner.Next()
		if err != nil {
			return nil, err
		}

		if line == nil {
			break Loop
		}

		switch line.Type {
		case scanner.T_SEP:
			break Loop

		case scanner.T_COMMENT:
			continue

		case scanner.T_FIELD:
			err = appendField()
			if err != nil {
				return nil, err
			}

			fieldName = line.Name
			fieldValue = strings.Builder{}
			_, err = fieldValue.WriteString(line.Value)
			if err != nil {
				return nil, err
			}

		case scanner.T_CONT:
			if fieldName == "" {
				panic("Unexpected T_CONT")
			}

			_, err = fieldValue.WriteRune('\n')
			if err != nil {
				return nil, err
			}

			_, err = fieldValue.WriteString(line.Value)
			if err != nil {
				return nil, err
			}
		}

	}

	err := appendField()
	if err != nil {
		return nil, err
	}

	return paragraph, nil
}

func (p *Parser) ConsumeAll() error {
	for {
		paragraph, err := p.NextParagraph()
		if err != nil {
			return err
		}

		if paragraph == nil {
			return nil
		}
	}
}
