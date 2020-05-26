package scanner

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/hanazuki/go-deb822"
)

type LineType int

const (
	T_SEP LineType = iota
	T_FIELD
	T_CONT
	T_COMMENT
)

type Scanner struct {
	Reader      *bufio.Reader
	inParagraph bool
	position    int
}

type Line struct {
	Type  LineType
	Name  string
	Value string
}

type ScanError struct {
	Message string
	Source  string
	Line    int
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("%s at line %d", e.Message, e.Line)
}

func New(source io.Reader) *Scanner {
	return &Scanner{
		Reader: bufio.NewReader(source),
	}
}

var (
	RE_EMPTY   = regexp.MustCompile(`\A[ \t]*\z`)
	RE_COMMENT = regexp.MustCompile(`\A#(?P<value>.*)\z`)
	RE_FIELD   = regexp.MustCompile(`\A(?P<name>` + deb822.FieldNamePattern + `):[ \t]*(?P<value>.*)\z`)
	RE_CONT    = regexp.MustCompile(`\A[ \t](?P<value>.*)\z`)
)

func (s *Scanner) Next() (*Line, error) {
	for {
		line, err := s.Reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		if line[len(line)-1] != '\n' {
			return nil, &ScanError{Line: s.position, Message: "Unexpected EOF", Source: line}
		}
		line = line[:len(line)-1] // chomp NL

		s.position += 1

		if RE_EMPTY.MatchString(line) {
			if s.inParagraph {
				s.inParagraph = false
				return &Line{
					Type: T_SEP,
				}, nil
			}
			continue
		}

		if m := RE_COMMENT.FindStringSubmatch(line); m != nil {
			return &Line{
				Type:  T_COMMENT,
				Value: m[1],
			}, nil
		}

		if m := RE_FIELD.FindStringSubmatch(line); m != nil {
			s.inParagraph = true
			return &Line{
				Type:  T_FIELD,
				Name:  m[1],
				Value: m[2],
			}, nil
		}

		if m := RE_CONT.FindStringSubmatch(line); m != nil {
			if !s.inParagraph {
				return nil, &ScanError{Line: s.position, Message: "Unexpected continuation line", Source: line}
			}

			return &Line{
				Type:  T_CONT,
				Value: m[1],
			}, nil
		}

		return nil, &ScanError{Line: s.position, Message: "Invalid deb822", Source: line}
	}

}
