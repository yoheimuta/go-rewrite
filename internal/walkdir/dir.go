package walkdir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Run walks the directory concurrently.
func Run(dir string, n *sync.WaitGroup, files chan<- string) {
	defer n.Done()

	for _, entry := range dirents(dir) {
		p := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			n.Add(1)
			go Run(p, n, files)
		} else {
			files <- p
		}
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dirents: %v\n", err)
		return nil
	}
	return entries
}
