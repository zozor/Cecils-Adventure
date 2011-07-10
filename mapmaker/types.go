package main

import (
	"sdl"
	"fmt"
	"math"
	"sdl/ttf"
	"strings"
	"json"
)
var saskjdasndakjsdn string = fmt.Sprint("kage")

type Buttons []*Button

//======================================
//CONTAINERS
//======================================
func (m Buttons) Blit(surface *sdl.Surface) {
	for i := 0; i < len(m); i++ {
		m[i].Blit(surface)
	}
}
func (m Buttons) Remove(i int) Buttons {
	return append(m[0:i], m[i+1:len(m)]...)
}

func (m Buttons) TryClick(r Rect) bool {
	for i := 0; i < len(m); i++ {
		if HitTest(r, m[i]) {
			m[i].OnClick()
			return true
		}
	}
	return false
}

type ScreenObjects []*ScreenObject

func (m ScreenObjects) Blit(surface *sdl.Surface) {
	for i := 0; i < len(m); i++ {
		m[i].Blit(surface)
	}
}

func (m ScreenObjects) Select(r Rect) *ScreenObject {
	selected = nil
	for i := 0; i < len(m); i++ {
		if HitTest(r, m[i]) {
			selected = m[i]
			return m[i]
		}
	}
	return nil
}

func (m ScreenObjects) Remove(i int) ScreenObjects {
	return append(m[0:i], m[i+1:len(m)]...)
}

func (m ScreenObjects) RemoveValue(v *ScreenObject) ScreenObjects {
	for i := 0; i < len(m); i++ {
		if m[i] == v {
			return m.Remove(i)
		}
	}
	return m
}

func (m ScreenObjects) String() string {
	s := ""
	for i := 0; i < len(m); i++ {
		s += m[i].String()+"\n"
	}
	return s
}

//======================================
//BUTTON
//======================================
type Button struct {
	X, Y, W, H float64
	Image *sdl.Surface
	OnClick	func()
}

func NewButton(x,y,w,h float64, text string, Func func()) *Button {
	b := new(Button)
	b.X, b.Y = x,y
	b.W, b.H = w,h
	b.OnClick = Func
	
	//Setup Image
	f := screen.Format
	b.Image = sdl.CreateRGBSurface(sdl.SWSURFACE, int(w), int(h), BPP, f.Rmask,f.Gmask,f.Bmask, f.Amask)
	b.Image.FillRect(&sdl.Rect{0,0,uint16(w),uint16(h)}, 0xFF00FF)
	tekst := ttf.RenderText_Blended(font, text, sdl.Color{255, 255, 255, 0})
	b.Image.Blit(&sdl.Rect{5,5,0,0}, tekst, nil)
	b.Image.Flip()
	
	
	return b
}

func (l *Button) GetRect() Rect {
	return Rect{l.X, l.Y, l.W, l.H}
}

func (l *Button) Blit(surface *sdl.Surface) {
	surface.Blit(&sdl.Rect{int16(l.X-l.W/2),int16(l.Y-l.H/2),0,0}, l.Image, nil);
}


//======================================
//SCREENOBJECT
//======================================
//ScreenObject.Type
const (
	SO_PLAYER = iota
	
	SO_GHOST
	SO_ROBOT
	SO_DEATHSHELL
	SO_MALBORO
	
	SO_SPIKE
	SO_TELEPORT_H
	SO_TELEPORT_V
)

func SOnavn(t int) string {
	switch t {
		case SO_PLAYER: return "Player"
		case SO_GHOST:  return "Ghost"
		case SO_ROBOT:  return "Robot"
		case SO_DEATHSHELL: return "Deathshell"
		case SO_MALBORO:	 return "Malboro"	
		case SO_SPIKE: return "Spike"
		case SO_TELEPORT_H: return "Teleport H"
		case SO_TELEPORT_V: return "Teleport V"
	}
	return ""
}

type ScreenObject struct {
	X, Y, W, H float64
	OldX, OldY float64
	Type int
	Waypoints []Vector64
}

func NewScreenObject() *ScreenObject {
	o := new(ScreenObject)
	o.X, o.Y = 100, 100
	o.Waypoints = []Vector64{}
	return o
}

func (o *ScreenObject) image() *sdl.Surface {
	var img *sdl.Surface
	switch o.Type {
	case SO_PLAYER:
		img = IM.Cecil[0][0]
	case SO_ROBOT:
		img = IM.Monsters[2][0]
	case SO_GHOST:
		img = IM.Monsters[0][0]
	case SO_DEATHSHELL:
		img = IM.Monsters[1][0]
	case SO_MALBORO:
		img = IM.Monsters[3][0]
	case SO_TELEPORT_H:
		img = IM.Obj[0]
	case SO_TELEPORT_V:
		img = IM.Obj[1]
	case SO_SPIKE:
		img = IM.Shots[2]
	}
	o.W = float64(img.W)
	o.H = float64(img.H)
	return img
}

