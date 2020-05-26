package parser

import (
	"strings"
	"testing"

	"github.com/hanazuki/go-deb822"
)

func expectFields(t *testing.T, paragraph *deb822.Paragraph, nvPairs ...string) {
	if len(nvPairs)%2 != 0 {
		panic("nvPairs must contain even number of strings")
	}

	if paragraph == nil {
		t.Fatalf("Paragraph expected")
	}

	if len(paragraph.Fields) != len(nvPairs)/2 {
		t.Fatalf("%d fields expected but actually %d", len(nvPairs)/2, len(paragraph.Fields))
	}

	for i, field := range paragraph.Fields {
		expectedName, expectedValue := nvPairs[i*2], nvPairs[i*2+1]

		if field.Name != expectedName {
			t.Errorf("Field name `%s' expected but actually `%s'", expectedName, field.Name)
		}
		if field.Value != expectedValue {
			t.Errorf("Field value `%s' expected but actually `%s'", expectedValue, field.Value)
		}
	}
}

func TestParser(t *testing.T) {
	s := `#comment
Package: pui
Version: 1.0
#comment
Description: a
 b
#comment

#comment

Package: pui
Version: 1.1
Description: a
 b
#comment
 .
 c

#comment
`

	parser := New(strings.NewReader(s))

	para, err := parser.NextParagraph()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	expectFields(t, para, "Package", "pui", "Version", "1.0", "Description", "a\nb")

	para, err = parser.NextParagraph()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	expectFields(t, para, "Package", "pui", "Version", "1.1", "Description", "a\nb\n.\nc")

	para, err = parser.NextParagraph()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if para != nil {
		t.Errorf("Unexpected paragraph")
	}
}
