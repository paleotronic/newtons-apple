# newtons-apple
A simple serial driven implementation of physics via an external computer and the apple ii.

This project contains multiple parts: 
* A `physics-server` in the `server` folder written in golang using chipmunk 2D.
* An assembly language API contained in `apple2\physics.s`
* Some examples in applesoft contained in the `applesoft` folder.
 * `blocks.bas` - a falling boxes physics demo.
 * `thrust.bas` - a simple space ship flying demo.
 * `batgame.bas` - a simple pong style bat and ball game.

## building the server
With go installed:
```
cd server
go build .
```
There is a Makefile to build a totally static executable with `musl-gcc`.
```
cd server
make build
```

### Using the server

```
Usage of ./physics-server:
  -baud-rate int
        Baud rate (default 115200)
  -data-bits int
        Data bits (5,6,7 or 8). (default 8)
  -list-ports
        List serial ports and exit.
  -max-delta int
        Max memory delta size (bytes) (default 64)
  -parity string
        Parity (E=even,O=odd,M=mark,S=space,N=none). (default "N")
  -serial
        Run on serial.
  -serial-port string
        Serial port to run service on. (default "/dev/pts/9")
  -stop-bits string
        Stop bits (1,1.5,2) (default "1")
  -telnet-port string
        TCP Port to run service on. (default "5555")
```

It will run on telnet by default, which allows bootstrapping it with emulators with telnet serial emulation. 

On some OSes, you may need to use `sudo` to get access to the serial hardware. 

The service has been validated on an Apple //c with `19200` as the baud rate. (This was because our serial card has the crystal bug and isn't stable at 115200).  

You can change `CONTROLVAL` in the assembly to set a different baud rate.

## building the assembly
Using merlin32:
```
cd apple2
merlin32 physics.s
```
The code will by default assemble to `$2000`. This engine proof of concept uses lo-res for brevity, so this memory is free. 

