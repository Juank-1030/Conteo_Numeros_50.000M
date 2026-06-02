<div align="center">

# 🔢 ContadorParalelo

### *Conteo de Números hasta 50.000 Millones con Hilos en Java*

> Programa Java que divide el rango `[1, límite]` en segmentos iguales y asigna cada uno
> a un hilo independiente, garantizando consistencia mediante operaciones atómicas CAS
> y midiendo el tiempo total de ejecución en paralelo.

---

### 🛠️ Tecnologías utilizadas

![Java](https://img.shields.io/badge/Java-21-007396?style=for-the-badge&logo=openjdk&logoColor=white)
![Concurrencia](https://img.shields.io/badge/Concurrencia-Threads-FF6B35?style=for-the-badge)
![AtomicLong](https://img.shields.io/badge/AtomicLong-CAS-blueviolet?style=for-the-badge)
![Streams](https://img.shields.io/badge/Java_Streams-IntStream%20%7C%20LongStream-009688?style=for-the-badge)

</div>

---

## 📑 Tabla de Contenidos

1. [🎯 Descripción del programa](#1--descripción-del-programa)
2. [⚙️ Tecnologías y APIs utilizadas](#2--tecnologías-y-apis-utilizadas)
3. [🧵 Concurrencia — cómo funciona](#3--concurrencia--cómo-funciona)
   - [3.1 División del rango](#31-división-del-rango)
   - [3.2 Creación y arranque de hilos](#32-creación-y-arranque-de-hilos)
   - [3.3 Conteo seguro con AtomicLong](#33-conteo-seguro-con-atomiclong)
   - [3.4 Sincronización de finalización](#34-sincronización-de-finalización)
4. [🔧 Métodos de la clase](#4--métodos-de-la-clase)
5. [📈 Consideraciones de rendimiento](#5--consideraciones-de-rendimiento)
6. [🗂️ Estructura del proyecto](#6--estructura-del-proyecto)
7. [🚀 Compilar y ejecutar](#7--compilar-y-ejecutar)
8. [🖥️ Uso del programa](#8--uso-del-programa)
9. [📋 Ejemplo de salida](#9--ejemplo-de-salida)

---

## 1. 🎯 Descripción del programa

El programa divide el rango `[1, límite]` en segmentos iguales y asigna cada segmento a un hilo independiente. Todos los hilos trabajan simultáneamente, y al finalizar se suma el total contado y se mide el tiempo total de ejecución.

<div align="center">

| ✅ **Qué hace** | 🧩 **Cómo lo logra** |
|:---------------|:---------------------|
| Cuenta números de 1 hasta 50.000 millones | Divide el rango entre N hilos en paralelo |
| Garantiza que no se pierda ni duplique ningún número | Usa `AtomicLong` con operaciones CAS |
| Ofrece dos modos de visualización | Modo resumen (rápido) y modo detallado (número a número) |
| Mide el tiempo de ejecución real | Usa `System.currentTimeMillis()` antes y después de los hilos |

</div>

---

## 2. ⚙️ Tecnologías y APIs utilizadas

<div align="center">

| **Tecnología / Clase** | **Paquete** | **Uso en el programa** |
|------------------------|-------------|------------------------|
| `Thread` | `java.lang` | Creación y gestión de hilos del sistema operativo |
| `AtomicLong` | `java.util.concurrent.atomic` | Contador compartido entre hilos sin necesidad de `synchronized` |
| `IntStream` | `java.util.stream` | Generación funcional del arreglo de hilos con `.mapToObj` y `.peek` |
| `LongStream.rangeClosed` | `java.util.stream` | Iteración número a número en modo detallado |
| `Scanner` | `java.util` | Lectura interactiva de parámetros por consola |
| `System.currentTimeMillis()` | `java.lang` | Medición del tiempo total de ejecución en milisegundos |

</div>

> 💡 No se usan dependencias externas. El programa es **100 % Java estándar (JDK)**.

---

## 3. 🧵 Concurrencia — cómo funciona

### 3.1 División del rango

El rango total `[1, límite]` se divide en `N` segmentos del mismo tamaño (`tamano = limite / hilos`). El último hilo siempre absorbe el remanente para no perder números por redondeo:

```
Hilo 1 → [1,          tamano]
Hilo 2 → [tamano+1,   2*tamano]
...
Hilo N → [...,        limite]   ← absorbe el sobrante
```

### 3.2 Creación y arranque de hilos

Se usa `IntStream.range(0, hilos)` para construir y arrancar todos los hilos en una sola expresión encadenada con `.peek(Thread::start)`. Cada hilo ejecuta el método `contar(...)`:

```java
Thread[] trabajadores = IntStream.range(0, hilos).mapToObj(i -> {
    long desde = i * tamano + 1;
    long hasta  = (i == hilos - 1) ? limite : (i + 1) * tamano;
    return new Thread(() -> contar(i + 1, desde, hasta, detall, c), "Hilo-" + (i + 1));
}).peek(Thread::start).toArray(Thread[]::new);
```

### 3.3 Conteo seguro con `AtomicLong`

En lugar de usar un `long` normal protegido con `synchronized`, se emplea `AtomicLong`, que internamente usa instrucciones **CAS** (*Compare-And-Swap*) del procesador. Esto evita bloqueos y es más eficiente en alta concurrencia:

```java
c.incrementAndGet();   // modo detallado: +1 por cada número
c.addAndGet(n);        // modo resumen:   +N de un golpe
```

<div align="center">

| Estrategia | Mecanismo | Ventaja |
|:-----------|:----------|:--------|
| `synchronized` | Bloqueo de monitor (mutex) | Simple pero genera contención |
| `AtomicLong` (CAS) | Instrucción atómica del CPU | Sin bloqueos, mayor throughput |

</div>

### 3.4 Sincronización de finalización

El hilo principal espera a que todos los trabajadores terminen con `t.join()` antes de imprimir el resultado final. Esto garantiza que el tiempo medido incluye el trabajo completo de todos los hilos.

---

## 4. 🔧 Métodos de la clase

<div align="center">

| **Método** | **Visibilidad** | **Descripción** |
|------------|:---------------:|-----------------|
| `leer(Scanner, String, long, long)` | `static` | Lee un `long` desde consola y valida que esté en `[min, max]`. Termina el programa si está fuera de rango. |
| `contar(int, long, long, boolean, AtomicLong)` | `static` | Lógica de conteo ejecutada por cada hilo. Soporta modo resumen y detallado. |
| `main(String[])` | `public static` | Punto de entrada. Lee parámetros, crea hilos, espera su fin e imprime el resumen. |

</div>

#### `contar(...)` — modos de operación

- **Modo resumen** (`detallado = false`): suma el bloque completo de una vez con `addAndGet`. Imprime una línea por hilo. Muy rápido, ideal para rangos grandes.
- **Modo detallado** (`detallado = true`): itera número a número con `LongStream`, incrementa el contador por cada valor e imprime cada número. Útil para verificar el comportamiento, pero lento para rangos grandes por el I/O de consola.

---

## 5. 📈 Consideraciones de rendimiento

<div align="center">

| Escenario | Comportamiento esperado |
|:----------|:------------------------|
| Modo resumen + pocos hilos | Milisegundos, incluso para 50.000 millones (no hay iteración real) |
| Modo resumen + muchos hilos | Beneficio real hasta `N hilos = N núcleos físicos` |
| Modo detallado | Tiempo crece linealmente con el límite (cuello de botella: I/O de consola) |
| Más hilos que núcleos | No garantiza mayor velocidad; puede haber sobrecarga de context-switching |

</div>

> ⚠️ El beneficio real del paralelismo se obtiene cuando `N hilos ≤ N núcleos físicos` del CPU.

---

## 6. 🗂️ Estructura del proyecto

```
Conteo_Numeros_50.000M/
├── src/
│   └── ContadorParalelo.java
└── README.md
```

---

## 7. 🚀 Compilar y ejecutar

```bash
# Compilar
javac src/ContadorParalelo.java -d out

# Ejecutar
java -cp out ContadorParalelo
```

---

## 8. 🖥️ Uso del programa

Al ejecutar, el programa solicita tres parámetros por consola:

| # | Parámetro | Descripción |
|:-:|-----------|-------------|
| 1 | **Número final** | Límite de conteo (entre `1` y `50.000.000.000`) |
| 2 | **Hilos** | Cantidad de hilos a usar en paralelo |
| 3 | **Modo** | `1` = resumen por hilo &nbsp;/&nbsp; `2` = número a número (detallado) |

---

## 9. 📋 Ejemplo de salida

#### Modo resumen (4 hilos, límite = 1000)

```
Número final (1-50.000.000.000): 1000
Hilos: 4
Modo (1=resumen  2=número a número): 1

Contando 1→1000 con 4 hilo(s)...

  Hilo-1: 1 → 250  (250 nums)
  Hilo-2: 251 → 500  (250 nums)
  Hilo-3: 501 → 750  (250 nums)
  Hilo-4: 751 → 1000  (250 nums)

--- Resultado ---
Total contado  : 1000
Tiempo de mora : 3 ms (0.003 s)
```

#### Modo detallado (2 hilos, límite = 6)

```
Número final (1-50.000.000.000): 6
Hilos: 2
Modo (1=resumen  2=número a número): 2

Contando 1→6 con 2 hilo(s)...

  [Hilo-1] = 1
  [Hilo-1] = 2
  [Hilo-1] = 3
  [Hilo-2] = 4
  [Hilo-2] = 5
  [Hilo-2] = 6

--- Resultado ---
Total contado  : 6
Tiempo de mora : 5 ms (0.005 s)
```