package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/reiver/go-telnet"
)

var (
	flTelnetPort = flag.String("telnet-port", "5555", "TCP Port to run service on.")
	flSerial     = flag.Bool("serial", false, "Run on serial.")
	flSerialPort = flag.String("serial-port", "/dev/pts/9", "Serial port to run service on.")
	flBaudRate   = flag.Int("baud-rate", 115200, "Baud rate")
	flStopBits   = flag.Int("stop-bits", 1, "Stop bits")
	flParity     = flag.String("parity", "N", "Parity.")
	flDataBits   = flag.Int("data-bits", 8, "Data bits.")
)

func main() {

	flag.Parse()

	if !*flSerial {
		var handler telnet.Handler = PhysicsService
		log.Printf("Starting Physics via TELNET (localhost:%s)", *flTelnetPort)
		err := telnet.ListenAndServe(fmt.Sprintf(":%s", *flTelnetPort), handler)
		if nil != err {
			log.Fatal(err)
		}
	} else {
		log.Printf("Starting Physics via SERIAL (%s -> %d %d%s%d)", *flSerialPort, *flBaudRate, *flDataBits, *flParity, *flStopBits)
		err := ListenAndServeSerial(
			context.Background(),
			PhysicsService,
			*flSerialPort,
			*flBaudRate,
			*flDataBits,
			*flStopBits,
			*flParity,
		)
		if nil != err {
			log.Fatal(err)
		}
	}
}
