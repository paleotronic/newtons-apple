package main

import (
	"fmt"
	"log"

	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

type caller struct{}

func (c caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	oi.LongWrite(w, []byte{0x01, 0x01, 0x00, 0x00})

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	for {
		n, err := r.Read(p)

		if n > 0 {
			log.Printf("ClientTest: received %d bytes ($%x)", n, byte(p[0]))
		}

		if nil != err {
			break
		}
	}
}

func main() {
	fmt.Printf("Dial to %s:%d\n", "localhost", 5555)
	err := telnet.DialToAndCall(fmt.Sprintf("%s:%d", "localhost", 5555), caller{})

	if err != nil {
		log.Fatal(err)
	}
}
