# Synchronization

## The basics
Synchronization is essential for concurrent programming. In Go,  the package `sync` does us a big favor when 
scaling goroutines. There are eight important types in the packages [1]:
- `Cond`. It represents a variable that defines a rendezvous point for goroutines waiting for or 
announcing the occurrence of an event. Each Cond has a `Locker L` (often a `*Mutex` or `*RWMutex`), which must be held when changing the 
condition and when calling the `Wait` function. 

- `Locker`. It represents an object that can be locked and unlocked.

- `Map`. It is structurally similar to `map[interface{}]interface{}` but is goroutines-safe without additional locking 
or coordination. *Loads, stores, and deletes run in amortized constant time.* In most cases, `map` is recommended to use 
 with `Mutex` or `RWMutex` thus it is easier to maintain other invariants along with the map content. But we should 
 consider use Map` if the situation follows either of the following:
    - when the entry for a given key is only ever written once but read many times, as in caches that only grows;
    - when multiple goroutines read, write, and overwrite entries for disjoint sets of keys. 
    
- `Once`. It represents an object that will perform exactly once action.

- `Mutex`. It is a mutual exclusion lock (互斥锁). The zero value for a Mutex is an unlocked mutex.

- `RWMutex`. It is a reader/writer mutual exclusion lock. The lock can be held by any number of readers or a single 
writer. If a goroutine holds a RWMutex for reading and another might call Lock, no goroutines should expect to be able 
to get a read lock until the initial read lock is released. This can avoid recursive read locking. 

- `pool`. It is a collection of temporary object that can be individually saved and retrieved. It is used to cache 
allocated but  unused items for later reuse, relieving pressure on GC. It makes easy to build efficient, goroutine-safe 
free lists. A proper use is to manage a group of temporary items silently shared among and potentially reused by 
concurrent independent clients of a package. E.g., package `fmt`.

- `WaitGroup`. It represents an object that waits for a collection of goroutines to finish. The main goroutine calls Add 
to set the number of goroutines to wait for. Then, each of the goroutines runs and calls `Done` when finished. At the same 
time, `Wait` can be used to block until all goroutines have finished.

## Example - Readers-Writers Problem
The main issues regarding Readers-Writers problem (RWP) are:
- Data integrity. How to ensure the data that each reader reads is accurate before and after writer writes.
- Performance. How to avoid dead locks while reader/writer waiting for its turn. 

`sync.RWMutex` provides Multi-Reader-Single-Writer mutual exclusion lock to solve one of the RWPs. Now let's look at 
a problem. The example is written and discussed by Michał Łowicki in his article [2]. The example does two things:
- performing reading and writing tasks, the task execution time is random (`sleep()`). 
- tracking current readers and writers (using channel) 
```
package main
import (
    "fmt"
    "math/rand"
    "strings"
    "sync"
    "time"
)
func init() {
    rand.Seed(time.Now().Unix())
}
func sleep() {
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
func reader(c chan int, m *sync.RWMutex, wg *sync.WaitGroup) {
    sleep()
    m.RLock()
    c <- 1
    sleep()
    c <- -1
    m.RUnlock()
    wg.Done()
}
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
    rsCh := make(chan int)
    wsCh := make(chan int)
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
```

Since we use `RLock()` from `sync.RWMutex`, readers will not block other readers, but writer is always exclusive to all 
the readers and other writers. The results looks like this:
```
R...RR...R...RR...R......W......R...RR...RRR...RRRR...RRRRR...RRRRRR...RRRRRRR...RRRRRR...RRRRR...RRRR...RRR...RR...R......W......W...
```

But if we change `RLock` to `Lock`, then every moment, readers and writers are mutually exclusive: 
```
R......R......R......W......R......R......R......R......W......R......R......R......W...
```

It can be seen that calling `Lock()` in `writer()` blocks all the new readers immediately. Then, writer waits for 
current readers to finish their jobs to start its job. As a result, we can see the number of readers decreases by one 
at a time. It should be noted that, when a writer is done, both readers and writers can be in the waiting list. Then, 
first readers will be unblocked and then writer. In this case, writer is required to wait for the current readers so 
neither readers or writers will get starved (doing nothing). Although, Go allows us to pend over a billion readers, we 
still need design good mechanism (e.g., set timeout for write task) so that the pending readers wont accumulated forever.


`atomic.AddInt32` is used to ensure the operation is performed at 'atomic level' so that other threads will not 
be inferred.

Recursive read locking can happen when a goroutine called `Lock`, other goroutines attempt to get a read lock. 
In this case, Go program throws a fatal error of deadlock (no one should starve at any moment). 

## Discussion
If we have many readers and only a few writers, RWMutex is a good choice, otherwise Mutex is just fine. Some also use channel to implement. There is no standard de facto answer for solving RW problems. We need pread out all the candidates and measure their performance.


Reference:
- [1] [sync package godoc](https://golang.org/pkg/sync)
- [2] [sync.RWMutex, Solving readers-writers problems](https://medium.com/golangspec/sync-rwmutex-ca6c6c3208a0)
