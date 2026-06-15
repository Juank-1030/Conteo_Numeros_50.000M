package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------
// Read and validate a number entered by the user
// ---------------------------------------------------------------
func pedirNumero(mensaje string, minimo, maximo int64) int64 {
	var valor int64
	fmt.Print(mensaje)
	_, err := fmt.Scan(&valor)

	if err != nil {
		fmt.Println("Error: ingresa un número válido.")
		os.Exit(1)
	}

	if valor < minimo || valor > maximo {
		fmt.Printf("Error: el valor debe estar entre %d y %d.\n", minimo, maximo)
		os.Exit(1)
	}

	return valor
}

// ---------------------------------------------------------------
// Work performed by each goroutine: count its assigned segment
//
// NOTE: A mutex-protected bufio.Writer is used instead of
// fmt.Printf directly to os.Stdout.
//
// Why?
// fmt.Printf in Go has no userspace buffer: each call makes
// an individual write() syscall to the OS. With millions of
// numbers printed, this generates millions of syscalls and is very slow.
//
// bufio.Writer accumulates text in a memory buffer (4 KB by
// default) and only calls the OS when the buffer fills or an
// explicit Flush() is made. This drastically reduces the number of syscalls.
// ---------------------------------------------------------------
func contarSegmento(numeroGoroutine int, inicio, fin int64,
	modoDetallado bool, contadorTotal *int64,
	escritor *bufio.Writer, mutexEscritor *sync.Mutex,
	grupo *sync.WaitGroup) {

	defer grupo.Done()

	cantidadNumeros := fin - inicio + 1

	if modoDetallado {
		// Build all output for this goroutine in memory
		// before writing, to minimize mutex contention.
		buf := make([]byte, 0, cantidadNumeros*30)
		for numero := inicio; numero <= fin; numero++ {
			atomic.AddInt64(contadorTotal, 1)
			buf = fmt.Appendf(buf, "  [Goroutine-%d] --> %d\n", numeroGoroutine, numero)
		}
		// Single write with the mutex held
		mutexEscritor.Lock()
		escritor.Write(buf)
		mutexEscritor.Unlock()
	} else {
		// Add the complete block at once (fast mode)
		atomic.AddInt64(contadorTotal, cantidadNumeros)
		mutexEscritor.Lock()
		fmt.Fprintf(escritor, "  Goroutine-%d: del %d al %d  (%d números)\n",
			numeroGoroutine, inicio, fin, cantidadNumeros)
		mutexEscritor.Unlock()
	}
}

// ---------------------------------------------------------------
// Program entry point
// ---------------------------------------------------------------
func main() {

	// 1. Read user parameters
	numeroFinal := pedirNumero("Número final (1 - 50.000.000.000): ", 1, 50_000_000_000)
	cantidadGoroutines := int(pedirNumero("Cantidad de goroutines: ", 1, 1_000_000))
	modo := pedirNumero("Modo (1 = resumen | 2 = número a número): ", 1, 2)
	modoDetallado := modo == 2

	// 2. Calculate the size of each segment
	tamanoSegmento := numeroFinal / int64(cantidadGoroutines)

	// 3. Shared counter between all goroutines (atomic operation)
	var contadorTotal int64 = 0

	// 4. Buffered writer: accumulates output in memory and reduces OS syscalls
	escritor := bufio.NewWriterSize(os.Stdout, 1024*1024) // 1 MB buffer
	var mutexEscritor sync.Mutex

	fmt.Printf("\nContando del 1 al %d usando %d goroutine(s)...\n\n",
		numeroFinal, cantidadGoroutines)

	// 5. Record the start time
	tiempoInicio := time.Now()

	// 6. Create and start each goroutine
	var grupo sync.WaitGroup

	for i := 0; i < cantidadGoroutines; i++ {
		grupo.Add(1)

		numeroGoroutine := i + 1
		inicio := int64(i)*tamanoSegmento + 1

		var fin int64
		if i == cantidadGoroutines-1 {
			fin = numeroFinal // Last goroutine takes the remainder
		} else {
			fin = int64(i+1) * tamanoSegmento
		}

		go contarSegmento(numeroGoroutine, inicio, fin, modoDetallado,
			&contadorTotal, escritor, &mutexEscritor, &grupo)
	}

	// 7. Wait for all goroutines to finish
	grupo.Wait()

	// 8. Flush the buffer: write all accumulated output to the OS
	escritor.Flush()

	// 9. Calculate total elapsed time
	duracion := time.Since(tiempoInicio)
	milisegundos := duracion.Milliseconds()
	segundos := duracion.Seconds()

	// 10. Show the final result
	fmt.Println()
	fmt.Println("--- Resultado ---")
	fmt.Printf("Total contado  : %d\n", contadorTotal)
	fmt.Printf("Tiempo de mora : %d ms (%.3f s)\n", milisegundos, segundos)
}
