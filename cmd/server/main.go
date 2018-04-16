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
	addr := ":80"
	app, err := app.New(addr, app.Options{
		Origins:     []string{"*"},
		AuthService: "127.0.0.1:8000",
		RedisHost:   "127.0.0.1:6379",
		RedisPass:   "daniyar",
	})
	if err != nil {
		util.Exit("%s", err.Error())
	}
	if err = app.Serve(); err != nil {
		util.Exit("%s", err.Error())
	}
}
