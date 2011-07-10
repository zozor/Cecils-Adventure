package main

import (
	"sdl"
	"fmt"
	"math"
)

func AOSDNOASDNOIASND() {
	fmt.Println("ASDOASDIOANSDSA")
}

type Monster interface {
	Recter
	Move()
	Damage(int) bool
	Blit(surface *sdl.Surface)
	Attack() (*EnemyLaser, bool)
}

//Group of monsters
type Monsters []Monster

func (m Monsters) Blit(surface *sdl.Surface) {
	for i := 0; i < len(m); i++ {
		m[i].Blit(surface)
	}
}

func (m Monsters) Move() {
	for i := 0; i < len(m); i++ {
		m[i].Move()
	}
}

func (m Monsters) Attack() []*EnemyLaser {
	l := []*EnemyLaser{}
	for i := 0; i < len(m); i++ {
		out, ok := m[i].Attack()
		if ok {
			l = append(l, out)
		}
	}
	return l
}

func (m Monsters) Remove(i int) Monsters {
	return append(m[0:i], m[i+1:len(m)]...)
}

//GHOST
type MonsterGhost struct {
	X float64
	Y float64
	W float64
	H float64

	HP    int
	Speed float64
	GotHit int64

	Waypoints []Vector64
	Waypoint  int
}

func NewMonsterGhost(x, y, speed float64, Waypoints []Vector64) *MonsterGhost {
	m := new(MonsterGhost)
	m.X = x
	m.Y = y
	m.HP = 10
	m.Speed = speed

	m.Waypoints = append([]Vector64{Vector64{x,y}}, Waypoints...)
	return m
}

func NewGhostJSON(m JSONmonster) *MonsterGhost {

	return NewMonsterGhost(m.X,m.Y,m.Speed,m.Waypoints)
}

func (m *MonsterGhost) Attack() (*EnemyLaser, bool) {
	return nil, false
}

func (m *MonsterGhost) Damage(dmg int) bool {
	m.HP -= dmg
	if m.HP <= 0 {
		return true
	}
	m.GotHit = now()+50e6
	return false
}

func (m *MonsterGhost) Blit(surface *sdl.Surface) {
	m.H = float64(IM.Monsters[0][0].H)
	m.W = float64(IM.Monsters[0][0].W)
	if difftime(m.GotHit) > 0 {
		surface.Blit(&sdl.Rect{int16(m.X - m.W/2), int16(m.Y - m.H/2), 0, 0}, IM.Monsters[0][0], nil)
	}
}

func (m *MonsterGhost) GetRect() Rect {
	return Rect{m.X, m.Y, m.W, m.H}
}

func (m *MonsterGhost) Move() {
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
			m.Waypoint = 0
		} else {
			m.Waypoint++
		}

	}
}

///DEATH SHELL
const (
	DS_SPEED = 400
)
const (
	DS_NORMAL = iota
	DS_ISBOSS
)

const (
	DS_LEFT = iota
	DS_RIGHT
)

type Deathshell struct {
	X float64
	Y float64
	W float64
	H float64

	HP    int
	GotHit int64
	
	CoolDown int64
	Reloading int
	
	//Boss Specifics
	Speed float64
	IsBoss int
	Direction int

}

func NewDeathshell(x, y, speed float64, boss int) *Deathshell {
	d := new(Deathshell)
	d.X = x
	d.Y = y
	d.Speed = speed
	d.IsBoss = boss
	
	if boss == DS_ISBOSS {
		d.HP = 50
	} else {
		d.HP = 20
	}
	

	return d
}

func NewDeathShellJSON(m JSONmonster) *Deathshell {
	return NewDeathshell(m.X, m.Y, m.Speed, m.IsBoss)
}

func (d *Deathshell) Damage(dmg int) bool {
	d.HP -= dmg
	if d.HP <= 0 {
		return true
	}
	d.GotHit = now()+50e6
	return false
}

func (l *Deathshell) Attack() (*EnemyLaser, bool) {
	if difftime(l.CoolDown) < 0 {
		return nil, false
	} 
	l.CoolDown = now()+300e6
	l.Reloading++
	if l.Reloading == 4 {
		l.Reloading = 0
		l.CoolDown = now()+3e9
	}
	diff := Vector64{l.X-player.X, l.Y-player.Y}
	c := diff.Length()
	
	
	if diff.Y < 0 {
		n := NewEnemyLaser(l.X, l.Y, Vector64{-diff.X/c*DS_SPEED, math.Sin(math.Acos(diff.X/c))*DS_SPEED})
		return n , true
	} else {
		n := NewEnemyLaser(l.X, l.Y, Vector64{-diff.X/c*DS_SPEED, -math.Sin(math.Acos(diff.X/c))*DS_SPEED})
		return n, true
	}
		
	return nil, false
}

