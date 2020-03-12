package editline

import (
	"io"
	"strings"
)

// NewWriter makes a Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

// Writer wraps another io.Writer and adds editing capabilities
// on a line by line basis.
type Writer struct {
	w io.Writer

	buf []byte
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	_, err := w.w.Write(w.buf)
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
		start, finish := it.Range()

		line := it.Line()
		suffix := "\n"

		if strings.HasSuffix(line, "\r") {
			line = line[0 : len(line)-1]
			suffix = "\r\n"
		}

		_, err := io.WriteString(w.w, line+suffix)
		if err != nil {
			return nn, err
		}

		nn += finish - start - offset
		offset = 0
	}

	it.buf = it.Remaining()

	return len(p), nil
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

func (it *lineIter) Range() (int, int) {
	return it.start, it.finish
}

func (it *lineIter) Remaining() []byte {
	return it.buf[it.start:]
}

func (it *lineIter) Next() bool {
	it.line = ""
	for i := it.finish; i < len(it.buf); i++ {
		if it.buf[i] == '\n' {
			it.start, it.finish = it.finish, i+1
			it.line = string(it.buf[it.start : it.finish-1])
			return true
		}
	}

	return false
}
