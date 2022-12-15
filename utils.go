package utils

import "time"

// Schedule запускает выполнение заданный процедуры с заданным интервалом времени
func Schedule(fn func(), interval time.Duration, done <-chan bool) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				fn()
			case <-done:
				return
			}
		}
	}()
	return ticker
}
