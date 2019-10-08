# Go routines

## Why goroutine
Goroutines aim to make concurrency easier, which tends to multiplex independently executing 
functions - coroutines - onto a set of threads [1]. The idea behind goroutines versus actual threads is 
about the computational resource concerns. 

Typically, when a coroutine blocks, the run-time automatically moves other coroutines on the same OS thread to a 
different, runnable thread so that they wont be blocked. Such operation is invisible to coders. The result, so called 
goroutines thus can be quite cheap: they have very little overhead beyond the memory for the stack, which is just a 
few KB [1].

To make the stack small, GO uses resizable, bounded stacks. A new goroutine is usually given a few KB, which is always 
enough. When it needs more, the run-time grows (and shrinks) the memory for storing the stack automatically, allowing 
many goroutines to live in a modest amount of memory. The CPU overhead averages about three cheap instructions per func 
call. Thus, it allows to create hundreds of thousands of goroutines in the same address space. If goroutines were 
threads, the system resources can be easily eaten up [1].

We can run goroutines simultaneously when multiple CPUs available. It can be regulated by the **GOMAXPROCS** env 
variable whose default value is the number of CPU cores available. The runtime can allocate more threads than the value 
of **GOMAXPROCS** to service multiple outstanding I/O requests [1].

Goroutines are anonymous, because GO want to enable the full GO language to be available when programming concurrent 
code, not just a library that enables one to do so [1].

## Issue 1 - Too soon to quit
Let's have a look at a classic example: 

```
package main

import "fmt"

func main() {
    for i := 0; i < 10; i++ {
        go fmt.Printf("Job: %v\n", i)
    }
}
```

If you get lucky, you may see some prints, e.g., Job 0, but you are unable to see the most expected ones simply because 
the main function ends too soon before those go routines are executed. In this case, **THERE IS NO WAY** to ensure every 
goroutine you start will be executed. 


### Make-A-Pause Solution
We can make sure every goroutines **WILL BE** executed by forcing a pause after 
every creation of goroutine, but it is not a reasonable practice for the real-world.

```
func main() {
    for i := 0; i < 10; i++ {
        go fmt.Printf("Job: %v\n", i)
        time.Sleep(time.Second * 1)
    }
}
```

### WaitGroup Solution
We can also use sync.WaitGroup. Let's have a look at the following example:
```
func main() {
    var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(jobID int) {
			defer wg.Done()
			fmt.Printf("This is job: %v\n", jobID)
		}(i)
	}
	wg.Wait()
}
```

In the above example, `wg.Add(1)` ensures all the goroutines will be executed before the main function quits.
We add a counter for each goroutine run. `defer wg.Done()` ensures we decrease the counter once the main function 
quits. Calling `wait` to make sure we wait for all the goroutine to finish. Although we can ensure all the goroutines 
are executed before the main quits, we are still `UNABLE` to know which goroutine is done first. As a result, you may 
see different result every run. That's pretty common in the world of concurrency. As we (coders) have no idea which 
goroutine will be 'blocked' and when. If we want to assemble the results of multiple goroutine, a typical approach is 
to use `channel`.

## Error Handling

### Capturing Error Without Stopping Goroutines
We may encounter error when executing goroutines, to find out which goroutine gave away the error, we can use `errgroup` 
from `golang.org/x/sync` (a package that provides Go concurrency primitives in addition to those provided by Golang and 
"sync" and "sync/atomic" packages [4]) to capture any error returned from executing goroutines [3]:

```
func job(ID int) error {
	if rand.Intn(10) == ID {
		return fmt.Errorf("job %v failed", ID)
	}

	fmt.Printf("Job %v completed!\n", ID)
	return nil
}

func main() {
    var eg errgroup.Group
	for i := 0; i < 10; i++ {
		jobID := i
		eg.Go(func() error {
			return job(jobID)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error captured %v\n", err)
	}

	fmt.Println("Mission completed!")
}
```

Calling `eg.Go()` starts a goroutine that has a function and then immediately returns so that the loop can be continued.
It is equivalent to calling `Add` and `Done` in the previous example. `jobID` receives a value copy of `i`. This is an 
essential use as `i` will be reused later. Note that in this case, we only capture the first error returned by 
`eg.Wait()`. If we would like to capture all the errors, we can also dump to, e.g.,  a channel.

### Capturing Error And Stop The Rest
Sometimes, we would like to let the caller or system signals to terminate goroutines. This can be done with `context`.
Let's see the following example mentioned in [3]:
```
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

func main() {
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
```

In the above example, a context created in `main` and is passed to each goroutine. We created a function called 
`jobWithContext` that reads from `Done` channel of the context. The channel is returned by calling `time.After` which 
performs a random sleep that can be interrupted by context cancellation. The use of `errgroup.WithContext` ensures that 
the context will be cancelled when there is an error encountered in the error group. In this case, if one goroutine 
returns an error, the program terminates the rest.


For the full example codes, please see in [go_routine.go](go_routine.go).



## Reference
1. [Go routine in FAQ](https://golang.org/doc/faq#goroutines).
2. [Coroutine](https://en.wikipedia.org/wiki/Coroutine).
3. [Michal Bock, The Startup, Medium, "Managing Groups of Goroutines in Go"](https://medium.com/swlh/managing-groups-of-gorutines-in-go-ee7523e3eaca).
4. [Go Sync package](https://github.com/golang/sync).