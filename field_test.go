package deb822

import "testing"

func TestEscapeFieldValue(t *testing.T) {
	test := func(input, expected string) {
		actual := EscapeFieldValue(input)

		if actual != expected {
			t.Errorf("Expected: `%s`, actual: `%s`", expected, actual)
		}
	}

	test("", "")
	test("\na", "\na")
	test("a", "a")
	test("a\n", "a\n.")
	test("a\nb", "a\nb")
	test("a\n\nb", "a\n.\nb")
	test("a\n.\nb", "a\n..\nb")
	test("a\n.x\nb", "a\n..x\nb")
}

func TestUnescapeFieldValue(t *testing.T) {
	test := func(input, expected string) {
		actual := UnescapeFieldValue(input)

		if actual != expected {
			t.Errorf("Expected: `%s`, actual: `%s`", expected, actual)
		}
	}

	test("", "")
	test(".", ".")
	test("a", "a")
	test(".a", ".a")
	test("a\nb", "a\nb")
	test("a\n.", "a\n")
	test("a\n.b", "a\nb")
	test("a\n..b", "a\n.b")
}
