package main

import "sdl"
import "fmt"

type Level struct {
    Grid	[15][20]int
    Surface *sdl.Surface
    
    TimeLimit	int64 //I Sekunder
}

func (l *Level) Blit(surface *sdl.Surface) {
	surface.Blit(&sdl.Rect{0,0,0,0},l.Surface, nil)
}

func (l *Level) CreateSurface(surface *sdl.Surface) {
	f := surface.Format
	l.Surface = sdl.CreateRGBSurface(sdl.SWSURFACE, SCREENWIDTH, SCREENHEIGHT, BPP, f.Rmask,f.Gmask,f.Bmask, f.Amask)
	var x,y int16
	for y=0; y<15; y++ {
		for x=0; x<20; x++ {
			l.Surface.Blit(&sdl.Rect{BOXSIZE*x, BOXSIZE*y, BOXSIZE, BOXSIZE}, IM.Tiles[l.Grid[y][x]], nil)
		}
	}
	l.Surface.Flip()
}

func strToArray(s string) [20]int {
	var k [20]int
	for i:=0;i<20;i++ {
		k[i] = int(s[i])
	}
	return k
}

func (l *Level) String() string {
	return fmt.Sprintf("%#v", l.Grid)
}

func Level_x() *Level {
	l := new(Level)
	
	//Setup Map
	l.Grid = [15][20]int{strToArray("WWWWWWWWWWWWWWWWWWQQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
		  				 strToArray("WWWWWWWWWWWWWWWWWWQQ"),}
    
    /*
    //Setup Monsters
    l.Monst = Monsters{NewMonsterGhost(200, 140, 80, []Vector64{Vector64{300,140}}),
    					NewMonsterGhost(240, 340, 80, []Vector64{Vector64{480,340}}),
    					NewMonsterGhost(680, 460, 80, []Vector64{}),
    }
    
    //Setup Objects
    t1 := NewTeleport(TileToVector(4, 2, Vector64{-10,0}), Vector64{1,0}, 2, TELE_VERTICAL, 1)
    t2 := NewTeleport(TileToVector(18, 4, Vector64{-10,0}), Vector64{-1,0}, 2, TELE_VERTICAL, 1)
    t1.Connect(t2)
    t2.Connect(t1)
    
    s1 := NewSpike(TileToVector(4, 3, Vector64{-10,0}), 80, true, []Vector64{Vector64{600,500}})
    
    l.Objs = Objects{t1,t2, s1}
    
    
    //Setup Misc
    player.X = 100
    player.Y = 100
    
    l.TimeLimit = 60
    */
	
	return l
}

/* CLEAN MAP
	l.Grid = [15][20]int{strToArray("WWWWWWWWWWWWWWWWWWQQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
						 strToArray("W.................QQ"),
		  				 strToArray("WWWWWWWWWWWWWWWWWWQQ"),}
*/

