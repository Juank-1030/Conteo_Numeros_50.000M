<div align="center">

# Contador Paralelo

### Contando números hasta 50 mil millones — Java & Go

> Implementación del mismo algoritmo en dos lenguajes: divide el rango `[1, límite]`
> en segmentos iguales, asigna cada segmento a un hilo o goroutine independiente,
> y garantiza un conteo correcto usando operaciones atómicas libres de cerrojos (CAS).

---

![Java](https://img.shields.io/badge/Java-21-007396?style=for-the-badge&logo=openjdk&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Threads](https://img.shields.io/badge/Concurrencia-Threads%20%7C%20Goroutines-FF6B35?style=for-the-badge)
![AtomicLong](https://img.shields.io/badge/Conteo-AtomicLong%20%7C%20sync%2Fatomic-blueviolet?style=for-the-badge)

</div>

---

## Índice de contenidos

1. [Descripción del programa](#1-descripción-del-programa)
2. [Tecnologías utilizadas](#2-tecnologías-utilizadas)
   - [2.1 Java](#21-java)
   - [2.2 Go](#22-go)
3. [Cómo funciona la concurrencia](#3-cómo-funciona-la-concurrencia)
   - [3.1 División del rango](#31-división-del-rango)
   - [3.2 Hilos vs Goroutines](#32-hilos-vs-goroutines)
   - [3.3 Conteo atómico sin cerrojos](#33-conteo-atómico-sin-cerrojos)
   - [3.4 Esperando a que todos los trabajadores terminen](#34-esperando-a-que-todos-los-trabajadores-terminen)
4. [Funciones en cada versión](#4-funciones-en-cada-versión)
5. [Tabla comparativa Java vs Go](#5-tabla-comparativa-java-vs-go)
6. [Estructura del proyecto](#6-estructura-del-proyecto)
7. [Compilación y ejecución](#7-compilación-y-ejecución)
   - [7.1 Script principal run.bat](#71-script-principal-runbat)
   - [7.2 Java manualmente](#72-java-manualmente)
   - [7.3 Go manualmente](#73-go-manualmente)
8. [Parámetros del programa](#8-parámetros-del-programa)
9. [Ejemplos de salida](#9-ejemplos-de-salida)
   - [9.1 Modo resumen](#91-modo-resumen)
   - [9.2 Modo detallado](#92-modo-detallado)
10. [Análisis profundo del rendimiento](#10-análisis-profundo-del-rendimiento)
    - [10.1 Modo resumen: diferencias mínimas](#101-modo-resumen-diferencias-mínimas)
    - [10.2 Modo detallado: el cuello de botella es la E/S](#102-modo-detallado-el-cuello-de-botella-es-la-es)
    - [10.3 Mediciones experimentales a escala](#103-mediciones-experimentales-a-escala)
    - [10.4 ¿Por qué Java puede superar a Go con muchos hilos?](#104-por-qué-java-puede-superar-a-go-con-muchos-hilos)
    - [10.5 La solución: buffer en Go con bufio.Writer](#105-la-solución-buffer-en-go-con-bufiowriter)
    - [10.6 Conclusión del análisis](#106-conclusión-del-análisis)
11. [¿Cuándo usar cada uno?](#11-cuándo-usar-cada-uno)

---

## 1. Descripción del programa

Este programa es un **contador paralelo** que cuenta números desde 1 hasta un límite especificado por el usuario, utilizando múltiples unidades de ejecución que trabajan simultáneamente. Está implementado en **Java** y en **Go** con el propósito de comparar el rendimiento, la sintaxis y el modelo de concurrencia de ambos lenguajes.

### ¿Qué problema resuelve?

Contar números secuencialmente del 1 al 50,000,000,000 (50 mil millones) sería extremadamente lento si se hiciera con un solo hilo de ejecución. Al dividir el trabajo entre múltiples hilos o goroutines, el tiempo total se reduce drásticamente porque cada unidad de ejecución procesa una parte del rango de forma independiente y en paralelo.

### Entradas del programa

El programa solicita tres parámetros al usuario:

| # | Parámetro | Rango de valores | Descripción |
|:-:|-----------|:-----------------:|-------------|
| 1 | **Número final** | 1 a 50,000,000,000 | El límite superior del rango a contar |
| 2 | **Número de hilos / goroutines** | 1 o más | Cuántos trabajadores en paralelo se usarán |
| 3 | **Modo** | 1 o 2 | `1` = modo resumen (una línea por hilo), `2` = modo detallado (número por número) |

### Modos de salida

- **Modo resumen (1):** Cada hilo imprime una sola línea indicando qué rango procesó y cuántos números contó. Es rápido porque no hay E/S intensiva.
- **Modo detallado (2):** Cada hilo imprime cada número individualmente. Es lento porque la salida por consola es una operación de E/S costosa, pero sirve para verificar visualmente que el reparto del trabajo es correcto.

### ¿Qué hace exactamente cada trabajador?

Cada hilo o goroutine:
1. Recibe un rango de números `[inicio, fin]` que debe procesar.
2. Itera sobre cada número en ese rango.
3. En modo resumen, suma la cantidad de números de una vez al contador global.
4. En modo detallado, incrementa el contador global de uno en uno y además imprime cada número por consola.
5. Al terminar, notifica al hilo principal que ha completado su trabajo.

### Resumen visual

| Aspecto | Cómo lo logra |
|---------|---------------|
| Contar del 1 al 50 mil millones | Divide el rango entre N hilos/goroutines en paralelo |
| Evitar que dos hilos escriban al mismo tiempo | Usa operaciones atómicas CAS (sin `synchronized` ni `mutex`) |
| Ofrecer dos modos de visualización | Modo resumen (una línea por hilo) y modo detallado (número por número) |
| Medir el tiempo real de ejecución | Toma el tiempo antes de que los hilos empiecen y después de que todos terminen |

---

## 2. Tecnologías utilizadas

### 2.1 Java

Java utiliza el modelo de **hilos a nivel de sistema operativo** (1:1). Cada objeto `Thread` crea un hilo nativo del SO, con su propia pila de llamadas y su propio contexto de ejecución.

| Clase / API | Paquete | Propósito |
|-------------|---------|-----------|
| `Thread` | `java.lang` | Crear e iniciar cada hilo del sistema operativo |
| `AtomicLong` | `java.util.concurrent.atomic` | Contador compartido entre hilos sin usar `synchronized` |
| `Scanner` | `java.util` | Leer los parámetros del usuario desde la consola |
| `System.currentTimeMillis()` | `java.lang` | Medir el tiempo de ejecución en milisegundos |

**AtomicLong** es la pieza clave: proporciona operaciones como `addAndGet()` que internamente usan la instrucción CAS (Compare-And-Swap) del procesador. Esto permite que múltiples hilos actualicen el mismo contador sin bloquearse mutuamente, logrando un rendimiento mucho mayor que con bloques `synchronized`.

> No utiliza dependencias externas — Solo emplea el JDK estándar.

### 2.2 Go

Go utiliza **goroutines**, que son unidades de ejecución ligeras gestionadas por el propio runtime de Go (modelo M:N). Múltiples goroutines se multiplexan sobre un número menor de hilos del SO.

| Paquete | Elemento | Propósito |
|---------|----------|-----------|
| `sync` | `WaitGroup` | Esperar a que todas las goroutines terminen antes de mostrar el resultado |
| `sync/atomic` | `AddInt64()` | Contador compartido entre goroutines sin usar `mutex` |
| `time` | `Now()`, `Since()` | Medir el tiempo de ejecución |
| `fmt` | `Scan()`, `Printf()` | Leer entrada del usuario e imprimir resultados |
| `os` | `Exit()` | Terminar el programa si el usuario ingresa un valor inválido |

**sync.WaitGroup** funciona como un contador interno: se incrementa con `Add()` por cada goroutine lanzada, cada goroutine llama a `Done()` al finalizar (decrementa el contador), y la goroutine principal llama a `Wait()` que bloquea la ejecución hasta que el contador llegue a cero.

> No utiliza dependencias externas — Solo emplea la biblioteca estándar de Go.

---

## 3. Cómo funciona la concurrencia

### 3.1 División del rango

El rango total `[1, límite]` se divide en N partes iguales. El último trabajador siempre se queda con el resto para asegurar que no falte ningún número:

```
Trabajador 1  -->  desde 1        hasta tamaño
Trabajador 2  -->  desde tamaño+1 hasta 2*tamaño
...
Trabajador N  -->  desde ...      hasta límite  (absorbe el resto)
```

#### Fórmula de división

```
tamaño = límite / cantidadHilos
resto  = límite % cantidadHilos
```

- Los primeros `(cantidadHilos - resto)` trabajadores procesan exactamente `tamaño` números.
- Los últimos `resto` trabajadores procesan `tamaño + 1` números.
- Si el límite es divisible exactamente, todos procesan la misma cantidad.

#### Ejemplo concreto

Con límite = 1000 y 4 hilos:
- Hilo 1: cuenta del 1 al 250 (250 números)
- Hilo 2: cuenta del 251 al 500 (250 números)
- Hilo 3: cuenta del 501 al 750 (250 números)
- Hilo 4: cuenta del 751 al 1000 (250 números)

#### Ejemplo con resto

Con límite = 1000 y 3 hilos:
- `tamaño = 1000 / 3 = 333`
- `resto = 1000 % 3 = 1`
- Hilo 1: del 1 al 333 (333 números)
- Hilo 2: del 334 al 666 (333 números)
- Hilo 3: del 667 al 1000 (334 números — absorbe el resto)

Este método garantiza que **ningún número quede sin contar** y que **ningún número se cuente dos veces**.

### 3.2 Hilos vs Goroutines

| Característica | Java — `Thread` | Go — `goroutine` |
|----------------|-----------------|-------------------|
| Modelo de concurrencia | 1:1 (un hilo Java = un hilo del SO) | M:N (M goroutines multiplexadas sobre N hilos del SO) |
| Huella de memoria | ~1 MB de pila por hilo | ~2 KB de pila inicial (crece si es necesario) |
| Cómo se crea | `new Thread(() -> ...).start()` | `go funcion(...)` |
| Quién los gestiona | La JVM delega en el sistema operativo | El runtime de Go tiene su propio planificador |
| Cómo esperar | `hilo.join()` | `waitGroup.Wait()` |
| Cuántos se pueden crear fácilmente | Miles (limitado por RAM) | Millones (muy ligeros) |
| Costo de creación | Relativamente alto (llamada al sistema) | Muy bajo (solo asignar ~2 KB en el heap de Go) |
| Planificación | El sistema operativo decide cuándo ejecutar cada hilo | El runtime de Go decide cuándo ejecutar cada goroutine (cooperativa en puntos de suspensión) |

#### Explicación detallada del modelo M:N de Go

El runtime de Go implementa un planificador **M:N**:
- **M** = número de goroutines (pueden ser millones)
- **N** = número de hilos del SO (normalmente igual a `GOMAXPROCS`, que por defecto es el número de núcleos de CPU)

El planificador de Go distribuye las M goroutines entre los N hilos del SO. Cuando una goroutine hace una operación que bloquea (como E/S o un llamado al sistema), el planificador automáticamente:
1. Suspende esa goroutine
2. Crea o reutiliza un hilo del SO para las goroutines restantes
3. Reanuda la goroutine bloqueada cuando la operación se completa

Esto hace que las goroutines sean mucho más eficientes que los hilos del SO para tareas con muchas unidades de trabajo concurrentes.

**En este programa**, la diferencia no es muy notable cuando se usan pocos trabajadores (por ejemplo, 4 u 8), pero en aplicaciones con cientos o miles de trabajadores concurrentes, Go escala mucho mejor.

### 3.3 Conteo atómico sin cerrojos

#### El problema de la condición de carrera

Cuando múltiples hilos intentan sumar al mismo contador simultáneamente, pueden ocurrir errores de concurrencia. Por ejemplo:

```java
// ESTO NO ES SEGURO EN MULTIHILO:
contador = contador + 1;
```

Esta operación aparentemente simple se descompone en tres pasos a nivel de CPU:
1. Leer `contador` de la memoria
2. Sumar 1 al valor leído
3. Escribir el resultado de vuelta a la memoria

Si dos hilos ejecutan esto al mismo tiempo:
- Hilo A: Lee contador = 5
- Hilo B: Lee contador = 5 (antes de que A escriba)
- Hilo A: Escribe 6
- Hilo B: Escribe 6 (¡debería ser 7!)

El resultado final es 6 en lugar de 7: se ha perdido una operación de incremento. Esto se llama **condición de carrera** (race condition).

#### Solución clásica: cerrojos

La solución tradicional es usar cerrojos (locks) para garantizar que solo un hilo acceda al contador a la vez:

```java
// Java: bloque synchronized
synchronized (contador) {
    contador++;
}
```

```go
// Go: mutex
mutex.Lock()
contador++
mutex.Unlock()
```

El problema de los cerrojos es que crean **contención**: si un hilo tiene el cerrojo, los demás tienen que esperar en una cola, lo que reduce el rendimiento.

#### Solución optimizada: operaciones atómicas CAS

En lugar de usar cerrojos, ambas versiones usan operaciones **CAS (Compare-And-Swap)**, que son instrucciones directas del procesador. CAS funciona así:

1. Leer el valor actual de la memoria
2. Calcular el nuevo valor
3. Intentar escribir el nuevo valor **solo si la memoria sigue teniendo el valor original**
4. Si otro hilo modificó la memoria mientras tanto, el CAS **falla** y se reintenta desde el paso 1

Este enfoque:
- Es **libre de cerrojos** (no hay colas de espera)
- Es **más rápido** que los cerrojos cuando la contención es baja
- Garantiza **seguridad en concurrencia** a nivel de hardware

```java
// Java — modo resumen (agrega todo el bloque de una vez)
contadorTotal.addAndGet(cantidadNumeros);

// Java — modo detallado (agrega de uno en uno)
contadorTotal.incrementAndGet();
```

```go
// Go — modo resumen
atomic.AddInt64(&contadorTotal, cantidadNumeros)

// Go — modo detallado
atomic.AddInt64(&contadorTotal, 1)
```

### 3.4 Esperando a que todos los trabajadores terminen

El hilo principal no puede mostrar el resultado hasta que todos los trabajadores hayan terminado. Cada lenguaje tiene su propia forma de hacerlo:

#### Java: `Thread.join()`

```java
// Java: espera hilo por hilo
for (Thread hilo : hilos) {
    hilo.join();  // El hilo principal se bloquea aquí hasta que 'hilo' termine
}
```

`join()` es un método que bloquea al hilo que lo llama hasta que el hilo objetivo muere. En el bucle, el hilo principal espera al hilo 1, luego al 2, y así sucesivamente. Como todos los hilos ya están en ejecución (se iniciaron antes del bucle), el tiempo total de espera es aproximadamente el tiempo del hilo más lento.

#### Go: `sync.WaitGroup`

```go
var grupo sync.WaitGroup

// Por cada goroutine que lanzamos:
grupo.Add(1)           // incrementa el contador interno
go funcion(...)        // lanza la goroutine

// Al inicio de cada goroutine:
defer grupo.Done()     // al terminar, decrementa el contador (ejecutado al salir)

// En main, después de lanzar todas:
grupo.Wait()           // bloquea hasta que el contador llegue a cero
```

`WaitGroup` es esencialmente un contador atómico:
- `Add(delta)` incrementa el contador en `delta`
- `Done()` decrementa el contador en 1 (equivalente a `Add(-1)`)
- `Wait()` bloquea hasta que el contador es 0

---

## 4. Funciones en cada versión

| Función | Java | Go | Propósito |
|---------|------|----|-----------|
| Leer número | `pedirNumero(Scanner, String, long, long)` | `pedirNumero(string, int64, int64)` | Lee un número de la consola y valida que esté dentro del rango permitido |
| Contar segmento | `contarSegmento(int, long, long, boolean, AtomicLong)` | `contarSegmento(int, int64, int64, bool, *int64, *WaitGroup)` | Lógica que ejecuta cada hilo/goroutine para contar su segmento |
| Punto de entrada | `main(String[])` | `main()` | Lee parámetros, crea los trabajadores, espera a que terminen y muestra el resultado |

### Descripción detallada de cada función

#### `pedirNumero` (ambos lenguajes)

Esta función es responsable de la interacción con el usuario:
1. Muestra un mensaje descriptivo pidiendo el valor
2. Lee el valor ingresado por el usuario
3. Verifica que el valor sea un número válido dentro del rango especificado
4. Si el valor es inválido, muestra un mensaje de error y termina el programa
5. Si el valor es válido, lo devuelve para su uso posterior

Parámetros:
- **Mensaje:** Texto descriptivo que se muestra al usuario
- **Mínimo:** Valor mínimo aceptable
- **Máximo:** Valor máximo aceptable

#### `contarSegmento` (ambos lenguajes)

Esta función contiene el núcleo del algoritmo paralelo y es ejecutada por cada trabajador:

1. Determina el rango `[inicio, fin]` que debe procesar basándose en su número de trabajador y la cantidad total de trabajadores
2. Itera desde `inicio` hasta `fin` (inclusive)
3. En **modo resumen**: simplemente suma la cantidad de números del segmento al contador global con una sola operación atómica
4. En **modo detallado**: por cada número, incrementa el contador global atómicamente y además imprime el número por consola
5. En **Go**: llama a `grupo.Done()` al finalizar para notificar al `WaitGroup`

Parámetros:
- **númeroTrabajador:** Identificador del trabajador (1, 2, 3, ...)
- **inicio:** Primer número del rango a procesar (inclusive)
- **fin:** Último número del rango a procesar (inclusive)
- **modoDetallado:** Booleano que indica si se debe imprimir cada número
- **contadorTotal:** Referencia al contador compartido (AtomicLong en Java, *int64 en Go)
- **grupo:** (solo Go) Referencia al WaitGroup para notificar cuando termine

#### `main` (ambos lenguajes)

Punto de entrada del programa:
1. Muestra un banner de bienvenida
2. Llama a `pedirNumero` tres veces para obtener los parámetros del usuario
3. Calcula cómo dividir el rango entre los trabajadores
4. Crea e inicia todos los trabajadores
5. Espera a que todos terminen
6. Muestra el resultado final (total contado y tiempo transcurrido)

---

## 5. Tabla comparativa Java vs Go

### Comparación técnica

| Aspecto | Java | Go |
|---------|------|----|
| Unidad de concurrencia | `Thread` (hilo del SO) | `goroutine` (hilo ligero del runtime) |
| Memoria por trabajador | ~1 MB de pila | ~2 KB de pila inicial |
| Operación atómica usada | `AtomicLong.addAndGet()` | `atomic.AddInt64()` |
| Mecanismo de espera | `Thread.join()` | `sync.WaitGroup` |
| Inicio de trabajadores | `hilo.start()` | `go funcion()` |
| Planificador | Sistema operativo | Runtime de Go (M:N, multiplexado) |
| Sintaxis concurrente | Verbosa (clase Thread, lambda) | Muy simple (`go` + función) |
| Tiempo de arranque | ~200-400 ms (calentamiento de JVM) | ~5-20 ms (binario nativo) |
| Compilación | Abytecode (JIT a código máquina en tiempo de ejecución) | A código máquina nativo en tiempo de compilación |
| Dependencias | JDK estándar (sin externas) | Biblioteca estándar de Go (sin externas) |

### Comparación de resultados — modo resumen

Los tiempos son aproximados y varían según el hardware. Las pruebas se ejecutaron con límite = 1,000,000,000 (mil millones).

| Número de hilos | Tiempo Java (ms) | Tiempo Go (ms) | Diferencia |
|:----------------:|:-----------------:|:---------------:|:----------:|
| 1 hilo | ~8 ms | ~3 ms | Go ~2.5x más rápido |
| 2 hilos | ~6 ms | ~2 ms | Go ~3x más rápido |
| 4 hilos | ~5 ms | ~2 ms | Go ~2.5x más rápido |
| 8 hilos | ~5 ms | ~2 ms | Go ~2.5x más rápido |

> **Nota importante:** En modo resumen no hay iteración real (cada hilo agrega todo su bloque de una vez con `addAndGet` y `AddInt64`), por lo que los tiempos son muy bajos en ambos lenguajes. La principal diferencia es el tiempo de arranque de la JVM.

### Comparación de resultados — modo detallado

En modo detallado hay una iteración real número por número. El cuello de botella es la E/S de consola, no la CPU.

| Número de hilos | Límite | Tiempo Java | Tiempo Go | Diferencia |
|:----------------:|:------:|:-----------:|:---------:|:----------:|
| 2 hilos | 1,000,000 | **101.503 s** | **14.483 s** | Go es **~7x más rápido** |

> Prueba ejecutada en modo detallado (número por número con salida a consola), 2 hilos, límite = 1,000,000. Resultados reales medidos en la misma máquina.

Agregar más hilos en modo detallado ayuda poco porque la E/S de consola es el verdadero cuello de botella, no el procesamiento.

### Mediciones experimentales — modo detallado (escala creciente)

Pruebas ejecutadas aumentando simultáneamente el número de hilos y el límite de conteo, todas en **modo detallado** (salida número por número por consola).

| Número de hilos | Límite | Tiempo Java (s) | Tiempo Go (s) | Observación |
|:----------------:|:------:|:----------------:|:--------------:|-------------|
| 2 | 10,000 | 0.585 | 0.009 | Go ~65x más rápido |
| 2 | 100,000 | 4.662 | 0.509 | Go ~9x más rápido |
| 10 | 1,000,000 | 48.585 | 5.820 | Go ~8x más rápido |
| 15 | 50,000,000 | ⚠️ No pudo completarse | 387.185 | Java colapsó el editor de código |
| 250 | 50,000,000 | ⚠️ No pudo completarse | 377.922 | Java colapsó el editor de código |

> ⚠️ **Nota sobre las entradas de Java:** Con límite = 50,000,000, Java no pudo completar la ejecución en ninguna de las dos configuraciones (15 y 250 hilos). En ambos casos, VS Code se cerró inesperadamente, probablemente debido a agotamiento de RAM o agotamiento de descriptores de archivo del SO, causado por el volumen de salida de consola y la carga de hilos del SO sostenida durante un período prolongado.

---

## 6. Estructura del proyecto

```
Conteo_Numeros_50.000M/
│
├── src/
│   ├── java/
│   │   └── ContadorParalelo.java    ← Implementación en Java
│   └── go/
│       └── contador_paralelo.go     ← Implementación en Go
│
├── run.bat                          ← Script para ejecutar ambas versiones
├── README.md                        ← Documentación en inglés
└── README.es.md                     ← Documentación en español (este archivo)
```

### Descripción de cada archivo

- **`src/java/ContadorParalelo.java`**: Implementación completa en Java. Contiene la clase `ContadorParalelo` con el método `main`, la función `pedirNumero` y la función `contarSegmento`. Utiliza `Thread`, `AtomicLong` y `Scanner` del JDK estándar.

- **`src/go/contador_paralelo.go`**: Implementación completa en Go. Contiene la función `main()`, la función `pedirNumero` y la función `contarSegmento`. Utiliza `goroutines`, `sync.WaitGroup`, `sync/atomic` y `bufio.Writer` de la biblioteca estándar de Go.

- **`run.bat`**: Script de Windows que muestra un menú interactivo para seleccionar qué versión ejecutar. Requiere que JDK y Go estén instalados y disponibles en el PATH del sistema.

---

## 7. Compilación y ejecución

### 7.1 Script principal run.bat

El archivo `run.bat` muestra un menú para elegir qué versión ejecutar:

```bash
# En la terminal de Windows:
run.bat
```

```
============================================
  ContadorParalelo -- Selector de Lenguaje
============================================

 [1]  Java   (Threads + AtomicLong)
 [2]  Go     (Goroutines + sync/atomic)
 [3]  Ambos  (Java primero, luego Go)
 [0]  Salir

Elige una opción:
```

> Requiere que el **JDK** y **Go** estén instalados y disponibles en el PATH del sistema.

### 7.2 Java manualmente

```bash
# Paso 1: compilar
javac src/java/ContadorParalelo.java -d out

# Paso 2: ejecutar
java -cp out ContadorParalelo
```

**Explicación:**
- `javac` compila el archivo `.java` y genera archivos `.class` con bytecode
- `-d out` especifica el directorio de salida para los archivos compilados
- `java -cp out` ejecuta la JVM usando `out` como classpath
- No es necesario especificar la extensión `.class` al ejecutar

### 7.3 Go manualmente

```bash
# Opción A: ejecutar directamente sin compilar (útil para pruebas rápidas)
go run src/go/contador_paralelo.go

# Opción B: compilar y luego ejecutar
go build -o contador src/go/contador_paralelo.go
./contador
```

**Explicación:**
- `go run` compila y ejecuta en un solo paso, sin dejar un binario
- `go build` genera un binario ejecutable nativo
- `-o contador` especifica el nombre del archivo de salida
- El binario generado es un ejecutable independiente que no necesita Go instalado para ejecutarse

---

## 8. Parámetros del programa

Ambas versiones solicitan las mismas tres entradas en tiempo de ejecución:

| # | Parámetro | Valores válidos | Descripción |
|:-:|-----------|:----------------:|-------------|
| 1 | **Número final** | 1 a 50,000,000,000 | Hasta qué número contar (límite superior inclusivo) |
| 2 | **Número de hilos / goroutines** | 1 o más | Cuántos trabajadores en paralelo se utilizarán |
| 3 | **Modo** | 1 o 2 | `1` = modo resumen (una línea por hilo), `2` = modo detallado (número por número) |

### Validación de parámetros

- **Número final:** Debe ser un entero entre 1 y 50,000,000,000. Si el usuario ingresa un valor fuera de este rango o un valor no numérico, el programa muestra un mensaje de error y termina.
- **Número de hilos:** Debe ser un entero positivo (1 o más). Valores como 0 o negativos son rechazados.
- **Modo:** Debe ser exactamente 1 (resumen) o 2 (detallado). Cualquier otro valor es rechazado.

---

## 9. Ejemplos de salida

### 9.1 Modo resumen

**Java — 4 hilos, límite = 1000**

```
Número final (1 - 50,000,000,000): 1000
Número de hilos: 4
Modo (1 = resumen | 2 = número por número): 1

Contando del 1 al 1000 usando 4 hilo(s)...

  Thread-1: del 1 al 250  (250 números)
  Thread-2: del 251 al 500  (250 números)
  Thread-3: del 501 al 750  (250 números)
  Thread-4: del 751 al 1000  (250 números)

--- Resultado ---
Total contado  : 1000
Tiempo total   : 3 ms (0.003 s)
```

**Go — 4 goroutines, límite = 1000**

```
Número final (1 - 50,000,000,000): 1000
Número de goroutines: 4
Modo (1 = resumen | 2 = número por número): 1

Contando del 1 al 1000 usando 4 goroutine(s)...

  Goroutine-1: del 1 al 250  (250 números)
  Goroutine-2: del 251 al 500  (250 números)
  Goroutine-3: del 501 al 750  (250 números)
  Goroutine-4: del 751 al 1000  (250 números)

--- Resultado ---
Total contado  : 1000
Tiempo total   : 1 ms (0.001 s)
```

### 9.2 Modo detallado

**Java — 2 hilos, límite = 6**

```
Número final (1 - 50,000,000,000): 6
Número de hilos: 2
Modo (1 = resumen | 2 = número por número): 2

Contando del 1 al 6 usando 2 hilo(s)...

  [Thread-1] --> 1
  [Thread-1] --> 2
  [Thread-1] --> 3
  [Thread-2] --> 4
  [Thread-2] --> 5
  [Thread-2] --> 6

--- Resultado ---
Total contado  : 6
Tiempo total   : 5 ms (0.005 s)
```

**Go — 2 goroutines, límite = 6**

```
Número final (1 - 50,000,000,000): 6
Número de goroutines: 2
Modo (1 = resumen | 2 = número por número): 2

Contando del 1 al 6 usando 2 goroutine(s)...

  [Goroutine-1] --> 1
  [Goroutine-1] --> 2
  [Goroutine-1] --> 3
  [Goroutine-2] --> 4
  [Goroutine-2] --> 5
  [Goroutine-2] --> 6

--- Resultado ---
Total contado  : 6
Tiempo total   : 2 ms (0.002 s)
```

---

## 10. Análisis profundo del rendimiento

### 10.1 Modo resumen: diferencias mínimas

En modo resumen, cada hilo ejecuta una única operación atómica (`addAndGet` o `AddInt64`) para sumar todos los números de su segmento de una sola vez. No hay iteración número por número ni E/S de consola intensiva.

```
Trabajador 1: contadorTotal.addAndGet(250)   ← una sola operación
Trabajador 2: contadorTotal.addAndGet(250)   ← una sola operación
Trabajador 3: contadorTotal.addAndGet(250)   ← una sola operación
Trabajador 4: contadorTotal.addAndGet(250)   ← una sola operación
```

El resultado es que ambos lenguajes completan la tarea en **milisegundos**, y la diferencia principal se debe al tiempo de arranque:
- **Java:** 200-400 ms de calentamiento de la JVM (carga de clases, compilación JIT inicial)
- **Go:** 5-20 ms de arranque (binario compilado a nativo)

En modo resumen, **Java y Go son prácticamente equivalentes** en rendimiento de procesamiento puro.

### 10.2 Modo detallado: el cuello de botella es la E/S

En modo detallado, cada hilo itera número por número y para cada uno:
1. Realiza un incremento atómico del contador global
2. Imprime el número por consola

El problema es que **imprimir por consola es una operación de E/S** que involucra:
- Una llamada al sistema (syscall) al sistema operativo
- Escritura en el búfer de la terminal
- Sincronización del flujo de salida

Con 1,000,000 de números para imprimir, se generan 1,000,000 de operaciones de E/S de consola. **El tiempo de CPU dedicado al conteo es insignificante comparado con el tiempo de E/S.**

Esto explica por qué:
- Agregar más hilos **no mejora significativamente** el tiempo en modo detallado (la E/S sigue siendo el cuello de botella)
- Go es consistentemente más rápido en este modo (como veremos a continuación)

### 10.3 Mediciones experimentales a escala

Las pruebas muestran tres patrones claros:

1. **Go es consistentemente más rápido en modo detallado.** La ventaja va de ~8x a ~65x dependiendo del límite. La diferencia es mayor con límites pequeños (el overhead de arranque de la JVM pesa más) y se estabiliza alrededor de 8-9x con límites grandes.

2. **Agregar goroutines apenas reduce el tiempo en Go con límites muy altos.** Con 50,000,000 de números, pasar de 15 a 250 goroutines solo redujo el tiempo de 387 s a 377 s (menos del 3%). Esto confirma que con tantos números para imprimir, el cuello de botella es la E/S de consola, no el procesamiento paralelo.

3. **Java no es viable para modo detallado a gran escala.** El modelo de hilos del SO consume memoria proporcional al número de hilos activos (~1 MB de pila cada uno) y mantiene miles de descriptores de salida abiertos. Combinado con 50,000,000 de líneas para imprimir, la presión sobre la JVM y el SO resulta en una falla catastrófica que cierra el entorno de desarrollo.

### 10.4 ¿Por qué Java puede superar a Go con muchos hilos?

Con configuraciones de **muchos trabajadores y pocos datos por trabajador** (por ejemplo, 10,000 hilos para 1,000,000 de números = 100 números por hilo), Java puede superar a Go. Esto parece contradictorio pero tiene una explicación técnica precisa.

#### El cuello de botella: llamadas al sistema de escritura en consola

El problema radica en cómo cada lenguaje escribe en la consola a nivel del sistema operativo:

| Aspecto | Java `System.out.println` | Go `fmt.Printf` (versión original) |
|---------|---------------------------|-------------------------------------|
| Búfer en espacio de usuario | **Sí** — 8 KB `BufferedOutputStream` | **No** — escribe directamente al SO |
| Llamadas al SO por 1,000,000 impresiones | Pocas (el búfer agrupa escrituras) | ~1,000,000 llamadas individuales |
| Contención entre hilos | Sincronizado sobre el búfer (eficiente) | Cada goroutine compite por stdout |

**Java usa un búfer interno de 8 KB:** las llamadas a `println` acumulan texto en memoria y solo hacen una llamada al SO cuando el búfer se llena. Con 10,000 hilos imprimiendo 100 números cada uno, el número total de llamadas al SO es muy bajo.

**Go (versión original) no tiene búfer:** cada `fmt.Printf` llama a `write()` directamente en el SO. Con 1,000,000 de números para imprimir, esto genera casi 1,000,000 de llamadas individuales al SO, lo que degrada severamente el rendimiento con muchos trabajadores.

#### Factores adicionales que favorecen a Java en este escenario

1. **Optimización JIT:** Con 10,000 hilos activos, la JVM tiene tiempo para compilar el código caliente de los hilos a código nativo optimizado. El bucle de impresión de 100 números se vuelve muy eficiente después del calentamiento.

2. **Planificador del SO para E/S:** Cuando un hilo de Java se bloquea esperando el búfer, el SO lo suspende eficientemente. Con 10,000 hilos del SO en cola, el kernel los gestiona bien para E/S serializada.

3. **Overhead de goroutines con trabajo trivial:** Con solo 100 números por goroutine, el overhead del planificador M:N de Go (crear, planificar y destruir 10,000 goroutines) puede superar el beneficio de su ligereza.

### 10.5 La solución: buffer en Go con bufio.Writer

Para solucionar la falta de búfer en Go y equiparar (o superar) el rendimiento de Java, se aplicaron tres cambios en el código Go:

#### Cambio 1: Nueva importación de `bufio`

```go
import (
    "bufio"    // <-- agregado
    "fmt"
    "os"
    "sync"
    "sync/atomic"
    "time"
)
```

#### Cambio 2: Nueva firma de `contarSegmento` con escritor y mutex

```go
// ANTES (sin búfer)
func contarSegmento(numeroGoroutine int, start, end int64,
    modoDetallado bool, contadorTotal *int64, grupo *sync.WaitGroup)

// DESPUÉS (con búfer)
func contarSegmento(numeroGoroutine int, start, end int64,
    modoDetallado bool, contadorTotal *int64,
    escritor *bufio.Writer, mutexEscritor *sync.Mutex,
    grupo *sync.WaitGroup)
```

El cuerpo en modo detallado ahora acumula en memoria y hace una sola escritura:

```go
// ANTES (una llamada al SO por número → ~1,000,000 llamadas)
for numero := start; numero <= end; numero++ {
    atomic.AddInt64(contadorTotal, 1)
    fmt.Printf("  [Goroutine-%d] --> %d\n", numeroGoroutine, numero)
}

// DESPUÉS (construye en []byte local, una escritura al final)
buf := make([]byte, 0, count*30)
for numero := start; numero <= end; numero++ {
    atomic.AddInt64(contadorTotal, 1)
    buf = fmt.Appendf(buf, "  [Goroutine-%d] --> %d\n", numeroGoroutine, numero)
}
mutexEscritor.Lock()
escritor.Write(buf)   // una sola llamada con el mutex tomado
mutexEscritor.Unlock()
```

#### Cambio 3: En `main` se crea el escritor compartido y se llama a `Flush()` al final

```go
// Escritor con búfer de 1 MB compartido entre todas las goroutines
escritor := bufio.NewWriterSize(os.Stdout, 1024*1024)
var mutexEscritor sync.Mutex

// Al lanzar cada goroutine, se pasa el escritor y el mutex:
go contarSegmento(numeroGoroutine, start, end, modoDetallado,
    &contadorTotal, escritor, &mutexEscritor, &grupo)

// Después de grupo.Wait(), se vacía el búfer al SO:
escritor.Flush()
```

Con estos tres cambios, el número de llamadas al SO se reduce de ~1,000,000 a unas pocas docenas, independientemente de cuántas goroutines haya.

### 10.6 Conclusión del análisis

La diferencia más impactante entre los dos lenguajes se observa en **modo detallado con salida a consola**:

- **Java tardó 101.503 segundos** en imprimir 1,000,000 de números con 2 hilos.
- **Go tardó 14.483 segundos** en la misma tarea (versión original) — aproximadamente **7 veces más rápido** con pocos hilos.
- **Con la corrección de bufio aplicada**, Go mantiene su ventaja incluso con 10,000 goroutines.

Por qué Java podía ganar con 10,000 hilos y la versión original de Go:

1. **Java tiene búfer de E/S por defecto** (8 KB `BufferedOutputStream` en `System.out`).
2. **Go original no tenía búfer:** cada `fmt.Printf` hacía una llamada individual al SO, generando ~1,000,000 de llamadas.
3. **La solución** es usar `bufio.Writer` + búfer por goroutine en Go para igualar o superar a Java.

En **modo resumen** (sin impresión número por número) la diferencia es mínima, porque el conteo se hace con una sola operación atómica por hilo y no hay E/S intensiva. Ahí Java y Go son prácticamente equivalentes.

**Conclusión general:** La ventaja de Go sobre Java no es automática; depende de usar correctamente las herramientas de E/S con búfer. Con `bufio.Writer`, Go recupera su ventaja en todos los escenarios. Además, las mediciones experimentales a escala (hasta 50,000,000 de números) mostraron que Java no es viable para modo detallado a gran escala: en ambas pruebas con ese límite, la JVM falló y cerró el editor, mientras que Go completó la tarea en ambos casos, aunque con tiempos altos (~6.5 minutos) dominados por la E/S más que por el procesamiento paralelo.

---

## 11. ¿Cuándo usar cada uno?

| Situación | Recomendado |
|-----------|-------------|
| Procesamiento con E/S intensiva de consola o archivos | **Go** |
| Proyecto académico o empresarial con ecosistema Java | **Java** |
| Se necesita el menor tiempo de arranque posible | **Go** |
| Muchos trabajadores concurrentes (miles) | **Go** (las goroutines son más ligeras) |
| Pocos hilos con lógica compleja | **Java o Go** (similares) |
| Equipo ya familiarizado con la JVM | **Java** |
| Aplicación que requiere desplegarse como un solo binario sin dependencias | **Go** |
| Necesidad de bibliotecas maduras para procesamiento de datos | **Java** (ecosistema más amplio) |
| Prototipado rápido de herramientas de línea de comandos | **Go** (compilación instantánea) |
| Sistemas críticos donde el rendimiento determinista es importante | **Java** (JIT con perfiles de calentamiento) |
