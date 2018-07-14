package rewrite

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/yoheimuta/go-rewrite/internal/walkdir"
)

// Rule is the rule for rewrite.
type Rule interface {
	// Filter filters the file using the filepath.
	Filter(filepath string) (isFilter bool, err error)
	// Mapping maps the file content with new one.
	Mapping(content []byte) (newContent []byte, isChanged bool, err error)
	// Output writes the content to the file.
	Output(filepath string, content []byte) error
}

// Run walks the rootPath and overwrites each file using the rule.
func Run(rootPath string, rule Rule) {
	files := make(chan string)
	var n sync.WaitGroup
	n.Add(1)
	go walkdir.Run(rootPath, &n, files)

	go func() {
		n.Wait()
		close(files)
	}()

loop:
	for {
		select {
		case file, ok := <-files:
			if !ok {
				break loop
			}
			err := rewrite(file, rule)
			if err != nil {
				fmt.Fprintf(os.Stderr, "process: %v\n", err)
			}
		}
	}
}

func rewrite(filepath string, rule Rule) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	ok, err := rule.Filter(filepath)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	newContent, changed, err := rule.Mapping(content)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	err = rule.Output(filepath, newContent)
	if err != nil {
		return err
	}
	fmt.Printf("overwrite: %s\n", filepath)
	return nil
}
