package editline

import (
	"io"
	"strings"
)

// NewWriter makes a Writer.
func NewWriter(w io.Writer, editors ...Editor) *Writer {
	return &Writer{
		w:       w,
		editors: buildPrefixTrie(editors),
	}
}

// Writer wraps another io.Writer and adds editing capabilities
// on a line by line basis.
type Writer struct {
	w       io.Writer
	editors *prefixTrie

	buf []byte
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	var err error
	if len(w.buf) > 0 {
		line, remove := w.editLine(string(w.buf))
		if !remove {
			_, err = io.WriteString(w.w, line)
		}
	}
	w.buf = nil
	return err
}

// Write writes the provided data to the underlying io.Writer after
// passing it through the various line editors.
func (w *Writer) Write(p []byte) (int, error) {
	offset := len(w.buf)
	it := lineIter{buf: append(w.buf, p...)}

	nn := 0

	for it.Next() {
		line := it.Line()
		suffix := "\n"

		if strings.HasSuffix(line, "\r") {
			line = line[0 : len(line)-1]
			suffix = "\r\n"
		}

		line, remove := w.editLine(line)
		if remove {
			continue
		}

		_, err := io.WriteString(w.w, line+suffix)
		if err != nil {
			return nn, err
		}

		nn += it.Len() - offset
		offset = 0
	}

	w.buf = it.Remaining()

	return len(p), nil
}

func (w *Writer) editLine(line string) (string, bool) {
	editor := Combine(w.editors.Get(line)...)
	line, action := editor.Edit(line)
	return line, action == RemoveAction
}

type lineIter struct {
	buf []byte

	start  int
	finish int

	line string
}

func (it *lineIter) Line() string {
	return it.line
}

func (it *lineIter) Len() int {
	return it.finish - it.start
}

func (it *lineIter) Remaining() []byte {
	return it.buf[it.finish:]
}

func (it *lineIter) Next() bool {
	it.line = ""
	for i := it.finish; i < len(it.buf); i++ {
		if it.buf[i] == '\n' {
			it.start, it.finish = it.finish, i+1
			it.line = string(it.buf[it.start:i])
			return true
		}
	}

	return false
}
