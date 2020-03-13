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

// Prefixer allows an editor to provide a prefix in order to allow
// the Writer to optimize which editors should be run on a given line.
type Prefixer interface {
	Prefix() string
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
	return &regexEditor{
		re:     re,
		editor: editor,
	}
}

// RegexpString is a helper for Regexp. It will panic if the pattern
// does not compile to a regexp.Regexp.
func RegexpString(pattern string, editor Editor) Editor {
	return &regexEditor{
		re:     regexp.MustCompile(pattern),
		editor: editor,
	}
}

type regexEditor struct {
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
	sre, _ := syntax.Parse(e.re.String(), syntax.Perl)
	return buildPrefixFromRegexp(sre.Simplify())
}

func buildPrefixFromRegexp(sre *syntax.Regexp) string {
	if sre.Op != syntax.OpConcat {
		return ""
	}

	println(sre.Flags)

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

// Replace returns an Editor that replaces a line with the provided string.
func Replace(replacement string) Editor {
	return replaceEditor(replacement)
}

// ReplaceEditor is an editor that always replaces a line with the replacement.
type replaceEditor string

// Edit implements the Editor interface.
func (e replaceEditor) Edit(string) (string, Action) {
	return string(e), ReplaceAction
}
