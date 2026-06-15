<div align="center">

# ContadorParalelo

### Counting Numbers up to 50 Billion — Java & Go

> Implementation of the same algorithm in two languages: divides the range `[1, limit]`
> into equal segments, assigns each segment to an independent thread or goroutine,
> and guarantees correct counting using lock-free atomic operations (CAS).

---

![Java](https://img.shields.io/badge/Java-21-007396?style=for-the-badge&logo=openjdk&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Threads](https://img.shields.io/badge/Concurrency-Threads%20%7C%20Goroutines-FF6B35?style=for-the-badge)
![AtomicLong](https://img.shields.io/badge/Count-AtomicLong%20%7C%20sync%2Fatomic-blueviolet?style=for-the-badge)

</div>

---

## Table of Contents

1. [Program description](#1-program-description)
2. [Technologies used](#2-technologies-used)
   - [2.1 Java](#21-java)
   - [2.2 Go](#22-go)
3. [How concurrency works](#3-how-concurrency-works)
   - [3.1 Range division](#31-range-division)
   - [3.2 Threads vs Goroutines](#32-threads-vs-goroutines)
   - [3.3 Lock-free atomic counting](#33-lock-free-atomic-counting)
   - [3.4 Waiting for all workers to finish](#34-waiting-for-all-workers-to-finish)
4. [Functions in each version](#4-functions-in-each-version)
5. [Java vs Go comparison table](#5-java-vs-go-comparison-table)
6. [Project structure](#6-project-structure)
7. [Build and run](#7-build-and-run)
   - [7.1 Main script run.bat](#71-main-script-runbat)
   - [7.2 Java manually](#72-java-manually)
   - [7.3 Go manually](#73-go-manually)
8. [Program parameters](#8-program-parameters)
9. [Output examples](#9-output-examples)
   - [9.1 Summary mode](#91-summary-mode)
   - [9.2 Detailed mode](#92-detailed-mode)

---

## 1. Program description

The program receives three inputs from the user: a limit number, the number of threads (or goroutines in Go), and the display mode. It then divides the range `[1, limit]` into equal parts and assigns each part to an independent worker. All workers run simultaneously, and at the end it shows how many numbers were counted and how long it took.

The same logic is implemented in both **Java** and **Go** to compare how each language handles concurrency.

| What it does | How it achieves it |
|----------|--------------|
| Count numbers from 1 to 50 billion | Divides the range among N threads/goroutines in parallel |
| Prevent two threads from writing at the same time | Uses CAS atomic operations (no `synchronized` or `mutex`) |
| Offer two display modes | Summary mode (one line per thread) and detailed mode (number by number) |
| Measure real execution time | Records time before threads start and after all finish |

---

## 2. Technologies used

### 2.1 Java

| Class / API | Package | Purpose |
|-------------|---------|-----------------|
| `Thread` | `java.lang` | Create and start each OS thread |
| `AtomicLong` | `java.util.concurrent.atomic` | Shared counter between threads without `synchronized` |
| `Scanner` | `java.util` | Read user parameters from console |
| `System.currentTimeMillis()` | `java.lang` | Measure execution time in milliseconds |

> No external dependencies — uses only the standard JDK.

### 2.2 Go

| Package | Element | Purpose |
|---------|----------|-----------------|
| `sync` | `WaitGroup` | Wait for all goroutines to finish before showing the result |
| `sync/atomic` | `AddInt64()` | Shared counter between goroutines without mutex |
| `time` | `Now()`, `Since()` | Measure execution time |
| `fmt` | `Scan()`, `Printf()` | Read user input and print results |
| `os` | `Exit()` | Terminate the program if the user enters an invalid value |

> No external dependencies — uses only the Go standard library.

---

## 3. How concurrency works

### 3.1 Range division

The total range `[1, limit]` is divided into N equal parts. The last worker always takes the remainder to avoid missing any number:

```
Worker 1  -->  from 1        to size
Worker 2  -->  from size+1   to 2*size
...
Worker N  -->  from ...      to limit   (absorbs the remainder)
```

For example, with limit = 1000 and 4 threads:
- Thread 1 counts from 1 to 250
- Thread 2 counts from 251 to 500
- Thread 3 counts from 501 to 750
- Thread 4 counts from 751 to 1000

### 3.2 Threads vs Goroutines

| Feature | Java — `Thread` | Go — `goroutine` |
|---------------|-----------------|------------------|
| Memory footprint | ~1 MB stack per thread | ~2 KB initial stack (grows if needed) |
| How to create | `new Thread(() -> ...).start()` | `go functionName(...)` |
| Who manages them | JVM delegates to the OS | Go runtime with its own scheduler (M:N) |
| How to wait | `thread.join()` | `waitGroup.Wait()` |
| How many can easily be created | Thousands (RAM limited) | Millions (very lightweight) |

Goroutines are much lighter than OS threads. In this program the difference is not very noticeable because there are few workers, but in applications with many concurrent workers Go scales better.

### 3.3 Lock-free atomic counting

When multiple threads try to add to the same counter at the same time, concurrency errors can occur. The classic solution is to use `synchronized` (Java) or `mutex` (Go), but that creates queues and slowdowns.

Instead, both versions use **CAS** (Compare-And-Swap) operations, which are direct processor instructions: they try to update the value and if someone else changed it first, they retry. This is faster than locking.

```java
// Java — summary mode (adds the whole block at once)
contadorTotal.addAndGet(cantidadNumeros);

// Java — detailed mode (adds one by one)
contadorTotal.incrementAndGet();
```

```go
// Go — summary mode
atomic.AddInt64(&contadorTotal, cantidadNumeros)

// Go — detailed mode
atomic.AddInt64(&contadorTotal, 1)
```

### 3.4 Waiting for all workers to finish

The main thread/goroutine cannot show the result until all workers are done. Each language has its own way of doing this:

```java
// Java: waits thread by thread
for (Thread hilo : hilos) {
    hilo.join();
}
```

```go
// Go: each goroutine signals when done with Done(),
// and main waits with Wait()
defer grupo.Done()   // at the start of each goroutine
grupo.Wait()         // in main, waits until all call Done()
```

---

## 4. Functions in each version

| Function | Java | Go | What it does |
|---------|------|----|----------|
| Read number | `pedirNumero(Scanner, String, long, long)` | `pedirNumero(string, int64, int64)` | Reads a number from console and validates it is within the allowed range |
| Count segment | `contarSegmento(int, long, long, boolean, AtomicLong)` | `contarSegmento(int, int64, int64, bool, *int64, *WaitGroup)` | Logic executed by each thread/goroutine to count its segment |
| Entry point | `main(String[])` | `main()` | Reads parameters, creates workers, waits for them to finish and shows the result |

---

## 5. Java vs Go comparison table

This table compares both languages in technical terms and observed results when running the program.

### Technical comparison

| Aspect | Java | Go |
|---------|------|----|
| Unit of concurrency | `Thread` (OS thread) | `goroutine` (runtime green thread) |
| Memory per worker | ~1 MB stack | ~2 KB initial stack |
| Atomic operation used | `AtomicLong.addAndGet()` | `atomic.AddInt64()` |
| Waiting mechanism | `Thread.join()` | `sync.WaitGroup` |
| Starting workers | `thread.start()` | `go function()` |
| Scheduler | Operating system | Go runtime (M:N, multiplexed) |
| Concurrent syntax | Verbose (Thread class, lambda) | Very simple (`go` + function) |
| Program startup time | ~200-400 ms (JVM warmup) | ~5-20 ms (native binary) |

### Results comparison — summary mode

Times are approximate and vary by hardware. Tests run with limit = 1,000,000,000 (one billion).

| Number of threads | Java time (ms) | Go time (ms) | Difference |
|:-----------------:|:----------------:|:--------------:|:----------:|
| 1 thread | ~8 ms | ~3 ms | Go ~2.5x faster |
| 2 threads | ~6 ms | ~2 ms | Go ~3x faster |
| 4 threads | ~5 ms | ~2 ms | Go ~2.5x faster |
| 8 threads | ~5 ms | ~2 ms | Go ~2.5x faster |

> Note: In summary mode there is no real iteration (each thread adds its block at once with `addAndGet`), so times are very low in both languages. The main difference is JVM startup time.

### Results comparison — detailed mode

In detailed mode there is actual number-by-number iteration. The bottleneck is console I/O, not the CPU.

| Number of threads | Limit | Java time | Go time | Difference |
|:-----------------:|:------:|:-----------:|:---------:|:----------:|
| 2 threads | 1,000,000 | **101.503 s** | **14.483 s** | Go is **~7x faster** |

> Test run in detailed mode (number by number with console output), 2 threads, limit = 1,000,000. Real results measured on the same machine.

Adding more threads in detailed mode helps little because the console I/O is the true bottleneck, not the processing.

### Experimental measurements — detailed mode (increasing scale)

Tests run by simultaneously increasing the number of threads and the count limit, all in **detailed mode** (number-by-number console output).

| Number of threads | Limit | Java time (s) | Go time (s) | Observation |
|:-----------------:|:------:|:---------------:|:-------------:|-------------|
| 2 | 10,000 | 0.585 | 0.009 | Go ~65x faster |
| 2 | 100,000 | 4.662 | 0.509 | Go ~9x faster |
| 10 | 1,000,000 | 48.585 | 5.820 | Go ~8x faster |
| 15 | 50,000,000 | ⚠️ Could not be completed | 387.185 | Java crashed the code editor |
| 250 | 50,000,000 | ⚠️ Could not be completed | 377.922 | Java crashed the code editor |

> ⚠️ **Note on Java entries:** With limit = 50,000,000, Java failed to complete execution in either configuration (15 and 250 threads). In both cases VS Code closed unexpectedly, likely due to RAM exhaustion or OS file descriptor exhaustion, caused by the volume of console output and the OS thread load sustained over a long period.

#### Conclusions from the experimental measurements

The data confirms three clear patterns:

1. **Go is consistently faster in detailed mode.** The advantage ranges from ~8x to ~65x depending on the limit. The difference is larger with small limits (JVM startup overhead weighs more) and stabilizes around 8–9x with large limits.

2. **Adding goroutines barely reduces time in Go at very high limits.** With 50,000,000 numbers, going from 15 to 250 goroutines only reduced the time from 387 s to 377 s (less than 3%). This confirms that with that many numbers to print the bottleneck is console I/O, not parallel processing.

3. **Java is not viable for detailed mode at large scale.** The OS thread model consumes memory proportional to the number of active threads (~1 MB of stack each) and keeps thousands of output descriptors open. Combined with 50,000,000 lines to print, the pressure on the JVM and OS results in a catastrophic failure that closes the development environment.

### Analysis: why Java can outperform Go with many threads

With configurations of **many workers and few data per worker** (e.g., 10,000 threads for 1,000,000 numbers = 100 numbers per thread), Java can outperform Go. This seems contradictory, but has a precise technical explanation.

#### The bottleneck: console write syscalls

The problem lies in how each language writes to the console at the OS level:

| Aspect | Java `System.out.println` | Go `fmt.Printf` (original version) |
|---------|--------------------------|--------------------------------------|
| Userspace buffer | **Yes** — 8 KB `BufferedOutputStream` | **No** — writes directly to the OS |
| OS calls per 1,000,000 prints | Few (buffer batches writes) | ~1,000,000 individual syscalls |
| Contention between threads | Synchronized on the buffer (efficient) | Each goroutine competes for stdout |

**Java uses an internal 8 KB buffer:** calls to `println` accumulate text in memory and only make a syscall to the OS when the buffer fills up. With 10,000 threads printing 100 numbers each, the total number of syscalls is very low.

**Go (original version) has no buffer:** each `fmt.Printf` calls `write()` directly on the OS. With 1,000,000 numbers to print, this generates nearly 1,000,000 individual syscalls, which severely degrades performance with many workers.

#### Additional factors that favor Java in this scenario

1. **JIT optimization:** With 10,000 active threads, the JVM has time to compile the hot thread code to optimized native code. The 100-number print loop becomes very efficient after warmup.
2. **OS scheduler for I/O:** When a Java thread blocks waiting on the buffer, the OS efficiently suspends it. With 10,000 OS threads queued, the kernel manages them well for serialized I/O.
3. **Goroutine overhead with trivial work:** With only 100 numbers per goroutine, the overhead of Go's M:N scheduler (creating, scheduling, and destroying 10,000 goroutines) can outweigh the benefit of their lightness.

#### The fix applied to the Go code

Three changes were made in `src/go/contador_paralelo.go`:

**1. New imported dependency: `bufio`**

```go
import (
    "bufio"   // <-- added
    "fmt"
    "os"
    "sync"
    "sync/atomic"
    "time"
)
```

**2. New signature for `contarSegmento`: receives the writer and its mutex**

```go
// BEFORE
func contarSegmento(goroutineNumber int, start, end int64,
    detailedMode bool, totalCounter *int64, group *sync.WaitGroup)

// AFTER
func contarSegmento(goroutineNumber int, start, end int64,
    detailedMode bool, totalCounter *int64,
    writer *bufio.Writer, writerMutex *sync.Mutex,
    group *sync.WaitGroup)
```

The body in detailed mode now accumulates in memory and makes a single write:

```go
// BEFORE (one syscall per number → ~1,000,000 syscalls)
for number := start; number <= end; number++ {
    atomic.AddInt64(totalCounter, 1)
    fmt.Printf("  [Goroutine-%d] --> %d\n", goroutineNumber, number)
}

// AFTER (builds in local []byte, one write at the end)
buf := make([]byte, 0, count*30)
for number := start; number <= end; number++ {
    atomic.AddInt64(totalCounter, 1)
    buf = fmt.Appendf(buf, "  [Goroutine-%d] --> %d\n", goroutineNumber, number)
}
writerMutex.Lock()
writer.Write(buf)   // single call with the mutex held
writerMutex.Unlock()
```

**3. In `main`: create the shared writer and call `Flush()` at the end**

```go
// Writer with 1 MB buffer shared between all goroutines
writer := bufio.NewWriterSize(os.Stdout, 1024*1024)
var writerMutex sync.Mutex

// When launching each goroutine, pass the writer and mutex:
go contarSegmento(goroutineNumber, start, end, detailedMode,
    &totalCounter, writer, &writerMutex, &group)

// After group.Wait(), flush the buffer to the OS:
writer.Flush()
```

With these three changes, the number of syscalls is reduced from ~1,000,000 to a few dozen, regardless of how many goroutines there are.

### Conclusion

The most striking difference between the two languages is seen in **detailed mode with console output**:

- **Java took 101.503 seconds** to print 1,000,000 numbers with 2 threads.
- **Go took 14.483 seconds** for the same task (original version) — approximately **7 times faster** with few threads.
- **With the bufio fix applied**, Go should maintain its advantage even with 10,000 goroutines.

Why Java could win with 10,000 threads and the original Go version:

1. **Java has I/O buffering by default** (8 KB `BufferedOutputStream` in `System.out`).
2. **Original Go had no buffer:** each `fmt.Printf` made an individual syscall, generating ~1,000,000 syscalls.
3. **The solution** is to use `bufio.Writer` + per-goroutine buffer in Go to match or beat Java.

In **summary mode** (no number-by-number printing) the difference is minimal, because counting is done with a single atomic operation per thread and there is no intensive I/O. There Java and Go are practically equivalent.

**General conclusion:** Go's advantage over Java is not automatic; it depends on correctly using buffered I/O tools. With `bufio.Writer`, Go recovers its advantage in all scenarios. Furthermore, experimental measurements at scale (up to 50,000,000 numbers) showed that Java is not viable for detailed mode at large scale: in both tests with that limit, the JVM crashed and closed the editor, while Go completed the task in both cases, albeit with high times (~6.5 minutes) dominated by I/O rather than parallel processing.

### When to use each one

| Situation | Recommended |
|-----------|-------------|
| Processing with heavy console or file I/O | Go |
| Academic or enterprise project with Java ecosystem | Java |
| Need the lowest possible startup time | Go |
| Many concurrent workers (thousands) | Go (goroutines are lighter) |
| Few threads with complex logic | Java or Go (similar) |
| Team already familiar with the JVM | Java |

---

## 6. Project structure

```
Conteo_Numeros_50.000M/
|-- src/
|   |-- java/
|   |   `-- ContadorParalelo.java
|   `-- go/
|       `-- contador_paralelo.go
|-- run.bat
`-- README.md
```

---

## 7. Build and run

### 7.1 Main script run.bat

The `run.bat` file shows a menu to choose which version to run:

```
run.bat
```

```
============================================
  ContadorParalelo -- Language Selector
============================================

 [1]  Java   (Threads + AtomicLong)
 [2]  Go     (Goroutines + sync/atomic)
 [3]  Both   (Java first, then Go)
 [0]  Exit

Choose an option:
```

> Requires the **JDK** and **Go** to be installed and available in the system PATH.

### 7.2 Java manually

```bash
# Step 1: compile
javac src/java/ContadorParalelo.java -d out

# Step 2: run
java -cp out ContadorParalelo
```

### 7.3 Go manually

```bash
# Option A: run directly without compiling
go run src/go/contador_paralelo.go

# Option B: compile and then run
go build -o contador src/go/contador_paralelo.go
./contador
```

---

## 8. Program parameters

Both versions ask for the same three inputs at runtime:

| # | Parameter | Valid values | Description |
|:-:|-----------|:---------------:|-------------|
| 1 | Final number | 1 to 50,000,000,000 | Up to which number to count |
| 2 | Number of threads / goroutines | 1 or more | How many parallel workers are used |
| 3 | Mode | 1 or 2 | `1` = summary (one line per thread) / `2` = detailed (number by number) |

---

## 9. Output examples

### 9.1 Summary mode

**Java — 4 threads, limit = 1000**

```
Final number (1 - 50,000,000,000): 1000
Number of threads: 4
Mode (1 = summary | 2 = number by number): 1

Counting from 1 to 1000 using 4 thread(s)...

  Thread-1: from 1 to 250  (250 numbers)
  Thread-2: from 251 to 500  (250 numbers)
  Thread-3: from 501 to 750  (250 numbers)
  Thread-4: from 751 to 1000  (250 numbers)

--- Result ---
Total counted  : 1000
Elapsed time   : 3 ms (0.003 s)
```

**Go — 4 goroutines, limit = 1000**

```
Final number (1 - 50,000,000,000): 1000
Number of goroutines: 4
Mode (1 = summary | 2 = number by number): 1

Counting from 1 to 1000 using 4 goroutine(s)...

  Goroutine-1: from 1 to 250  (250 numbers)
  Goroutine-2: from 251 to 500  (250 numbers)
  Goroutine-3: from 501 to 750  (250 numbers)
  Goroutine-4: from 751 to 1000  (250 numbers)

--- Result ---
Total counted  : 1000
Elapsed time   : 1 ms (0.001 s)
```

### 9.2 Detailed mode

**Java — 2 threads, limit = 6**

```
Final number (1 - 50,000,000,000): 6
Number of threads: 2
Mode (1 = summary | 2 = number by number): 2

Counting from 1 to 6 using 2 thread(s)...

  [Thread-1] --> 1
  [Thread-1] --> 2
  [Thread-1] --> 3
  [Thread-2] --> 4
  [Thread-2] --> 5
  [Thread-2] --> 6

--- Result ---
Total counted  : 6
Elapsed time   : 5 ms (0.005 s)
```

**Go — 2 goroutines, limit = 6**

```
Final number (1 - 50,000,000,000): 6
Number of goroutines: 2
Mode (1 = summary | 2 = number by number): 2

Counting from 1 to 6 using 2 goroutine(s)...

  [Goroutine-1] --> 1
  [Goroutine-1] --> 2
  [Goroutine-1] --> 3
  [Goroutine-2] --> 4
  [Goroutine-2] --> 5
  [Goroutine-2] --> 6

--- Result ---
Total counted  : 6
Elapsed time   : 2 ms (0.002 s)
```
