package editline_test

import (
	"bytes"
	"testing"

	"github.com/craiggwilson/editline/pkg/editline"
)

func TestWriter(t *testing.T) {
	testCases := []struct {
		input   string
		editors []editline.Editor
		output  string
	}{
		{
			input:  "one two three",
			output: "one two three",
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
			output: "one\r\ntwo\nthree",
		},
		{
			input: "one\r\ntwo\nthree",
			editors: []editline.Editor{
				editline.Remove(),
			},
			output: "",
		},
		{
			input: "one\r\ntwo\nthree",
			editors: []editline.Editor{
				editline.ReplaceLiteral("yes"),
				editline.ReplaceLiteral("no"),
			},
			output: "no\r\nno\nno",
		},
		{
			input: "one\r\ntwo\nthree",
			editors: []editline.Editor{
				editline.Prefix("t", editline.Remove()),
			},
			output: "one\r\n",
		},
		{
			input: "one\r\ntwo\nthree",
			editors: []editline.Editor{
				editline.ReplaceRegexpString("o", "ee"),
				editline.ReplaceRegexpString("ee", "oo"),
			},
			output: "oone\r\ntwoo\nthroo",
		},
		{
			input: "one\r\ntwo\nthree",
			editors: []editline.Editor{
				editline.RegexpString("^t", editline.EditorFunc(func(line string) (string, editline.Action) {
					return "yeah", editline.ReplaceAction
				})),
			},
			output: "one\r\nyeah\nyeah",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input+"_"+tc.output, func(t *testing.T) {
			var out bytes.Buffer
			w := editline.NewWriter(&out, tc.editors...)
			n, err := w.Write([]byte(tc.input))
			if err != nil {
				t.Fatalf("got an error writing: %v", err)
			}

			err = w.Flush()
			if err != nil {
				t.Fatalf("got an error flushing: %v", err)
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
