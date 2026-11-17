package game

import (
	"context"
	"fmt"
	"math"
	"time"
)

const (
	TickRate    = 20
	TickDur     = time.Second / TickRate
	MapSize     = 1000.0
	Speed       = 200.0
	InputBuffer = 256 // 输入队列容量
)

type Player struct {
	ID         int
	X, Y       float64
	DirX, DirY float64
}

type Input struct {
	ID         int
	DirX, DirY float64
}

// normalizeInput 把方向归一化为单位向量。
// 零向量、NaN、无穷大一律视为"停下"，不信任客户端传来的数据。
func normalizeInput(in Input) Input {
	length := math.Hypot(in.DirX, in.DirY)
	if length == 0 || math.IsNaN(length) || math.IsInf(length, 0) {
		return Input{ID: in.ID}
	}
	return Input{ID: in.ID, DirX: in.DirX / length, DirY: in.DirY / length}
}

type Room struct {
	players map[int]*Player // 按 ID 查找
	ticks   int
	inputs  chan Input
}

func NewRoom() *Room {
	return &Room{
		players: make(map[int]*Player),
		inputs:  make(chan Input, InputBuffer),
	}
}

// SendInput 不阻塞：队列满了返回 false。
func (r *Room) SendInput(in Input) bool {
	select {
	case r.inputs <- in:
		return true
	default:
		return false
	}
}

// AddPlayer 只能在 Run 之前调用，Run 开始后 players 只允许房间自己的 goroutine 修改。
func (r *Room) AddPlayer(p *Player) {
	r.players[p.ID] = p
}

func (r *Room) applyInput(in Input) {
	p, ok := r.players[in.ID]
	if !ok {
		return // 未知玩家，忽略
	}
	in = normalizeInput(in)
	p.DirX, p.DirY = in.DirX, in.DirY
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

// update 每个 tick 调用一次：让所有玩家按方向移动
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
		case in := <-r.inputs:
			r.applyInput(in)
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
