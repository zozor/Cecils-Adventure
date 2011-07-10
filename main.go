/*
TODO liste:
-Mere grafik (Sådan det ser pænere ud)
-Lav en masse baner :D
-lyd
*/

package main

import (
	"fmt"
	"sdl"
	"math"
	"time"
	//"sdl/ttf"
	"sdl/mixer"
	"os"
)

//=================================
//Konstanter
//=================================
const (
	SCREENWIDTH  = 800
	SCREENHEIGHT = 600
	BPP          = 32
	BOXSIZE      = 40
)

const (
	G_EXITPROGRAM = iota
	G_NEXTLEVEL
	G_GAMEOVER
	G_ENDGAME
)

//const GRIDX = 20
//const GRIDY = 15


//=================================
//Typer
//=================================

type Bliter interface {
	Blit(surface *sdl.Surface)
}
type Blites []Bliter

//Vector
type Vector struct {
	X int
	Y int
}

func (v Vector) Length() int {
	return int(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vector) Add(vec Vector) Vector {
	v.X += vec.X
	v.Y += vec.Y
	return v
}

func (v Vector) Subtract(vec Vector) Vector {
	v.X -= vec.X
	v.Y -= vec.Y
	return v
}

func (v Vector) Scalar(scal int) Vector {
	v.X *= scal
	v.Y *= scal
	return v
}

//Vector with floats!
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

//Rect
type Rect struct {
	X, Y, W, H float64
}

func GetRect(rect sdl.Rect) (r Rect) {
	r.X = float64(rect.X)
	r.W = float64(rect.W)
	r.Y = float64(rect.Y)
	r.H = float64(rect.H)
	return
}

type Recter interface {
	GetRect() Rect
}

//=================================
//Globals
//=================================
var (
	screen *sdl.Surface
	background *sdl.Surface
	debug *sdl.Surface
	currentLvl *Level
	player *Player
	DeltaTime int64 = now()
	checkpoint int = 0
)

//=================================
//Praktiske Funktioner
//=================================

//Time Funktioner
func now() int64 {
	return time.Nanoseconds()
}

//Returns fps
func fps() float64 {
	return float64(1e9/(now()-DeltaTime))
}

//Returns fps()^2
func fps2() float64 {
	return fps()*fps()
}

func difftime(t int64) int64 {
	return now()-t
}



func Refresh(surface *sdl.Surface) {
	surface.FillRect(&sdl.Rect{0, 0, SCREENWIDTH, SCREENHEIGHT}, 0x000000)
}

func abs(i float64) float64 {
	if i < 0 {
		return -i
	}
	return i
}

//Loads and converts Image
func loadImage(path string) *sdl.Surface {
	loaded := sdl.Load(path)
	
	if loaded == nil {
		panic(sdl.GetError())
	}
	defer loaded.Free()
	return sdl.DisplayFormat(loaded)
}

func loadMusic(path string) *mixer.Music {
	load :=  mixer.LoadMUS(path)
	if load == nil {
		panic(sdl.GetError())
	}
	return load
}

func loadSound(path string) *mixer.Chunk {
	load := mixer.LoadWAV(path)
	if load == nil {
		panic(sdl.GetError())
	}
	return load
}

func HitTest(obj1, obj2 Recter) bool {
	a := obj1.GetRect()
	b := obj2.GetRect()

	if a.Y+a.H/2 > b.Y-b.H/2 && a.Y-a.H/2 < b.Y+b.H/2 {
		if a.X-a.W/2 < b.X+b.W/2 && a.X+a.W/2 > b.X-b.W/2 {
			return true
		}
	}
	return false
}

func BoxCollide(a Recter, l *Level) ([]Vector, int) {
	rect := a.GetRect()

	//OutOfBounds
	if rect.X < 0 || rect.Y < 0 || rect.Y+rect.H > SCREENHEIGHT || rect.X+rect.W > SCREENWIDTH {
		return nil, 0
	}

	x := make([]int, 0, 4)
	y := make([]int, 0, 4)

	min := int(rect.X - rect.W/2) / BOXSIZE
	max := int(rect.X + rect.W/2) / BOXSIZE
	for ; min <= max; min += 1 {
		x = append(x, min)
	}

	min = int(rect.Y - rect.H/2) / BOXSIZE
	max = int(rect.Y + rect.H/2) / BOXSIZE
	for ; min <= max; min += 1 {
		y = append(y, int(min))
	}

	response := 0
	V := make([]Vector, 0, 10)

	for i := 0; i < len(y); i++ {
		for j := 0; j < len(x); j++ {
			response = response << 1
			if l.HitWall(y[i], x[j]) {
				response = response | 1
				V = append(V, Vector{x[j], y[i]})
			}
		}
	}
	if len(y) == 2 && len(x) == 1 && len(V) != 0 {
		response = response | 16
	} else if len(V) == 1 && len(y) != 1 {
		response = response | 32
	}

	return V, response
}

func TileToVector(x, y float64, offset Vector64) Vector64 {
	t := Vector64{x * BOXSIZE, y * BOXSIZE}
	return offset.Add(t)
}

func ExitProgram() {
	sdl.Quit()
	os.Exit(0)
	fmt.Println("kage")
}

//=================================
//Initialize Lortet
//=================================
func main() {
	//INIT SDL
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}

	screen = sdl.SetVideoMode(SCREENWIDTH, SCREENHEIGHT, BPP, sdl.SWSURFACE)
	if screen == nil {
		panic(sdl.GetError())
	}

	sdl.WM_SetCaption("KAGE!", "")

	/*if ttf.Init() != 0 {
	    panic(sdl.GetError())
	}

	font := ttf.OpenFont("Fontin Sans.otf", 15)

	if font == nil {
		panic(sdl.GetError())
	}*/

	//debug = ttf.RenderText_Blended(font, "Test (with music)", sdl.Color{255, 255, 255, 0})
	
	if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 1096) != 0 {
		panic(sdl.GetError())
	}

	sdl.ShowCursor(sdl.DISABLE)

	//Load Media
	IM = LoadImages()
	MU = LoadSound()
	
	//Load Maps
	baner := NewJSONlevel("levels")
	for _, v := range baner {
		MapList = append(MapList, v.Create())
	}

	MainLoop()
	
	ExitProgram()
}