func (d *Deathshell) Blit(surface *sdl.Surface) {
	var img *sdl.Surface
	if d.IsBoss == DS_ISBOSS {
		img = IM.Monsters[1][1+d.Direction]
	} else {
		img = IM.Monsters[1][0]
	}
	
	d.H = float64(img.H)
	d.W = float64(img.W)
	
	if difftime(d.GotHit) > 0 {
		surface.Blit(&sdl.Rect{int16(d.X - d.W/2), int16(d.Y - d.H/2), 0, 0}, img, nil)
	}
}

func (d *Deathshell) GetRect() Rect {
	return Rect{d.X, d.Y, d.W, d.H}
}

func (d *Deathshell) Move() {
	if d.IsBoss == DS_NORMAL {
		return
	}
	if player.X > d.X {
		d.X += d.Speed/fps()
		d.Direction = DS_RIGHT
	}
	if player.X < d.X {
		d.X -= d.Speed/fps()
		d.Direction = DS_LEFT
	}
	if player.Y > d.Y{
		d.Y += d.Speed/fps()
	}
	if player.Y < d.Y{
		d.Y -= d.Speed/fps()
	}
}


//ROBOT
const (
	ROBOT_ISBOSS = 2
	ROBOT_NORMAL = 0
)

const (
	ROBOT_LEFT = iota
	ROBOT_RIGHT
)

type Robot struct {
	X float64
	Y float64
	W float64
	H float64

	HP    int
	Speed float64
	GotHit int64
	
	IsBoss int
	Direction int

}

func NewRobot(x, y, speed float64, boss int) *Robot {
	d := new(Robot)
	d.X = x
	d.Y = y
	d.IsBoss = boss
	d.Speed = speed
	
	
	if boss == ROBOT_ISBOSS {
		d.HP = 50
	} else {
		d.HP = 15
	}

	return d
}

func NewRobotJSON(m JSONmonster) *Robot {
	return NewRobot(m.X,m.Y,m.Speed,m.IsBoss)
}

func (r *Robot) Damage(dmg int) bool {
	r.HP -= dmg
	if r.HP <= 0 {
		return true
	}
	r.GotHit = now()+50e6
	return false
}

func (l *Robot) Attack() (*EnemyLaser, bool) {
	return nil, false
}

func (r *Robot) Blit(surface *sdl.Surface) {
	img := IM.Monsters[2][r.Direction+r.IsBoss]
	r.H = float64(img.H)
	r.W = float64(img.W)
	if difftime(r.GotHit) > 0 {
		surface.Blit(&sdl.Rect{int16(r.X - r.W/2), int16(r.Y - r.H/2), 0, 0}, img, nil)
	}
}

func (r *Robot) GetRect() Rect {
	return Rect{r.X, r.Y, r.W, r.H}
}

func (r *Robot) Move() {
	if player.X > r.X {
		r.X += r.Speed/fps()
		r.Direction = ROBOT_RIGHT
	}
	if player.X < r.X {
		r.X -= r.Speed/fps()
		r.Direction = ROBOT_LEFT
	}
	if player.Y > r.Y{
		r.Y += r.Speed/fps()
	}
	if player.Y < r.Y{
		r.Y -= r.Speed/fps()
	}
}


//MALBORO
const MAL_SPEED = 400

type Malboro struct {
	X float64
	Y float64
	W float64
	H float64

	HP    int
	GotHit int64
	
	CoolDown int64
	Reloading int
	Direction float64 //radianer
}

func NewMalboro(x, y float64) *Malboro {
	d := new(Malboro)
	d.X = x
	d.Y = y
	d.HP = 30
	
	return d
}
func NewMalboroJSON(m JSONmonster) *Malboro {
	return NewMalboro(m.X,m.Y)
}

func (d *Malboro) Damage(dmg int) bool {
	d.HP -= dmg
	if d.HP <= 0 {
		return true
	}
	d.GotHit = now()+50e6
	return false
}

func (l *Malboro) Attack() (*EnemyLaser, bool) {
	if difftime(l.CoolDown) < 0 {
		return nil, false
	} 
	l.CoolDown = now()+300e6
	l.Reloading++
	if l.Reloading == 11 {
		l.Reloading = 0
		l.CoolDown = now()+3e9
	}

	n := NewEnemyLaser(l.X, l.Y, Vector64{MAL_SPEED*math.Cos(l.Direction), MAL_SPEED*math.Sin(l.Direction)})
		
	return n, true
}

func (d *Malboro) Blit(surface *sdl.Surface) {
	img := IM.Monsters[3][0]	
	d.H = float64(img.H)
	d.W = float64(img.W)
	
	if difftime(d.GotHit) > 0 {
		surface.Blit(&sdl.Rect{int16(d.X - d.W/2), int16(d.Y - d.H/2), 0, 0}, img, nil)
	}
}

func (d *Malboro) GetRect() Rect {
	return Rect{d.X, d.Y, d.W, d.H}
}

func (d *Malboro) Move() {
	d.Direction += math.Pi/2/fps()
}
