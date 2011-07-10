package main

import (
	"sdl"
	"fmt"
	"math"
)

func ASOIDAIS() {
	fmt.Println("LASD")
}

const (
	PLAYER_LIVES = 8
	PLAYER_TOPSPEED = 200
	PLAYER_MAXGRAVITY = 500
	
	//TEMP
	MAP_FRICTIONWALL = 500
	MAP_FRICTIONGROUND = 200
)

//Player
type Player struct {
	X                  float64
	Y                  float64
	W                  float64
	H                  float64
	GravityMaxVelocity float64
	Topspeed           float64

	InAir bool
	
	Lives int
	Immortal int64

	Direction bool
	Lvl       int
	Cooldown  int64
	BigCool   int
	Shots     int

	Velocity Vector64
	Force    Vector64

	TimeFrame	int64
	Frame     int
	Animation int
}

func NewPlayer() *Player {
	out := new(Player)
	out.W, out.H = 24, 36

	out.Lives = PLAYER_LIVES

	return out
}

func (b *Player) Animate() (out *sdl.Surface) {
	i := IM.Cecil
	anim := b.Animation & 3
	switch anim {
	case 0,2:
		if b.Direction {
			out = i[anim][1]
		} else {
			out = i[anim][0]
		}
	case 1:
		if b.Direction {
			out = i[anim][b.Frame+2]
		} else {
			out = i[anim][b.Frame]
		}
		
		if difftime(b.TimeFrame) > 0 {
			b.Frame++
			b.TimeFrame = now()+100e6
		}
		if b.Frame == 2 {
			b.Frame = 0
		}
	}
	if b.Animation & 4 != 0 {
		out.SetAlpha(sdl.SRCALPHA, 125)
	} else {
		out.SetAlpha(0, sdl.ALPHA_OPAQUE)
	}
	
	
	b.H = float64(out.H)
	b.W = float64(out.W)
	return out
}


func (b *Player) Blit(surface *sdl.Surface) {
	surface.Blit(&sdl.Rect{int16(b.X - b.W/2), int16(b.Y - b.H/2), 0, 0}, b.Animate(), nil)
	
	surface.Blit(&sdl.Rect{19*BOXSIZE+10, 10, 0, 0}, IM.Shots[b.Lvl], nil)
}

func (b *Player) SetAnimation(ani int) {
	if ani != b.Animation {
		b.Animation = ani
		b.Frame = 0
	}
	if b.IsImmortal() {
		b.Animation |= 4
	}
}

func (b *Player) IsImmortal() bool {
	if difftime(b.Immortal) > 0 {
		return false
	}
	return true
}

func (b *Player) Events() {
	key := sdl.GetKeyState()
	kright := key[sdl.K_RIGHT]
	kleft := key[sdl.K_LEFT]
	
	//Set Maximums
	b.Topspeed = PLAYER_TOPSPEED/fps()
	b.GravityMaxVelocity = PLAYER_MAXGRAVITY/fps()
	b.Force = Vector64{PLAYER_TOPSPEED*2/fps2(), PLAYER_TOPSPEED*2/fps2()}

	//Apply Force
	if key[sdl.K_UP] != 0 && b.InAir == false {
		b.Velocity.Y = -math.Sqrt(210*gravity.Y)
		b.InAir = true
	}
	if !(kright != 0 && kleft != 0) {	
		if kright != 0 && b.Velocity.X < b.Topspeed {
			b.Velocity.X += b.Force.X
			b.Direction = true
		}
	
		if kleft != 0 && b.Velocity.X > -b.Topspeed {
			b.Velocity.X -= b.Force.X
			b.Direction = false
			b.SetAnimation(1)
		}
	}

	if b.Velocity.Y < b.GravityMaxVelocity {
		b.Velocity.Y += gravity.Y
	}
	if b.Velocity.Y > b.GravityMaxVelocity {
		b.Velocity.Y = b.GravityMaxVelocity
	}
	
	/*DEBUG
	if key[sdl.K_UP] != 0 && b.Velocity.Y > -b.Topspeed {
		b.Velocity.Y -= b.Force.Y
	}
	if key[sdl.K_DOWN] != 0 && b.Velocity.Y < b.Topspeed {
		b.Velocity.Y += b.Force.Y
	}*/

	b.InAir = true
	
	//Animation
	if (kleft != 0 || kright != 0) && !(kright != 0 && kleft != 0) {
		b.SetAnimation(1)
	} else {
		b.SetAnimation(0)
	}
}

func (b *Player) WasHit() {
	b.Immortal = now()+3e9
	b.Lives--
}

func (b *Player) Shoot() *Laser {
	if now()-b.Cooldown < 0 {
		return nil
	}
	b.BigCool++
	b.Cooldown = now()+200*1e6
	if b.BigCool == 5 {
		b.BigCool = 0
		b.Cooldown = now()+2000*1e6 //2 sekunders cooldown
	}
	
	b.SetAnimation(2)
	if b.Direction {
		return NewLaser(b.X, b.Y, b.Lvl, 300)
	}
	return NewLaser(b.X, b.Y, b.Lvl, -300)
}

