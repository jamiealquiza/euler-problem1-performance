package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

var (
	target     float64
	goroutines int
)

func init() {
	flag.Float64Var(&target, "target", 1000, "Eueler Project problem #1 target limit")
	flag.IntVar(&goroutines, "goroutines", 1, "Number of Goroutines")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func moder(wg *sync.WaitGroup, portion int, numbers chan float64, results chan float64) {
	var nums = make([]float64, 0, portion)

	for n := range numbers {
		if math.Mod(n, 3) == 0 || math.Mod(n, 5) == 0 {
			nums = append(nums, n)
		}
	}

	var sum float64
	for _, val := range nums {
		sum += val
	}

	results <- sum
	defer wg.Done()
}

func main() {
	processStart := time.Now()

	numbers := make(chan float64, int(target)+1)
	for i := float64(0); i < target+float64(1); i++ {
		numbers <- i
	}
	fmt.Println(len(numbers))
	close(numbers)
	fmt.Printf("Populated %d numbers in %s\n",
		int(target), time.Since(processStart))

	fmt.Printf("Starting %d Goroutines\n", goroutines)

	results := make(chan float64, int(target))
	wg := &sync.WaitGroup{}
	wg.Add(goroutines)

	findStart := time.Now()
	for i := 0; i < goroutines; i++ {
		go moder(wg, int(target)/goroutines, numbers, results)
	}

	wg.Wait()
	fmt.Printf("Found values in %s\n", time.Since(findStart))
	close(results)

	sumStart := time.Now()
	var sum float64
	for n := range results {
		sum += n
	}

	fmt.Printf("Calculated sum %d in %s\n", int(sum), time.Since(sumStart))
	fmt.Printf("Total run time %s\n", time.Since(processStart))
}