func (l *ScreenObject) Blit(surface *sdl.Surface) {
	if l == selected {
		surface.FillRect(&sdl.Rect{int16(l.X-l.W/2),int16(l.Y-l.H/2),uint16(l.W),uint16(l.H)}, 0xFFFFFF);
	}
	surface.Blit(&sdl.Rect{int16(l.X-l.W/2),int16(l.Y-l.H/2),0,0}, l.image(), nil);
}

func (l *ScreenObject) GetRect() Rect {
	return Rect{l.X, l.Y, l.W, l.H}
}

func (o *ScreenObject) String() string {
	s := ""
	switch o.Type {
	case SO_PLAYER:
		s = fmt.Sprintf("player.X = %d\nplayer.Y = %d", int(o.X), int(o.Y))
	case SO_ROBOT:
		s = fmt.Sprintf("NewRobot(%f, %f, speed, boss)", o.X, o.Y)
	case SO_GHOST:
		s = fmt.Sprintf("NewMonsterGhost(%f, %f, speed, %#v)", o.X, o.Y, o.Waypoints)
	case SO_DEATHSHELL:
		s = fmt.Sprintf("NewDeathshell(%f, %f, speed, boss)", o.X, o.Y)
	case SO_MALBORO:
		s = fmt.Sprintf("NewMalboro(%f, %f)", o.X, o.Y)
	case SO_TELEPORT_H:
		s = fmt.Sprintf("NewTeleport(Vector64{%f, %f}, Vector64{-1,0}, size, TELE_VERTICAL, 1)", o.X, o.Y)
	case SO_TELEPORT_V:
		s = fmt.Sprintf("NewTeleport(Vector64{%f, %f}, Vector64{-1,0}, size, TELE_HORIZONTAL, 1)", o.X, o.Y)
	case SO_SPIKE:
		s = fmt.Sprintf("NewSpike(Vector64{%f, %f}, speed, jumpback, %#v)", o.X, o.Y, o.Waypoints)
	}
	return strings.Replace(s, "main.", "", -1)
}

//======================================
//VECTOR
//======================================
type Vector64 struct {
	X float64
	Y float64
}

func (v Vector64) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector64) Add(vec Vector64) Vector64 {
	v.X += vec.X
	v.Y += vec.Y
	return v
}

func (v Vector64) Subtract(vec Vector64) Vector64 {
	v.X -= vec.X
	v.Y -= vec.Y
	return v
}

func (v Vector64) Scalar(scal float64) Vector64 {
	v.X *= scal
	v.Y *= scal
	return v
}

//======================================
//Rect
//======================================
type Rect struct {
	X,Y,W,H float64
}

func (r Rect) GetRect() Rect {
	return r
}

//================================
//JSON
//================================
type JSONlevel struct {
	Grid string
	Timelimit int64
	JSONplayer Vector64
	JSONmonsters []map[string]interface{}
	JSONobjects []map[string]interface{}
	JSONteleports []map[string]interface{}
}

func NewJSONlevel(obj ScreenObjects) *JSONlevel {
	j := new(JSONlevel)
	j.Grid = "INSERT GRID HERE"
	j.Timelimit = 60
	j.JSONmonsters = []map[string]interface{}{}
	j.JSONobjects = []map[string]interface{}{}
	j.JSONteleports = []map[string]interface{}{}
	for i, v := range obj {
		switch v.Type {
		case SO_PLAYER:
			j.JSONplayer = Vector64{v.X,v.Y}
		case SO_DEATHSHELL, SO_GHOST, SO_MALBORO, SO_ROBOT:
			input := map[string]interface{}{
				"X": v.X,
				"Y": v.Y,
				"Type": v.Type,
				"Waypoints": v.Waypoints,
				"Speed": 80,
				"IsBoss": 0,
				"Name": SOnavn(v.Type),
			}
			j.JSONmonsters = append(j.JSONmonsters, input)
		case SO_SPIKE:
			input := map[string]interface{}{
				"X": v.X,
				"Y": v.Y,
				"Type": v.Type,
				"Waypoints": v.Waypoints,
				"Speed": 80,
				"Restart":false,
				"Name": SOnavn(v.Type),
			}
			j.JSONobjects = append(j.JSONobjects, input)			
		
		case SO_TELEPORT_H, SO_TELEPORT_V:
			input := map[string]interface{}{
				"X": v.X,
				"Y": v.Y,
				"Type": v.Type,
				"Way": Vector64{-1,0},
				"Size": 1,
				"Id": i,
				"Boost": 1,
				"ConnectWith": -1,
				"Name": SOnavn(v.Type),
			}
			j.JSONteleports = append(j.JSONteleports, input)
		}
	}
	return j
}

func (j *JSONlevel) Print(lvl *Level) {
	out, err := json.Marshal(lvl.Grid)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	
	out, err = json.MarshalIndent(j, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
