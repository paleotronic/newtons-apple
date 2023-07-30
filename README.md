# newtons-apple
A simple serial driven implementation of physics via an external computer and the apple ii.

This project contains multiple parts: 
* A `physics-server` in the `server` folder written in golang using chipmunk 2D.
* An assembly language API contained in `apple2\serialtest.s`
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

## building the assembly
Using merlin32:
```
cd apple2
merlin32 serialtest.s
```

