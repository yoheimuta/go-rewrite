package rewrite

import "os"

// Config configures how to rewrite.
type Config struct {
	dryrun   bool
	infoFile *os.File
	errFile  *os.File
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

// WithInfoFile is the option to set infoFile.
// infoFile is used to log the information. Default is os.Stdout.
func WithInfoFile(infoFile *os.File) ConfigOption {
	return func(c *Config) {
		c.infoFile = infoFile
	}
}

// WithErrFile is the option to set errFile.
// errFile is used to log the error. Default is os.Stderr.
func WithErrFile(errFile *os.File) ConfigOption {
	return func(c *Config) {
		c.errFile = errFile
	}
}

func newConfig(opts ...ConfigOption) *Config {
	config := &Config{
		dryrun:   false,
		infoFile: os.Stdout,
		errFile:  os.Stderr,
	}

	for _, opt := range opts {
		opt(config)
	}
	return config
}
