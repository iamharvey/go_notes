// +build darwin

package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

)

func main() {
	// In example1, goroutines may not be executed because the main func ends to soon.
	// example1()

	// In example2, we take care goroutine execution by let the main sleeps for a while so that it will only quits after
	// the last goroutine executes
	// example2()

	// In example 3, we use WaitGroup to address the end-to-soon issue in Example 1
	// example3()

	// In example 4, we use errgroup from x/sync so that error that returns from executing goroutines can be captured
	// example4()

	// In example 5, we use errgroup with context, so that when the parent context is cancelled, its child contexts will
	// be cancelled too with proper goroutine termination
	example5()
}

func example1() {
	for i := 0; i < 10; i++ {
		go fmt.Printf("Job: %v\n", i)
	}
}

func example2() {
	for i := 0; i < 10; i++ {
		go fmt.Printf("Job: %v\n", i)
		time.Sleep(time.Second * 1)
	}
}

func example3() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(ID int) {
			defer wg.Done()
			fmt.Printf("This is job: %v\n", ID)
		}(i)
	}
	wg.Wait()
}

func job(ID int) error {
	if rand.Intn(10) == ID {
		fmt.Printf("===> oops! job %v failed\n", ID)
		return fmt.Errorf("job %v failed", ID)
	}

	fmt.Printf("Job %v completed!\n", ID)
	return nil
}

func example4() {
	var eg errgroup.Group
	for i := 0; i < 10; i++ {
		ID := i
		eg.Go(func() error {
			return job(ID)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error captured: %v.\n", err)
	}

	fmt.Println("Mission completed!")
}

// NewContext() creates a new context, if can be terminated when it receives termination signal from OS
func NewContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<- sigChan
		cancel()
	}()

	return ctx
}

func jobWithContext(ctx context.Context, ID int) error {
	select {
		case <-ctx.Done():
			fmt.Printf("Context cancelled! Terminating job %v.\n", ID)
			return nil
		case <-time.After(time.Second * time.Duration(rand.Intn(3))):
	}

	if rand.Intn(10) == ID {
		fmt.Printf("===> oops! job %v failed\n", ID)
		return fmt.Errorf("job %v failed", ID)
	}

	fmt.Printf("Job %v completed!\n", ID)
	return nil
}

func example5() {
	eg, ctx := errgroup.WithContext(NewContext())

	for i := 0; i < 10; i++ {
		ID := i
		eg.Go(func() error {
			return jobWithContext(ctx, ID)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error captured: %v.\n", err)
	}

	fmt.Println("Mission completed!")
}



