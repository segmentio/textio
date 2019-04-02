package textio

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTreeWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewTreeWriter(b)
	w.WriteString(".")

	w1 := NewTreeWriter(w)
	w2 := NewTreeWriter(w)
	w3 := NewTreeWriter(w)

	w1.WriteString("A")
	w2.WriteString("B\n(newline)\n")
	w3.WriteString("C")

	w21 := NewTreeWriter(w2)
	w21.WriteString("Hello World!")

	w.Close()

	const expected = `.
├── A
├── B
│   (newline)
│   └── Hello World!
└── C`

	found := b.String()

	if expected != found {
		t.Error("content mismatch")
		t.Logf("expected: %s", expected)
		t.Logf("found:    %s", found)
	}

}

func TestTreeWriterBase(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewTreeWriter(b)
	w1 := NewTreeWriter(w)

	if x := Base(w1); x != w {
		t.Error("Base(w1) != w")
	}

	if x := Base(w); x != b {
		t.Error("Base(w) != b")
	}

	if x := Base(b); x != b {
		t.Error("Base(b) != b")
	}
}

func TestTreeWriterParent(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewTreeWriter(b)
	w1 := NewTreeWriter(w)

	if x := Parent(w1); x != w {
		t.Error("Parent(w1) != w")
	}

	if x := Parent(w); x != b {
		t.Error("Parent(w) != w")
	}

	if x := Parent(b); x != b {
		t.Error("Parent(b) != b")
	}
}

func TestTreeWriterRoot(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewTreeWriter(b)
	w1 := NewTreeWriter(w)

	if x := Root(w1); x != w {
		t.Error("Root(w1) != w")
	}

	if x := Root(w); x != w {
		t.Error("Root(w) != w")
	}

	if x := Root(b); x != b {
		t.Error("Root(b) != b")
	}
}

func ExampleNewTreeWriter() {
	var ls func(io.Writer, string)

	ls = func(w io.Writer, path string) {
		tree := NewTreeWriter(w)
		tree.WriteString(filepath.Base(path))
		defer tree.Close()

		files, _ := ioutil.ReadDir(path)

		for _, f := range files {
			if f.Mode().IsDir() {
				ls(tree, filepath.Join(path, f.Name()))
			}
		}

		for _, f := range files {
			if !f.Mode().IsDir() {
				io.WriteString(NewTreeWriter(tree), f.Name())
			}
		}
	}

	ls(os.Stdout, "examples")
	// Output: examples
	// ├── A
	// │   ├── 1
	// │   └── 2
	// └── message
}
