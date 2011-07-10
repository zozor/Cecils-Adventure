include $(GOROOT)/src/Make.inc

TARG=game
GOFILES=\
	maps.go\
	main.go\
	monsters.go\
	player.go\
	misc.go\
	media.go\
	game.go\
	objects.go\

include $(GOROOT)/src/Make.cmd
