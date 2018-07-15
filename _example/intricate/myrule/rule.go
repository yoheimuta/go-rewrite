package myrule

import (
	"bytes"
	"io/ioutil"
	"strings"
)

func containsMulti(b []byte, subslices [][]byte) bool {
	for _, subslice := range subslices {
		if bytes.Contains(b, subslice) {
			return true
		}
	}
	return false
}

func mappingRecursive(n int, lines [][]byte) [][]byte {
	beforeLine := lines[n]
	if bytes.Contains(beforeLine, []byte("/// ")) {
		return lines
	}

	if bytes.Contains(beforeLine, []byte("// ")) {
		beforeLine = bytes.Replace(beforeLine, []byte("// "), []byte("/// "), 1)
	}

	if bytes.Equal(lines[n], beforeLine) {
		return lines
	}

	lines[n] = beforeLine
	return mappingRecursive(n-1, lines)
}

// Rule implements the Rewrite.Rule interface.
type Rule struct{}

// Filter filters the file using the filepath.
func (*Rule) Filter(filepath string) (bool, error) {
	if !strings.HasSuffix(filepath, ".swift") {
		return false, nil
	}
	return true, nil
}

// Mapping maps the file content with new one.
func (*Rule) Mapping(content []byte) ([]byte, bool, error) {
	newLine := []byte("\n")
	lines := bytes.Split(content, newLine)
	for i, line := range lines {
		if containsMulti(line, [][]byte{[]byte("func"), []byte("class"), []byte("var")}) {
			lines = mappingRecursive(i-1, lines)
		}
	}
	newContent := bytes.Join(lines, newLine)
	return newContent, !bytes.Equal(content, newContent), nil
}

// Output writes the content to the file.
func (*Rule) Output(filepath string, content []byte) error {
	return ioutil.WriteFile(filepath, content, 0644)
}
