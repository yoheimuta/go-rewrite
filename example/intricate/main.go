package main

import (
	"flag"

	"github.com/yoheimuta/go-rewrite/example/intricate/myrule"
	"github.com/yoheimuta/go-rewrite/rewrite"
)

var (
	root   = flag.String("root", ".", "root path")
	dryrun = flag.Bool("dryrun", true, "the flag whether to overwrite")
)

func main() {
	flag.Parse()

	rule := &myrule.Rule{}
	config := rewrite.Config{
		Dryrun: *dryrun,
	}
	rewrite.Run(*root, config, rule)
}
