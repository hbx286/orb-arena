package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("exit, tick times:", count)
			return
		case <-ticker.C:
			count++
			if count%20 == 0 {
				fmt.Println("count:", count)
			}
		}

	}
}
