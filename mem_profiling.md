# Memory Profiling - Understanding The Impact Of `Escape To Heap`

**Escape analysis** is used to find out is a variable escapes to heap. That is part of the compiler optimisation 
procedure. Because GC is involved in heap management. What GC algorithm does is finding out those are no longer referenced 
and wipe them out. Such process can 'stop the world' for a tiny while. If there are too many values in heap, longer latency 
can be introduced. Thus, sometimes, when our program becomes slow-running, it is worth knowing if values will escape to 
heap or not by using `gcflags`. 

```
GC is a complex topic. Rick Hudson mentioned in his 2015's talk about GC's phases:
- **GC off** - pointer writes are just memory writes: *slot=pointer
- **Stack scan** - collect pointers from globals and goroutine stacks. Stacks scanned at preemption points
- **Mark** - Mark objects and follow pointers until pointer queue is empty. Write barrier tracks pointer changes yb mutator
- **Mark termination** - Rescan global/changed stacks, finish marking, shrink stacks etc.
- **Sweep** - Reclaim unmarked object as needed. Adjust GC pacing for next cycle
- **GC off** - Rinse and repeat

There is also so called pacing algorithm that balances the trade-off between the heap growth and CPU utilized. I will not 
go into it in details as this is not my focus for this note. 
```

**NOTE** that small heap size may increase the frequency of GO's stop-the-world behaviour which increases the latency.

So when 'escaping' occurs. There are several common scenarios:

