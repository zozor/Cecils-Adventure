package main

import "sdl"
import "json"
import "os"
import "fmt"

type Level struct {
    Monst Monsters
    Objs	Objects
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

func (l *Level) HitWall(y,x int) bool {
	if l.Grid[y][x] != '.' {
		return true
	}
	return false
}

func strToArray(s string) [20]int {
	var k [20]int
	for i:=0;i<20;i++ {
		k[i] = int(s[i])
	}
	return k
}




type Maps []func()(*Level)
var MapList Maps = Maps{
	func()(*Level) {
		l := new(Level)
		
		//Setup Map
		l.Grid = [15][20]int{strToArray("WWWWWWWWWWWWWWWWWWQQ"),
							 strToArray("W.................QQ"),
							 strToArray("W.................QQ"),
							 strToArray("W..W..............QQ"),
							 strToArray("W..WWWWW..........QQ"),
							 strToArray("W...............WWQQ"),
							 strToArray("W.................QQ"),
							 strToArray("W..............W..QQ"),
							 strToArray("W............W....QQ"),
							 strToArray("W....WWWW.........QQ"),
							 strToArray("W.................QQ"),
							 strToArray("W...........W.....QQ"),
							 strToArray("W.......WW..W...WWQQ"),
							 strToArray("W...........W...WWQQ"),
			  				 strToArray("WWWWWWWWWWWWWWWWWWQQ"),}
	    
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
	    
	    player.X = 100
	    player.Y = 100
	    
	    l.TimeLimit = 60
		
		return l
	},
	
	func ()(*Level) {
		l := new(Level)
		l.Grid = [15][20]int{[20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 87, 87, 46, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 87, 87, 87, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 87, 87, 46, 46, 46, 46, 46, 46, 46, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}}
		
		player.X = 79
		player.Y = 118
		
		l.Monst = Monsters{
			NewMonsterGhost(358.000000, 85.000000, 80, []Vector64{Vector64{X:263, Y:86}, Vector64{X:283, Y:200}, Vector64{X:461, Y:193}, Vector64{X:451, Y:86}}),
			NewMonsterGhost(658.000000, 80.000000, 80, []Vector64{Vector64{X:661, Y:255}}),
			NewMonsterGhost(412.000000, 243.000000, 80, []Vector64{Vector64{X:297, Y:244}}),
			NewMonsterGhost(666.000000, 523.000000, 80, []Vector64{}),
			NewMonsterGhost(107.000000, 452.000000, 80, []Vector64{}),
		}
		
		l.Objs = Objects{
			NewSpike(Vector64{287.000000, 403.000000}, 80, false, []Vector64{Vector64{280,141}}),
			NewSpike(Vector64{140.000000, 346.000000}, 80, false, []Vector64{Vector64{143,462}}),
			NewSpike(Vector64{396.000000, 428.000000}, 80, false, []Vector64{Vector64{443,404},Vector64{439,341}}),
		}
			
		l.TimeLimit = 45
		return l
	},
	
	func ()(*Level) {
		l := new(Level)
		l.Grid = [15][20]int{[20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 87, 46, 46, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 87, 87, 87, 87, 87, 87, 87, 46, 46, 46, 46, 46, 87, 87, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 46, 87, 87, 87, 87, 87, 87, 87, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 46, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}}
	
		player.X = 179
		player.Y = 524
		
		
		
		t1 := NewTeleport(Vector64{288.000000, 523.000000}, Vector64{1,0}, 2, TELE_VERTICAL, 1)
		t2 := NewTeleport(Vector64{713.000000, 83.000000}, Vector64{-1,0}, 2, TELE_VERTICAL, 1)
		t1.Connect(t2)
		t2.Connect(t1)
		
	    speed := 80.0
	    jumpback := false
	    
		l.Monst = Monsters{
			NewMonsterGhost(681.000000, 525.000000, speed, []Vector64{Vector64{X:327, Y:523}}),
			NewMonsterGhost(664.000000, 252.000000, speed, []Vector64{Vector64{X:593, Y:253}}),
			NewMonsterGhost(57.000000, 531.000000, speed, []Vector64{}),
			NewMonsterGhost(141.000000, 131.000000, speed, []Vector64{}),
			NewMonsterGhost(357.000000, 245.000000, speed, []Vector64{Vector64{X:361, Y:83}}),
		}
		
		l.Objs = Objects{
			NewSpike(Vector64{302.000000, 423.000000}, speed, jumpback, []Vector64{Vector64{X:302, Y:342}}),
			NewSpike(Vector64{397.000000, 341.000000}, speed, jumpback, []Vector64{Vector64{X:396, Y:424}}),
			NewSpike(Vector64{192.000000, 138.000000}, speed, jumpback, []Vector64{Vector64{X:380, Y:141}}),
			t1,
			t2,
		}
							
		l.TimeLimit = 60
		return l
	},
	
	func ()(*Level) {
		l := new(Level)
		l.Grid = [15][20]int{[20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 87, 87, 87, 46, 46, 81, 81}, [20]int{87, 87, 46, 46, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 46, 46, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 87, 87, 87, 87, 87, 87, 87, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 87, 81, 81}, [20]int{87, 87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 87, 81, 81}, [20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}}
		
	    speed := 80.0
	    jumpback := false
	    
		l.Monst = Monsters{
			NewMonsterGhost(170.000000, 91.000000, speed, []Vector64{Vector64{X:71, Y:86}}),
			NewMonsterGhost(590.000000, 90.000000, speed, []Vector64{Vector64{X:695, Y:88}}),
			NewMonsterGhost(562.000000, 454.000000, speed, []Vector64{Vector64{X:561, Y:321}}),
			NewMonsterGhost(199.000000, 457.000000, speed, []Vector64{Vector64{X:196, Y:323}}),
			NewMonsterGhost(372.000000, 323.000000, speed, []Vector64{}),
			NewMonsterGhost(378.000000, 214.000000, speed, []Vector64{}),
		}
		
		l.Objs = Objects{
			NewSpike(Vector64{337.000000, 164.000000}, speed, jumpback, []Vector64{Vector64{X:419, Y:165}}),
			NewSpike(Vector64{420.000000, 61.000000}, speed, jumpback, []Vector64{Vector64{X:419, Y:143}}),
			NewSpike(Vector64{341.000000, 60.000000}, speed, jumpback, []Vector64{Vector64{X:339, Y:144}}),			
			NewSpike(Vector64{703.000000, 382.000000}, speed, jumpback, []Vector64{Vector64{X:706, Y:292}}),
			NewSpike(Vector64{60.000000, 382.000000}, speed, jumpback, []Vector64{Vector64{X:61, Y:297}}),
		}
		
		
		player.X = 379
		player.Y = 526
							
		l.TimeLimit = 60
		return l
	},
	
	func ()(*Level) {
		l := new(Level)
		l.TimeLimit = 60
		
		l.Grid = [15][20]int{[20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 87, 87, 87, 87, 46, 46, 46, 87, 87, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 46, 46, 46, 46, 46, 46, 87, 46, 46, 46, 46, 46, 46, 87, 87, 81, 81}, [20]int{87, 46, 46, 46, 87, 46, 46, 46, 46, 87, 46, 46, 46, 46, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 87, 87, 46, 46, 46, 46, 46, 46, 46, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 87, 87, 87, 87, 46, 46, 46, 87, 87, 87, 87, 46, 46, 46, 81, 81}, [20]int{87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 81, 81}, [20]int{87, 87, 87, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 87, 87, 81, 81}, [20]int{87, 46, 87, 87, 46, 46, 87, 46, 46, 46, 46, 46, 87, 46, 46, 87, 87, 46, 81, 81}, [20]int{87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 87, 81, 81}}
	
		speed := 80.0
	    
	    t1 := NewTeleport(Vector64{562.000000, 550.000000}, Vector64{0,-1}, 2, TELE_HORIZONTAL, 2)
		t2 := NewTeleport(Vector64{200.000000, 551.000000}, Vector64{0,-1}, 2, TELE_HORIZONTAL, 2)
		t3 := NewTeleport(Vector64{409.000000, 81.000000}, Vector64{1,0}, 2, TELE_VERTICAL, 1)
		t4 := NewTeleport(Vector64{351.000000, 82.000000}, Vector64{-1,0}, 2, TELE_VERTICAL, 1)
	    
	    t1.ConnectBoth(t3)
	    t2.ConnectBoth(t4)
	    
		l.Monst = Monsters{
			NewRobot(385.000000, 180.000000, speed, ROBOT_ISBOSS),
		}
		
		l.Objs = Objects{
			t1,t2,t3,t4,
		}
		
		
		player.X = 375
		player.Y = 500
		
		return l
		
	},
	/*
	func ()(*Level) {
		
	},
	*/
}

//================================
//JSON
//================================
type JSONlevel struct {
	Grid [15][20]int
	Timelimit int64
	JSONplayer Vector64
	JSONmonsters []JSONmonster
	JSONobjects []JSONobject
	JSONteleports []JSONteleport
}

type JSONmonster struct {
	X,Y float64
	Type int
	Waypoints []Vector64
	Speed float64
	IsBoss int
}

type JSONobject struct {
	X,Y float64
	Type int
	Waypoints []Vector64
	Speed float64
	Restart bool
}

type JSONteleport struct {
	X,Y float64
	Type int
	Way Vector64
	Size int
	Id int
	Boost float64
	ConnectWith int
}

func NewJSONlevel(path string) []*JSONlevel {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d := json.NewDecoder(file)
	
	var out []*JSONlevel
	err = d.Decode(&out)
	if err != nil {
		panic(err)
		return nil
	}
	return out
}

func (j *JSONlevel) Create() (func ()(*Level)) {
	return func ()(*Level) {
		l := new(Level)
		l.TimeLimit = j.Timelimit
		
		player.X = j.JSONplayer.X
		player.Y = j.JSONplayer.Y
		
		l.Grid = j.Grid 
			    	    
		l.Monst = Monsters{}
		l.Objs = Objects{}
		
		
		for i:=0;i<len(j.JSONmonsters);i++ {
			l.Monst = append(l.Monst, JSONgetMonster(j.JSONmonsters, i))
		}
		
		for i:=0;i<len(j.JSONobjects);i++ {
			l.Objs = append(l.Objs, JSONgetObject(j.JSONobjects, i))
		}
		
		teleports := make(map[int]*Teleport)
		
		
		for _, v := range j.JSONteleports {
			teleports[v.Id] = NewTeleportJSON(v)
		}
		
		var to *Teleport
		var ok bool
		for _, v := range j.JSONteleports {
			to, ok = teleports[v.ConnectWith]
			if ok {
				teleports[v.Id].Connect(to)
			}
			l.Objs = append(l.Objs, teleports[v.Id])
		}
				
		return l
	}
}
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

func JSONgetMonster(m []JSONmonster, index int) Monster {

	var out Monster
	switch m[index].Type {
	case SO_GHOST: out = NewGhostJSON(m[index])
	case SO_ROBOT: out = NewRobotJSON(m[index])
	case SO_DEATHSHELL: out = NewDeathShellJSON(m[index])
	case SO_MALBORO: out = NewMalboroJSON(m[index])
	}
	return out
}

func JSONgetObject(m []JSONobject, index int) Object {
	var out Object
	switch m[index].Type {
	case SO_SPIKE: out = NewSpikeJSON(m[index])
	}
	return out
}
/*
	l.TimeLimit = 60
	speed := 80.0
    jumpback := false
    
	l.Monst = Monsters{
	}
	
	l.Objs = Objects{
	}
	
*/

