package rewrite

import (
	"io/ioutil"
	"sync"

	"github.com/yoheimuta/go-rewrite/internal/logger"
	"github.com/yoheimuta/go-rewrite/internal/walkdir"
)

// Run walks the rootPath and overwrites each file using the rule.
func Run(rootPath string, rule Rule, opts ...ConfigOption) {
	config := newConfig(opts...)
	logger := logger.NewClient(config.infoFile, config.errFile)
	rewriter := newRewriter(config, logger)

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
			err := rewriter.run(file, rule)
			if err != nil {
				logger.Errorf("rewrite: %v\n", err)
			}
		}
	}
}

type rewriter struct {
	c *Config
	l *logger.Client
}

func newRewriter(config *Config, logger *logger.Client) *rewriter {
	return &rewriter{
		c: config,
		l: logger,
	}
}

func (r *rewriter) run(filepath string, rule Rule) error {
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

	if r.c.dryrun {
		r.l.Infof("dryrun: %s\n", filepath)
	} else {
		err = rule.Output(filepath, newContent)
		if err != nil {
			return err
		}
		r.l.Infof("overwrite: %s\n", filepath)
	}
	return nil
}
