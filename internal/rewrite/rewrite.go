package rewrite

import (
	"fmt"
	"io/ioutil"
)

// Rule is the rule for rewrite.
type Rule interface {
	filter(file string) (isFilter bool, err error)
	mapping(content []byte) (newContent []byte, isChanged bool, err error)
	output(file string, content []byte) error
}

// Rewrite overwrites the file using the rule.
func Rewrite(file string, rule Rule) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	ok, err := rule.filter(file)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	newContent, changed, err := rule.mapping(content)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	err = rule.output(file, newContent)
	if err != nil {
		return err
	}
	fmt.Printf("overwrite: %s\n", file)
	return nil
}
