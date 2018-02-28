//
// automated test services
//
// @author darryl.west <darwest@ebay.com>
// @created 2018-02-27 16:17:24

package main

import "./proxy"

func main() {
	hub.CreateLogger()
	cfg := hub.ParseArgs()
	if cfg == nil {
		hub.ShowHelp()
		return
	}

	service, err := hub.NewService(cfg)
	if err != nil {
		panic(err)
	}

	err = service.Start()
	if err != nil {
		println(err.Error())
	}
}
