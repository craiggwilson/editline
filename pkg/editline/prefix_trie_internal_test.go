package editline

import (
	"reflect"
	"testing"
)

func TestBuildPrefixTrie(t *testing.T) {
	editors := []Editor{
		testPrefixEditor("aaa"),
		testPrefixEditor(""),
		testPrefixEditor("aab"),
		testPrefixEditor("baab"),
		testPrefixEditor("baabc"),
	}

	trie := buildPrefixTrie(editors)

	testCases := []struct {
		input  string
		result []Editor
	}{
		{
			input:  "aaacdsdffe",
			result: []Editor{editors[0], editors[1]},
		},
		{
			input:  "aa",
			result: []Editor{editors[1]},
		},
		{
			input:  "baabc",
			result: []Editor{editors[1], editors[3], editors[4]},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := trie.Get(tc.input)
			if !reflect.DeepEqual(result, tc.result) {
				t.Fatalf("expected editors %v, but got %v", tc.result, result)
			}
		})
	}
}

type testPrefixEditor string

func (e testPrefixEditor) Edit(line string) (string, EditMode) {
	panic("nah")
}

func (e testPrefixEditor) Prefix() string {
	return string(e)
}
