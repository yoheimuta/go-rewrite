package walkdir_test

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/yoheimuta/go-rewrite/internal/setting_test"
	"github.com/yoheimuta/go-rewrite/internal/walkdir"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		inputDir     string
		wantFiles    []string
		wantErrCount int
	}{
		{
			name:         "Empty input generates empty output and 1 error",
			inputDir:     "",
			wantFiles:    []string{},
			wantErrCount: 1,
		},
		{
			name:     "Test input directory has 4 files",
			inputDir: setting_test.TestDataPath("testdir"),
			wantFiles: []string{
				setting_test.TestDataPath("testdir", "test.pl"),
				setting_test.TestDataPath("testdir", "test.txt"),
				setting_test.TestDataPath("testdir", "test123.swift"),
				setting_test.TestDataPath("testdir", "insidedir", "test.swift"),
			},
		},
	}

	for _, test := range tests {
		files := make(chan string)
		errs := make(chan error)
		var n sync.WaitGroup
		n.Add(1)
		go walkdir.Run(test.inputDir, &n, files, errs)

		go func() {
			n.Wait()
			close(files)
			close(errs)
		}()

		outputFiles := []string{}
		outputErrs := []error{}

	loop:
		for {
			select {
			case file, ok := <-files:
				if !ok {
					continue
				}
				outputFiles = append(outputFiles, file)
			case err, ok := <-errs:
				if !ok {
					break loop
				}
				outputErrs = append(outputErrs, err)
			}
		}

		sort.Strings(outputFiles)
		sort.Strings(test.wantFiles)
		if !reflect.DeepEqual(outputFiles, test.wantFiles) {
			t.Errorf("[%s] got %v, but want %v", test.name, outputFiles, test.wantFiles)
		}
		if len(outputErrs) != test.wantErrCount {
			t.Errorf(`[%s] got %v, but want %v, error="%v"`, test.name, len(outputErrs), test.wantErrCount, outputErrs)
		}
	}
}