func (b *Player) Move() {
	//Handle Wall Collision
	oldy := b.Y
	b.X += b.Velocity.X
	b.Y += b.Velocity.Y
	CollList, hit := BoxCollide(b, currentLvl)

	if hit == 0 {
		return
	}

	//1-hit-4-squars
	if hit == 40 { //TOPLEFT
		if oldy-b.H/2 >= float64(CollList[0].Y*BOXSIZE+BOXSIZE) {
			y := float64(CollList[0].Y*BOXSIZE + BOXSIZE)
			y = y - (b.Y - b.H/2)
			b.RespondY(y)
		} else {
			x := float64(CollList[0].X*BOXSIZE + BOXSIZE)
			x = x - (b.X - b.W/2)
			b.RespondX(x)
		}
	}
	if hit == 34 { //BOTTOMLEFT
		if oldy+b.H/2 <= float64(CollList[0].Y*BOXSIZE) {
			y := float64(CollList[0].Y * BOXSIZE)
			y = (b.Y + b.H/2) - y
			b.RespondY(-y)
			b.InAir = false
		} else {
			x := float64(CollList[0].X*BOXSIZE + BOXSIZE)
			x = x - (b.X - b.W/2)
			b.RespondX(x)
		}
	}
	if hit == 36 { //TOPRIGHT
		if oldy-b.H/2 >= float64(CollList[0].Y*BOXSIZE+BOXSIZE) {
			y := float64(CollList[0].Y*BOXSIZE + BOXSIZE)
			y = y - (b.Y - b.H/2)
			b.RespondY(y)
		} else {
			x := float64(CollList[0].X * BOXSIZE)
			x = (b.X + b.W/2) - x + 1
			b.RespondX(-x)
		}
	}
	if hit == 33 { //BOTTOMRIGHT
		if oldy+b.H/2 <= float64(CollList[0].Y*BOXSIZE) {
			y := float64(CollList[0].Y * BOXSIZE)
			y = (b.Y + b.H/2) - y
			b.RespondY(-y)
			b.InAir = false
		} else {
			x := float64(CollList[0].X * BOXSIZE)
			x = (b.X + b.W/2) - x + 1
			b.RespondX(-x)
		}
	}

	//X-axis
	if hit == 13 || hit == 5 || hit == 7 || hit == 1 { //Going right
		x := 0.0
		if hit == 7 || hit == 13 {
			x = float64(CollList[len(CollList)-1].X * BOXSIZE)
		} else {
			x = float64(CollList[0].X * BOXSIZE)
		}
		x = ((b.X + b.W/2) - x) + 1
		b.RespondX(-x)
	}
	if hit == 10 || hit == 11 || hit == 14 || hit == 2 { //Going left
		x := 0.0
		if hit == 14 || hit == 11 {
			x = float64(CollList[0].X*BOXSIZE + BOXSIZE)
		} else {
			x = float64(CollList[len(CollList)-1].X*BOXSIZE + BOXSIZE)
		}
		x = (x - (b.X - b.W/2))
		b.RespondX(x)
	}

	//Y-axis
	if hit == 7 || hit == 3 || hit == 11 || hit == 17 { //Going Down
		y := 0.0
		if hit == 7 || hit == 11 {
			y = float64(CollList[len(CollList)-1].Y * BOXSIZE)
		} else {
			y = float64(CollList[0].Y * BOXSIZE)
		}
		y = ((b.Y + b.H/2) - y)
		b.RespondY(-y)
		b.InAir = false
	}
	if hit == 12 || hit == 13 || hit == 14 || hit == 18 { //Going Up
		y := 0.0
		if hit == 14 || hit == 13 {
			y = float64(CollList[0].Y*BOXSIZE + BOXSIZE)
		} else {
			y = float64(CollList[len(CollList)-1].Y*BOXSIZE + BOXSIZE)
		}
		y = (y - (b.Y - b.H/2))
		b.RespondY(y)
	}
}

func (b *Player) RespondX(offset float64) {
	b.X += offset
	b.Velocity.X = 0
	b.Friction(MAP_FRICTIONWALL/fps2(), 'y')
}

func (b *Player) RespondY(offset float64) {
	b.Y += offset
	b.Velocity.Y = 0
	b.Friction(MAP_FRICTIONGROUND/fps2(), 'x')
}

func (b *Player) Friction(u float64, xy int) {

	if xy == 'x' { //Apply Friction to X velocity
		if abs(b.Velocity.X) < u {
			b.Velocity.X = 0
			return
		}
		if b.Velocity.X > 0 {
			b.Velocity.X -= u
		}
		if b.Velocity.X < 0 {
			b.Velocity.X += u
		}
	} else { //Apply friction to Y velocity
		if abs(b.Velocity.Y) < u {
			b.Velocity.Y = 0
			return
		}
		if b.Velocity.Y > 0 {
			b.Velocity.Y -= u
		}
		if b.Velocity.Y < 0 {
			b.Velocity.Y += u
		}
	}
}

//Interface Methods
func (b *Player) GetRect() Rect {
	return Rect{b.X, b.Y, b.W, b.H}
}
