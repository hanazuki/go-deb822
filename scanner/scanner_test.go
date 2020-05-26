package scanner

import (
	"strings"
	"testing"
)

func expect(t *testing.T, scanner *Scanner, expected Line) {
	actual, err := scanner.Next()
	if err != nil {
		t.Errorf("Scanner error: %s", err.Error())
	}

	if actual == nil {
		t.Errorf("Unexpected eof")
		return
	}

	if *actual != expected {
		t.Errorf("Unexpected result. Expected: %v, actual: %v", expected, actual)
	}
}

func expectEof(t *testing.T, scanner *Scanner) {
	actual, err := scanner.Next()
	if err != nil {
		t.Errorf("Scanner error: %s", err.Error())
	}

	if actual != nil {
		t.Errorf("Expected eof: %v", actual)
	}
}

func TestScannerSingle(t *testing.T) {
	s := `Package: pui
Version: 1.1
Description: a 
  b
 c 
`

	scanner := New(strings.NewReader(s))

	expect(t, scanner, Line{Type: T_FIELD, Name: "Package", Value: "pui"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Version", Value: "1.1"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Description", Value: "a"})
	expect(t, scanner, Line{Type: T_CONT, Value: "b"})
	expect(t, scanner, Line{Type: T_CONT, Value: "c"})

	expectEof(t, scanner)
}

func TestScannerMulti(t *testing.T) {
	s := `Package: pui
Version: 1.1

Package: pui
Version: 1.2
`

	scanner := New(strings.NewReader(s))

	expect(t, scanner, Line{Type: T_FIELD, Name: "Package", Value: "pui"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Version", Value: "1.1"})
	expect(t, scanner, Line{Type: T_SEP})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Package", Value: "pui"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Version", Value: "1.2"})

	expectEof(t, scanner)
}

func TestScannerComment(t *testing.T) {
	s := `#beginning of file
Package: pui
#in paragraph
Version: 1.1
#end of paragraph

#between paragraphs

Package: pui
Version: 1.2
#end of file
`

	scanner := New(strings.NewReader(s))

	expect(t, scanner, Line{Type: T_COMMENT, Value: "beginning of file"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Package", Value: "pui"})
	expect(t, scanner, Line{Type: T_COMMENT, Value: "in paragraph"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Version", Value: "1.1"})
	expect(t, scanner, Line{Type: T_COMMENT, Value: "end of paragraph"})
	expect(t, scanner, Line{Type: T_SEP})
	expect(t, scanner, Line{Type: T_COMMENT, Value: "between paragraphs"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Package", Value: "pui"})
	expect(t, scanner, Line{Type: T_FIELD, Name: "Version", Value: "1.2"})
	expect(t, scanner, Line{Type: T_COMMENT, Value: "end of file"})

	expectEof(t, scanner)
}
