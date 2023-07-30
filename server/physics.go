package main

import (
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/jakecoffman/cp"
)

const (
	maxObjects = 8
)

type ShapeType byte

const (
	stPoint  ShapeType = 0x00
	stRect   ShapeType = 0x01
	stCircle ShapeType = 0x02
)

type PhysicsObject struct {
	body         *cp.Body
	radius       float64
	width        float64
	height       float64
	kind         ShapeType
	lastPubX     float64
	lastPubY     float64
	lastPubH     float64
	color        int
	mass         int
	bodyType     int
	collided     bool
	collidedWith int
	elasticity   float64
}

const rotatedBoxes = true

func (po *PhysicsObject) undraw(screen *LoResBuffer, color int) {
	if po.lastPubX == -1 {
		return
	}
	switch po.kind {
	case stPoint:
		screen.WithDeltasDo(func(screen *LoResBuffer) {
			screen.Plot(po.lastPubX, po.lastPubY, byte(color))
		})
	case stCircle:
		screen.WithDeltasDo(func(screen *LoResBuffer) {
			screen.DrawCircle(po.lastPubX, po.lastPubY, po.radius, byte(color))
		})
	case stRect:
		if rotatedBoxes {
			screen.DrawRotatedBox(po.lastPubX, po.lastPubY, po.lastPubH, po.width, po.height, byte(color))
		} else {
			screen.WithDeltasDo(func(screen *LoResBuffer) {
				x1 := po.lastPubX - po.width/2
				x2 := po.lastPubX + po.width/2
				y1 := po.lastPubY - po.height/2
				y2 := po.lastPubY + po.height/2
				screen.DrawBox(x1, y1, x2, y2, byte(color))
			})
		}
	}
}

func radToDeg(r float64) float64 {
	return 360 - r*180/math.Pi
}

func (po *PhysicsObject) draw(screen *LoResBuffer, color int) {
	var pos = po.body.Position()
	var heading = po.body.Angle()
	var cx, cy = pos.X, pos.Y
	// use current pos
	switch po.kind {
	case stPoint:
		screen.WithDeltasDo(func(screen *LoResBuffer) {
			screen.Plot(cx, cy, byte(color))
		})
	case stCircle:
		screen.WithDeltasDo(func(screen *LoResBuffer) {
			screen.DrawCircle(cx, cy, po.radius/2, byte(color))
		})
	case stRect:
		if rotatedBoxes {
			screen.DrawRotatedBox(cx, cy, heading, po.width, po.height, byte(color))
		} else {
			screen.WithDeltasDo(func(screen *LoResBuffer) {
				x1 := cx - po.width/2
				x2 := cx + po.width/2
				y1 := cy - po.height/2
				y2 := cy + po.height/2
				screen.DrawBox(x1, y1, x2, y2, byte(color))
			})
		}
	}
	po.lastPubX = cx
	po.lastPubY = cy
	po.lastPubH = heading
}

type PhysicsEngine struct {
	sync.Mutex
	minBounds cp.Vector
	maxBounds cp.Vector
	objects   [maxObjects]*PhysicsObject
	space     *cp.Space
	bounds    *cp.Shape
	interval  time.Duration
	running   bool
	screen    *LoResBuffer
	// reportFunc func(delta [][2]int)
}

func NewPhysicsEngine(x0, y0, x1, y1 float64, interval time.Duration) *PhysicsEngine {
	s := cp.NewSpace()
	p := &PhysicsEngine{
		minBounds: cp.Vector{X: x0, Y: y0},
		maxBounds: cp.Vector{X: x1, Y: y1},
		space:     s,
		interval:  interval,
		screen:    &LoResBuffer{},
		// reportFunc: reportFunc,
	}
	handler := s.NewCollisionHandler(1, 1)
	handler.BeginFunc = p.BeginCollision
	handler.SeparateFunc = p.EndCollision
	//p.createBoundingBox()
	return p
}

func (p *PhysicsEngine) SetForce(f float64, heading float64) {
	fv := headingToVector(heading).Mult(f)
	p.space.SetGravity(fv)
}