func MainLoop() {
	outcode := -1
	option := 0
	lvl := 0

	Refresh(screen)
	screen.Blit(&sdl.Rect{0, 0, 0, 0}, IM.BG[0], nil)
	screen.Flip()
	
	//MU.Music[0].PlayMusic(-1)

	for running := true; running; {
		for ev := sdl.WaitEvent(); ev!=nil; ev = sdl.WaitEvent() {
			Refresh(screen)
			screen.Blit(&sdl.Rect{0, 0, 0, 0}, IM.BG[0], nil)
			if option == 0 {
				screen.FillRect(&sdl.Rect{320, 212, 10, 10}, 0xFF00FF)
			} else {
				screen.FillRect(&sdl.Rect{320, 233, 10, 10}, 0xFF00FF)
			}
			
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				running = false
				return
			case *sdl.KeyboardEvent:
				if e.Type != sdl.KEYDOWN {
					break
				}
				switch e.Keysym.Sym {
				case sdl.K_SPACE:
					if option == 0 {
						goto game
					} else {
						screen.FillRect(&sdl.Rect{100, 200, 100, 10}, 0xFFFFFF)
						screen.FillRect(&sdl.Rect{100, 220, 100, 10}, 0xFFFFFF)
						screen.FillRect(&sdl.Rect{100, 200, 10, 20}, 0xFFFFFF)
						screen.FillRect(&sdl.Rect{200, 200, 10, 30}, 0xFFFFFF)
						screen.Flip()
						var ok bool
						lvl, ok = Password()
						if ok {
							goto game
						}
					}
				case sdl.K_UP:
					option = 0
				case sdl.K_DOWN:
					option = 1
				}
			}
			screen.Flip()
		}

		continue

		//Gameing
	game:
		outcode = MainGame(lvl)

		if outcode == G_EXITPROGRAM { //Exit Program
			running = false
			return
		}

		if outcode == G_NEXTLEVEL { //NextLevel
			screen.Blit(&sdl.Rect{350, 250, 0, 0}, IM.Misc[2], nil)
			screen.Flip()
			sdl.Delay(1000)
			for PressKey() != sdl.K_SPACE {
			}
			lvl++
			if lvl % 4 == 0 {
				checkpoint = lvl
			}
			goto game
				
		}

		if outcode == G_GAMEOVER { //Gameover
			screen.Blit(&sdl.Rect{350, 250, 0, 0}, IM.Misc[1], nil)
			screen.Flip()
			for PressKey() != sdl.K_SPACE {
			}
			player = nil
			lvl = checkpoint
			goto game
			
		}

		if outcode == G_ENDGAME { //Game Completed
			screen.Blit(&sdl.Rect{350, 250, 0, 0}, IM.Misc[3], nil)
			screen.Flip()
			for PressKey() != sdl.K_SPACE {
			}
			lvl = 0
		}
	}
}

func PressKey() int {
	for ev := sdl.WaitEvent(); ev!=nil; ev = sdl.WaitEvent() {
		switch e := ev.(type) {
		case *sdl.QuitEvent:
			ExitProgram()
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				return int(e.Keysym.Sym)
			}
		}
	}
	return 0
}

func Password() (int, bool) {
	input := []byte{}
	for {
		in := byte(PressKey())
		if in >= 48 && in <= 57 {
			screen.FillRect(&sdl.Rect{115 + int16(len(input))*10, 212, 5, 5}, 0xFFFFFF)
			screen.Flip()
			input = append(input, in)
		} else {
			return 0, false
		}
		
		if len(input) == 8 {
			break
		}
	}
	switch string(input) {
	case "11111111":
		return 1, true
	case "22222222":
		return 2, true
	case "33333333":
		return 3, true
	case "44444444":
		return 4, true
	case "55555555":
		return 5, true
	case "66666666":
		return 6, true
	case "77777777":
		return 7, true
	case "88888888":
		return 8, true
	case "99999999":
		return 9, true
	}
	return 0, false
}
