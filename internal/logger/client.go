package logger

import (
	"fmt"
	"os"
)

// Client is a logger.
type Client struct {
	infoFile *os.File
	errFile  *os.File
}

// NewClient generates a new Client.
func NewClient(infoFile, errFile *os.File) *Client {
	return &Client{
		infoFile: infoFile,
		errFile:  errFile,
	}
}

// Infof formats according to a format specifier and writes to infoFile.
// A format specifier is same with fmt.Sprintf.
func (c *Client) Infof(format string, a ...interface{}) {
	fmt.Fprintf(c.infoFile, format, a...)
}

// Errorf formats according to a format specifier and writes to errFile.
// A format specifier is same with fmt.Sprintf.
func (c *Client) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(c.errFile, format, a...)
}
