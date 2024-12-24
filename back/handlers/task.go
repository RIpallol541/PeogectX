package handlers

import (
	"fmt"
	"time"
)

var rateLimiter = make(chan struct{}, 5) // Semaphore with a limit of 5 concurrent operations

// PerformParallelTask performs a task with rate limiting
func PerformParallelTask() {
	rateLimiter <- struct{}{} // Acquire semaphore
	defer func() { <-rateLimiter }() // Release semaphore

	fmt.Println("Starting parallel task")
	time.Sleep(5 * time.Second) // Simulate a long-running task
	fmt.Println("Parallel task completed")
}