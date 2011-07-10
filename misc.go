package main

import (
	"sdl"
	"fmt"
)

func AOSDNOASDNOIASN2D() {
	fmt.Println("ASDOASDIOANSDSA")
}

type Lasers []Shots

func (l Lasers) Blit(surface *sdl.Surface) {
	for _, v := range l {
		v.Blit(surface)
	}
}

func (l Lasers) Move() {
	for _, v := range l {
		v.Move()
	}
}

func (l Lasers) Remove(i int) Lasers {
	return append(l[0:i], l[i+1:len(l)]...)
}

func (l Lasers) Clean() Lasers {
	for i:=0; i < len(l); i++ {
		r := l[i].GetRect()
		if r.X+r.W/2 < 0 || r.Y+r.H/2 < 0 {
			l = l.Remove(i)
			i--
		} else if r.X-r.W/2 > SCREENWIDTH || r.Y-r.H/2 > SCREENHEIGHT {
			l = l.Remove(i)
			i--
		}
	}
	return l
}

func (l Lasers) HitWall(lvl *Level) Lasers {
	for i:=0; i<len(l); i++ {
		if l[i].HitWall(lvl) {
			l = l.Remove(i)
			i--
		}
	}
	return l
}

type Shots interface {
	Recter
	Blit(*sdl.Surface)
	Move()
	HitWall(lvl *Level) bool
	Damage() int
}


type Laser struct {
	X float64
    Y float64
    W float64
    H float64
    
    Dmg int //(Damage+1)*3
    Speed float64
}

func NewLaser(x,y float64,dmg int, speed float64) *Laser {
	m := new(Laser)
	m.X = x
	m.Y = y
	m.Speed = speed
	m.Dmg = dmg

	//Sound
	
	MU.Sound["shot"].PlayChannel(-1,0)

	return m
}

func (l *Laser) Damage() int {
	return (l.Dmg+1)*3
}

func (l *Laser) GetRect() Rect {
	return Rect{l.X, l.Y, l.W, l.H}
}

func (l *Laser) HitWall(lvl *Level) bool {
	_, i := BoxCollide(l, lvl)
	if i != 0 {
		return true
	}
	return false
}

func (l *Laser) Blit(surface *sdl.Surface) {
	l.H = float64(IM.Shots[l.Dmg].H)
	l.W = float64(IM.Shots[l.Dmg].W)
	surface.Blit(&sdl.Rect{int16(l.X-l.W/2),int16(l.Y-l.H/2),0,0}, IM.Shots[l.Dmg], nil);
}

func (l *Laser) Move() {
	l.X += l.Speed/fps()
}

//Enemy Laser
type EnemyLaser struct {
	X float64
    Y float64
    W float64
    H float64
    
    Speed Vector64
}

func NewEnemyLaser(x,y float64, speed Vector64) *EnemyLaser {
	m := new(EnemyLaser)
	m.X = x
	m.Y = y
	m.Speed = speed

	//Sound
	MU.Sound["enemyshot"].PlayChannel(-1,0)

	return m
}

func (l *EnemyLaser) Damage() int {
	return 0
}

func (l *EnemyLaser) GetRect() Rect {
	return Rect{l.X, l.Y, l.W, l.H}
}

func (l *EnemyLaser) HitWall(lvl *Level) bool {
	return false
}

func (l *EnemyLaser) Blit(surface *sdl.Surface) {
	l.H = float64(IM.Shots[0].H)
	l.W = float64(IM.Shots[0].W)
	surface.Blit(&sdl.Rect{int16(l.X-l.W/2),int16(l.Y-l.H/2),0,0}, IM.Shots[0], nil);
}

func (l *EnemyLaser) Move() {
	l.X += l.Speed.X/fps()
	l.Y += l.Speed.Y/fps()
}

//Right Menu
type GameTime struct {
	Start int64
	End	int64
}

func BlitHeartsAndTimer(surface *sdl.Surface, antal int, t GameTime) {
	x := 18*BOXSIZE+10
	y := 5
	for i:=0;i<antal;i++ {
		surface.Blit(&sdl.Rect{int16(x),int16(i*30+y),0,0}, IM.Misc[0], nil)
	}
	
	max := t.End-t.Start
	nu := t.End-now()
	if nu <= 0 {
		return
	}
	surface.FillRect(&sdl.Rect{int16(x+BOXSIZE), 5+BOXSIZE, 15, uint16((200*nu)/max)}, 0xFFFFFF)
}