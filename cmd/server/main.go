package main

import (
	"flag"

	"github.com/danikarik/constantinople/pkg/app"
	"github.com/danikarik/constantinople/pkg/util"
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()
}

func main() {
	addr := ":3000"
	app, err := app.New(addr, app.Options{})
	if err != nil {
		util.Exit("[server] %s", err.Error())
	}
	if err = app.Serve(); err != nil {
		util.Exit("[server] %s", err.Error())
	}
}
