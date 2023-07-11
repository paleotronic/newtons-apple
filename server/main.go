package main

import (
	"flag"
	"fmt"

	"github.com/reiver/go-telnet"
)

var (
	flPort = flag.String("port", "5555", "Port to run service on.")
)

func main() {

	flag.Parse()

	var handler telnet.Handler = PhysicsService
	
	err := telnet.ListenAndServe(fmt.Sprintf(":%s", *flPort), handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}
