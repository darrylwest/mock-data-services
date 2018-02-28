//
// automated test services
//
// @author darryl.west <darwest@ebay.com>
// @created 2018-02-27 16:17:24

package main

import "./proxy"

func main() {
	proxy.CreateLogger()
	cfg := proxy.ParseArgs()
	if cfg == nil {
		proxy.ShowHelp()
		return
	}

	service, err := proxy.NewService(cfg)
	if err != nil {
		panic(err)
	}

	err = service.Start()
	if err != nil {
		println(err.Error())
	}
}
