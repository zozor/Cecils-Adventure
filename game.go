package main

import (
	"sdl"
	"fmt"
)

func AOISDAOISD() {
	fmt.Println("KAS")
}

var gravity Vector64
const GRAVITYCONSTANT = 600

func MainGame(lvl int) int {
	//=================================
    //Init Stuff
    //=================================
   	if player == nil {
    	player = NewPlayer()
    }
    currentLvl = MapList[lvl]()
    currentLvl.CreateSurface(screen)
    lasers := Lasers{}
    elasers := Lasers{}
    timer := GameTime{now(), now()+int64(currentLvl.TimeLimit*1e9)}
    var WinDelay int64
    
    DeltaTime = now()-50e6
    for running := true; running; {
        Refresh(screen)
        //=================================
        //Events
        //=================================	    
	    for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
	        switch e := ev.(type) {
	        case *sdl.QuitEvent:
	            return G_EXITPROGRAM
	        /*
	        case *sdl.KeyboardEvent:
	        	if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_p {
	        		d := difftime(DeltaTime)
	        		for PressKey() != sdl.K_p {
	        		}
	        		DeltaTime = now()-d
	        	}*/
	        }
	    }
	    
	    if fps() < 20 {
	    	DeltaTime = now()
	    	continue
	    }
	    
	    //=================================
        //Object Events
        //=================================
        gravity.Y = GRAVITYCONSTANT/fps2()
        
        player.Events()
        lasers = lasers.Clean()
        elasers = elasers.Clean()
        currentLvl.Objs = currentLvl.Objs.Clean()
	    
	    if sdl.GetKeyState()[sdl.K_SPACE] != 0 {
	    	t := player.Shoot()
	    	if t != nil {
	    		lasers = append(lasers, t)
	    	}
	    }
	    
	    for _, v := range currentLvl.Monst.Attack() {
	    	elasers = append(elasers, v)
	    }
	    
	    //=================================
        //Moving Objects
        //=================================
	    player.Move()
	    currentLvl.Monst.Move()
	    currentLvl.Objs.Move()
	    lasers.Move()
	    elasers.Move()
        
	    //=================================
        //Collission
        //=================================
        for im:=0; im<len(currentLvl.Monst); im++ {
        	mon := currentLvl.Monst[im]
        	for il:=0; il<len(lasers); il++ {        		
        		las := lasers[il]
        		if HitTest(mon, las) {
        			if mon.Damage(las.Damage()) { //if dies
        				//Spoils
        				spoil := NewSpoils(mon.GetRect().X,mon.GetRect().Y)
        				if !spoil.IsDead() {
        					currentLvl.Objs = append(currentLvl.Objs, spoil)
        				}
        				//Removing
        				currentLvl.Monst = currentLvl.Monst.Remove(im)
        				im--
        				
        			}
        			lasers = lasers.Remove(il)
        			il--
        		}
        	}
        }
        
        if !player.IsImmortal() {
	        for _,v := range currentLvl.Monst {
	        	if HitTest(player, v) {
	        		player.WasHit()
	        	}
	        }
	        for _, v := range elasers {
	        	if HitTest(player, v) {
	        		player.WasHit()
	        	}
	        }
	    }
        
        lasers = lasers.HitWall(currentLvl)
        
        currentLvl.Objs.Hit(player)
	    
	    DeltaTime = now()
	    //=================================
        //Blitting
        //=================================
        currentLvl.Blit(screen)
        player.Blit(screen)
        currentLvl.Monst.Blit(screen)
        currentLvl.Objs.Blit(screen)
        lasers.Blit(screen)
        elasers.Blit(screen)
        BlitHeartsAndTimer(screen, player.Lives, timer)
        
        //=================================
        //Win Lose Conditions
        //=================================
        if WinDelay == 0 {
	        if player.Lives == 0 || difftime(timer.End) > 0 {
	        	if len(currentLvl.Monst) != 0 {
	        		return G_GAMEOVER
	        	}
	        }
	        if len(currentLvl.Monst) == 0 {
	        	WinDelay = now()+2e9
	        }
	    } else if now()-WinDelay > 0 {
        	if lvl == len(MapList)-1 {
        		return G_ENDGAME
        	}
        	return G_NEXTLEVEL
        }
        //=================================
        //FrameSetting
        //=================================
        screen.Flip()
        f := now()-DeltaTime
        if f/1e6 < 30 {
        	sdl.Delay(uint32(30-(f/1e6)))
        }
	}
	return G_NEXTLEVEL
}