func (p *PhysicsEngine) BeginCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
	a, b := arb.Bodies()
	if aID, ok := a.UserData.(int); ok {
		if bID, ok := b.UserData.(int); ok {
			p.objects[aID].collided = true
			p.objects[aID].collidedWith = bID
			p.objects[bID].collided = true
			p.objects[bID].collidedWith = aID
		}
	}
	return true
}

func (p *PhysicsEngine) GetCollidedWith(id int) (bool, int) {
	o := p.objects[id%maxObjects]
	if o == nil {
		return false, 0
	}
	defer func() {
		o.collided = false
	}()
	c, w := o.collided, o.collidedWith
	return c, w
}

func (p *PhysicsEngine) EndCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) {
}

func (p *PhysicsEngine) createBoundingBox() {
	var KinematicBoxBox = p.space.AddBody(cp.NewKinematicBody())

	a := cp.Vector{p.minBounds.X, p.minBounds.Y}
	b := cp.Vector{p.minBounds.X, p.maxBounds.Y}
	c := cp.Vector{p.maxBounds.X, p.maxBounds.Y}
	d := cp.Vector{p.maxBounds.X, p.minBounds.Y}

	shape := p.space.AddShape(cp.NewSegment(KinematicBoxBox, a, b, 0))
	shape.SetElasticity(1)
	shape.SetFriction(1)

	shape = p.space.AddShape(cp.NewSegment(KinematicBoxBox, b, c, 0))
	shape.SetElasticity(1)
	shape.SetFriction(1)

	shape = p.space.AddShape(cp.NewSegment(KinematicBoxBox, c, d, 0))
	shape.SetElasticity(1)
	shape.SetFriction(1)

	shape = p.space.AddShape(cp.NewSegment(KinematicBoxBox, d, a, 0))
	shape.SetElasticity(1)
	shape.SetFriction(1)

	p.bounds = shape
}

func headingToVector(h float64) cp.Vector {
	h = 360 - h
	m := mgl64.QuatRotate(mgl64.DegToRad(h), mgl64.Vec3{0, 0, 1}).Mat4()
	v := mgl64.TransformNormal(mgl64.Vec3{0, -1, 0}, m)
	return cp.Vector{X: v[0], Y: v[1]}
}

func (p *PhysicsEngine) addRect(slot int, width, height float64, mass float64, pos cp.Vector, vel cp.Vector, color int, bodyType int, replace bool) {
	var elasticity float64 = 1
	if b := p.objects[slot%maxObjects]; b != nil {
		if replace {
			elasticity = b.elasticity
			mass = b.body.Mass()
			vel = b.body.Velocity()
			pos = b.body.Position()
			color = b.color
		}
		// p.space.EachShape(func(s *cp.Shape) {
		// 	if s.Body() == b.body {
		// 		b.body.RemoveShape(s)
		// 	}
		// })
		// p.space.RemoveBody(b.body)
	}
	var b = p.space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, width, height)))
	b.SetMass(mass)
	b.SetAngle(0)
	b.SetPosition(pos)
	var s = p.space.AddShape(cp.NewBox(b, width, height, 0))
	s.SetElasticity(elasticity)
	s.SetFriction(0)
	s.SetMass(mass)
	b.SetType(bodyType)
	b.UserData = slot % maxObjects
	b.SetVelocity(vel.X, vel.Y)
	s.SetCollisionType(1)
	p.objects[slot%maxObjects] = &PhysicsObject{elasticity: elasticity, color: color, body: b, width: width, height: height, kind: stRect, lastPubX: -1, lastPubY: -1}
}

func (p *PhysicsEngine) RemoveObject(id int) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	if o.lastPubX != -1 {
		o.undraw(p.screen, 0)
	}
	p.space.EachShape(func(s *cp.Shape) {
		if s.Body() == o.body {
			o.body.RemoveShape(s)
		}
	})
	p.space.RemoveBody(o.body)
	// p.objects[id%maxObjects] = nil
}

