/*
Copyright(C)Michał Łowicki 2018. The original code can be found
https://play.golang.org/p/xoiqW0RQQE9. The discussion of the code
can be found https://medium.com/golangspec/sync-rwmutex-ca6c6c3208a0.
 */

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// init sets seed for random number generation
func init() {
	rand.Seed(time.Now().Unix())
}

// sleep mimics a task execution with random duration ~ (0 - 1000 milliseconds)
func sleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

// reader mimics a reading task
func reader(c chan int, m *sync.RWMutex, wg *sync.WaitGroup) {
	sleep()
	m.RLock()
	c <- 1
	sleep()
	c <- -1
	m.RUnlock()
	wg.Done()
}

// writer mimics a writing task
func writer(c chan int, m *sync.RWMutex, wg *sync.WaitGroup) {
	sleep()
	m.Lock()
	c <- 1
	sleep()
	c <- -1
	m.Unlock()
	wg.Done()
}

func main() {
	var m sync.RWMutex
	var rs, ws int

	// make two channels for concurrent tracking readers and writers
	rsCh := make(chan int)
	wsCh := make(chan int)

	// print current readers and writers
	go func() {
		for {
			select {
			case n := <-rsCh:
				rs += n
			case n := <-wsCh:
				ws += n
			}
			fmt.Printf("%s%s...", strings.Repeat("R", rs), strings.Repeat("W", ws))
		}
	}()

	// perform reading and writing tasks
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go reader(rsCh, &m, &wg)
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go writer(wsCh, &m, &wg)
	}
	wg.Wait()
}
