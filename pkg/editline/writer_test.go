package editline_test

import (
	"bytes"
	"testing"

	"github.com/craiggwilson/editline/pkg/editline"
)

func TestWriter(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "one two three",
			output: "",
		},
		{
			input:  "one two three\n",
			output: "one two three\n",
		},
		{
			input:  "one two three\r\n",
			output: "one two three\r\n",
		},
		{
			input:  "one\r\ntwo\nthree",
			output: "one\r\ntwo\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input+"_"+tc.output, func(t *testing.T) {
			var out bytes.Buffer
			w := editline.NewWriter(&out)
			n, err := w.Write([]byte(tc.input))
			if err != nil {
				t.Fatalf("got an error writing input: %v", err)
			}

			if n != len(tc.input) {
				t.Fatalf("expected %d bytes to be written, but got %d", len(tc.input), n)
			}

			if out.String() != tc.output {
				t.Fatalf("expected output to be %q, but was %q", tc.output, out.String())
			}
		})
	}
}
