// counting prime numbers with vs without concurrency
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const NUM_ITERATIONS int = 10000000

func main() {
	primeCount := 0
	start := time.Now()
	countPrimes(&primeCount, NUM_ITERATIONS)
	fmt.Println("there are ", primeCount, " prime numbers.")
	fmt.Println("it took ", time.Since(start), " without concurrency")

	jobs := make(chan int)
	res := make(chan bool)

	primeCount = 0
	start = time.Now()
	// fill jobs channel
	go func() {
		for i := 0; i < NUM_ITERATIONS; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// create wait group and queue up worker routines
	var wg sync.WaitGroup
	W := runtime.NumCPU()
	wg.Add(W)
	for i := 0; i < W; i++ {
		go worker(jobs, res, &wg)
	}

	// close the res channel down here
	go func() {
		wg.Wait()
		close(res)
	}()

	// tally up the primes based on the bools in the res channel
	for b := range res {
		if b {
			primeCount++
		}
	}

	fmt.Println("there are ", primeCount, " prime numbers.")
	fmt.Println("it took ", time.Since(start), " with concurrency")
}

func worker(jobs <-chan int, res chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for n := range jobs {
		b := isPrime(n)
		res <- b
	}
}

func countPrimes(primeCount *int, n int) {
	for i := 0; i <= n; i++ {
		if isPrime(i) {
			*primeCount++
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
