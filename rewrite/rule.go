package rewrite

// Rule is the rule for rewrite.
type Rule interface {
	// Filter filters the file using the filepath.
	Filter(filepath string) (isFilter bool, err error)
	// Mapping maps the file content with new one.
	Mapping(content []byte) (newContent []byte, isChanged bool, err error)
	// Output writes the content to the file.
	Output(filepath string, content []byte) error
}
