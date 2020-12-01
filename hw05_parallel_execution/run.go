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
	var numberOfWorkers int = n
	var numberOfErrors int = 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	chanWithTasks := make(chan Task, len(tasks))

	if n <= 0 {
		numberOfWorkers = 1
	}

	go func() {
		for _, task := range tasks {
			chanWithTasks <- task
		}
		close(chanWithTasks)
	}()

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localNumberOfErrors := 0
			for task := range chanWithTasks {
				if err := task(); err != nil {
					if m <= 0 {
						continue
					}
					mu.Lock()
					numberOfErrors++
					localNumberOfErrors = numberOfErrors
					mu.Unlock()
					if localNumberOfErrors >= m {
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	if numberOfErrors >= m && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
