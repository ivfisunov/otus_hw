package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	// Place your code here
	var wg sync.WaitGroup
	var mu sync.Mutex
	var numberOfWorkers int = n
	tasksCount := 0
	errorsCount := 0

	if n <= 0 {
		numberOfWorkers = 1
	}

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				currentTaskIndex := tasksCount
				tasksCount++
				mu.Unlock()
				if currentTaskIndex >= len(tasks) {
					return
				}
				if err := tasks[currentTaskIndex](); err != nil {
					if m <= 0 { // ignore errors
						continue
					}
					mu.Lock()
					errorsThreshold := errorsCount
					errorsCount++
					mu.Unlock()
					if errorsThreshold >= m {
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	if m <= 0 {
		return nil
	}
	if errorsCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
