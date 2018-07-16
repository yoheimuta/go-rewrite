package rewrite

import (
	"io"
	"os"
)

// Config configures how to rewrite.
type Config struct {
	dryrun      bool
	infoWriter  io.Writer
	errWriter   io.Writer
	concurrency int
}

// ConfigOption is used to set the argument to the config.
type ConfigOption func(*Config)

// WithDryrun is the option to set dryrun.
// dryrun is used whether to skip output. Default is false.
func WithDryrun(dryrun bool) ConfigOption {
	return func(c *Config) {
		c.dryrun = dryrun
	}
}

// WithInfoWriter is the option to set infoWriter.
// infoWriter is used to log the information. Default is os.Stdout.
func WithInfoWriter(infoWriter io.Writer) ConfigOption {
	return func(c *Config) {
		c.infoWriter = infoWriter
	}
}

// WithErrWriter is the option to set errWriter.
// errWriter is used to log the error. Default is os.Stderr.
func WithErrWriter(errWriter io.Writer) ConfigOption {
	return func(c *Config) {
		c.errWriter = errWriter
	}
}

// WithConcurrency is the option to set concurrency.
// concurrency is used to determine the number of rewrites at the same time. Default is 10.
func WithConcurrency(concurrency int) ConfigOption {
	return func(c *Config) {
		c.concurrency = concurrency
	}
}

func newConfig(opts ...ConfigOption) *Config {
	config := &Config{
		dryrun:      false,
		infoWriter:  os.Stdout,
		errWriter:   os.Stderr,
		concurrency: 10,
	}

	for _, opt := range opts {
		opt(config)
	}
	return config
}
