package rewrite_test

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/yoheimuta/go-rewrite/internal/setting_test"
	"github.com/yoheimuta/go-rewrite/rewrite"
)

type rule struct {
	filter  func(filepath string) (isFilter bool, err error)
	mapping func(content []byte) (newContent []byte, isChanged bool, err error)

	outputFiles    []string
	outputContents []string
}

func newRule(
	filter func(filepath string) (isFilter bool, err error),
	mapping func(content []byte) (newContent []byte, isChanged bool, err error),
) *rule {
	return &rule{
		filter:         filter,
		mapping:        mapping,
		outputFiles:    []string{},
		outputContents: []string{},
	}
}

func (r *rule) Filter(filepath string) (bool, error) {
	return r.filter(filepath)
}

func (r *rule) Mapping(content []byte) ([]byte, bool, error) {
	return r.mapping(content)
}

func (r *rule) Output(filepath string, content []byte) error {
	r.outputFiles = append(r.outputFiles, filepath)
	r.outputContents = append(r.outputContents, string(content))
	return nil
}

func TestRun(t *testing.T) {
	tests := []struct {
		name               string
		inputRoot          string
		inputFilter        func(string) (bool, error)
		inputMapping       func([]byte) ([]byte, bool, error)
		inputDryrun        bool
		wantOutputFiles    []string
		wantOutputContents []string
		wantInfoCount      int
		wantErrCount       int
	}{
		{
			name:               "Empty rootPath generates an error",
			wantOutputFiles:    []string{},
			wantOutputContents: []string{},
			wantErrCount:       1,
		},
		{
			name:      "Filters files",
			inputRoot: setting_test.TestDataPath("testdir"),
			inputFilter: func(filepath string) (bool, error) {
				if strings.HasSuffix(filepath, ".swift") {
					return true, nil
				}
				return false, nil
			},
			inputMapping: func(in []byte) ([]byte, bool, error) {
				return []byte("rewrite"), true, nil
			},
			wantOutputFiles: []string{
				setting_test.TestDataPath("testdir", "test123.swift"),
				setting_test.TestDataPath("testdir", "insidedir", "test.swift"),
			},
			wantOutputContents: []string{
				"rewrite",
				"rewrite",
			},
			wantInfoCount: 2,
		},
		{
			name:      "Skip to call Output because of dryrun",
			inputRoot: setting_test.TestDataPath("testdir"),
			inputFilter: func(filepath string) (bool, error) {
				return true, nil
			},
			inputMapping: func(in []byte) ([]byte, bool, error) {
				return in, true, nil
			},
			inputDryrun:        true,
			wantOutputFiles:    []string{},
			wantOutputContents: []string{},
			wantInfoCount:      4,
		},
	}

	for _, test := range tests {
		mockRule := newRule(
			test.inputFilter,
			test.inputMapping,
		)
		wantInfo := &bytes.Buffer{}
		wantErr := &bytes.Buffer{}

		rewrite.Run(
			test.inputRoot,
			mockRule,
			rewrite.WithDryrun(test.inputDryrun),
			rewrite.WithInfoWriter(wantInfo),
			rewrite.WithErrWriter(wantErr),
		)

		sort.Strings(mockRule.outputFiles)
		sort.Strings(test.wantOutputFiles)
		if !reflect.DeepEqual(mockRule.outputFiles, test.wantOutputFiles) {
			t.Errorf("[%s] got %v, but want %v", test.name, mockRule.outputFiles, test.wantOutputFiles)
		}

		sort.Strings(mockRule.outputContents)
		sort.Strings(test.wantOutputContents)
		if !reflect.DeepEqual(mockRule.outputContents, test.wantOutputContents) {
			t.Errorf("[%s] got %v, but want %v", test.name, mockRule.outputContents, test.wantOutputContents)
		}

		infoLogs := strings.Split(wantInfo.String(), "\n")
		infoLogLen := len(infoLogs) - 1
		if infoLogLen != test.wantInfoCount {
			t.Errorf(`[%s] got %d, but want %d, info "%v"`, test.name, infoLogLen, test.wantInfoCount, infoLogs)
		}

		errLogs := strings.Split(wantErr.String(), "\n")
		errLogLen := len(errLogs) - 1
		if errLogLen != test.wantErrCount {
			t.Errorf(`[%s] got %d, but want %d, err "%v"`, test.name, errLogLen, test.wantErrCount, errLogs)
		}
	}
}
