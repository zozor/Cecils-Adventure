package main

import (
	"sdl"
	"fmt"
	"rand"
)

func AOSDNOASDNOIASND2() {
	fmt.Println("ASDOASDIOANSDSA")
}

type Object interface {
	Recter
	HitPlayer(*Player)
	Blit(*sdl.Surface)
	Move()
	IsDead() bool
}

type Objects []Object

func (o Objects) Blit(surface *sdl.Surface) {
	for i := 0; i < len(o); i++ {
		o[i].Blit(surface)
	}
}

func (o Objects) Move() {
	for i := 0; i < len(o); i++ {
		o[i].Move()
	}
}

func (o Objects) Remove(i int) Objects {
	return append(o[0:i], o[i+1:len(o)]...)
}

func (o Objects) Hit(p *Player) {
	for i := 0; i < len(o); i++ {
		o[i].HitPlayer(p)
	}
}

func (o Objects) Clean() Objects {
	for i := 0; i < len(o); i++ {
		if o[i].IsDead() {
			o = o.Remove(i)
			i--
		}
	}
	return o
}




//Teleport
const (
	TELE_VERTICAL = iota
	TELE_HORIZONTAL
)

type Teleport struct {
	X, Y float64
	W, H float64
	S, O int //Size, retning
	Way	Vector64
	Boost float64
	
	Conn *Teleport
	
	Image	*sdl.Surface
	
}

func NewTeleport(xy, way Vector64, s, o int, boost float64) *Teleport {
	t := new(Teleport)
	t.X = xy.X
	t.Y = xy.Y
	t.Way = way
	t.S = s
	t.O = o
	t.Boost = boost
	
	//Draw Stuff
	if o == TELE_HORIZONTAL {
		t.Image = IM.Obj[1]
		t.H = float64(t.Image.H)
		t.W = float64(s*BOXSIZE)
	} else {
		t.Image = IM.Obj[0]
		t.W = float64(t.Image.W)
		t.H = float64(s*BOXSIZE)
	}
	
	return t
}

func NewTeleportJSON(m JSONteleport) *Teleport {
	o := 0
	if m.Type == SO_TELEPORT_V {
		o = TELE_HORIZONTAL
	} else {
		o = TELE_VERTICAL
	}
	return NewTeleport(Vector64{m.X,m.Y}, m.Way, m.Size, o, m.Boost)
}

func (t *Teleport) Blit(surface *sdl.Surface) {
	if t.O == TELE_HORIZONTAL {
		for x:=0; x<t.S*2; x++ {
			surface.Blit(&sdl.Rect{int16(t.X - t.W/2)+int16(20*x), int16(t.Y - t.H/2), 0, 0}, t.Image, nil)
		}
	} else {
		for x:=0; x<t.S*2; x++ {
			surface.Blit(&sdl.Rect{int16(t.X - t.W/2) , int16(t.Y - t.H/2)+int16(20*x), 0, 0}, t.Image, nil)
		}
	}
}

func (t *Teleport) Connect(tele *Teleport) {
	t.Conn = tele
}

func (t *Teleport) ConnectBoth(tele *Teleport) {
	t.Conn = tele
	tele.Conn = t
}

func (t *Teleport) GetRect() Rect {
	return Rect{t.X, t.Y, t.W, t.H}
}

func (t *Teleport) HitPlayer(p *Player) {
	if HitTest(t, p) {
		if t.Conn == nil {
			t.MoveTo(t, p)
		} else {
			t.MoveTo(t.Conn, p)
		}
	}
}

func (t *Teleport) Move() { //implement Object interface
}

func (t *Teleport) MoveTo(to *Teleport, player *Player) {
	w := to.Way
	p := &player.Velocity
	
	//X-axis
	switch {
	case p.X <= 0 && w.X > 0:
		p.X = -(p.X*t.Boost)
	case p.X >= 0 && w.X > 0:
		p.X = (p.X*t.Boost)
	case p.X <= 0 && w.X < 0:
		p.X = p.X*t.Boost
	case p.X >= 0 && w.X < 0:
		p.X = -(p.X*t.Boost)
	}
	
	//Y-Axis
	switch {
	case p.Y <= 0 && w.Y > 0:
		p.Y = -(p.Y*t.Boost)
	case p.Y >= 0 && w.Y > 0:
		p.Y = (p.Y*t.Boost)
	case p.Y <= 0 && w.Y < 0:
		p.Y = p.Y*t.Boost
	case p.Y >= 0 && w.Y < 0:
		p.Y = -(p.Y*t.Boost)
	}

	player.X = to.X + to.Way.X*to.W + to.Way.X*player.W/2
	player.Y = to.Y + to.Way.Y*to.H + to.Way.Y*player.H/2
}

func (s *Teleport) IsDead() bool { //Implements interface Objects
	return false
}

//Spikes

type Spike struct {
	X, Y float64
	W, H float64
	Speed float64
	
	Waypoint int
	Waypoints []Vector64
	Restart bool //Teleport to start
	
}

