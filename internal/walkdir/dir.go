package walkdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Run walks the directory concurrently.
func Run(dir string, n *sync.WaitGroup, files chan<- string, errs chan<- error) {
	defer n.Done()

	entries, err := dirents(dir)
	if err != nil {
		errs <- err
		return
	}

	for _, entry := range entries {
		p := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			n.Add(1)
			go Run(p, n, files, errs)
		} else {
			files <- p
		}
	}
}

func dirents(dir string) ([]os.FileInfo, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
