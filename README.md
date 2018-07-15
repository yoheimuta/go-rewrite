# go-rewrite [![GoDoc](https://godoc.org/github.com/yoheimuta/go-rewrite/rewrite?status.svg)](https://godoc.org/github.com/yoheimuta/go-rewrite/rewrite) [![Build Status](https://travis-ci.org/yoheimuta/go-rewrite.svg?branch=master)](https://travis-ci.org/yoheimuta/go-rewrite)

go-rewrite is a thin go package which helps replacing files.

You can use this package...

- To focus on coding to filter files and overwriting a content.
- To use an intricate searching instead of grep.
- To use an intricate overwriting instead of sed.
- To traverse a directory much faster with utilizing goroutines.

### Motivation

For example, if you want to replace comments on multiple lines of `func or class or var declaration`
from `//` to `///` for files with extension .swift.

It is easier to write Go's code than to do it with grep and sed.

### Installation

```
go get github.com/yoheimuta/go-rewrite
```

### Usage

See `_example/simple` and `_example/intricate` in detail.

```go
func main() {
	rule := &myrule{}
	rewrite.Run(".", rule)
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
	return ioutil.WriteFile(filepath, content, 0644)
}
```
