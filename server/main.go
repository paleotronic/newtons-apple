package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/reiver/go-telnet"
	"go.bug.st/serial"
)

var (
	flTelnetPort = flag.String("telnet-port", "5555", "TCP Port to run service on.")
	flSerial     = flag.Bool("serial", false, "Run on serial.")
	flSerialPort = flag.String("serial-port", "/dev/pts/9", "Serial port to run service on.")
	flBaudRate   = flag.Int("baud-rate", 115200, "Baud rate")
	flStopBits   = flag.String("stop-bits", "1", "Stop bits (1,1.5,2)")
	flParity     = flag.String("parity", "N", "Parity (E=even,O=odd,M=mark,S=space,N=none).")
	flDataBits   = flag.Int("data-bits", 8, "Data bits (5,6,7 or 8).")
	flListPorts  = flag.Bool("list-ports", false, "List serial ports and exit.")
)

func main() {

	flag.Parse()

	if *flListPorts {
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			fmt.Println("No serial ports found!")
		} else {
			for _, port := range ports {
				fmt.Printf("Found port: %v\n", port)
			}
		}
		return
	}

	if !*flSerial {
		var handler telnet.Handler = PhysicsService
		log.Printf("Starting Physics via TELNET (localhost:%s)", *flTelnetPort)
		err := telnet.ListenAndServe(fmt.Sprintf(":%s", *flTelnetPort), handler)
		if nil != err {
			log.Fatal(err)
		}
	} else {
		log.Printf("Starting Physics via SERIAL (%s -> %d %d%s%s)", *flSerialPort, *flBaudRate, *flDataBits, *flParity, *flStopBits)
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
