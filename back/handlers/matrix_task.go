package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/RIpallol541/PeogectX/models"
)

var matrixTaskSemaphore = make(chan struct{}, 5) // Limit concurrent matrix tasks

// PerformMatrixMultiplicationTask executes a parallel matrix multiplication task
func PerformMatrixMultiplicationTask(conn *net.UDPConn, addr *net.UDPAddr, params map[string]uint32) {
	matrixTaskSemaphore <- struct{}{} // Acquire semaphore
	defer func() { <-matrixTaskSemaphore }() // Release semaphore

	sizeMatrix := params["sizeMatrix"]
	maxDimension := params["maxDimension"]

	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	results := []map[string]interface{}{}

	for idxDimension := uint32(2); idxDimension <= maxDimension; idxDimension++ {
		for idxLambda := uint32(0); idxLambda <= idxDimension; idxLambda++ {
			for idxMu := uint32(0); idxMu < idxDimension-idxLambda; idxMu++ {
				if (2*idxDimension)-idxLambda-(2*idxMu) > maxDimension {
					continue
				}

				wg.Add(1)
				go func(lambda, mu, dimension uint32) {
					defer wg.Done()

					lhs := models.NewRandomMatrix(sizeMatrix, dimension)
					rhs := models.NewRandomMatrix(sizeMatrix, dimension)

					start := time.Now()
					result := lhs.Multiplication(lambda, mu, rhs)
					duration := time.Since(start)

					log.Printf("Task completed: Lambda=%d, Mu=%d, Dimension=%d in %v\n", lambda, mu, dimension, duration)
					results = append(results, map[string]interface{}{
						"Lambda":    lambda,
						"Mu":        mu,
						"Dimension": dimension,
						"Duration":  duration.String(),
						"Result":    result.Data,
					})
				}(idxLambda, idxMu, idxDimension)
			}
		}
	}

	wg.Wait()
	log.Println("Matrix multiplication task completed")

	response, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error marshaling response: %v\n", err)
		return
	}
	conn.WriteToUDP(response, addr)
}