## Memory profiling
Let's see an example that mimics the experiment codes from 
[William Kennedy's talk on memory profile](https://www.ardanlabs.com/blog/2017/06/language-mechanics-on-memory-profiling.html).
The code finds the number of given word in a text file (see below):
```
package mem_profiling

import (
	"bytes"
	"io/ioutil"
)

func findInFile(word string) (count int, err error) {

	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return count, err
	}

	// create an input stream (a reader) for the content
	input := bytes.NewReader(data)

	// we just need the same size of word for reading bytes
	size := len(word)

	// make the buf
	buf := make([]byte, size)
	end := size - 1

	// read in an initial number of bytes we need to get started
	if n, err := input.Read(buf[:end]); err != nil || n < end {
		return count, err
	}

	for {
		// read in one byte from the input stream.
		if _, err := input.Read(buf[end:]); err != nil {
			return count, err
		}

		// if bytes matches, count +1
		if bytes.Compare(buf, []byte(word)) == 0 {
			count ++
		}

		// remove one that has been read
		copy(buf, buf[1:])

	}

	return count, err
}
```

Let's write a benchmark test for the above code:
```
package mem_profiling

import (
	"testing"
)

func BenchmarkFindInFile(b *testing.B) {
	b.ResetTimer()
	count := 0
	for i := 0; i < b.N; i++ {
		count, _ = findInFile("Fermi paradox")
	}
	println("keyword:Fermi, count:", count)
}
```

The collective result can be seen after running the test:
```
$ go test -run none -bench FindInFile -benchtime 3s -benchmem -memprofile mem.out
keyword:Fermi, count: 5
goos: darwin
goarch: amd64
BenchmarkFindInFile-8           keyword:Fermi, count: 5
keyword:Fermi, count: 5
keyword:Fermi, count: 5
   20000            198615 ns/op            9817 B/op          6 allocs/op
PASS
ok      _/Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling   5.998s
```

The above go test command generates two files:
- mem.out - it contains the profile data
- memcpu.test - it contains a test binary we need to have access to symbols when looking at the profile data.

It can be seen that `findInFile` function is allocating 6 values worth a total 9817 bytes per operation.
Now we use 'pprof' to find out what lines of code causes those allocations:

```
$ go tool pprof -alloc_space mem_profiling.test mem.out
File: mem_profiling.test
Type: alloc_space
Time: Jul 29, 2019 at 6:24pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list findInFile
Total: 297.55MB
ROUTINE ======================== _/Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling.findInFile in /Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling/file_utils.go
    1.50MB   297.55MB (flat, cum)   100% of Total
         .          .      5:   "io/ioutil"
         .          .      6:)
         .          .      7:
         .          .      8:func findInFile(word string) (count int, err error) {
         .          .      9:
         .   296.05MB     10:   data, err := ioutil.ReadFile("test.txt")
         .          .     11:   if err != nil {
         .          .     12:           return count, err
         .          .     13:   }
         .          .     14:
         .          .     15:   // create an input stream (a reader) for the content
         .          .     16:   input := bytes.NewReader(data)
         .          .     17:
         .          .     18:   // we just need the same size of word for reading bytes
         .          .     19:   size := len(word)
         .          .     20:
         .          .     21:   // make the buf
    1.50MB     1.50MB     22:   buf := make([]byte, size)
         .          .     23:   end := size - 1
         .          .     24:
         .          .     25:   // read in an initial number of bytes we need to get started
         .          .     26:   if n, err := input.Read(buf[:end]); err != nil || n < end {
         .          .     27:           return count, err
(pprof) 
```

It can be seen that it is the line 22  - constructing a byte slice - causing the allocation on heap. 
Let's also check the Benchmark code:
```
(pprof) list Benchmark
Total: 297.55MB
ROUTINE ======================== _/Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling.BenchmarkFindInFile in /Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling/file_utils_test.go
         0   297.55MB (flat, cum)   100% of Total
         .          .      6:
         .          .      7:func BenchmarkFindInFile(b *testing.B) {
         .          .      8:   b.ResetTimer()
         .          .      9:   count := 0
         .          .     10:   for i := 0; i < b.N; i++ {
         .   297.55MB     11:           count, _ = findInFile("Fermi paradox")
         .          .     12:   }
         .          .     13:   println("keyword:Fermi, count:", count)
         .          .     14:}
(pprof) 
```

It also can be seen that there is no direct allocation in the benchmark code. Thus, all those allocations come from the 
`findInFile` function. To 'dive into the water', let's make use of the GC flag:
```
$ go build -gcflags "-m -m"
# _/Users/xiali/GO_PROJECTS/github.com/iamharvey/go_notes/mem_profiling
./file_utils.go:8:6: cannot inline findInFile: unhandled op FOR
./file_utils.go:16:26: inlining call to bytes.NewReader func([]byte) *bytes.Reader { return &bytes.Reader literal }
./file_utils.go:26:25: inlining call to bytes.(*Reader).Read method(*bytes.Reader) func([]byte) (int, error) { if bytes.r.i >= int64(len(bytes.r.s)) { return int(0), io.EOF }; bytes.r.prevRune = int(-1); bytes.n = copy(bytes.b, bytes.r.s[bytes.r.i:]); bytes.r.i += int64(bytes.n); return  }
./file_utils.go:32:26: inlining call to bytes.(*Reader).Read method(*bytes.Reader) func([]byte) (int, error) { if bytes.r.i >= int64(len(bytes.r.s)) { return int(0), io.EOF }; bytes.r.prevRune = int(-1); bytes.n = copy(bytes.b, bytes.r.s[bytes.r.i:]); bytes.r.i += int64(bytes.n); return  }
./file_utils.go:37:19: inlining call to bytes.Compare func([]byte, []byte) int { return bytealg.Compare(bytes.a, bytes.b) }
./file_utils.go:22:13: make([]byte, size) escapes to heap
./file_utils.go:22:13:  from make([]byte, size) (non-constant size) at ./file_utils.go:22:13
./file_utils.go:8:17: findInFile word does not escape
./file_utils.go:16:26: findInFile &bytes.Reader literal does not escape
./file_utils.go:37:31: findInFile ([]byte)(word) does not escape
```

Examining the result, we can observe that `buf` escapes to heap due to `size`. Because size is determined at run time. 
The compiler has no idea how big this byte slice might be thus it 'worries' that stack might not be big enough for the value, 
thus moves buf to heap. The codes are available [here](mem_profiling)
    
## Summary
- Running benchmark test tells us how many bytes are allocated
- pprof (with `-alloc_space` rather than default `-inuse_space`) helps us to find out the exact location of the 
memory allocations, which are usually caused by values escape to heap
- `run` or `build` with `-gcflags` allows us to find out **WHY** such escapes happen.

Mastering this tool chain is able to help us analyse the memory use of our program, and make you a better coder.
