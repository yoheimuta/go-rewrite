package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/yoheimuta/go-rewrite/rewrite"
)

var (
	root   = flag.String("root", ".", "root path")
	dryrun = flag.Bool("dryrun", true, "the flag whether to overwrite")
)

func main() {
	flag.Parse()

	rule := &myrule{}
	rewrite.Run(*root, rule, rewrite.WithDryrun(*dryrun))
}

type myrule struct{}

// Filter filters the file using the filepath.
func (*myrule) Filter(filepath string) (bool, error) {
	if !strings.HasSuffix(filepath, ".txt") {
		return false, nil
	}
	return true, nil
}

// Mapping maps the file content with new one.
func (*myrule) Mapping(content []byte) ([]byte, bool, error) {
	content = bytes.Replace(content, []byte("hoge"), []byte("fuga"), -1)
	return content, true, nil
}

// Output writes the content to the file.
func (*myrule) Output(_ string, content []byte) error {
	fmt.Printf("%s", string(content))
	return nil
}
