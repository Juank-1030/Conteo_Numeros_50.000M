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
// Lectura y validación de un número ingresado por el usuario
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
// Trabajo que realiza cada goroutine: contar su segmento asignado
//
// NOTA: Se usa un bufio.Writer protegido por mutex en lugar de
// fmt.Printf directo a os.Stdout.
//
// ¿Por qué?
// fmt.Printf en Go no tiene buffer en userspace: cada llamada hace
// una syscall write() individual al sistema operativo. Con millones
// de números impresos, eso genera millones de syscalls y es muy lento.
//
// bufio.Writer acumula los textos en un buffer de memoria (4 KB por
// defecto) y solo llama al SO cuando el buffer se llena o se hace
// Flush() explícito. Esto reduce drásticamente el número de syscalls.
// ---------------------------------------------------------------
func contarSegmento(numeroGoroutine int, inicio, fin int64,
	modoDetallado bool, contadorTotal *int64,
	escritor *bufio.Writer, mutexEscritor *sync.Mutex,
	grupo *sync.WaitGroup) {

	defer grupo.Done()

	cantidadNumeros := fin - inicio + 1

	if modoDetallado {
		// Construir toda la salida de esta goroutine en memoria
		// antes de escribir, para minimizar la contención del mutex.
		buf := make([]byte, 0, cantidadNumeros*30)
		for numero := inicio; numero <= fin; numero++ {
			atomic.AddInt64(contadorTotal, 1)
			buf = fmt.Appendf(buf, "  [Goroutine-%d] --> %d\n", numeroGoroutine, numero)
		}
		// Una sola escritura con el mutex tomado
		mutexEscritor.Lock()
		escritor.Write(buf)
		mutexEscritor.Unlock()
	} else {
		// Suma el bloque completo de una sola vez (modo rápido)
		atomic.AddInt64(contadorTotal, cantidadNumeros)
		mutexEscritor.Lock()
		fmt.Fprintf(escritor, "  Goroutine-%d: del %d al %d  (%d números)\n",
			numeroGoroutine, inicio, fin, cantidadNumeros)
		mutexEscritor.Unlock()
	}
}

// ---------------------------------------------------------------
// Punto de entrada del programa
// ---------------------------------------------------------------
func main() {

	// 1. Leer parámetros del usuario
	numeroFinal := pedirNumero("Número final (1 - 50.000.000.000): ", 1, 50_000_000_000)
	cantidadGoroutines := int(pedirNumero("Cantidad de goroutines: ", 1, 1_000_000))
	modo := pedirNumero("Modo (1 = resumen | 2 = número a número): ", 1, 2)
	modoDetallado := modo == 2

	// 2. Calcular el tamaño de cada segmento
	tamanoSegmento := numeroFinal / int64(cantidadGoroutines)

	// 3. Contador compartido entre todas las goroutines (operación atómica)
	var contadorTotal int64 = 0

	// 4. Escritor con buffer: acumula salida en memoria y reduce syscalls al SO
	escritor := bufio.NewWriterSize(os.Stdout, 1024*1024) // buffer de 1 MB
	var mutexEscritor sync.Mutex

	fmt.Printf("\nContando del 1 al %d usando %d goroutine(s)...\n\n",
		numeroFinal, cantidadGoroutines)

	// 5. Registrar el tiempo de inicio
	tiempoInicio := time.Now()

	// 6. Crear y arrancar cada goroutine
	var grupo sync.WaitGroup

	for i := 0; i < cantidadGoroutines; i++ {
		grupo.Add(1)

		numeroGoroutine := i + 1
		inicio := int64(i)*tamanoSegmento + 1

		var fin int64
		if i == cantidadGoroutines-1 {
			fin = numeroFinal // El último toma el remanente
		} else {
			fin = int64(i+1) * tamanoSegmento
		}

		go contarSegmento(numeroGoroutine, inicio, fin, modoDetallado,
			&contadorTotal, escritor, &mutexEscritor, &grupo)
	}

	// 7. Esperar a que todas las goroutines terminen
	grupo.Wait()

	// 8. Vaciar el buffer: escribir al SO todo lo que quedó acumulado
	escritor.Flush()

	// 9. Calcular el tiempo total transcurrido
	duracion := time.Since(tiempoInicio)
	milisegundos := duracion.Milliseconds()
	segundos := duracion.Seconds()

	// 10. Mostrar el resultado final
	fmt.Println()
	fmt.Println("--- Resultado ---")
	fmt.Printf("Total contado  : %d\n", contadorTotal)
	fmt.Printf("Tiempo de mora : %d ms (%.3f s)\n", milisegundos, segundos)
}
