package logger

import (
	"fmt"
	"io"
	"sync"
)

// Client is a logger.
type Client struct {
	infom    sync.Mutex
	infoFile io.Writer

	errm    sync.Mutex
	errFile io.Writer
}

// NewClient generates a new Client.
func NewClient(infoFile, errFile io.Writer) *Client {
	return &Client{
		infoFile: infoFile,
		errFile:  errFile,
	}
}

// Infof formats according to a format specifier and writes to infoFile.
// A format specifier is same with fmt.Sprintf.
func (c *Client) Infof(format string, a ...interface{}) {
	c.infom.Lock()
	defer c.infom.Unlock()
	fmt.Fprintf(c.infoFile, format, a...)
}

// Errorf formats according to a format specifier and writes to errFile.
// A format specifier is same with fmt.Sprintf.
func (c *Client) Errorf(format string, a ...interface{}) {
	c.errm.Lock()
	defer c.errm.Unlock()
	fmt.Fprintf(c.errFile, format, a...)
}