func NewSpike(xy Vector64, speed float64, restart bool, wp []Vector64) *Spike {
	s := new(Spike)
	s.X = xy.X
	s.Y = xy.Y
	s.Speed = speed
	s.Restart = restart
	
	s.Waypoints = append([]Vector64{xy}, wp...)
	
	s.W = float64(IM.Shots[2].W)
	s.H = float64(IM.Shots[2].H)
	
	return s
}

func NewSpikeJSON(m JSONobject) *Spike {
	return NewSpike(Vector64{m.X,m.Y}, m.Speed, m.Restart, m.Waypoints)
}

func (s *Spike) Blit(surface *sdl.Surface) {
	surface.Blit(&sdl.Rect{int16(s.X - s.W/2), int16(s.Y - s.H/2), 0, 0}, IM.Shots[2], nil)
}

func (s *Spike) GetRect() Rect {
	return Rect{s.X, s.Y, s.W, s.H}
}

func (s *Spike) HitPlayer(p *Player) {
	if !player.IsImmortal() && HitTest(s, p) {
		p.WasHit()
	}
}

func (m *Spike) Move() {
	if len(m.Waypoints) < 2 {
		return
	}
	
	if m.X < m.Waypoints[m.Waypoint].X {
		if m.Waypoints[m.Waypoint].X - m.X < m.Speed/fps() {
			m.X += m.Waypoints[m.Waypoint].X - m.X
		} else {
			m.X += m.Speed/fps()
		}
	}
	if m.X > m.Waypoints[m.Waypoint].X {
		if -m.Waypoints[m.Waypoint].X + m.X < m.Speed/fps() {
			m.X -= -m.Waypoints[m.Waypoint].X + m.X
		} else {
			m.X -= m.Speed/fps()
		}
	}
	if m.Y < m.Waypoints[m.Waypoint].Y {
		if m.Waypoints[m.Waypoint].Y - m.Y < m.Speed/fps() {
			m.Y += m.Waypoints[m.Waypoint].Y - m.Y
		} else {
			m.Y += m.Speed/fps()
		}
			}
	if m.Y > m.Waypoints[m.Waypoint].Y {
		if m.Y - m.Waypoints[m.Waypoint].Y < m.Speed/fps() {
			m.Y -= m.Y - m.Waypoints[m.Waypoint].Y
		} else {
			m.Y -= m.Speed/fps()
		}
	}
	if m.Y == m.Waypoints[m.Waypoint].Y && m.X == m.Waypoints[m.Waypoint].X {
		if m.Waypoint == len(m.Waypoints)-1 {
			if m.Restart {
				m.X = m.Waypoints[0].X
				m.Y = m.Waypoints[0].Y
				m.Waypoint = 0
			} else {
				m.Waypoint = 0
			}
			
			
		} else {
			m.Waypoint++
		}

	}
}

func (s *Spike) IsDead() bool { //Implements interface Objects
	return false
}

//Spoils of War
const (
	SP_HEART = iota
	SP_UPGRADE
)

type Spoils struct {
	X, Y float64
	W, H float64
	Velocity Vector64
	
	Type int
	Frame int
	
	Dead bool
}

type UpgradeAnimation struct {
	frame int
}

func NewSpoils(x,y float64) *Spoils {
	s := new(Spoils)
	if rand.Float64() > 0.2 {
		s.X, s.Y = -100,-100
		s.Dead = true
		return s
	}
	
	s.X = x
	s.Y = y
	s.Type = rand.Intn(2)
	
	s.Velocity = Vector64{(rand.Float64()*10-5)*20,(rand.Float64()*10-5)*20}	
	
	switch s.Type {
	case SP_HEART:
		s.W = float64(IM.Misc[0].W)
		s.H = float64(IM.Misc[0].H)
	case SP_UPGRADE:
		s.W = 10
		s.H = 10
	}
	
	return s
}

func (s *Spoils) Blit(surface *sdl.Surface) {
	switch s.Type {
	case SP_HEART:
		surface.Blit(&sdl.Rect{int16(s.X - s.W/2), int16(s.Y - s.H/2), 0, 0}, IM.Misc[0], nil)
	case SP_UPGRADE:
		surface.FillRect(&sdl.Rect{int16(s.X - s.W/2), int16(s.Y - s.H/2), 10, 10},uint32(s.Frame))
		s.Frame = s.Frame*s.Frame+1
		if s.Frame > 0xFFFFFF {
			s.Frame = 0
		}
	}
}

func (s *Spoils) GetRect() Rect {
	return Rect{s.X, s.Y, s.W, s.H}
}

func (s *Spoils) HitPlayer(p *Player) {
	if HitTest(s, p) {
		switch s.Type {
		case SP_HEART:
			p.Lives++
		case SP_UPGRADE:
			if p.Lvl < 3 {
				p.Lvl++
			}
		}
		s.Dead = true
	}
}

func (s *Spoils) Move() {
	s.X += s.Velocity.X/fps()
	s.Y += s.Velocity.Y/fps()
	if s.X > SCREENWIDTH+50 || s.X < -50 {
		s.Dead = true
	}
	if s.Y > SCREENHEIGHT+50 || s.Y < -50 {
		s.Dead = true
	}
}

func (s *Spoils) IsDead() bool {
	return s.Dead
}