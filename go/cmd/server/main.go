package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/hbx286/orb-arena/internal/game"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	room := game.NewRoom()
	room.AddPlayer(&game.Player{ID: 1, X: 100, Y: 100, DirX: 1, DirY: 0})
	room.AddPlayer(&game.Player{ID: 2, X: 900, Y: 500, DirX: -0.6, DirY: 0.8})
	// 在 Run 之前起一个并发任务，因为 Run 会一直占住当前 goroutine
	go func() {
		time.Sleep(2 * time.Second)
		ok := room.SendInput(game.Input{ID: 1, DirX: 0, DirY: 1}) // Y 增大 = 向下（和画布坐标一致）
		fmt.Println("input sent:", ok)
	}()

	room.Run(ctx)
}
