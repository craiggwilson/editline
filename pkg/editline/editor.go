package editline

// EditMode returns the mode that should be performed from the result
// of an edit.
type EditMode byte

const (
	// ReplaceEditMode indicates that the edit should be used to replace
	// the line.
	ReplaceEditMode EditMode = iota
	// RemoveEditMode indicates that the edit should be used to remove
	// the line.
	RemoveEditMode
)

// Editor is an editor for a line.
type Editor interface {
	// Edit edits and returns the provided line as well as a bool when, if true, would
	// remove the line from output.
	Edit(line string) (string, EditMode)
}

// PrefixEditor is an editor that provides the prefix it's looking for in order to allow
// the Writer to optimize which editors should be run on a given line.
type PrefixEditor interface {
	Prefix() string
}