func (p *PhysicsEngine) addCircle(slot int, radius float64, mass float64, pos cp.Vector, vel cp.Vector, color int, bodyType int, replace bool) {
	var elasticity float64 = 1
	if b := p.objects[slot%maxObjects]; b != nil {
		if replace {
			elasticity = b.elasticity
			mass = b.body.Mass()
			vel = b.body.Velocity()
			pos = b.body.Position()
			color = b.color
		}
		// p.space.EachShape(func(s *cp.Shape) {
		// 	if s.Body() == b.body {
		// 		b.body.RemoveShape(s)
		// 	}
		// })
		// p.space.RemoveBody(b.body)
	}
	var b = p.space.AddBody(cp.NewBody(mass, cp.MomentForCircle(mass, 0, radius, cp.Vector{})))
	b.SetMass(mass)
	b.SetAngle(0)
	b.SetPosition(pos)
	var s = p.space.AddShape(cp.NewCircle(b, radius, cp.Vector{}))
	s.SetElasticity(elasticity)
	s.SetFriction(0)
	s.SetMass(mass)
	b.SetType(bodyType)
	b.SetVelocity(vel.X, vel.Y)
	s.SetCollisionType(1)
	b.UserData = slot % maxObjects
	p.objects[slot%maxObjects] = &PhysicsObject{elasticity: elasticity, color: color, body: b, radius: radius, kind: stCircle, lastPubX: -1, lastPubY: -1}
}

func (p *PhysicsEngine) Start() {
	if p.running {
		return
	}
	p.Lock()
	defer p.Unlock()
	go func(p *PhysicsEngine) {
		p.running = true
		for p.running {
			time.Sleep(p.interval)
			var dt = float64(p.interval) / float64(time.Second)
			// log.Printf("Stepping with dt of %f", dt)
			p.Step(dt)
		}
	}(p)
}

func (p *PhysicsEngine) GetDeltasWithBase(base int) [][2]int {
	p.reportDeltas()
	return p.screen.GetDeltasWithBase(base)
}

func (p *PhysicsEngine) Step(dt float64) {
	p.Lock()
	defer p.Unlock()
	p.space.Step(dt)
	//p.reportDeltas()
}

func (p *PhysicsEngine) Stop() {
	p.Lock()
	defer p.Unlock()
	if !p.running {
		return
	}
	p.running = false
	time.Sleep(2 * p.interval)
	for idx, o := range p.objects {
		if o == nil {
			continue
		}
		p.space.EachShape(func(s *cp.Shape) {
			if s.Body() == o.body {
				o.body.RemoveShape(s)
			}
		})
		p.space.RemoveBody(o.body)
		p.objects[idx] = nil
	}
}

func (p *PhysicsEngine) GetObjectColor(id int) int {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return 0
	}
	pos := o.color
	return pos
}

func (p *PhysicsEngine) GetObjectPos(id int) (int, int) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return 0, 0
	}
	pos := o.body.Position()
	return int(pos.X), int(pos.Y)
}

func (p *PhysicsEngine) GetAnyOOB() (int, bool) {
	p.Lock()
	defer p.Unlock()
	for idx, o := range p.objects {
		if o == nil {
			return idx, true
		}
		pos := o.body.Position()
		// log.Printf("OOB check: body %d pos = %+v", idx, pos)
		if math.IsNaN(pos.X) || math.IsNaN(pos.Y) {
			return idx, true
		}
		if pos.X < p.minBounds.X || pos.X > p.maxBounds.X || pos.Y < p.minBounds.Y || pos.Y > p.maxBounds.Y {
			return idx, true
		}
	}
	return 0, false
}

func (p *PhysicsEngine) GetObjectOOB(id int) int {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return 1
	}
	pos := o.body.Position()
	if math.IsNaN(pos.X) || math.IsNaN(pos.Y) {
		return 1
	}
	if pos.X < p.minBounds.X || pos.X > p.maxBounds.X || pos.Y < p.minBounds.Y || pos.Y > p.maxBounds.Y {
		return 1
	}
	return 0
}

