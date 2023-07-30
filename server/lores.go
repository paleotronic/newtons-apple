package main

import (
	"math"
	"sync"

	"github.com/jakecoffman/cp"
)

type LoResBuffer struct {
	sync.Mutex
	Data    [1024]byte
	Prev    [1024]byte
	Changed [1024]bool
}

func (lrb *LoResBuffer) Clear() {
	for idx := range lrb.Data {
		lrb.Data[idx] = 0
		lrb.Prev[idx] = 0
		lrb.Changed[idx] = false
	}
}

func (lrb *LoResBuffer) GetIndex(x, y int) int {
	my := y / 2
	offset := ((my % 8) * 128) + ((my / 8) * 40) + x
	return offset
}

func (lrb *LoResBuffer) IndexToXY(index int) (int, int) {
	mx := (index % 128) % 40
	my := ((index % 128) / 40 * 8) + ((index / 128) % 8)
	return mx, my
}

func (lrb *LoResBuffer) Init() {
	for idx, _ := range lrb.Data {
		lrb.Data[idx] = 0
		lrb.Changed[idx] = false
	}
}

func (lrb *LoResBuffer) Write(idx int, value byte) {
	idx = idx % 1024
	if lrb.Data[idx] == value {
		return
	}
	lrb.Data[idx] = value
	lrb.Changed[idx] = true
}

func (lrb *LoResBuffer) Read(idx int) byte {
	idx = idx % 1024
	return lrb.Data[idx]
}

const xratio = 1

func (lrb *LoResBuffer) plot(x, y int, c byte) {
	if x < 0 || x >= 40 {
		return
	}
	if y < 0 || y >= 40 {
		return
	}
	x = x % 40
	y = y % 48
	idx := lrb.GetIndex(x, y)
	cidx := y % 2
	b := lrb.Read(idx)
	switch cidx {
	case 0:
		b = (b & 0xf0) | (c & 0x0f)
	case 1:
		b = (b & 0x0f) | ((c & 0x0f) << 4)
	}
	lrb.Write(idx, b)
}

func (lrb *LoResBuffer) Plot(x, y float64, c byte) {
	lrb.Lock()
	defer lrb.Unlock()
	lrb.plot(int(math.Round(x)), int(math.Round(y)), c)
}

func (lrb *LoResBuffer) DrawRotatedBox(cx, cy, heading, width, height float64, c byte) {
	lrb.Lock()
	defer lrb.Unlock()
	hv := cp.ForAngle(heading)
	cen := cp.Vector{X: cx, Y: cy}
	for y := -height / 2; y <= height/2; y += 0.5 {
		for x := -width / 2; x <= width/2; x += 0.5 {
			pos := cp.Vector{X: x, Y: y}
			pos = pos.Rotate(hv)
			pos = pos.Add(cen)
			lrb.plot(int(math.Round(pos.X*xratio)), int(math.Round(pos.Y)), c)
		}
	}
}

func (lrb *LoResBuffer) DrawBox(x1, y1, x2, y2 float64, c byte) {
	lrb.Lock()
	defer lrb.Unlock()
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			lrb.plot(int(math.Round(x*xratio)), int(math.Round(y)), c)
		}
	}
}

func (lrb *LoResBuffer) DrawCircle(cx, cy float64, r float64, c byte) {
	lrb.Lock()
	defer lrb.Unlock()
	x1 := cx - r
	x2 := cx + r
	y1 := cy - r
	y2 := cy + r
	center := cp.Vector{X: float64(cx), Y: float64(cy)}
	var p cp.Vector
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			p = cp.Vector{X: float64(x), Y: float64(y)}
			if p.Sub(center).Length() <= float64(r) {
				lrb.plot(int(math.Round(x*xratio)), int(math.Round(y)), c)
			}
		}
	}
}

func (lrb *LoResBuffer) GetDeltasWithBase(base int) [][2]int {
	lrb.Lock()
	defer lrb.Unlock()
	out := [][2]int{}
	curr := lrb.Data
	for idx, newVal := range curr {
		if newVal != lrb.Prev[idx] {
			out = append(out, [2]int{base + idx, int(newVal)})
		}
		lrb.Changed[idx] = false
	}
	lrb.Prev = lrb.Data
	return out
}

func (lrb *LoResBuffer) WithDeltasDo(f func(b *LoResBuffer)) {
	old := lrb.Data
	f(lrb)
	new := lrb.Data
	for idx, v := range new {
		lrb.Changed[idx] = (v != old[idx])
	}
}
