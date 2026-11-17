package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/hbx286/orb-arena/internal/game"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	room := game.NewRoom()
	room.AddPlayer(&game.Player{ID: 1, X: 100, Y: 100, DirX: 1, DirY: 0})
	room.AddPlayer(&game.Player{ID: 2, X: 900, Y: 500, DirX: -0.6, DirY: 0.8})

	room.Run(ctx)
}
