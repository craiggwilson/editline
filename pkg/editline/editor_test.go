package editline_test

import (
	"testing"

	"github.com/craiggwilson/editline/pkg/editline"
)

func TestRegexp_Edit(t *testing.T) {
	testCases := []struct {
		pattern string
		input   string
		output  string
	}{
		{
			pattern: ".*funny.*",
			input:   "life is funny",
			output:  "YAY!",
		},
		{
			pattern: ".*funny.*",
			input:   "life is funn",
			output:  "life is funn",
		},
		{
			pattern: "(?m)^life",
			input:   "life is funny",
			output:  "YAY!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			editor := editline.RegexpString(tc.pattern, editline.ReplaceLiteral("YAY!"))
			output, _ := editor.Edit(tc.input)
			if output != tc.output {
				t.Fatalf("expected output %q, but got %q", tc.output, output)
			}
		})
	}
}

func TestRegexp_Prefix(t *testing.T) {
	testCases := []struct {
		pattern string
		prefix  string
	}{
		{
			pattern: ".*funny.*",
			prefix:  "",
		},
		{
			pattern: "life",
			prefix:  "",
		},
		{
			pattern: "^life",
			prefix:  "life",
		},
		{
			pattern: "(?m)^life",
			prefix:  "life",
		},
		{
			pattern: "(?i)^life",
			prefix:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			editor := editline.RegexpString(tc.pattern, editline.Remove())
			prefix := editor.(editline.Prefixer).Prefix()
			if prefix != tc.prefix {
				t.Fatalf("expected prefix %q, but got %q", tc.prefix, prefix)
			}
		})
	}
}
