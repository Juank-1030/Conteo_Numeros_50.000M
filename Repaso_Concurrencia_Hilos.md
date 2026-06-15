# 📋 INFORME TÉCNICO: Concurrencia e Hilos en Java y Go

**Autores:** Luis Daniel Benavides Navarro | Rodrigo Humberto Gualtero Martínez  
**Fecha:** 20/06/2024  
**Institución:** Escuela Colombiana de Ingeniería Julio Garavito

---

## 📑 TABLA DE CONTENIDO

1. [Concurrencia](#1-concurrencia)
2. [Estrategias de Concurrencia](#2-estrategias-de-concurrencia)
   - [Paso de Mensajes](#21-paso-de-mensajes)
   - [Memoria Compartida](#22-memoria-compartida)
3. [Java Concurrente](#3-java-concurrente)
   - [Creación de Hilos](#31-creación-de-hilos)
   - [Join](#32-join---esperando-un-hilo)
   - [Métodos Sincronizados](#33-métodos-sincronizados)
   - [Barreras CyclicBarrier](#34-barreras-cyclicbarrier)
   - [Seguros (Locks)](#35-seguros-y-sincronización)
   - [Liveness y Deadlock](#36-liveness-y-deadlock)
   - [Bloques con Guarda](#37-bloques-con-guarda)
   - [wait y notifyAll](#38-wait-y-notifyall)
   - [Objetos Inmutables](#39-objetos-inmutables)
   - [Concurrencia de Alto Nivel](#310-concurrencia-de-alto-nivel)
   - [Pool de Hilos](#311-pool-de-hilos)
   - [Callable y Future](#312-callable-y-future)
4. [Go (Golang)](#4-go-golang)
   - [Goroutines](#41-goroutines)
   - [Canales (Channels)](#42-canales-channels)
   - [Select](#43-select)
5. [Glosario](#-glosario-de-términos)

---

## 1. Concurrencia

### 📌 Definición
La **concurrencia** es la capacidad de un sistema para manejar **múltiples operaciones o transacciones de manera simultánea**. No necesariamente significa que ocurran al mismo tiempo físicamente (eso sería paralelismo), sino que el sistema puede gestionar varias tareas en progreso al mismo tiempo, alternando entre ellas.

### 🔍 Diferencia entre Concurrencia y Paralelismo

| Concepto | Descripción |
|----------|-------------|
| **Concurrencia** | Múltiples tareas progresan simultáneamente (puede ser en un solo core) |
| **Paralelismo** | Múltiples tareas se ejecutan físicamente al mismo tiempo (multi-core) |

---

## 2. Estrategias de Concurrencia

El documento presenta **dos estrategias principales** para manejar la concurrencia:

---

### 2.1 Paso de Mensajes

#### 📌 Definición
Los hilos o procesos se **comunican enviando mensajes** entre sí. Cada proceso tiene su **propio espacio de memoria aislado**, por lo que no comparten datos directamente.

#### ✅ Ventajas
- **Elimina condiciones de carrera**: Al no compartir memoria, no hay interferencia.
- **Escalabilidad**: Escala fácilmente en sistemas distribuidos y redes de computadoras.
- **Aislamiento**: Cada proceso tiene su propio espacio de memoria.

#### ❌ Desventajas
- **Sobrecarga de comunicación**: Es necesario serializar y deserializar mensajes, lo que tiene un costo.
- **Latencia**: La comunicación entre procesos puede ser más lenta comparada con memoria compartida.

---

### 2.2 Memoria Compartida

#### 📌 Definición
Múltiples hilos o procesos **acceden y modifican datos** almacenados en un **espacio de memoria común**. Requiere mecanismos de sincronización para coordinar el acceso concurrente.

#### ✅ Ventajas
- **No hay sobrecarga de comunicación**: Generalmente más rápido.
- **Ideal** para sistemas con hilos ligeros donde la comunicación es frecuente y de baja latencia.

#### ❌ Desventajas
- **Condiciones de carrera**: Los hilos pueden interferir entre sí, causando inconsistencias en los datos.
- **Complejidad en la sincronización**: Uso de mecanismos de sincronización puede generar **deadlocks**.

---

## 3. Java Concurrente

Java utiliza una estrategia de **memoria compartida a través de Threads**. Provee los siguientes mecanismos de sincronización:
- **Monitores**
- **Semáforos**
- **CyclicBarrier** (Barrera cíclica)

---

### 3.1 Creación de Hilos

#### 📌 Definición
En Java existen **dos formas principales** de crear un hilo:

1. **Extendiendo la clase `Thread`**
2. **Implementando la interfaz `Runnable`**

#### 💻 Ejemplo Detallado — Método 1: Extendiendo `Thread`

```java
// Se crea una clase que EXTIENDE Thread
// Esto convierte nuestra clase directamente en un hilo ejecutable
public class HelloThread extends Thread {

    // El método run() es el CUERPO del hilo
    // Todo lo que esté aquí se ejecutará en un hilo separado
    @Override
    public void run() {
        System.out.println("Hello from a thread!");
        // Este mensaje lo imprimirá el hilo nuevo, NO el hilo principal (main)
    }

    public static void main(String[] args) {
        // Se crea una instancia de HelloThread
        // y se llama .start() para INICIAR el hilo
        // IMPORTANTE: No se llama run() directamente, sino start()
        // .start() → crea un nuevo hilo y luego llama run() en ese nuevo hilo
        // .run()   → simplemente ejecutaría el método en el hilo actual (no crea hilo nuevo)
        (new HelloThread()).start();
    }
}
```

#### 💻 Ejemplo Detallado — Método 2: Implementando `Runnable`

```java
// Se crea una clase que IMPLEMENTA la interfaz Runnable
// Esta es la forma PREFERIDA porque:
// 1. No "gasta" la herencia (Java no permite herencia múltiple)
// 2. Separa la lógica de la tarea del manejo del hilo
public class HelloRunnable implements Runnable {

    // Implementamos el método run() obligatorio de la interfaz Runnable
    @Override
    public void run() {
        System.out.println("Hello from a thread!");
    }

    public static void main(String[] args) {
        // Para ejecutarlo, se crea un objeto Thread
        // pasándole el Runnable como argumento
        // Thread → es el "vehículo" que ejecuta la tarea (Runnable)
        (new Thread(new HelloRunnable())).start();
    }
}
```

#### 💻 Ejemplo con Múltiples Hilos (Problema de Orden)

```java
// Este ejemplo muestra el PROBLEMA de la concurrencia:
// el orden de ejecución es NO DETERMINÍSTICO
public class Multi extends Thread {

    public void run() {
        // Este mensaje puede aparecer en cualquier orden
        System.out.println("Thread is running...");
    }

    public static void main(String[] args) {
        Multi t1 = new Multi(); // Creamos primer hilo
        Multi t2 = new Multi(); // Creamos segundo hilo

        t1.start(); // Iniciamos t1 → el SO decide cuándo ejecutarlo
        t2.start(); // Iniciamos t2 → el SO decide cuándo ejecutarlo

        // La salida puede ser:
        // "Thread is running..." (t1)
        // "Thread is running..." (t2)
        // O también al revés, t2 antes que t1
        System.out.println("end main execution");
        // Incluso este mensaje puede aparecer ANTES que los de t1 o t2
    }
}
```

---

### 3.2 Join — Esperando un Hilo

#### 📌 Definición
El método **`join()`** permite que un hilo **espere** a que otro termine antes de continuar su ejecución. Es útil cuando necesitamos el resultado de un hilo antes de continuar.

#### 💻 Ejemplo Detallado con `join()`

```java
import java.time.LocalDateTime;
import java.util.logging.Level;
import java.util.logging.Logger;

public class HelloRunnableWithJoin implements Runnable {

    @Override
    public void run() {
        try {
            // Imprimimos la hora de inicio del hilo
            System.out.println("Hello from a thread! time: " + LocalDateTime.now());

            // Simulamos trabajo pesado: dormimos 4 segundos
            // Thread.sleep() pausa el hilo actual sin liberar locks
            Thread.sleep(4000); // 4000 ms = 4 segundos

        } catch (InterruptedException ex) {
            // Si el hilo es interrumpido mientras duerme, capturamos la excepción
            Logger.getLogger(HelloRunnableWithJoin.class.getName())
                  .log(Level.SEVERE, null, ex);
        }
    }

    public static void main(String[] args) {
        // Creamos el hilo con nuestro Runnable
        Thread t = new Thread(new HelloRunnableWithJoin());

        // Iniciamos el hilo → comienza a ejecutar run() en paralelo
        t.start();

        try {
            // JOIN: el hilo principal se BLOQUEA aquí
            // y espera hasta que el hilo 't' termine completamente
            t.join();
            // Sin join(), el main continuaría sin esperar a 't'

        } catch (InterruptedException ex) {
            Logger.getLogger(HelloRunnableWithJoin.class.getName())
                  .log(Level.SEVERE, null, ex);
        }

        // Este mensaje SOLO se imprime DESPUÉS de que el hilo 't' haya terminado
        System.out.println("Hello from main thread! time: " + LocalDateTime.now());
        // La diferencia de tiempo entre ambos mensajes será ~4 segundos
    }
}
```

---

### 3.3 Métodos Sincronizados

#### 📌 Definición
La **sincronización** garantiza que solo **un hilo a la vez** pueda ejecutar métodos sincronizados de un objeto. Usa la palabra clave `synchronized`.

#### 💻 Ejemplo Detallado — Deadlock con Métodos Sincronizados

```java
// Este ejemplo ilustra el problema del DEADLOCK con métodos sincronizados
public class Friend {
    private final String name;

    public Friend(String name) {
        this.name = name;
    }

    public String getName() {
        return this.name;
    }

    // Método SINCRONIZADO: solo un hilo puede ejecutarlo a la vez
    // Cuando un hilo entra aquí, adquiere el "lock" (seguro) del objeto
    public synchronized void bow(Friend bower) {
        System.out.format("%s: %s has bowed to me!%n",
                this.name, bower.getName());

        // PROBLEMA: Mientras tenemos el lock de THIS,
        // intentamos llamar bowBack() en BOWER
        // Si bower también tiene su lock, hay DEADLOCK
        bower.bowBack(this);
    }

    // Otro método SINCRONIZADO del mismo objeto
    public synchronized void bowBack(Friend bower) {
        System.out.format("%s: %s has bowed back to me!%n",
                this.name, bower.getName());
    }
}
```

---

### 3.4 Barreras CyclicBarrier

#### 📌 Definición
`CyclicBarrier` es una barrera de sincronización que permite que **un conjunto de hilos se esperen mutuamente** hasta que todos lleguen a un punto común (la barrera). Después todos continúan juntos.

#### 💻 Ejemplo Detallado con `CyclicBarrier`

```java
import java.util.concurrent.BrokenBarrierException;
import java.util.concurrent.CyclicBarrier;

// Clase que representa un "trabajador" (hilo)
class WorkerThread extends Thread {
    private CyclicBarrier barrier; // Referencia a la barrera compartida
    private int id;                // Identificador del trabajador

    public WorkerThread(CyclicBarrier barrier, int id) {
        this.barrier = barrier;
        this.id = id;
    }

    @Override
    public void run() {
        try {
            // Simulamos 3 fases de trabajo
            for (int i = 1; i <= 3; i++) {

                System.out.println("Worker " + id + " está trabajando en la fase " + i);

                // Simulamos tiempo de trabajo variable según el id del worker
                // Worker 1 tarda 1s, Worker 2 tarda 2s, Worker 3 tarda 3s
                Thread.sleep(1000 * id);

                System.out.println("Worker " + id + " ha terminado la fase " + i);

                // BARRERA: El hilo se DETIENE aquí hasta que TODOS los hilos
                // hayan llamado barrier.await()
                // Ningún worker avanza a la siguiente fase hasta que todos terminen la actual
                barrier.await();

                System.out.println("Worker " + id + " ha cruzado la barrera de la fase " + i);
            }
        } catch (InterruptedException | BrokenBarrierException e) {
            e.printStackTrace();
        }
    }
}

public class CyclicBarrierExample {
    public static void main(String[] args) {
        // Número de hilos que deben llegar a la barrera antes de que se abra
        final int NUM_WORKERS = 3;

        // Creamos la barrera con:
        // - NUM_WORKERS: cantidad de hilos que deben llegar
        // - Runnable: acción que se ejecuta cuando TODOS llegan a la barrera
        CyclicBarrier barrier = new CyclicBarrier(NUM_WORKERS, new Runnable() {
            @Override
            public void run() {
                // Este código se ejecuta UNA VEZ cuando todos los hilos llegaron
                System.out.println("¡Todos los trabajadores han alcanzado la barrera!");
            }
        });

        // Creamos e iniciamos los 3 workers
        for (int i = 1; i <= NUM_WORKERS; i++) {
            new WorkerThread(barrier, i).start();
        }

        // Flujo esperado:
        // Fase 1: Los 3 workers trabajan → todos llegan a la barrera → continúan
        // Fase 2: Los 3 workers trabajan → todos llegan a la barrera → continúan
        // Fase 3: Los 3 workers trabajan → todos llegan a la barrera → continúan
    }
}
```

---

### 3.5 Seguros y Sincronización

#### 📌 Definición
Los **seguros intrínsecos (intrinsic locks o monitors)** son el mecanismo base de sincronización en Java. Cada objeto tiene un seguro. Un hilo debe **adquirirlo** para ejecutar métodos sincronizados. Si no puede adquirirlo, **espera**.

#### 💻 Ejemplo Detallado — Bloques Sincronizados con Locks Separados

```java
// Este ejemplo muestra cómo MEJORAR la concurrencia
// usando locks separados para datos independientes
public class SynchronizedStatement {

    private long c1 = 0; // Contador 1
    private long c2 = 0; // Contador 2

    // Dos objetos que actúan como "llaves" (locks) independientes
    // Esto permite que c1 y c2 sean modificados CONCURRENTEMENTE
    // porque cada uno tiene su propio lock
    private Object lock1 = new Object();
    private Object lock2 = new Object();

    // Método que incrementa c1
    // Solo necesita el lock1, no bloquea acceso a c2
    public void inc1() {
        // Bloque sincronizado: adquiere lock1
        synchronized (lock1) {
            c1++;
            // Al salir del bloque, libera lock1
        }
    }

    // Método que incrementa c2
    // Solo necesita el lock2, no bloquea acceso a c1
    public void inc2() {
        // Bloque sincronizado: adquiere lock2
        synchronized (lock2) {
            c2++;
            // Al salir del bloque, libera lock2
        }
    }

    // COMPARACIÓN:
    // Si usáramos "synchronized" en el método completo (synchronized void inc1())
    // solo UN hilo podría ejecutar cualquiera de los dos métodos a la vez
    // Con locks separados, inc1() e inc2() pueden ejecutarse SIMULTÁNEAMENTE
}
```

---

### 3.6 Liveness y Deadlock

#### 📌 Definiciones
- **Liveness**: La habilidad de una aplicación de ejecutarse de manera temporalmente correcta, sin quedarse bloqueada indefinidamente.
- **Deadlock**: Situación donde dos o más hilos se bloquean mutuamente esperando recursos que el otro tiene.

#### 💻 Ejemplo Detallado de Deadlock

```java
// EJEMPLO CLÁSICO DE DEADLOCK
// Dos amigos que se reverencian mutuamente → ninguno puede terminar
public class DeadlockExample {

    static class Friend {
        private final String name;

        public Friend(String name) { this.name = name; }
        public String getName() { return this.name; }

        // Hilo A: adquiere lock de "alphonse", luego intenta lock de "gaston"
        // Hilo B: adquiere lock de "gaston", luego intenta lock de "alphonse"
        // → DEADLOCK: cada hilo espera el lock del otro
        public synchronized void bow(Friend bower) {
            System.out.format("%s: %s has bowed to me!%n",
                    this.name, bower.getName());
            // Aquí intentamos acceder al método sincronizado de OTRO objeto
            // Si ese objeto ya tiene su lock tomado por otro hilo → DEADLOCK
            bower.bowBack(this);
        }

        public synchronized void bowBack(Friend bower) {
            System.out.format("%s: %s has bowed back to me!%n",
                    this.name, bower.getName());
        }
    }

    public static void main(String[] args) {
        final Friend alphonse = new Friend("Alphonse");
        final Friend gaston = new Friend("Gaston");

        // Hilo 1: alphonse.bow(gaston)
        // → adquiere lock de alphonse
        // → intenta adquirir lock de gaston (bloqueado por Hilo 2)
        new Thread(new Runnable() {
            public void run() { alphonse.bow(gaston); }
        }).start();

        // Hilo 2: gaston.bow(alphonse)
        // → adquiere lock de gaston
        // → intenta adquirir lock de alphonse (bloqueado por Hilo 1)
        new Thread(new Runnable() {
            public void run() { gaston.bow(alphonse); }
        }).start();

        // RESULTADO: Ambos hilos esperan indefinidamente → DEADLOCK
    }
}
```

---

### 3.7 Bloques con Guarda

#### 📌 Definición
Un **bloque con guarda (guarded block)** es un patrón donde un hilo **espera que una condición se cumpla** antes de continuar. La guarda ineficiente hace un loop continuo (busy-waiting), desperdiciando CPU.

#### 💻 Ejemplo Detallado — Guarda Ineficiente vs. Eficiente

```java
// ❌ FORMA INEFICIENTE: Busy-waiting (polling loop)
// El hilo consume CPU continuamente revisando la condición
public void guardedJoy_MALO() {
    // Este loop revisa la condición millones de veces por segundo
    // DESPERDICIA tiempo de procesador
    while (!joy) {
        // No hace nada útil, solo sigue chequeando
    }
    System.out.println("Joy has been achieved!");
}

// ✅ FORMA EFICIENTE: Usando wait() y notifyAll()
// El hilo SE DUERME y es despertado solo cuando la condición puede haber cambiado
public synchronized void guardedJoy_BUENO() {
    // El while (no if) es importante: hay que reverificar la condición
    // porque pueden ocurrir "spurious wakeups" (despertares falsos)
    while (!joy) {
        try {
            // wait() hace tres cosas importantes:
            // 1. Libera el lock del objeto (permite que otros hilos accedan)
            // 2. Suspende la ejecución del hilo actual
            // 3. Cuando recibe notificación, vuelve a adquirir el lock y continúa
            wait();
        } catch (InterruptedException e) {
            // Manejar la interrupción
        }
    }
    System.out.println("Joy and efficiency have been achieved!");
}
```

---

### 3.8 wait() y notifyAll()

#### 📌 Definición
- **`Object.wait()`**: Suspende la ejecución del hilo y **libera el lock**, esperando una notificación.
- **`Object.notifyAll()`**: Despierta **a todos los hilos** que están esperando en ese objeto.
- **`Object.notify()`**: Despierta **a un hilo aleatorio** que está esperando.

#### 💻 Ejemplo Detallado — Patrón Productor-Consumidor con wait/notifyAll

```java
// Clase que implementa una cola sincronizada usando wait/notifyAll
public class SyncQueue {
    private boolean joy = false; // La condición de guarda

    // CONSUMIDOR: espera hasta que joy sea true
    public synchronized void guardedJoy() {
        // SIEMPRE usar while, nunca if
        // Razón: pueden ocurrir "spurious wakeups"
        while (!joy) {
            try {
                // 1. Libera el lock de este objeto
                // 2. Suspende este hilo
                // 3. Cuando sea notificado, adquiere el lock y re-chequea la condición
                wait();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
        System.out.println("Joy and efficiency have been achieved!");
    }

    // PRODUCTOR: cambia la condición y notifica a todos los que esperan
    public synchronized void notifyJoy() {
        // Cambiamos la condición de guarda
        joy = true;

        // Notificamos a TODOS los hilos que están en wait()
        // Ellos se despertarán, re-adquirirán el lock, y evaluarán la condición
        notifyAll();

        // Diferencia entre notify() y notifyAll():
        // notify()    → despierta a UN hilo aleatorio (puede causar inanición)
        // notifyAll() → despierta a TODOS los hilos (más seguro, preferido)
    }
}
```

---

### 3.9 Objetos Inmutables

#### 📌 Definición
Un **objeto inmutable** es aquel cuyo **estado no puede cambiar** después de ser construido. Son inherentemente seguros en entornos concurrentes porque no necesitan sincronización.

#### 💻 Ejemplo Detallado — Clase Inmutable

```java
// La clase es FINAL para evitar que subclases la modifiquen
public final class ImmutablePoint {

    // Todos los campos son FINAL y PRIVATE
    // final → no pueden reasignarse después del constructor
    // private → no se exponen directamente
    private final int x;
    private final int y;
    private final String label;

    // Constructor: único momento donde se asignan los valores
    public ImmutablePoint(int x, int y, String label) {
        this.x = x;
        this.y = y;
        // Si el campo referencia a un objeto mutable, hacemos copia defensiva
        this.label = label; // String ya es inmutable en Java, OK directo
    }

    // Solo GETTERS, nunca setters
    public int getX() { return x; }
    public int getY() { return y; }
    public String getLabel() { return label; }

    // Para "modificar", devolvemos una NUEVA instancia
    // El objeto original no cambia jamás
    public ImmutablePoint moverA(int nuevoX, int nuevoY) {
        return new ImmutablePoint(nuevoX, nuevoY, this.label);
    }

    // Múltiples hilos pueden leer este objeto sin sincronización
    // porque nunca cambia → Thread-safe por naturaleza
}

// EJEMPLO DE USO
class Main {
    public static void main(String[] args) {
        ImmutablePoint p1 = new ImmutablePoint(0, 0, "Origen");

        // Múltiples hilos pueden usar p1 de forma segura
        Runnable tarea = () -> {
            // No hay riesgo de condición de carrera
            // porque p1 nunca cambia
            System.out.println("X: " + p1.getX() + ", Y: " + p1.getY());
        };

        new Thread(tarea).start();
        new Thread(tarea).start();
    }
}
```

---

### 3.10 Concurrencia de Alto Nivel

#### 📌 Definición
Java provee un **API de alto nivel** (`java.util.concurrent`) con herramientas más sofisticadas que los hilos básicos:

| Herramienta | Descripción |
|-------------|-------------|
| **Executors** | Gestores de pools de hilos |
| **Objetos Lock** | Locks más avanzados que `synchronized` |
| **Colecciones Concurrentes** | Estructuras thread-safe |
| **Variables Atómicas** | Variables que solo un hilo modifica a la vez |

#### 💻 Ejemplo — Interfaz Executor

```java
import java.util.concurrent.Executor;

public class SimpleExecutorExample {
    public static void main(String[] args) {

        // Implementamos la interfaz Executor de forma anónima
        // Executor es la interfaz más básica: solo tiene el método execute()
        Executor executor = new Executor() {
            @Override
            public void execute(Runnable command) {
                // Cada vez que llamamos execute(), creamos un hilo nuevo
                // En producción usaríamos un pool, no new Thread() cada vez
                new Thread(command).start();
            }
        };

        // Enviamos tareas al executor
        // Las lambdas (() -> ...) son implementaciones de Runnable
        executor.execute(() ->
            System.out.println("Hola desde otro hilo - "
                + Thread.currentThread().getName())
        );

        executor.execute(() ->
            System.out.println("Otra tarea ejecutándose - "
                + Thread.currentThread().getName())
        );

        // Resultado posible:
        // "Hola desde otro hilo - Thread-0"
        // "Otra tarea ejecutándose - Thread-1"
        // (El orden puede variar)
    }
}
```

---

### 3.11 Pool de Hilos

#### 📌 Definición
Un **pool de hilos (Thread Pool)** es un patrón que administra un **conjunto limitado y reutilizable de hilos** para evitar el costo de crear y destruir hilos constantemente.

#### Tipos de Pools disponibles en Java

| Método | Descripción |
|--------|-------------|
| `newFixedThreadPool(n)` | Pool de tamaño fijo con `n` hilos |
| `newCachedThreadPool()` | Pool dinámico para tareas cortas y paralelas |
| `newSingleThreadExecutor()` | Un solo hilo a la vez, serializado |

#### 💻 Ejemplo Detallado — Pool de Hilos con `newFixedThreadPool`

```java
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class ThreadPoolExample {
    public static void main(String[] args) {

        // Creamos un pool de exactamente 3 hilos
        // Solo habrá 3 hilos activos simultáneamente
        // Las tareas restantes esperan en una cola interna
        ExecutorService executor = Executors.newFixedThreadPool(3);

        // Enviamos 6 tareas al pool (más tareas que hilos disponibles)
        for (int i = 1; i <= 6; i++) {
            final int taskId = i; // Variable efectivamente final para la lambda

            executor.execute(() -> {
                String threadName = Thread.currentThread().getName();

                // Las tareas 1,2,3 se ejecutan primero (3 hilos disponibles)
                // Las tareas 4,5,6 esperan a que algún hilo quede libre
                System.out.println("Ejecutando tarea " + taskId + " en " + threadName);

                try {
                    // Simulamos 2 segundos de trabajo
                    TimeUnit.SECONDS.sleep(2);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }

                System.out.println("Tarea " + taskId + " finalizada en " + threadName);
            });
        }

        // shutdown(): ya no acepta más tareas, pero termina las existentes
        // Es importante llamarlo para liberar recursos
        executor.shutdown();

        try {
            // Esperamos hasta 10 segundos para que terminen todas las tareas
            if (!executor.awaitTermination(10, TimeUnit.SECONDS)) {
                // Si en 10s no terminaron, las cancelamos forzosamente
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
        }

        // Flujo esperado:
        // t=0s: Tareas 1,2,3 inician en pool-1-thread-1, 2, 3
        // t=2s: Tareas 1,2,3 terminan → Tareas 4,5,6 inician
        // t=4s: Tareas 4,5,6 terminan → programa termina
    }
}
```

---

### 3.12 Callable y Future

#### 📌 Definición
- **`Callable<V>`**: Similar a `Runnable`, pero **puede retornar un valor** y lanzar excepciones.
- **`Future<V>`**: Representa el **resultado futuro** de una tarea asíncrona. Permite obtener el resultado cuando esté disponible.

#### 💻 Ejemplo Detallado — Callable + Future

```java
import java.util.concurrent.*;

public class CallableFutureExample {
    public static void main(String[] args) {

        // Pool de un solo hilo para este ejemplo
        ExecutorService executor = Executors.newSingleThreadExecutor();

        // Callable<String>: tarea que retorna un String
        // A diferencia de Runnable.run(), Callable.call() PUEDE retornar valor
        Callable<String> tarea = () -> {
            System.out.println("Tarea iniciada en: " + Thread.currentThread().getName());
            TimeUnit.SECONDS.sleep(2); // Simulamos trabajo de 2 segundos
            return "Resultado de la tarea"; // Valor que retornará
        };

        // submit() envía la tarea y retorna un Future
        // La tarea se ejecuta en SEGUNDO PLANO
        // Future es como un "ticket" que nos permite recoger el resultado después
        Future<String> future = executor.submit(tarea);

        // IMPORTANTE: El main continúa ejecutándose mientras la tarea trabaja en paralelo
        System.out.println("Tarea enviada. Haciendo otras cosas en main...");

        try {
            // future.get() BLOQUEA el hilo actual hasta que la tarea termine
            // y retorna el resultado
            // Si la tarea lanzó una excepción, aquí se lanza ExecutionException
            String resultado = future.get();
            System.out.println("El resultado es: " + resultado);

            // También podemos usar: future.get(timeout, TimeUnit)
            // para no esperar indefinidamente

        } catch (InterruptedException e) {
            // El hilo fue interrumpido mientras esperaba
            Thread.currentThread().interrupt();
        } catch (ExecutionException e) {
            // La tarea lanzó una excepción durante su ejecución
            e.printStackTrace();
        } finally {
            // Siempre cerrar el executor
            executor.shutdown();
        }
    }
}
```

---

## 4. Go (Golang)

### 📌 Definición
**Go** es un lenguaje compilado, concurrente y de código abierto, diseñado por Google en 2007. Es eficiente, simple y rápido. Usa el modelo de **paso de mensajes** a través de **canales (channels)** y **goroutines** (hilos ultraligeros).

> **Filosofía de Go**: *"No comuniques compartiendo memoria; comparte memoria comunicando"*

---

### 4.1 Goroutines

#### 📌 Definición
Una **goroutine** es una función que se ejecuta **concurrentemente** con otras goroutines. Son extremadamente ligeras (comienzan con ~2KB de stack) comparadas con los threads de Java.

#### 💻 Ejemplo Detallado — Goroutines Básicas

```go
package main

import (
    "fmt"
    "time"
)

// Función normal que imprime un texto
func say(text string) {
    fmt.Println(text)
}

func main() {
    // La palabra clave "go" lanza say() como una GOROUTINE
    // El programa NO espera a que termine → continúa inmediatamente
    go say("Hola")   // Goroutine 1
    go say("desde")  // Goroutine 2
    go say("Go!")    // Goroutine 3

    // PROBLEMA: Si main() termina antes que las goroutines,
    // las goroutines se eliminan sin ejecutarse
    // Solución temporal: dormir el programa
    time.Sleep(1 * time.Second)

    // NOTA: El orden de salida es NO DETERMINÍSTICO
    // Puede imprimir: "Go!", "Hola", "desde"
    // O cualquier otra combinación
}
```

---

### 4.2 Canales (Channels)

#### 📌 Definición
Los **canales** son "tuberías" por donde las goroutines **envían y reciben valores**. Permiten sincronización y comunicación segura entre goroutines. Se crean con `make(chan tipo)`.

#### Operaciones básicas:
```
ch <- valor   // ENVIAR valor al canal (bloqueante hasta que alguien reciba)
valor := <-ch // RECIBIR valor del canal (bloqueante hasta que alguien envíe)
```

#### 💻 Ejemplo Detallado — Canal Básico (1 Productor, 1 Consumidor)

```go
package main

import "fmt"

// worker recibe mensajes del canal ch
func worker(ch chan string) {
    // <- ch BLOQUEA la goroutine hasta que haya un mensaje disponible
    msg := <-ch
    fmt.Println("Recibido:", msg)
}

func main() {
    // Creamos un canal que transporta strings
    // make(chan string) → canal sin buffer (sincrónico)
    ch := make(chan string)

    // Lanzamos worker como goroutine
    // worker está esperando recibir algo del canal
    go worker(ch)

    // Enviamos un mensaje al canal
    // ch <- "..." BLOQUEA main hasta que worker reciba el mensaje
    // Esto SINCRONIZA main y la goroutine worker
    ch <- "mensaje desde main"

    // Flujo:
    // 1. main lanza worker (goroutine)
    // 2. worker espera en: msg := <-ch
    // 3. main envía: ch <- "mensaje desde main"
    // 4. worker recibe y continúa
    // 5. main continúa (o termina)
}
```

#### 💻 Ejemplo Detallado — Canal con Múltiples Productores

```go
package main

import "fmt"

// Cada productor envía UN mensaje al canal compartido
func productor(id int, ch chan string) {
    // Enviamos un mensaje formateado con el ID del productor
    // Esta goroutine se bloquea hasta que alguien reciba el mensaje
    ch <- fmt.Sprintf("mensaje del productor %d", id)
}

func main() {
    // Canal compartido por todos los productores
    ch := make(chan string)

    // Lanzamos 3 goroutines productoras
    for i := 1; i <= 3; i++ {
        go productor(i, ch) // Cada una envía independientemente
    }

    // Recibimos exactamente 3 mensajes del canal
    for i := 1; i <= 3; i++ {
        // <-ch bloquea hasta que haya un mensaje disponible
        fmt.Println(<-ch)
    }

    // El ORDEN de los mensajes puede variar en cada ejecución
    // Posible salida 1: "mensaje del productor 2", "1", "3"
    // Posible salida 2: "mensaje del productor 1", "3", "2"
}
```

---

### 4.3 Select

#### 📌 Definición
`select` es como un `switch` pero para **canales**. Espera en múltiples canales y procesa el primero que tenga datos disponibles. Si múltiples canales están listos, elige uno **aleatoriamente**.

#### 💻 Ejemplo Detallado — Select con Múltiples Canales

```go
package main

import "fmt"

func main() {
    // Creamos DOS canales independientes
    ch1 := make(chan string)
    ch2 := make(chan string)

    // Goroutine 1: envía mensaje a ch1
    // go func() { ... }() → función anónima lanzada como goroutine
    go func() {
        ch1 <- "Hola desde ch1"
    }()

    // Goroutine 2: envía mensaje a ch2
    go func() {
        ch2 <- "Hola desde ch2"
    }()

    // select: espera al primer canal que tenga datos
    // Si ch1 está listo antes → ejecuta case msg1
    // Si ch2 está listo antes → ejecuta case msg2
    // Si ambos están listos → elige UNO aleatoriamente
    select {
    case msg1 := <-ch1:
        fmt.Println(msg1) // "Hola desde ch1"
    case msg2 := <-ch2:
        fmt.Println(msg2) // "Hola desde ch2"
    }

    // NOTA: select SIN default es bloqueante (espera hasta que un canal tenga datos)
    // select CON default (case default:) es no bloqueante

    // Ejemplo con default (no bloqueante):
    select {
    case msg := <-ch1:
        fmt.Println("Recibido de ch1:", msg)
    default:
        fmt.Println("Ningún canal listo aún") // Se ejecuta si no hay datos
    }
}
```

---

## 📊 Comparación Java vs Go

| Característica | Java | Go |
|----------------|------|----|
| **Unidad de concurrencia** | Thread | Goroutine |
| **Estrategia principal** | Memoria compartida | Paso de mensajes (canales) |
| **Sincronización** | `synchronized`, `Lock`, `Semaphore` | Canales (`chan`) |
| **Peso** | ~1MB por hilo | ~2KB por goroutine |
| **Complejidad** | Mayor (deadlocks, condiciones de carrera) | Menor (canales evitan muchos problemas) |
| **Pool de hilos** | `ExecutorService` | Gestionado automáticamente por el runtime |
| **Resultado futuro** | `Future<T>` | Canales con valor de retorno |

---

## 📚 GLOSARIO DE TÉRMINOS

| Término | Definición | Enlace de Referencia |
|---------|------------|----------------------|
| **Concurrencia** | Capacidad de manejar múltiples operaciones simultáneamente | [Oracle Java Concurrency](https://docs.oracle.com/javase/tutorial/essential/concurrency/) |
| **Thread (Hilo)** | Unidad básica de ejecución dentro de un proceso | [Java Thread Docs](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/lang/Thread.html) |
| **Runnable** | Interfaz funcional que define la tarea a ejecutar en un hilo | [Java Runnable](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/lang/Runnable.html) |
| **Memoria Compartida** | Estrategia donde múltiples hilos acceden al mismo espacio de memoria | [Wikipedia - Shared Memory](https://en.wikipedia.org/wiki/Shared_memory) |
| **Paso de Mensajes** | Estrategia donde procesos se comunican enviando mensajes | [Wikipedia - Message Passing](https://en.wikipedia.org/wiki/Message_passing) |
| **Deadlock** | Bloqueo mutuo donde dos o más hilos esperan indefinidamente | [Oracle - Deadlock](https://docs.oracle.com/javase/tutorial/essential/concurrency/deadlock.html) |
| **Race Condition** | Condición de carrera: comportamiento no determinístico por acceso concurrente | [Wikipedia - Race Condition](https://en.wikipedia.org/wiki/Race_condition) |
| **Liveness** | Propiedad que garantiza que un programa eventualmente progresa | [Oracle - Liveness](https://docs.oracle.com/javase/tutorial/essential/concurrency/liveness.html) |
| **synchronized** | Palabra clave de Java para sincronizar acceso a métodos o bloques | [Java synchronized](https://docs.oracle.com/javase/tutorial/essential/concurrency/syncmeth.html) |
| **CyclicBarrier** | Mecanismo que sincroniza múltiples hilos en un punto común | [Java CyclicBarrier](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/CyclicBarrier.html) |
| **Semaphore** | Contador que controla el acceso concurrente a recursos | [Java Semaphore](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/Semaphore.html) |
| **Monitor** | Mecanismo de sincronización que combina mutex y variables de condición | [Wikipedia - Monitor](https://en.wikipedia.org/wiki/Monitor_(synchronization)) |
| **wait()** | Método que suspende un hilo y libera el lock del objeto | [Java Object.wait()](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/lang/Object.html#wait()) |
| **notifyAll()** | Despierta todos los hilos que están en wait() sobre el objeto | [Java Object.notifyAll()](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/lang/Object.html#notifyAll()) |
| **join()** | Espera a que un hilo termine su ejecución | [Java Thread.join()](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/lang/Thread.html#join()) |
| **ExecutorService** | Interfaz para gestionar pools de hilos y tareas asíncronas | [Java ExecutorService](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/ExecutorService.html) |
| **Callable** | Interfaz como Runnable pero que puede retornar un valor | [Java Callable](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/Callable.html) |
| **Future** | Representa el resultado futuro de una tarea asíncrona | [Java Future](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/Future.html) |
| **Thread Pool** | Conjunto reutilizable de hilos para ejecutar tareas | [Java Thread Pool](https://docs.oracle.com/javase/tutorial/essential/concurrency/pools.html) |
| **Objeto Inmutable** | Objeto cuyo estado no puede cambiar después de construirse | [Oracle - Immutable Objects](https://docs.oracle.com/javase/tutorial/essential/concurrency/immutable.html) |
| **Goroutine** | Hilo ultraligero de Go gestionado por el runtime de Go | [Go Goroutines](https://go.dev/tour/concurrency/1) |
| **Channel (Canal)** | Conducto tipado para comunicación entre goroutines en Go | [Go Channels](https://go.dev/tour/concurrency/2) |
| **select (Go)** | Sentencia que espera en múltiples operaciones de canal | [Go Select](https://go.dev/tour/concurrency/5) |
| **WaitGroup** | Mecanismo en Go para esperar a que un conjunto de goroutines termine | [Go WaitGroup](https://pkg.go.dev/sync#WaitGroup) |
| **Mutex** | Mecanismo de exclusión mutua para proteger secciones críticas | [Go Mutex](https://pkg.go.dev/sync#Mutex) |
| **Busy-waiting** | Técnica ineficiente donde un hilo verifica continuamente una condición | [Wikipedia - Busy Waiting](https://en.wikipedia.org/wiki/Busy_waiting) |
| **Guarded Block** | Bloque de código que solo se ejecuta cuando una condición se cumple | [Oracle - Guarded Blocks](https://docs.oracle.com/javase/tutorial/essential/concurrency/guardmeth.html) |
| **Intrinsic Lock** | Lock intrínseco de Java asociado a cada objeto | [Oracle - Intrinsic Locks](https://docs.oracle.com/javase/tutorial/essential/concurrency/locksync.html) |
| **AtomicInteger** | Variable entera que soporta operaciones atómicas sin sincronización explícita | [Java AtomicInteger](https://docs.oracle.com/en/java/se/17/docs/api/java.base/java/util/concurrent/atomic/AtomicInteger.html) |

---

> 📌 **Conclusión**: El documento aborda la concurrencia desde dos perspectivas: **Java**, que usa memoria compartida con sincronización explícita (`synchronized`, `wait/notify`, `CyclicBarrier`, `ExecutorService`), y **Go**, que usa paso de mensajes mediante **goroutines y canales**, resultando en código más simple y menos propenso a errores de sincronización. Ambos enfoques tienen sus casos de uso ideales según la naturaleza del problema a resolver.