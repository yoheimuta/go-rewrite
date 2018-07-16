package rewrite

import (
	"io/ioutil"
	"sync"

	"github.com/yoheimuta/go-rewrite/internal/logger"
	"github.com/yoheimuta/go-rewrite/internal/walkdir"
)

// Run walks the rootPath and overwrites each file using the rule.
func Run(rootPath string, rule Rule, opts ...ConfigOption) {
	files := make(chan string)
	errs := make(chan error)
	var n sync.WaitGroup
	n.Add(1)
	go walkdir.Run(rootPath, &n, files, errs)

	go func() {
		n.Wait()
		close(files)
		close(errs)
	}()

	config := newConfig(opts...)
	logger := logger.NewClient(config.infoWriter, config.errWriter)
	rewriter := newRewriter(config, logger)
	var w sync.WaitGroup
	for i := 0; i < config.concurrency; i++ {
		w.Add(1)
		go func() {
			defer w.Done()

			for {
				select {
				case file, ok := <-files:
					if !ok {
						continue
					}
					err := rewriter.run(file, rule)
					if err != nil {
						logger.Errorf("failed to rewrite %s: %v\n", file, err)
					}
				case err, ok := <-errs:
					if !ok {
						return
					}
					if err != nil {
						logger.Errorf("%v\n", err)
					}
				}
			}
		}()
	}
	w.Wait()
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
	ok, err := rule.Filter(filepath)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
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
