package main

import (
	"sdl"
	"sdl/ttf"
	"os"
	"fmt"
)

//=================================
//Constants
//=================================
const (
	SCREENWIDTH  = 800
	SCREENHEIGHT = 600
	BPP          = 32
	BOXSIZE      = 40
)

//=================================
//Praktiske Funktioner
//=================================
func ExitProgram() {
	sdl.Quit()
	os.Exit(0)
	fmt.Println("kage")
}

func Refresh(surface *sdl.Surface) {
	surface.FillRect(&sdl.Rect{0, 0, SCREENWIDTH, SCREENHEIGHT}, 0x000000)
}

func loadImage(path string) *sdl.Surface {
	loaded := sdl.Load(path)
	
	if loaded == nil {
		panic(sdl.GetError())
	}
	defer loaded.Free()
	return sdl.DisplayFormat(loaded)
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

func TileToVector(x, y float64, offset Vector64) Vector64 {
	t := Vector64{x * BOXSIZE, y * BOXSIZE}
	return offset.Add(t)
}

//=================================
//Interfaces
//=================================
type Recter interface {
	GetRect() Rect
}

//=================================
//Globals
//=================================

//States
const (
	SELECT = 1 << iota
	FLIPTILE
	DRAGGING
	ADDWAYPOINT
	REMOVEOBJECT
	CLEARWAYPOINTS
	PRINT
)

var (
	selected *ScreenObject
	screen *sdl.Surface
	state int
	font *ttf.Font
	objects ScreenObjects
)

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

	sdl.WM_SetCaption("MapMaker", "")

	if ttf.Init() != 0 {
	    panic(sdl.GetError())
	}

	font = ttf.OpenFont("../Fontin Sans.otf", 15)

	if font == nil {
		panic(sdl.GetError())
	}
	
	objects = ScreenObjects{}

	//Load Media
	IM = LoadImages()

	MainLoop()
	
	ExitProgram()
}

func MainLoop() {
	//Initializing More stuff
	lvl := Level_x()
	lvl.CreateSurface(screen)
	But := SetButtons()
	state = SELECT
	fmt.Println("state: Select")
	
	for ev := sdl.WaitEvent(); ev != nil; ev = sdl.WaitEvent() {
		switch e := ev.(type) {
		case *sdl.QuitEvent:
			return
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN {
				mouse := Rect{float64(e.X), float64(e.Y), 0,0}
				
				if state != DRAGGING || state != DRAGGING | ADDWAYPOINT {
					if But.TryClick(mouse) {
						if state != PRINT {
							continue
						}
						//Print
						state = SELECT
						j := NewJSONlevel(objects)
						j.Print(lvl)
						continue
					}
				}
				
				switch state {
				case FLIPTILE:
					x := int(mouse.X/BOXSIZE)
					y := int(mouse.Y/BOXSIZE)
					if x >= 18 {break}
					
					if lvl.Grid[y][x] == 'W' {
						lvl.Grid[y][x] = '.'
					} else {
						lvl.Grid[y][x] = 'W'
					}
					lvl.CreateSurface(screen)
				case SELECT:
					old := selected
					so := objects.Select(mouse)
					if so == nil {break}
					if old == so {
						fmt.Println("State: DRAGGING")
						state = DRAGGING
					}
				case DRAGGING:
					fmt.Println("State: SELECT")
					state = SELECT
				
				case ADDWAYPOINT:
					so := objects.Select(mouse)
					if so == nil {break}
					fmt.Println("State: DRAGGING | ADDWAYPOINT")
					so.OldX, so.OldY = so.X, so.Y
					state = DRAGGING | ADDWAYPOINT
					
				case DRAGGING | ADDWAYPOINT:
					fmt.Println("State: SELECT")
					state = SELECT
					if selected != nil {
						selected.Waypoints = append(selected.Waypoints, Vector64{mouse.X, mouse.Y})
						selected.X = selected.OldX
						selected.Y = selected.OldY
						
						fmt.Println(selected)
					}
					
				case REMOVEOBJECT:
					fmt.Println("State: SELECT")
					state = SELECT
					so := objects.Select(mouse)
					if so == nil {break}
					objects = objects.RemoveValue(so)
					selected = nil
					
				case CLEARWAYPOINTS:
					fmt.Println("State: SELECT")
					state = SELECT
					so := objects.Select(mouse)
					if so == nil {break}
					so.Waypoints = []Vector64{}
					fmt.Println(selected)			
				}					
			}
			
		case *sdl.MouseMotionEvent:
			mouse := Rect{float64(e.X), float64(e.Y), 0,0}
			switch state {
			case DRAGGING, DRAGGING | ADDWAYPOINT:
				if selected != nil {
					selected.X = mouse.X
					selected.Y = mouse.Y
				}					
			}
			
		case *sdl.KeyboardEvent:
			if selected == nil {
				break
			}
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_1: selected.Type = SO_PLAYER
				case sdl.K_2: selected.Type = SO_GHOST
				case sdl.K_3: selected.Type = SO_MALBORO
				case sdl.K_4: selected.Type = SO_DEATHSHELL
				case sdl.K_5: selected.Type = SO_ROBOT
				case sdl.K_6: selected.Type = SO_SPIKE
				case sdl.K_7: selected.Type = SO_TELEPORT_H
				case sdl.K_8: selected.Type = SO_TELEPORT_V
				}
			}
		}
		
		//Blitting
		lvl.Blit(screen)
		But.Blit(screen)
		objects.Blit(screen)
		
		screen.Flip()	
	}
}

func SetButtons() Buttons {
	b1 := NewButton(760, 50, 60, 30, "Flip Tile", func() {
		if state != FLIPTILE {
			state = FLIPTILE
			fmt.Println("State: FlipTile")
		} else {
			state = SELECT
			fmt.Println("State: SELECT")
		}
	})
	
	b2 := NewButton(760, 100, 60, 30, "Add Obj", func() {
		objects = append(objects, NewScreenObject())
	})
	
	b3 := NewButton(760, 150, 60, 30, "Rm Obj", func() {
		if state != REMOVEOBJECT {
			state = REMOVEOBJECT
			fmt.Println("State: Remove Object")
		} else {
			state = SELECT
			fmt.Println("State: SELECT")
		}
	})
	
	b4 := NewButton(760, 200, 60, 30, "Add wp", func() {
		if state != ADDWAYPOINT {
			state = ADDWAYPOINT
			fmt.Println("State: Add Waypoint")
		} else {
			state = SELECT
			fmt.Println("State: SELECT")
		}
	})
	
	b5 := NewButton(760, 250, 60, 30, "Clear wp", func() {
		if state != CLEARWAYPOINTS {
			state = CLEARWAYPOINTS
			fmt.Println("State: Clear Waypoints")
		} else {
			state = SELECT
			fmt.Println("State: SELECT")
		}
	})
	
	b6 := NewButton(760, 350, 60, 30, "Print", func() {
		if state != PRINT {
			state = PRINT
		}
	})
	
	
	
	return Buttons{b1,b2,b3,b4,b5, b6}
}
