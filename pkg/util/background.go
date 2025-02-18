package util

import (
	"sync"

	"github.com/thaian1234/green_light/pkg/logger"
)

func Background(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.Error("handle panic in goroutine", "msg", err)
			}
		}()

		fn()
	}()
}
