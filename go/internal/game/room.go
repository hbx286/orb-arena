package game

import (
	"context"
	"fmt"
	"time"
)

const (
	TickRate = 20
	TickDur  = time.Second / TickRate
	MapSize  = 1000.0
	Speed    = 200.0
)

type Player struct {
	ID         int
	X, Y       float64
	DirX, DirY float64
}

type Room struct {
	players []*Player
	ticks   int
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) AddPlayer(player *Player) {
	r.players = append(r.players, player)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// update 每个tick调用一次: 让所有的玩家按方向移动
func (r *Room) update() {
	dt := TickDur.Seconds()
	for _, p := range r.players {
		p.X += p.DirX * Speed * dt
		p.Y += p.DirY * Speed * dt
		p.X = clamp(p.X, 0, MapSize)
		p.Y = clamp(p.Y, 0, MapSize)
	}
}

func (r *Room) Run(ctx context.Context) {
	ticker := time.NewTicker(TickDur)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("room stopped, ticks:", r.ticks)
			return
		case <-ticker.C:
			r.ticks++
			r.update()
			if r.ticks%TickRate == 0 {
				for _, p := range r.players {
					fmt.Printf("player %d at (%.1f, %.1f)\n", p.ID, p.X, p.Y)
				}
			}
		}
	}
}
