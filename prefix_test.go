package textio

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrefixWriter(t *testing.T) {
	b := &bytes.Buffer{}
	b.WriteByte('\n')

	w1 := NewPrefixWriter(b, "\t")
	w2 := NewPrefixWriter(w1, "\t- ")

	fmt.Fprint(w1, "hello:\n")

	fmt.Fprint(w2, "value: 1")
	fmt.Fprint(w2, "\n")
	fmt.Fprint(w2, "value: 2\nvalue: 3\n")

	w2.Flush()
	w1.Flush()

	const expected = `
	hello:
		- value: 1
		- value: 2
		- value: 3
`

	found := b.String()

	if expected != found {
		t.Error("content mismatch")
		t.Log("expected:", expected)
		t.Log("found:", found)
	}
}
