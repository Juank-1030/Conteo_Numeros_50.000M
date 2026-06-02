<div align="center">

# ContadorParalelo

### Conteo de Números hasta 50.000 Millones — Java & Go

> Implementación del mismo algoritmo en dos lenguajes: divide el rango `[1, límite]`
> en segmentos iguales, asigna cada segmento a un hilo o goroutine independiente,
> y garantiza un conteo correcto usando operaciones atómicas sin bloqueos (CAS).

---

![Java](https://img.shields.io/badge/Java-21-007396?style=for-the-badge&logo=openjdk&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Threads](https://img.shields.io/badge/Concurrencia-Threads%20%7C%20Goroutines-FF6B35?style=for-the-badge)
![AtomicLong](https://img.shields.io/badge/Conteo-AtomicLong%20%7C%20sync%2Fatomic-blueviolet?style=for-the-badge)

</div>

---

## Tabla de Contenidos

1. [Descripcion del programa](#1-descripcion-del-programa)
2. [Tecnologias utilizadas](#2-tecnologias-utilizadas)
   - [2.1 Java](#21-java)
   - [2.2 Go](#22-go)
3. [Como funciona la concurrencia](#3-como-funciona-la-concurrencia)
   - [3.1 Division del rango](#31-division-del-rango)
   - [3.2 Hilos vs Goroutines](#32-hilos-vs-goroutines)
   - [3.3 Conteo atomico sin bloqueos](#33-conteo-atomico-sin-bloqueos)
   - [3.4 Esperar a que todos terminen](#34-esperar-a-que-todos-terminen)
4. [Funciones de cada version](#4-funciones-de-cada-version)
5. [Tabla comparativa Java vs Go](#5-tabla-comparativa-java-vs-go)
6. [Estructura del proyecto](#6-estructura-del-proyecto)
7. [Compilar y ejecutar](#7-compilar-y-ejecutar)
   - [7.1 Script principal run.bat](#71-script-principal-runbat)
   - [7.2 Java manualmente](#72-java-manualmente)
   - [7.3 Go manualmente](#73-go-manualmente)
8. [Parametros del programa](#8-parametros-del-programa)
9. [Ejemplos de salida](#9-ejemplos-de-salida)
   - [9.1 Modo resumen](#91-modo-resumen)
   - [9.2 Modo detallado](#92-modo-detallado)

---

## 1. Descripcion del programa

El programa recibe tres datos del usuario: un número límite, la cantidad de hilos (o goroutines en Go) y el modo de visualización. Con eso, divide el rango `[1, límite]` en partes iguales y le asigna cada parte a un trabajador independiente. Todos corren al mismo tiempo y al final se muestra cuántos números se contaron y cuánto tardó.

La misma lógica está implementada en **Java** y en **Go** para poder comparar cómo cada lenguaje maneja la concurrencia.

| Qué hace | Cómo lo logra |
|----------|--------------|
| Contar números de 1 hasta 50.000 millones | Divide el rango entre N hilos/goroutines en paralelo |
| Evitar que dos hilos escriban al mismo tiempo | Usa operaciones atómicas CAS (sin `synchronized` ni `mutex`) |
| Ofrecer dos modos de visualización | Modo resumen (una línea por hilo) y modo detallado (número a número) |
| Medir el tiempo real de ejecución | Registra el tiempo antes de arrancar los hilos y al terminar todos |

---

## 2. Tecnologias utilizadas

### 2.1 Java

| Clase / API | Paquete | Para qué se usa |
|-------------|---------|-----------------|
| `Thread` | `java.lang` | Crear y arrancar cada hilo del sistema operativo |
| `AtomicLong` | `java.util.concurrent.atomic` | Contador compartido entre hilos sin necesidad de `synchronized` |
| `Scanner` | `java.util` | Leer los parámetros que ingresa el usuario por consola |
| `System.currentTimeMillis()` | `java.lang` | Medir cuánto tiempo tarda la ejecución en milisegundos |

> Sin dependencias externas — usa únicamente el JDK estándar de Java.

### 2.2 Go

| Paquete | Elemento | Para qué se usa |
|---------|----------|-----------------|
| `sync` | `WaitGroup` | Esperar a que todas las goroutines terminen antes de mostrar el resultado |
| `sync/atomic` | `AddInt64()` | Contador compartido entre goroutines sin necesidad de mutex |
| `time` | `Now()`, `Since()` | Medir cuánto tiempo tarda la ejecución |
| `fmt` | `Scan()`, `Printf()` | Leer entradas del usuario e imprimir resultados |
| `os` | `Exit()` | Terminar el programa si el usuario ingresa un valor inválido |

> Sin dependencias externas — usa únicamente la librería estándar de Go.

---

## 3. Como funciona la concurrencia

### 3.1 Division del rango

El rango total `[1, límite]` se divide en N partes del mismo tamaño. El último trabajador siempre se queda con el sobrante para no perder ningún número:

```
Trabajador 1  -->  del 1        al tamano
Trabajador 2  -->  del tamano+1 al 2*tamano
...
Trabajador N  -->  del ...      al limite   (absorbe el sobrante)
```

Por ejemplo, con límite = 1000 y 4 hilos:
- Hilo 1 cuenta del 1 al 250
- Hilo 2 cuenta del 251 al 500
- Hilo 3 cuenta del 501 al 750
- Hilo 4 cuenta del 751 al 1000

### 3.2 Hilos vs Goroutines

| Característica | Java — `Thread` | Go — `goroutine` |
|---------------|-----------------|------------------|
| Peso en memoria | ~1 MB de stack por hilo | ~2 KB de stack inicial (crece si necesita) |
| Cómo se crea | `new Thread(() -> ...).start()` | `go nombreFuncion(...)` |
| Quién los gestiona | La JVM delega al sistema operativo | El runtime de Go con scheduler propio (M:N) |
| Cómo se espera | `hilo.join()` | `waitGroup.Wait()` |
| Cuántos se pueden crear fácilmente | Miles (limitado por RAM) | Millones (muy ligeras) |

Las goroutines son mucho más ligeras que los hilos del SO. En este programa la diferencia no es tan notoria porque son pocos trabajadores, pero en aplicaciones con muchos trabajadores concurrentes Go escala mejor.

### 3.3 Conteo atomico sin bloqueos

Cuando varios hilos intentan sumar al mismo contador al mismo tiempo, puede haber errores de concurrencia. La solución clásica es usar `synchronized` (Java) o `mutex` (Go), pero eso genera colas y lentitud.

En cambio, ambas versiones usan operaciones **CAS** (Compare-And-Swap), que son instrucciones directas del procesador: intentan actualizar el valor y si alguien más lo cambió antes, reintenta. Esto es más rápido que bloquear.

```java
// Java — modo resumen (suma todo el bloque de una vez)
contadorTotal.addAndGet(cantidadNumeros);

// Java — modo detallado (suma de a uno)
contadorTotal.incrementAndGet();
```

```go
// Go — modo resumen
atomic.AddInt64(&contadorTotal, cantidadNumeros)

// Go — modo detallado
atomic.AddInt64(&contadorTotal, 1)
```

### 3.4 Esperar a que todos terminen

El hilo/goroutine principal no puede mostrar el resultado hasta que todos los trabajadores terminen. Cada lenguaje tiene su forma de hacer esto:

```java
// Java: espera hilo por hilo
for (Thread hilo : hilos) {
    hilo.join();
}
```

```go
// Go: cada goroutine avisa cuando termina con Done(),
// y el main espera con Wait()
defer grupo.Done()   // al inicio de cada goroutine
grupo.Wait()         // en el main, espera hasta que todas llamen Done()
```

---

## 4. Funciones de cada version

| Función | Java | Go | Qué hace |
|---------|------|----|----------|
| Pedir número | `pedirNumero(Scanner, String, long, long)` | `pedirNumero(string, int64, int64)` | Lee un número por consola y valida que esté en el rango permitido |
| Contar segmento | `contarSegmento(int, long, long, boolean, AtomicLong)` | `contarSegmento(int, int64, int64, bool, *int64, *WaitGroup)` | Lógica que ejecuta cada hilo/goroutine para contar su segmento |
| Principal | `main(String[])` | `main()` | Lee parámetros, crea los trabajadores, espera que terminen y muestra el resultado |

---

## 5. Tabla comparativa Java vs Go

Esta tabla compara los dos lenguajes en términos técnicos y de resultados observados al ejecutar el programa.

### Comparacion tecnica

| Aspecto | Java | Go |
|---------|------|----|
| Unidad de concurrencia | `Thread` (hilo del SO) | `goroutine` (hilo verde del runtime) |
| Memoria por trabajador | ~1 MB de stack | ~2 KB de stack inicial |
| Operación atómica usada | `AtomicLong.addAndGet()` | `atomic.AddInt64()` |
| Mecanismo de espera | `Thread.join()` | `sync.WaitGroup` |
| Arranque de trabajadores | `hilo.start()` | `go funcion()` |
| Scheduler | Sistema operativo | Runtime de Go (M:N, multiplexado) |
| Sintaxis concurrente | Verbosa (clase Thread, lambda) | Muy simple (`go` + función) |
| Tiempo de inicio del programa | ~200-400 ms (JVM warmup) | ~5-20 ms (binario nativo) |

### Comparacion de resultados — modo resumen

Los tiempos son aproximados y varían según el hardware. Pruebas realizadas con límite = 1.000.000.000 (mil millones).

| Cantidad de hilos | Tiempo Java (ms) | Tiempo Go (ms) | Diferencia |
|:-----------------:|:----------------:|:--------------:|:----------:|
| 1 hilo | ~8 ms | ~3 ms | Go ~2.5x más rápido |
| 2 hilos | ~6 ms | ~2 ms | Go ~3x más rápido |
| 4 hilos | ~5 ms | ~2 ms | Go ~2.5x más rápido |
| 8 hilos | ~5 ms | ~2 ms | Go ~2.5x más rápido |

> Nota: En modo resumen no hay iteración real (cada hilo suma su bloque de una vez con `addAndGet`), por lo que los tiempos son muy bajos en ambos lenguajes. La diferencia principal es el tiempo de inicio de la JVM.

### Comparacion de resultados — modo detallado

En modo detallado sí hay iteración número a número. El cuello de botella es la escritura en consola (I/O), no el CPU.

| Cantidad de hilos | Límite | Tiempo Java | Tiempo Go | Diferencia |
|:-----------------:|:------:|:-----------:|:---------:|:----------:|
| 2 hilos | 1.000.000 | **101.503 s** | **14.483 s** | Go es **~7x más rápido** |

> Prueba ejecutada en modo detallado (número a número con impresión por consola), 2 hilos, límite = 1.000.000. Resultados reales medidos en la misma máquina.

Aumentar hilos en modo detallado ayuda poco porque el I/O de consola es el verdadero cuello de botella, no el procesamiento.

### Conclusion

La diferencia más llamativa entre los dos lenguajes se ve en el **modo detallado con impresión por consola**:

- **Java tardó 101.503 segundos** para imprimir 1.000.000 de números con 2 hilos.
- **Go tardó 14.483 segundos** para el mismo trabajo — aproximadamente **7 veces más rápido**.

Esta diferencia se explica principalmente por tres razones:

1. **I/O de consola más eficiente en Go:** El runtime de Go maneja la salida estándar de forma más liviana que la JVM. Cada llamada a `fmt.Printf` genera menos overhead que `System.out.println` en Java.
2. **Menor costo de las goroutines:** Aunque en este caso solo son 2 trabajadores, el scheduler de Go gestiona el I/O de forma más eficiente al ser parte del propio runtime, sin depender del scheduler del sistema operativo.
3. **Ausencia de JVM warmup:** Go compila a binario nativo, mientras que Java necesita tiempo de arranque de la JVM y compilación JIT antes de alcanzar su velocidad óptima.

En el **modo resumen** (sin impresión número a número) la diferencia es mínima en ambos lenguajes, porque no hay I/O intensivo y el conteo se hace con una sola operación atómica por hilo. Ahí Java y Go son prácticamente equivalentes.

**Conclusión general:** Si el programa necesita imprimir o procesar grandes volúmenes de datos de forma concurrente, **Go ofrece una ventaja real y medible** sobre Java. Si el trabajo es principalmente cálculo sin I/O intensivo, ambos lenguajes tienen rendimiento similar.

### Cuando usar cada uno

| Situación | Recomendado |
|-----------|-------------|
| Procesamiento con mucho I/O de consola o archivos | Go |
| Proyecto académico o empresarial con ecosistema Java | Java |
| Necesitas el menor tiempo de inicio posible | Go |
| Muchos trabajadores concurrentes (miles) | Go (goroutines son más ligeras) |
| Pocos hilos con lógica compleja | Java o Go (similar) |
| Equipo ya familiarizado con la JVM | Java |

---

## 6. Estructura del proyecto

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

## 7. Compilar y ejecutar

### 7.1 Script principal run.bat

El archivo `run.bat` muestra un menú para elegir qué versión ejecutar:

```
run.bat
```

```
============================================
  ContadorParalelo -- Selector de lenguaje
============================================

 [1]  Java   (Threads + AtomicLong)
 [2]  Go     (Goroutines + sync/atomic)
 [3]  Ambos  (Java primero, luego Go)
 [0]  Salir

Elige una opcion:
```

> Requiere tener el **JDK** y **Go** instalados y disponibles en el PATH del sistema.

### 7.2 Java manualmente

```bash
# Paso 1: compilar
javac src/java/ContadorParalelo.java -d out

# Paso 2: ejecutar
java -cp out ContadorParalelo
```

### 7.3 Go manualmente

```bash
# Opcion A: ejecutar directo sin compilar
go run src/go/contador_paralelo.go

# Opcion B: compilar y luego ejecutar
go build -o contador src/go/contador_paralelo.go
./contador
```

---

## 8. Parametros del programa

Ambas versiones piden los mismos tres datos al ejecutarse:

| # | Parámetro | Valores válidos | Descripción |
|:-:|-----------|:---------------:|-------------|
| 1 | Número final | 1 a 50.000.000.000 | Hasta qué número se va a contar |
| 2 | Cantidad de hilos / goroutines | 1 en adelante | Cuántos trabajadores paralelos se usan |
| 3 | Modo | 1 o 2 | `1` = resumen (una línea por hilo) / `2` = detallado (número a número) |

---

## 9. Ejemplos de salida

### 9.1 Modo resumen

**Java — 4 hilos, límite = 1000**

```
Número final (1 - 50.000.000.000): 1000
Cantidad de hilos: 4
Modo (1 = resumen | 2 = número a número): 1

Contando del 1 al 1000 usando 4 hilo(s)...

  Hilo-1: del 1 al 250  (250 números)
  Hilo-2: del 251 al 500  (250 números)
  Hilo-3: del 501 al 750  (250 números)
  Hilo-4: del 751 al 1000  (250 números)

--- Resultado ---
Total contado  : 1000
Tiempo de mora : 3 ms (0.003 s)
```

**Go — 4 goroutines, límite = 1000**

```
Número final (1 - 50.000.000.000): 1000
Cantidad de goroutines: 4
Modo (1 = resumen | 2 = número a número): 1

Contando del 1 al 1000 usando 4 goroutine(s)...

  Goroutine-1: del 1 al 250  (250 números)
  Goroutine-2: del 251 al 500  (250 números)
  Goroutine-3: del 501 al 750  (250 números)
  Goroutine-4: del 751 al 1000  (250 números)

--- Resultado ---
Total contado  : 1000
Tiempo de mora : 1 ms (0.001 s)
```

### 9.2 Modo detallado

**Java — 2 hilos, límite = 6**

```
Número final (1 - 50.000.000.000): 6
Cantidad de hilos: 2
Modo (1 = resumen | 2 = número a número): 2

Contando del 1 al 6 usando 2 hilo(s)...

  [Hilo-1] --> 1
  [Hilo-1] --> 2
  [Hilo-1] --> 3
  [Hilo-2] --> 4
  [Hilo-2] --> 5
  [Hilo-2] --> 6

--- Resultado ---
Total contado  : 6
Tiempo de mora : 5 ms (0.005 s)
```

**Go — 2 goroutines, límite = 6**

```
Número final (1 - 50.000.000.000): 6
Cantidad de goroutines: 2
Modo (1 = resumen | 2 = número a número): 2

Contando del 1 al 6 usando 2 goroutine(s)...

  [Goroutine-1] --> 1
  [Goroutine-1] --> 2
  [Goroutine-1] --> 3
  [Goroutine-2] --> 4
  [Goroutine-2] --> 5
  [Goroutine-2] --> 6

--- Resultado ---
Total contado  : 6
Tiempo de mora : 2 ms (0.002 s)
```
