package editline

import "sort"

func buildPrefixTrie(editors []Editor) *prefixTrie {
	root := &prefixTrie{}
	for i, editor := range editors {
		prefix := ""
		if pe, ok := editor.(Prefixer); ok {
			prefix = pe.Prefix()
		}

		root.add(prefix, prefixTrieItem{idx: i, editor: editor})
	}

	return root
}

type prefixTrieItem struct {
	idx    int
	editor Editor
}

type prefixTrie struct {
	items    []prefixTrieItem
	children map[byte]*prefixTrie
}

func (t *prefixTrie) Get(input string) []Editor {
	items := t.get(input)

	sort.Slice(items, func(i, j int) bool {
		return items[i].idx < items[j].idx
	})

	editors := make([]Editor, 0, len(items))
	for _, item := range items {
		editors = append(editors, item.editor)
	}

	return editors
}

func (t *prefixTrie) get(input string) []prefixTrieItem {
	items := t.items[:]

	if len(input) == 0 || t.children == nil {
		return items
	}

	child, ok := t.children[input[0]]
	if !ok {
		return items
	}

	return append(items, child.get(input[1:])...)
}

func (t *prefixTrie) add(prefix string, item prefixTrieItem) {
	if len(prefix) == 0 {
		t.items = append(t.items, item)
		return
	}

	if t.children == nil {
		t.children = make(map[byte]*prefixTrie)
	}

	trie, ok := t.children[prefix[0]]
	if !ok {
		trie = &prefixTrie{}
		t.children[prefix[0]] = trie
	}

	trie.add(prefix[1:], item)
}