func (p *PhysicsEngine) SetObjectVelocityHeading(id int, v float64, heading float64) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	fv := headingToVector(heading).Mult(v)
	o.body.SetVelocity(fv.X, fv.Y)
}

func (p *PhysicsEngine) AddObjectVelocityHeading(id int, v float64, heading float64) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	fv := headingToVector(heading).Mult(v)
	ov := o.body.Velocity()
	fv = fv.Add(ov)
	o.body.SetVelocity(fv.X, fv.Y)
}

func (p *PhysicsEngine) SetObjectRect(id int, w, h int) {
	p.Lock()
	defer p.Unlock()
	w = w % 40
	h = h % 48
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	o.undraw(p.screen, 0)
	p.addRect(
		id,
		float64(w),
		float64(h),
		float64(o.mass),
		o.body.Position(),
		o.body.Velocity(),
		o.color,
		o.body.GetType(),
		true,
	)
	o.draw(p.screen, o.color)
}

func (p *PhysicsEngine) SetObjectPos(id int, x, y int) {
	p.Lock()
	defer p.Unlock()
	x = x % 40
	y = y % 48
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d pos to %d, %d", id, x, y)
	o.body.SetPosition(cp.Vector{X: float64(x), Y: float64(y)})
}

func (p *PhysicsEngine) SetObjectType(id int, kind int) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d type to %d", id, kind&0x1)
	o.bodyType = kind & 0x01
	switch kind & 0x01 {
	case 0:
		o.body.SetType(cp.BODY_DYNAMIC)
	case 1:
		o.body.SetType(cp.BODY_KINEMATIC)
	}
}

func (p *PhysicsEngine) SetObjectElasticity(id int, e float64) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d elasticity to %f", id, e)
	o.body.EachShape(func(s *cp.Shape) {
		s.SetElasticity(e)
	})
	o.elasticity = e
}

func (p *PhysicsEngine) SetObjectMass(id int, mass int) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d mass to %d", id, mass)
	o.mass = mass
	o.body.SetMass(float64(mass))
	// if !p.running {
	// 	log.Printf("Starting physics engine...")
	// 	p.Start()
	// }
}

func (p *PhysicsEngine) SetObjectHeading(id int, heading float64) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d heading to %f", id, heading)
	a := mgl64.DegToRad(360 - heading)
	o.body.SetAngle(a)
}

func (p *PhysicsEngine) GetObjectHeading(id int) int {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return 0
	}
	a := 360 - mgl64.RadToDeg(o.body.Angle())
	return int(math.Round(a))
}

func (p *PhysicsEngine) SetObjectColor(id int, color int) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d color to %d", id, color)
	o.color = color
	// if !p.running {
	// 	log.Printf("Starting physics engine...")
	// 	p.Start()
	// }
}

func (p *PhysicsEngine) SetObjectVelocity(id int, velX, velY float64) {
	p.Lock()
	defer p.Unlock()
	o := p.objects[id%maxObjects]
	if o == nil {
		return
	}
	log.Printf("Setting object %d velocity to %f, %f", id, velX, velY)
	o.body.SetVelocity(velX, velY)
	// if !p.running {
	// 	log.Printf("Starting physics engine...")
	// 	p.Start()
	// }
}

func (p *PhysicsEngine) reportDeltas() {
	var tmp = []*PhysicsObject{}
	for _, b := range p.objects {
		if b == nil {
			continue
		}
		tmp = append(tmp, b)
	}
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].body.GetType() < tmp[j].body.GetType()
	})
	for _, b := range tmp {
		if b == nil {
			continue
		}
		// var pos = b.body.Position()
		// var cx, cy = int(pos.X), int(pos.Y)
		// log.Printf("Body at %d, %d", cx, cy)
		// if cx != b.lastPubX || cy != b.lastPubY {
		//p.screen.WithDeltasDo(func(lrb *LoResBuffer) {
		b.undraw(p.screen, 0)
		// p.screen.Plot(cx, cy, byte(b.color))
		b.draw(p.screen, b.color)
		//})
		// log.Printf("Body %d: moved to %d, %d", idx, cx, cy)
		// }
	}
}
