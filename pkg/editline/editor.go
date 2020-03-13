package editline

import (
	"regexp"
	"regexp/syntax"
	"strings"
)

// Action returns the action that should be performed from the result
// of an edit.
type Action byte

const (
	// ReplaceEditMode indicates that the edit should be used to replace
	// the line.
	ReplaceAction Action = iota
	// RemoveAction indicates that the edit should be used to remove
	// the line.
	RemoveAction
)

// Editor is an editor for a line.
type Editor interface {
	// Edit edits and returns the provided line as well as a bool when, if true, would
	// remove the line from output.
	Edit(line string) (string, Action)
}

// EditorFunc is a functional implementation of an Editor.
type EditorFunc func(string) (string, Action)

// Edit implements the Editor interface.
func (f EditorFunc) Edit(line string) (string, Action) {
	return f(line)
}

// Prefixer allows an editor to provide a prefix in order to allow
// the Writer to optimize selecting the editors to run on a given line.
type Prefixer interface {
	Prefix() string
}

// Combine combines multiple editors together.
func Combine(editors ...Editor) Editor {
	return &combineEditor{editors: editors}
}

type combineEditor struct {
	editors []Editor
}

// Edit implements the Editor interface.
func (e *combineEditor) Edit(line string) (string, Action) {
	for _, editor := range e.editors {
		replacement, action := editor.Edit(line)
		switch action {
		case ReplaceAction:
			line = replacement
		case RemoveAction:
			return replacement, action
		default:
			panic("unsupported action")
		}
	}

	return line, ReplaceAction
}

// Prefix returns an Editor that conditionally executes based on whether
// the line has a prefix matching the specified prefix.
func Prefix(prefix string, editor Editor) Editor {
	return &prefixEditor{
		prefix: prefix,
		editor: editor,
	}
}

type prefixEditor struct {
	prefix string
	editor Editor
}

// Edit implements the Editor interface.
func (e *prefixEditor) Edit(line string) (string, Action) {
	if strings.HasPrefix(line, e.prefix) {
		return e.editor.Edit(line)
	}

	return line, ReplaceAction
}

// Prefix implements the PrefixEditor interface.
func (e *prefixEditor) Prefix() string {
	return e.prefix
}

// Regexp returns an Editor that conditionally executes based on whether
// the line matches the regular expression.
func Regexp(re *regexp.Regexp, editor Editor) Editor {
	sre, _ := syntax.Parse(re.String(), syntax.Perl)
	return &regexEditor{
		prefix: buildPrefixFromRegexp(sre),
		re:     re,
		editor: editor,
	}
}

// RegexpString is a helper for Regexp. It will panic if the pattern
// does not compile to a regexp.Regexp.
func RegexpString(pattern string, editor Editor) Editor {
	sre, _ := syntax.Parse(pattern, syntax.Perl)
	return &regexEditor{
		prefix: buildPrefixFromRegexp(sre),
		re:     regexp.MustCompile(pattern),
		editor: editor,
	}
}

type regexEditor struct {
	prefix string
	re     *regexp.Regexp
	editor Editor
}

// Edit implements the Editor interface.
func (e *regexEditor) Edit(line string) (string, Action) {
	if e.re.MatchString(line) {
		return e.editor.Edit(line)
	}

	return line, ReplaceAction
}

// Prefix implements the PrefixEditor interface.
func (e *regexEditor) Prefix() string {
	return e.prefix
}

// Replace returns an Editor that replaces a line with the provided string.
func Remove() Editor {
	return removeEditor{}
}

// removeEditor is an editor that always removes a line.
type removeEditor struct{}

// Edit implements the Editor interface.
func (removeEditor) Edit(line string) (string, Action) {
	return line, RemoveAction
}

// ReplaceLiteral returns an Editor that replaces a line with the provided string.
func ReplaceLiteral(replacement string) Editor {
	return replaceLiteralEditor(replacement)
}

// ReplaceEditor is an editor that always replaces a line with the replacement.
type replaceLiteralEditor string

// Edit implements the Editor interface.
func (e replaceLiteralEditor) Edit(string) (string, Action) {
	return string(e), ReplaceAction
}

// ReplaceLiteral returns an Editor that replaces a line with the provided string.
func ReplaceRegexp(re *regexp.Regexp, replacement string) Editor {
	sre, _ := syntax.Parse(re.String(), syntax.Perl)
	return &replaceRegexpEditor{
		prefix:      buildPrefixFromRegexp(sre),
		re:          re,
		replacement: replacement,
	}
}

// ReplaceLiteral returns an Editor that replaces a line with the provided string.
func ReplaceRegexpString(pattern string, replacement string) Editor {
	sre, _ := syntax.Parse(pattern, syntax.Perl)
	return &replaceRegexpEditor{
		prefix:      buildPrefixFromRegexp(sre),
		re:          regexp.MustCompile(pattern),
		replacement: replacement,
	}
}

// ReplaceEditor is an editor that always replaces a line with the replacement.
type replaceRegexpEditor struct {
	prefix      string
	re          *regexp.Regexp
	replacement string
}

// Edit implements the Editor interface.
func (e *replaceRegexpEditor) Edit(line string) (string, Action) {
	return e.re.ReplaceAllString(line, e.replacement), ReplaceAction
}

func (e *replaceRegexpEditor) Prefix() string {
	return e.prefix
}

func buildPrefixFromRegexp(sre *syntax.Regexp) string {
	if sre.Op != syntax.OpConcat {
		return ""
	}

	subs := sre.Sub
	if len(subs) <= 1 {
		return ""
	}

	if subs[0].Op != syntax.OpBeginLine && subs[0].Op != syntax.OpBeginText {
		return ""
	}

	if subs[1].Op != syntax.OpLiteral || (subs[1].Flags&syntax.FoldCase) != 0 {
		return ""
	}

	return subs[1].String()
}
