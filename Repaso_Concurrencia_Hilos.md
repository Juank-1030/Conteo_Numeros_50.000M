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

---

## 5. Caso de la vida real: Sistema de Restaurante Concurrente

**Escenario:** Un restaurante donde varios meseros toman pedidos de clientes y la cocina los prepara. Cada pedido tiene un número de recibo único, pasa por varias etapas y al final se genera una cuenta.

Este ejemplo integra **todos los conceptos de concurrencia** en un solo sistema realista.

### 5.1 Diagrama del flujo

```
                    ┌──────────────────┐
                    │    CAJEROS       │ (productores)
                    │  (toman pedidos) │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │                  │
                    │  COLA DE        │ ← BlockingQueue / chan
                    │  PEDIDOS        │
                    │                  │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │                  │
                    │   COCINA        │ ← Pool de workers
                    │  (cocineros)    │
                    │                  │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │                  │
                    │  FACTURACIÓN     │ ← Escritura de recibos
                    │  (recibo final)  │
                    └──────────────────┘
```

### 5.2 Componentes del sistema

| Componente | Rol en concurrencia | Concepto que demuestra |
|------------|---------------------|------------------------|
| **`Pedido`** | Dato inmutable que viaja entre hilos | Objeto inmutable (thread-safe por naturaleza) |
| **`GeneradorRecibos`** | Asigna número único a cada pedido | `AtomicLong` / `atomic.AddInt64` |
| **`Cajero`** | Toma pedidos y los encola | Hilo productor |
| **`Cocina`** | Pool de cocineros que preparan pedidos | Pool de workers, consumidores |
| **`SistemaRestaurante`** | Coordina todo el flujo y espera | `join()` / `WaitGroup` |

### 5.3 Paso a paso del flujo completo

```
1. El sistema inicia: crea cocineros (pool de workers)
2. Cada cajero (hilo) genera pedidos de clientes:
   a. Crea un Pedido (inmutable)
   b. Le asigna un número de recibo único (atómico)
   c. Lo pone en la cola de pedidos (productor)
3. Los cocineros (workers) compiten por los pedidos:
   a. Toman el siguiente pedido de la cola (consumidor)
   b. "Preparan" el pedido (simulado con sleep)
   c. Generan el recibo final (escriben en archivo compartido)
4. El sistema espera a que todos los pedidos terminen
5. Se muestra el resumen final
```

---

### 5.4 Implementación en Java

#### Clase Pedido (inmutable)

```java
// CLASE INMUTABLE: una vez creada, no puede cambiar
// Esto la hace inherentemente THREAD-SAFE
import java.util.List;
import java.util.ArrayList;
import java.util.Collections;

public final class Pedido {

    private final int numeroRecibo;          // Número único del recibo
    private final String cliente;            // Nombre del cliente
    private final List<String> platillos;    // Lista de platillos (inmutable)
    private final long tiempoCreacion;       // Timestamp de creación

    public Pedido(int numeroRecibo, String cliente, List<String> platillos) {
        this.numeroRecibo = numeroRecibo;
        this.cliente = cliente;
        // COPIA DEFENSIVA: evitamos que modifiquen la lista original
        this.platillos = Collections.unmodifiableList(new ArrayList<>(platillos));
        this.tiempoCreacion = System.currentTimeMillis();
    }

    // Solo Getters — no hay setters (el objeto nunca cambia)
    public int getNumeroRecibo()       { return numeroRecibo; }
    public String getCliente()         { return cliente; }
    public List<String> getPlatillos() { return platillos; }
    public long getTiempoCreacion()    { return tiempoCreacion; }

    @Override
    public String toString() {
        return String.format("Recibo #%d | Cliente: %s | Platillos: %s",
                numeroRecibo, cliente, String.join(", ", platillos));
    }
}
```

#### Clase GeneradorRecibos (contador atómico)

```java
// CONTADOR ATÓMICO: asigna números de recibo únicos sin usar cerrojos
import java.util.concurrent.atomic.AtomicInteger;

public class GeneradorRecibos {

    private static final AtomicInteger contador = new AtomicInteger(1000);
    // Empieza en 1000 para simular recibos realistas

    // incrementAndGet() usa CAS (Compare-And-Swap):
    // 1. Lee el valor actual
    // 2. Suma 1
    // 3. Si nadie cambió el valor entre lectura y escritura → escribe
    // 4. Si alguien lo cambió → reintenta desde el paso 1
    public static int siguienteNumero() {
        return contador.incrementAndGet();
    }

    // Para consultar el total sin modificarlo
    public static int totalRecibosEmitidos() {
        return contador.get() - 1000;
    }
}
```

#### Clase Cajero (hilo productor)

```java
// HILO PRODUCTOR: genera pedidos y los pone en la cola
import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;

public class Cajero extends Thread {

    private final String nombre;
    private final BlockingQueue<Pedido> colaPedidos;
    private final List<String> menu;
    private final int pedidosAGenerar;

    public Cajero(String nombre, BlockingQueue<Pedido> colaPedidos,
                  List<String> menu, int pedidosAGenerar) {
        this.nombre = nombre;
        this.colaPedidos = colaPedidos;
        this.menu = menu;
        this.pedidosAGenerar = pedidosAGenerar;
    }

    @Override
    public void run() {
        try {
            for (int i = 0; i < pedidosAGenerar; i++) {
                // 1. Elegir platillo aleatorio del menú
                String platillo = menu.get((int)(Math.random() * menu.size()));

                // 2. Crear el pedido (inmutable) con número de recibo único
                Pedido pedido = new Pedido(
                    GeneradorRecibos.siguienteNumero(),  // ← atómico
                    "Cliente-" + nombre + "-" + (i + 1),
                    List.of(platillo)
                );

                // 3. Poner el pedido en la cola compartida
                // put() se bloquea si la cola está llena
                colaPedidos.put(pedido);

                System.out.printf("  👤 %s tomó pedido: %s%n",
                        nombre, pedido);

                // 4. Simular tiempo entre clientes
                TimeUnit.MILLISECONDS.sleep(100 + (long)(Math.random() * 200));
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
```

#### Clase Cocinero (worker consumidor)

```java
// WORKER CONSUMIDOR: toma pedidos de la cola y los prepara
import java.io.BufferedWriter;
import java.io.IOException;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class Cocinero extends Thread {

    private final String nombre;
    private final BlockingQueue<Pedido> colaPedidos;
    private final BufferedWriter archivoRecibos;
    private final AtomicInteger pedidosPreparados;
    private volatile boolean ejecutando = true;

    public Cocinero(String nombre, BlockingQueue<Pedido> colaPedidos,
                    BufferedWriter archivoRecibos,
                    AtomicInteger pedidosPreparados) {
        this.nombre = nombre;
        this.colaPedidos = colaPedidos;
        this.archivoRecibos = archivoRecibos;
        this.pedidosPreparados = pedidosPreparados;
    }

    // Método para detener el worker de forma cooperativa
    public void detener() {
        this.ejecutando = false;
    }

    @Override
    public void run() {
        try {
            while (ejecutando) {
                // poll() con timeout: intenta tomar un pedido por 1 segundo
                // Si no hay pedidos en 1s y ya se cerró la cocina → termina
                Pedido pedido = colaPedidos.poll(1, TimeUnit.SECONDS);

                if (pedido == null) {
                    continue;  // No había pedido, reintentar
                }

                // SIMULACIÓN: preparar el pedido (tarda 200-500ms)
                System.out.printf("    🍳 %s cocinando: Recibo #%d (%s)%n",
                        nombre, pedido.getNumeroRecibo(),
                        pedido.getPlatillos().get(0));
                TimeUnit.MILLISECONDS.sleep(200 + (long)(Math.random() * 300));

                // ESCRIBIR RECIBO EN ARCHIVO (sección crítica con cerrojo)
                // El archivo es compartido entre todos los cocineros → synchronized
                synchronized (archivoRecibos) {
                    archivoRecibos.write(String.format(
                        "================================%n"));
                    archivoRecibos.write(String.format(
                        "  RESTAURANTE CONCURRENTE%n"));
                    archivoRecibos.write(String.format(
                        "================================%n"));
                    archivoRecibos.write(String.format(
                        "  Recibo #%d%n", pedido.getNumeroRecibo()));
                    archivoRecibos.write(String.format(
                        "  Cliente: %s%n", pedido.getCliente()));
                    archivoRecibos.write(String.format(
                        "  Platillo: %s%n", pedido.getPlatillos().get(0)));
                    archivoRecibos.write(String.format(
                        "  Atendido por: %s%n", nombre));
                    archivoRecibos.write(String.format(
                        "================================%n%n"));
                    archivoRecibos.flush();
                }

                // Incrementar contador de pedidos preparados (atómico)
                pedidosPreparados.incrementAndGet();

                System.out.printf("    ✅ %s terminó: Recibo #%d%n",
                        nombre, pedido.getNumeroRecibo());
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } catch (IOException e) {
            System.err.println("Error escribiendo recibo: " + e.getMessage());
        }
    }
}
```

#### Clase SistemaRestaurante (orquestador principal)

```java
// ORQUESTADOR PRINCIPAL: coordina todo el sistema concurrente
import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

public class SistemaRestaurante {

    public static void main(String[] args) throws InterruptedException, IOException {

        System.out.println("===================================");
        System.out.println("  🏪 RESTAURANTE CONCURRENTE");
        System.out.println("===================================\n");

        // ============================================================
        // CONFIGURACIÓN DEL SISTEMA
        // ============================================================

        // Menú del restaurante
        List<String> menu = List.of(
            "Hamburguesa", "Pizza", "Ensalada", "Tacos", "Pasta"
        );

        int numCajeros = 3;         // Hilos productores
        int pedidosPorCajero = 4;   // Cada cajero toma 4 pedidos
        int numCocineros = 2;       // Workers consumidores
        String archivoRecibosPath = "recibos_restaurante.txt";

        // Contadores atómicos compartidos
        AtomicInteger pedidosPreparados = new AtomicInteger(0);

        // Cola compartida entre cajeros (productores) y cocineros (consumidores)
        // ArrayBlockingQueue: cola con capacidad máxima (evita desbordamiento)
        BlockingQueue<Pedido> colaPedidos = new ArrayBlockingQueue<>(10);

        // Archivo de recibos (compartido entre todos los cocineros)
        BufferedWriter archivoRecibos = new BufferedWriter(
            new FileWriter(archivoRecibosPath)
        );

        // ============================================================
        // INICIAR COCINEROS (workers consumidores)
        // ============================================================

        List<Cocinero> cocineros = new ArrayList<>();
        for (int i = 1; i <= numCocineros; i++) {
            Cocinero cocinero = new Cocinero(
                "Cocinero-" + i, colaPedidos, archivoRecibos, pedidosPreparados
            );
            cocineros.add(cocinero);
            cocinero.start();  // Inicia el hilo (espera pedidos en la cola)
        }

        // ============================================================
        // INICIAR CAJEROS (hilos productores)
        // ============================================================

        List<Cajero> cajeros = new ArrayList<>();
        for (int i = 1; i <= numCajeros; i++) {
            Cajero cajero = new Cajero(
                "Cajero-" + i, colaPedidos, menu, pedidosPorCajero
            );
            cajeros.add(cajero);
            cajero.start();  // Inicia el hilo (genera pedidos)
        }

        // ============================================================
        // ESPERAR A QUE LOS CAJEROS TERMINEN
        // ============================================================

        for (Cajero cajero : cajeros) {
            cajero.join();  // join(): el main espera a que este cajero termine
        }
        System.out.println("\n  📋 Todos los pedidos fueron tomados.\n");

        // ============================================================
        // DETENER LA COCINA (forma cooperativa)
        // ============================================================

        // Esperar un poco para que los cocineros terminen los pedidos en cola
        Thread.sleep(2000);

        for (Cocinero cocinero : cocineros) {
            cocinero.detener();  // Señal de que ya no hay más pedidos
        }

        // Esperar a que los cocineros terminen
        for (Cocinero cocinero : cocineros) {
            cocinero.join();
        }

        // ============================================================
        // CERRAR RECURSOS Y MOSTRAR RESUMEN
        // ============================================================

        archivoRecibos.close();

        System.out.println("\n===================================");
        System.out.println("  📊 RESUMEN DEL DÍA");
        System.out.println("===================================");
        System.out.printf("  Total recibos emitidos: %d%n",
                GeneradorRecibos.totalRecibosEmitidos());
        System.out.printf("  Total pedidos preparados: %d%n",
                pedidosPreparados.get());
        System.out.printf("  Archivo de recibos: %s%n", archivoRecibosPath);
        System.out.println("===================================");
        System.out.println("  🎉 ¡Restaurante cerrado exitosamente!");
        System.out.println("===================================");
    }
}
```

---

### 5.5 Implementación en Go

```go
// SISTEMA DE RESTAURANTE CONCURRENTE EN GO
package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "sync"
    "sync/atomic"
    "time"
)

// ============================================================
// PEDIDO (estructura inmutable)
// ============================================================

// Pedido representa un pedido de cliente.
// Todos los campos son exportados pero nunca se modifican después de creados.
type Pedido struct {
    NumeroRecibo  int
    Cliente       string
    Platillo      string
    TiempoCreacion time.Time
}

// Método String para imprimir el pedido
func (p Pedido) String() string {
    return fmt.Sprintf("Recibo #%d | Cliente: %s | Platillo: %s",
        p.NumeroRecibo, p.Cliente, p.Platillo)
}

// ============================================================
// GENERADOR DE RECIBOS (contador atómico)
// ============================================================

// GeneradorRecibos asigna números de recibo únicos usando atomic
type GeneradorRecibos struct {
    contador int64
}

func NewGeneradorRecibos(inicio int64) *GeneradorRecibos {
    return &GeneradorRecibos{contador: inicio}
}

// SiguienteNumero incrementa y devuelve atómicamente
func (g *GeneradorRecibos) SiguienteNumero() int64 {
    return atomic.AddInt64(&g.contador, 1)
}

func (g *GeneradorRecibos) TotalEmitidos() int64 {
    return atomic.LoadInt64(&g.contador) - 1000
}

// ============================================================
// CAJERO (goroutine productora)
// ============================================================

// Cajero genera pedidos y los envía a la cocina
func Cajero(id int, menu []string, pedidosPorCajero int,
    colaPedidos chan<- Pedido, generador *GeneradorRecibos,
    wg *sync.WaitGroup) {

    defer wg.Done()  // Al terminar, avisar al WaitGroup

    for i := 0; i < pedidosPorCajero; i++ {
        // 1. Elegir platillo aleatorio
        platillo := menu[rand.Intn(len(menu))]

        // 2. Crear pedido con número de recibo único
        pedido := Pedido{
            NumeroRecibo:   int(generador.SiguienteNumero()),
            Cliente:        fmt.Sprintf("Cliente-Cajero%d-%d", id, i+1),
            Platillo:       platillo,
            TiempoCreacion: time.Now(),
        }

        // 3. Enviar pedido a la cola (canal)
        // El envío se bloquea si la cola está llena
        colaPedidos <- pedido

        fmt.Printf("  👤 Cajero-%d tomó pedido: %s\n", id, pedido)

        // 4. Simular tiempo entre clientes
        time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
    }
}

// ============================================================
// COCINERO (goroutine consumidora / worker)
// ============================================================

// Cocinero toma pedidos de la cola y los prepara
func Cocinero(id int, colaPedidos <-chan Pedido,
    archivoRecibos *bufio.Writer, mutexArchivo *sync.Mutex,
    pedidosPreparados *int64, wgCocineros *sync.WaitGroup) {

    defer wgCocineros.Done()

    for pedido := range colaPedidos {
        // SIMULACIÓN: preparar el pedido (200-500ms)
        fmt.Printf("    🍳 Cocinero-%d cocinando: Recibo #%d (%s)\n",
            id, pedido.NumeroRecibo, pedido.Platillo)
        time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)

        // ESCRIBIR RECIBO EN ARCHIVO (sección protegida con mutex)
        mutexArchivo.Lock()
        fmt.Fprintf(archivoRecibos, "================================\n")
        fmt.Fprintf(archivoRecibos, "  RESTAURANTE CONCURRENTE\n")
        fmt.Fprintf(archivoRecibos, "================================\n")
        fmt.Fprintf(archivoRecibos, "  Recibo #%d\n", pedido.NumeroRecibo)
        fmt.Fprintf(archivoRecibos, "  Cliente: %s\n", pedido.Cliente)
        fmt.Fprintf(archivoRecibos, "  Platillo: %s\n", pedido.Platillo)
        fmt.Fprintf(archivoRecibos, "  Atendido por: Cocinero-%d\n", id)
        fmt.Fprintf(archivoRecibos, "================================\n\n")
        archivoRecibos.Flush()
        mutexArchivo.Unlock()

        // Incrementar contador de pedidos preparados (atómico)
        atomic.AddInt64(pedidosPreparados, 1)

        fmt.Printf("    ✅ Cocinero-%d terminó: Recibo #%d\n",
            id, pedido.NumeroRecibo)
    }
}

// ============================================================
// MAIN: ORQUESTADOR PRINCIPAL
// ============================================================

func main() {
    rand.Seed(time.Now().UnixNano())

    fmt.Println("===================================")
    fmt.Println("  🏪 RESTAURANTE CONCURRENTE")
    fmt.Println("===================================")
    fmt.Println()

    // ============================================================
    // CONFIGURACIÓN
    // ============================================================

    menu := []string{"Hamburguesa", "Pizza", "Ensalada", "Tacos", "Pasta"}

    numCajeros := 3
    pedidosPorCajero := 4
    numCocineros := 2
    archivoRecibosPath := "recibos_restaurante.txt"

    var pedidosPreparados int64 = 0
    generador := NewGeneradorRecibos(1000)

    // CANAL: cola compartida entre cajeros y cocineros
    // Canal con búfer de 10 posiciones (como ArrayBlockingQueue)
    colaPedidos := make(chan Pedido, 10)

    // ============================================================
    // ABRIR ARCHIVO DE RECIBOS
    // ============================================================

    archivo, err := os.Create(archivoRecibosPath)
    if err != nil {
        fmt.Println("Error creando archivo:", err)
        return
    }
    defer archivo.Close()

    archivoRecibos := bufio.NewWriter(archivo)
    var mutexArchivo sync.Mutex

    // ============================================================
    // INICIAR COCINEROS (workers)
    // ============================================================

    var wgCocineros sync.WaitGroup
    for i := 1; i <= numCocineros; i++ {
        wgCocineros.Add(1)
        go Cocinero(i, colaPedidos, archivoRecibos,
            &mutexArchivo, &pedidosPreparados, &wgCocineros)
    }

    // ============================================================
    // INICIAR CAJEROS (productores)
    // ============================================================

    var wgCajeros sync.WaitGroup
    for i := 1; i <= numCajeros; i++ {
        wgCajeros.Add(1)
        go Cajero(i, menu, pedidosPorCajero, colaPedidos,
            generador, &wgCajeros)
    }

    // ============================================================
    // ESPERAR A QUE LOS CAJEROS TERMINEN
    // ============================================================

    wgCajeros.Wait()
    fmt.Println("\n  📋 Todos los pedidos fueron tomados.\n")

    // CERRAR LA COLA: los cocineros dejarán de recibir pedidos
    // cuando el canal se cierre y se acaben los pendientes
    close(colaPedidos)

    // ============================================================
    // ESPERAR A QUE LOS COCINEROS TERMINEN
    // ============================================================

    wgCocineros.Wait()
    archivoRecibos.Flush()

    // ============================================================
    // MOSTRAR RESUMEN
    // ============================================================

    fmt.Println()
    fmt.Println("===================================")
    fmt.Println("  📊 RESUMEN DEL DÍA")
    fmt.Println("===================================")
    fmt.Printf("  Total recibos emitidos: %d\n", generador.TotalEmitidos())
    fmt.Printf("  Total pedidos preparados: %d\n", pedidosPreparados)
    fmt.Printf("  Archivo de recibos: %s\n", archivoRecibosPath)
    fmt.Println("===================================")
    fmt.Println("  🎉 ¡Restaurante cerrado exitosamente!")
    fmt.Println("===================================")
}
```

---

### 5.6 Conceptos de concurrencia aplicados en este sistema

| Concepto | Dónde se usa en el sistema |
|----------|---------------------------|
| **Hilo / Goroutine** | Cada `Cajero` y cada `Cocinero` es un hilo/goroutine independiente |
| **Condición de carrera** | Si dos cocineros escribieran al archivo sin `synchronized`/`Mutex`, los recibos se mezclarían |
| **Cerrojo (`synchronized`/`Mutex`)** | Protege la escritura al archivo de recibos (sección crítica) |
| **Operación atómica (`AtomicInteger`/`atomic.AddInt64`)** | El `GeneradorRecibos` asigna números únicos sin cerrojos |
| **Cola compartida (`BlockingQueue`/`chan`)** | Los cajeros ponen pedidos, los cocineros los toman (productor-consumidor) |
| **Espera a que terminen (`join()`/`WaitGroup`)** | El main espera a que cajeros terminen de tomar pedidos y cocineros terminen de preparar |
| **Objeto inmutable (`Pedido`)** | Los pedidos no cambian después de creados, seguros entre hilos sin sincronización |
| **Detención cooperativa (`volatile`/`close(chan)`)** | Los cocineros se detienen limpiamente cuando ya no hay más pedidos |
| **Pool de workers** | Los `Cocineros` forman un pool fijo que procesa pedidos de una cola común |
| **Productor-Consumidor** | Cajeros = productores, Cocineros = consumidores, Cola = buffer |

### 5.7 Cómo ejecutar

**Java:**
```bash
javac SistemaRestaurante.java Pedido.java GeneradorRecibos.java Cajero.java Cocinero.java
java SistemaRestaurante
```

**Go:**
```bash
go run restaurante.go
```

### 5.8 Salida esperada

```
===================================
  🏪 RESTAURANTE CONCURRENTE
===================================

  👤 Cajero-1 tomó pedido: Recibo #1001 | Cliente: Cliente-Cajero1-1 | Platillo: Pizza
  👤 Cajero-2 tomó pedido: Recibo #1002 | Cliente: Cliente-Cajero2-1 | Platillo: Tacos
    🍳 Cocinero-1 cocinando: Recibo #1001 (Pizza)
    🍳 Cocinero-2 cocinando: Recibo #1002 (Tacos)
  👤 Cajero-3 tomó pedido: Recibo #1003 | Cliente: Cliente-Cajero3-1 | Platillo: Ensalada
  ...
    ✅ Cocinero-1 terminó: Recibo #1001
    ✅ Cocinero-2 terminó: Recibo #1002
  ...

  📋 Todos los pedidos fueron tomados.

===================================
  📊 RESUMEN DEL DÍA
===================================
  Total recibos emitidos: 12
  Total pedidos preparados: 12
  Archivo de recibos: recibos_restaurante.txt
===================================
  🎉 ¡Restaurante cerrado exitosamente!
===================================
```

### 5.9 Ejercicios de extensión

Una vez que entiendas el sistema base, intenta:

1. **Agregar un nuevo tipo de worker** (Ej: `Repartidor`) que toma pedidos listos y los entrega
2. **Implementar prioridad de pedidos** (los pedidos urgentes se procesan antes)
3. **Limitar el tiempo de preparación** con un timeout
4. **Agregar un monitor** que muestre en tiempo real cuántos pedidos hay en la cola
5. **Convertir el sistema a usar `ExecutorService`** en Java (en lugar de hilos manuales)
6. **En Go, agregar un `select`** para manejar múltiples canales (pedidos nuevos + señales de cancelación)

---

## 6. Caso de la vida real: Sistema de Turnos para Bienestar Universitario

**Escenario:** Una oficina de Bienestar Universitario ofrece múltiples servicios (Caja, Información, Consulta). Los estudiantes llegan, toman un ticket de la máquina dispensadora para el servicio que necesitan, y esperan a que su ticket sea llamado por la ventanilla correspondiente.

### 6.1 Diagrama del flujo

```
                        ┌──────────────────────┐
                        │  MÁQUINA DISPENSADORA │
                        │  (contador atómico)   │
                        └──────┬───────────┬────┘
                               │           │
              ┌────────────────┤           ├────────────────┐
              ▼                ▼           ▼                ▼
        ┌──────────┐    ┌──────────┐    ┌──────────┐
        │  CAJA    │    │INFORMACIÓN│    │ CONSULTA │
        │  Cola 1  │    │  Cola 2   │    │  Cola 3   │
        └────┬─────┘    └────┬──────┘    └────┬──────┘
             │               │                │
        ┌────▼─────┐    ┌────▼──────┐    ┌────▼──────┐
        │Ventanilla│    │Ventanilla │    │Ventanilla  │
        │ Caja-1   │    │ Info-1    │    │ Consulta-1 │
        │ Caja-2   │    │ Info-2    │    │ Consulta-2 │
        └──────────┘    └───────────┘    └────────────┘
             │               │                │
             └───────────────┼────────────────┘
                             ▼
                    ┌──────────────────┐
                    │   TABLERO        │
                    │  (pantalla LCD)  │
                    │  "Ahora: C-005"  │
                    └──────────────────┘
```

### 6.2 Componentes del sistema

| Componente | Rol en concurrencia | Concepto que demuestra |
|------------|---------------------|------------------------|
| **`Ticket`** | Dato inmutable con tipo de servicio y número | Objeto inmutable |
| **`Dispensador`** | Asigna números secuenciales por tipo de servicio | Múltiples contadores atómicos |
| **`Estudiante`** | Llega, toma ticket y espera ser atendido | Hilo con llegada aleatoria |
| **`Ventanilla`** | Atiende estudiantes de un tipo de servicio específico | Worker con cola dedicada |
| **`Tablero`** | Muestra en tiempo real el ticket actual en cada ventanilla | Monitor compartido con cerrojo |
| **`SistemaBienestar`** | Coordina todo y genera reporte final | Orquestador principal |

### 6.3 Paso a paso del flujo completo

```
1. El sistema inicia: abre las ventanillas de cada servicio
2. Cada estudiante (hilo) llega en un tiempo aleatorio:
   a. Elige un servicio al azar (Caja, Información o Consulta)
   b. Toma un ticket del dispensador para ese servicio (atómico)
   c. Espera en la cola del servicio correspondiente
   d. Cuando su ticket es llamado, es atendido (simulado con sleep)
   e. Sale de la oficina
3. Cada ventanilla (worker):
   a. Toma el siguiente ticket de su cola
   b. Actualiza el tablero compartido (cerrojo)
   c. Atiende al estudiante (simulado)
   d. Registra la atención
4. Después de un tiempo, la oficina cierra:
   a. No se aceptan más estudiantes nuevos
   b. Se espera a que terminen los que están en cola
   c. Se genera el reporte del día
```

### 6.4 Implementación en Java

#### Clase Ticket (inmutable)

```java
// CLASE INMUTABLE: un ticket no puede cambiar después de emitido
public final class Ticket {

    private final String codigo;          // Ej: "C-005", "I-012", "CON-003"
    private final String tipoServicio;    // "Caja", "Informacion", "Consulta"
    private final int numero;             // Número secuencial dentro del tipo
    private final long timestamp;         // Hora de emisión

    public Ticket(String codigo, String tipoServicio, int numero) {
        this.codigo = codigo;
        this.tipoServicio = tipoServicio;
        this.numero = numero;
        this.timestamp = System.currentTimeMillis();
    }

    public String getCodigo()        { return codigo; }
    public String getTipoServicio()  { return tipoServicio; }
    public int getNumero()           { return numero; }
    public long getTimestamp()       { return timestamp; }

    @Override
    public String toString() {
        return String.format("[%s] %s-%03d", tipoServicio, codigo, numero);
    }
}
```

#### Clase Dispensador (múltiples contadores atómicos)

```java
// DISPENSADOR: múltiples contadores atómicos (uno por tipo de servicio)
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Dispensador {

    // Mapa de contadores: "Caja" → AtomicInteger, "Informacion" → AtomicInteger, etc.
    private final ConcurrentHashMap<String, AtomicInteger> contadores;

    public Dispensador(String... tiposServicio) {
        this.contadores = new ConcurrentHashMap<>();
        for (String tipo : tiposServicio) {
            // Cada tipo de servicio tiene su propio contador independiente
            contadores.put(tipo, new AtomicInteger(0));
        }
    }

    // Genera un ticket con el siguiente número para el tipo de servicio indicado
    public Ticket emitirTicket(String tipoServicio) {
        AtomicInteger contador = contadores.get(tipoServicio);
        if (contador == null) {
            throw new IllegalArgumentException("Servicio no válido: " + tipoServicio);
        }

        int numero = contador.incrementAndGet();  // ← CAS, atómico por tipo
        String prefijo = switch (tipoServicio) {
            case "Caja" -> "C";
            case "Informacion" -> "I";
            case "Consulta" -> "CON";
            default -> "X";
        };
        String codigo = String.format("%s-%03d", prefijo, numero);

        return new Ticket(codigo, tipoServicio, numero);
    }

    // Reporte de cuántos tickets se emitieron por servicio
    public void imprimirReporte() {
        System.out.println("\n  📊 Tickets emitidos por servicio:");
        for (String tipo : contadores.keySet()) {
            int total = contadores.get(tipo).get();
            System.out.printf("     %-15s → %d tickets%n", tipo, total);
        }
    }
}
```

#### Clase Estudiante (hilo que llega aleatoriamente)

```java
// HILO ESTUDIANTE: llega a la universidad, toma ticket y espera ser atendido
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

public class Estudiante extends Thread {

    private static int contadorEstudiantes = 0;  // Solo para dar nombre único

    private final String nombre;
    private final Dispensador dispensador;
    private final ConcurrentHashMap<String, BlockingQueue<Ticket>> colas;
    private final String[] servicios;

    public Estudiante(Dispensador dispensador,
                      ConcurrentHashMap<String, BlockingQueue<Ticket>> colas,
                      String[] servicios) {
        super("Estudiante-" + (++contadorEstudiantes));
        this.nombre = this.getName();
        this.dispensador = dispensador;
        this.colas = colas;
        this.servicios = servicios;
    }

    @Override
    public void run() {
        try {
            // 1. Simular que el estudiante camina hacia la oficina (50-300ms)
            TimeUnit.MILLISECONDS.sleep(50 + (long)(Math.random() * 250));

            // 2. Elegir un servicio al azar
            String servicio = servicios[(int)(Math.random() * servicios.length)];

            // 3. Tomar un ticket del dispensador (atómico)
            Ticket ticket = dispensador.emitirTicket(servicio);

            System.out.printf("  🧑‍🎓 %s tomó ticket %s para %s%n",
                    nombre, ticket.getCodigo(), servicio);

            // 4. Poner el ticket en la cola del servicio correspondiente
            BlockingQueue<Ticket> cola = colas.get(servicio);
            cola.put(ticket);  // Se bloquea si la cola está llena

            // 5. El ticket ya está en cola, el estudiante espera pasivamente
            //    La ventanilla lo tomará cuando le toque su turno
            //    (En este sistema, el estudiante "espera" representado por
            //     el ticket en la cola. En un sistema real, el estudiante
            //     estaría sentado viendo el tablero.)

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
```

#### Clase Ventanilla (worker que atiende estudiantes)

```java
// WORKER VENTANILLA: atiende estudiantes de un tipo de servicio específico
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class Ventanilla extends Thread {

    private final String nombre;
    private final String tipoServicio;
    private final BlockingQueue<Ticket> cola;
    private final Tablero tablero;
    private final AtomicInteger totalAtendidos;
    private volatile boolean abierta = true;

    public Ventanilla(String nombre, String tipoServicio,
                      BlockingQueue<Ticket> cola, Tablero tablero,
                      AtomicInteger totalAtendidos) {
        this.nombre = nombre;
        this.tipoServicio = tipoServicio;
        this.cola = cola;
        this.tablero = tablero;
        this.totalAtendidos = totalAtendidos;
    }

    public void cerrar() {
        this.abierta = false;
    }

    @Override
    public void run() {
        try {
            while (abierta) {
                // Intentar tomar el siguiente ticket de la cola (espera 1s)
                Ticket ticket = cola.poll(1, TimeUnit.SECONDS);

                if (ticket == null) {
                    continue;  // No hay estudiantes, seguir esperando
                }

                // 1. Actualizar el tablero (sección crítica)
                tablero.mostrarAtendiendo(tipoServicio, ticket.getCodigo());

                // 2. Simular la atención (300-800ms según el servicio)
                long tiempoAtencion = switch (tipoServicio) {
                    case "Caja" -> 200 + (long)(Math.random() * 400);   // Rápido
                    case "Informacion" -> 300 + (long)(Math.random() * 500);
                    case "Consulta" -> 500 + (long)(Math.random() * 1000); // Más lento
                    default -> 300;
                };

                System.out.printf("    🪟 %s atendiendo %s...%n",
                        nombre, ticket.getCodigo());
                TimeUnit.MILLISECONDS.sleep(tiempoAtencion);

                // 3. Registrar atención
                totalAtendidos.incrementAndGet();
                System.out.printf("    ✅ %s atendió %s (%s)%n",
                        nombre, ticket.getCodigo(),
                        tipoServicio);
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
```

#### Clase Tablero (monitor compartido con cerrojo)

```java
// TABLERO COMPARTIDO: muestra el ticket actual siendo atendido en cada servicio
// El cerrojo evita que dos ventanillas actualicen al mismo tiempo
import java.util.HashMap;
import java.util.Map;

public class Tablero {

    private final Map<String, String> atendiendo = new HashMap<>();

    // synchronized: solo una ventanilla puede actualizar el tablero a la vez
    public synchronized void mostrarAtendiendo(String servicio, String codigoTicket) {
        atendiendo.put(servicio, codigoTicket);
        System.out.printf("\n  📺 TABLERO ─ %s: Atendiendo %s%n",
                servicio, codigoTicket);
    }

    // synchronized: solo una ventanilla puede leer el tablero a la vez
    public synchronized String consultarAtendiendo(String servicio) {
        return atendiendo.getOrDefault(servicio, "---");
    }

    // Mostrar estado completo del tablero (útil para resumen)
    public synchronized void mostrarEstadoCompleto() {
        System.out.println("\n  📺 TABLERO ─ Estado actual:");
        for (String servicio : atendiendo.keySet()) {
            System.out.printf("    %-15s → Atendiendo: %s%n",
                    servicio, atendiendo.get(servicio));
        }
    }
}
```

#### Clase SistemaBienestar (orquestador principal)

```java
// ORQUESTADOR PRINCIPAL: Sistema de Turnos para Bienestar Universitario
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

public class SistemaBienestar {

    public static void main(String[] args) throws InterruptedException {

        System.out.println("===============================================");
        System.out.println("  🏫 BIENESTAR UNIVERSITARIO — Sistema de Turnos");
        System.out.println("===============================================\n");

        // ================================================================
        // CONFIGURACIÓN DEL SISTEMA
        // ================================================================

        String[] servicios = {"Caja", "Informacion", "Consulta"};
        int aforoMaximo = 5;  // Máximo de estudiantes esperando por servicio

        // Map: tipo de servicio → cola de tickets (una cola independiente por servicio)
        ConcurrentHashMap<String, BlockingQueue<Ticket>> colas = new ConcurrentHashMap<>();
        for (String servicio : servicios) {
            // ArrayBlockingQueue: capacidad limitada (controla aforo)
            colas.put(servicio, new ArrayBlockingQueue<>(aforoMaximo));
        }

        // Dispensador con contadores atómicos independientes por servicio
        Dispensador dispensador = new Dispensador(servicios);

        // Tablero compartido (con cerrojo interno)
        Tablero tablero = new Tablero();

        // ================================================================
        // ABRIR VENTANILLAS (workers)
        // ================================================================

        // Configuración: qué ventanillas atienden qué servicios
        Map<String, Integer> ventanillasPorServicio = new HashMap<>();
        ventanillasPorServicio.put("Caja", 2);         // 2 ventanillas de caja
        ventanillasPorServicio.put("Informacion", 2);  // 2 de información
        ventanillasPorServicio.put("Consulta", 1);     // 1 de consulta

        List<Ventanilla> ventanillas = new ArrayList<>();
        for (String servicio : servicios) {
            int cantidad = ventanillasPorServicio.get(servicio);
            for (int i = 1; i <= cantidad; i++) {
                String nombre = servicio.substring(0, 4) + "-" + i;
                Ventanilla v = new Ventanilla(
                    nombre, servicio, colas.get(servicio), tablero,
                    new AtomicInteger(0)  // contador de atendidos (indiv. por ventanilla)
                );
                ventanillas.add(v);
                v.start();  // La ventanilla empieza a atender
                System.out.printf("  🪟 Ventanilla %s abierta (%s)%n",
                        nombre, servicio);
            }
        }

        // ================================================================
        // LLEGADA DE ESTUDIANTES (ola de 15 estudiantes)
        // ================================================================

        System.out.println("\n  📢 ¡Comienza la jornada! Llegada de estudiantes...\n");

        int totalEstudiantes = 15;
        List<Estudiante> estudiantes = new ArrayList<>();

        for (int i = 0; i < totalEstudiantes; i++) {
            Estudiante e = new Estudiante(dispensador, colas, servicios);
            estudiantes.add(e);
            e.start();  // El estudiante llega, toma ticket y se pone en cola
            Thread.sleep(80 + (long)(Math.random() * 120));  // Llegada escalonada
        }

        // ================================================================
        // ESPERAR A QUE TODOS LOS ESTUDIANTES TOMEN TICKET
        // ================================================================

        for (Estudiante e : estudiantes) {
            e.join();  // Esperar a que todos hayan tomado ticket
        }

        System.out.println("\n  📋 Todos los estudiantes tienen ticket.");
        System.out.println("  ⏳ Esperando a que las ventanillas terminen...\n");

        // ================================================================
        // CERRAR VENTANILLAS (esperar a que terminen las colas)
        // ================================================================

        // Dar tiempo para que se atiendan los tickets pendientes
        Thread.sleep(5000);

        for (Ventanilla v : ventanillas) {
            v.cerrar();  // Señal cooperativa de cierre
        }

        for (Ventanilla v : ventanillas) {
            v.join();  // Esperar a que termine su último atendido
        }

        // ================================================================
        // REPORTE FINAL
        // ================================================================

        System.out.println("\n===============================================");
        System.out.println("  📊 REPORTE FINAL DEL DÍA");
        System.out.println("===============================================");

        dispensador.imprimirReporte();

        int totalVentanillas = ventanillas.size();
        System.out.printf("\n  Ventanillas abiertas: %d%n", totalVentanillas);
        System.out.printf("  Estudiantes atendidos: %d%n", totalEstudiantes);

        System.out.println("\n===============================================");
        System.out.println("  ✅ ¡Jornada finalizada exitosamente!");
        System.out.println("===============================================");
    }
}
```

---

### 6.5 Implementación en Go

```go
// SISTEMA DE TURNOS PARA BIENESTAR UNIVERSITARIO EN GO
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

// ============================================================
// TICKET (estructura inmutable)
// ============================================================

type Ticket struct {
    Codigo       string
    TipoServicio string
    Numero       int
    Timestamp    time.Time
}

func (t Ticket) String() string {
    return fmt.Sprintf("[%s] %s", t.TipoServicio, t.Codigo)
}

// ============================================================
// DISPENSADOR (múltiples contadores atómicos)
// ============================================================

type Dispensador struct {
    contadores map[string]*int64
    mu         sync.Mutex  // Solo para proteger el mapa, no los contadores
}

func NewDispensador(tipos ...string) *Dispensador {
    d := &Dispensador{contadores: make(map[string]*int64)}
    for _, tipo := range tipos {
        var cero int64 = 0
        d.contadores[tipo] = &cero
    }
    return d
}

func (d *Dispensador) EmitirTicket(tipoServicio string) Ticket {
    d.mu.Lock()
    contador := d.contadores[tipoServicio]
    d.mu.Unlock()

    // Incremento atómico independiente por tipo de servicio
    numero := atomic.AddInt64(contador, 1)

    var prefijo string
    switch tipoServicio {
    case "Caja":
        prefijo = "C"
    case "Informacion":
        prefijo = "I"
    case "Consulta":
        prefijo = "CON"
    default:
        prefijo = "X"
    }
    codigo := fmt.Sprintf("%s-%03d", prefijo, numero)

    return Ticket{
        Codigo:       codigo,
        TipoServicio: tipoServicio,
        Numero:       int(numero),
        Timestamp:    time.Now(),
    }
}

func (d *Dispensador) ImprimirReporte() {
    fmt.Println("\n  📊 Tickets emitidos por servicio:")
    d.mu.Lock()
    defer d.mu.Unlock()
    for tipo, contador := range d.contadores {
        total := atomic.LoadInt64(contador)
        fmt.Printf("     %-15s → %d tickets\n", tipo, total)
    }
}

// ============================================================
// TABLERO (monitor compartido con mutex)
// ============================================================

type Tablero struct {
    mu        sync.Mutex
    atendiendo map[string]string // servicio → codigoTicket
}

func NewTablero() *Tablero {
    return &Tablero{atendiendo: make(map[string]string)}
}

func (t *Tablero) MostrarAtendiendo(servicio, codigoTicket string) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.atendiendo[servicio] = codigoTicket
    fmt.Printf("\n  📺 TABLERO ─ %s: Atendiendo %s\n", servicio, codigoTicket)
}

func (t *Tablero) ConsultarAtendiendo(servicio string) string {
    t.mu.Lock()
    defer t.mu.Unlock()
    if codigo, ok := t.atendiendo[servicio]; ok {
        return codigo
    }
    return "---"
}

// ============================================================
// VENTANILLA (goroutine worker)
// ============================================================

func Ventanilla(id int, nombre, tipoServicio string,
    cola <-chan Ticket, tablero *Tablero,
    totalAtendidos *int64, wg *sync.WaitGroup) {

    defer wg.Done()

    for ticket := range cola {
        // 1. Actualizar tablero
        tablero.MostrarAtendiendo(tipoServicio, ticket.Codigo)

        // 2. Simular atención (tiempo según servicio)
        var tiempoAtencion time.Duration
        switch tipoServicio {
        case "Caja":
            tiempoAtencion = time.Duration(200+rand.Intn(400)) * time.Millisecond
        case "Informacion":
            tiempoAtencion = time.Duration(300+rand.Intn(500)) * time.Millisecond
        case "Consulta":
            tiempoAtencion = time.Duration(500+rand.Intn(1000)) * time.Millisecond
        }

        fmt.Printf("    🪟 %s atendiendo %s...\n", nombre, ticket.Codigo)
        time.Sleep(tiempoAtencion)

        // 3. Registrar atención
        atomic.AddInt64(totalAtendidos, 1)
        fmt.Printf("    ✅ %s atendió %s (%s)\n",
            nombre, ticket.Codigo, tipoServicio)
    }
}

// ============================================================
// MAIN: ORQUESTADOR PRINCIPAL
// ============================================================

type servicioConfig struct {
    nombre         string
    ventanillas    int
    aforo          int
}

func main() {
    rand.Seed(time.Now().UnixNano())

    fmt.Println("===============================================")
    fmt.Println("  🏫 BIENESTAR UNIVERSITARIO — Sistema de Turnos")
    fmt.Println("===============================================")
    fmt.Println()

    // ================================================================
    // CONFIGURACIÓN
    // ================================================================

    servicios := []servicioConfig{
        {nombre: "Caja",        ventanillas: 2, aforo: 5},
        {nombre: "Informacion", ventanillas: 2, aforo: 5},
        {nombre: "Consulta",    ventanillas: 1, aforo: 5},
    }

    totalEstudiantes := 15
    nombresServicios := make([]string, len(servicios))
    for i, s := range servicios {
        nombresServicios[i] = s.nombre
    }

    // ================================================================
    // INICIALIZAR COMPONENTES COMPARTIDOS
    // ================================================================

    dispensador := NewDispensador(nombresServicios...)
    tablero := NewTablero()

    // Crear un canal (cola) por cada servicio
    canales := make(map[string]chan Ticket)
    for _, s := range servicios {
        // Canal con búfer = aforo máximo por servicio
        canales[s.nombre] = make(chan Ticket, s.aforo)
    }

    // ================================================================
    // ABRIR VENTANILLAS
    // ================================================================

    var totalAtendidos int64 = 0
    var wgVentanillas sync.WaitGroup
    contadorVent := 0

    for _, s := range servicios {
        for i := 1; i <= s.ventanillas; i++ {
            contadorVent++
            nombre := fmt.Sprintf("%s-%d", s.nombre[:4], i)
            wgVentanillas.Add(1)
            go Ventanilla(contadorVent, nombre, s.nombre,
                canales[s.nombre], tablero, &totalAtendidos, &wgVentanillas)
            fmt.Printf("  🪟 Ventanilla %s abierta (%s)\n", nombre, s.nombre)
        }
    }

    // ================================================================
    // LLEGADA DE ESTUDIANTES
    // ================================================================

    fmt.Println("\n  📢 ¡Comienza la jornada! Llegada de estudiantes...\n")

    var wgEstudiantes sync.WaitGroup
    estudianteID := 0

    for i := 0; i < totalEstudiantes; i++ {
        wgEstudiantes.Add(1)
        estudianteID++

        go func(id int) {
            defer wgEstudiantes.Done()

            // Simular caminata hacia la oficina
            time.Sleep(time.Duration(50+rand.Intn(250)) * time.Millisecond)

            // Elegir servicio al azar
            servicio := nombresServicios[rand.Intn(len(nombresServicios))]

            // Tomar ticket (atómico)
            ticket := dispensador.EmitirTicket(servicio)

            fmt.Printf("  🧑‍🎓 Estudiante-%d tomó ticket %s para %s\n",
                id, ticket.Codigo, servicio)

            // Poner ticket en la cola del servicio
            canales[servicio] <- ticket

        }(estudianteID)

        // Llegada escalonada
        time.Sleep(time.Duration(80+rand.Intn(120)) * time.Millisecond)
    }

    // ================================================================
    // ESPERAR A QUE TODOS TOMEN TICKET
    // ================================================================

    wgEstudiantes.Wait()
    fmt.Println("\n  📋 Todos los estudiantes tienen ticket.")
    fmt.Println("  ⏳ Esperando a que las ventanillas terminen...\n")

    // Cerrar los canales para que las ventanillas terminen
    // cuando acaben los tickets pendientes
    for _, s := range servicios {
        close(canales[s.nombre])
    }

    // ================================================================
    // ESPERAR A QUE LAS VENTANILLAS TERMINEN
    // ================================================================

    wgVentanillas.Wait()

    // ================================================================
    // REPORTE FINAL
    // ================================================================

    fmt.Println("\n===============================================")
    fmt.Println("  📊 REPORTE FINAL DEL DÍA")
    fmt.Println("===============================================")

    dispensador.ImprimirReporte()

    fmt.Printf("\n  Estudiantes atendidos: %d\n", totalAtendidos)

    fmt.Println("\n===============================================")
    fmt.Println("  ✅ ¡Jornada finalizada exitosamente!")
    fmt.Println("===============================================")
}
```

### 6.6 Archivos del proyecto

```
SistemaTurnos/
├── Java/
│   ├── SistemaBienestar.java     ← Orquestador principal
│   ├── Ticket.java               ← Objeto inmutable
│   ├── Dispensador.java          ← Contadores atómicos
│   ├── Estudiante.java           ← Hilo productor
│   ├── Ventanilla.java           ← Worker consumidor
│   └── Tablero.java              ← Monitor compartido
│
└── Go/
    └── bienestar_universitario.go  ← Todo en un archivo
```

### 6.7 Conceptos de concurrencia aplicados

| Concepto | Dónde se usa en el sistema |
|----------|---------------------------|
| **Múltiples colas independientes** | Cada servicio (Caja, Información, Consulta) tiene su propia `BlockingQueue`/`chan` |
| **Contadores atómicos separados** | El `Dispensador` mantiene un `AtomicInteger`/`int64` independiente por cada tipo de servicio |
| **Monitor compartido con cerrojo** | El `Tablero` se actualiza con `synchronized`/`Mutex` para evitar lecturas inconsistentes |
| **Hilos productores** | Los `Estudiante` generan tickets y los ponen en la cola del servicio elegido |
| **Workers consumidores** | Las `Ventanilla` toman tickets de su cola y atienden (pool por tipo de servicio) |
| **Llegada aleatoria escalonada** | Los estudiantes llegan en intervalos de tiempo variables (simula realidad) |
| **Aforo limitado** | Las colas tienen capacidad máxima (`ArrayBlockingQueue`, buffer de canal) para evitar desbordamiento |
| **Cierre graceful** | Las ventanillas se cierran de forma cooperativa cuando ya no hay más tickets |
| **Ticket como objeto inmutable** | Los `Ticket` no cambian después de creados — seguros entre hilos |
| **switch para tiempos distintos** | Cada tipo de servicio tiene un tiempo de atención diferente (caja es más rápido que consulta) |

### 6.8 Salida esperada

```
===============================================
  🏫 BIENESTAR UNIVERSITARIO — Sistema de Turnos
===============================================

  🪟 Ventanilla Caja-1 abierta (Caja)
  🪟 Ventanilla Caja-2 abierta (Caja)
  🪟 Ventanilla Info-1 abierta (Informacion)
  🪟 Ventanilla Info-2 abierta (Informacion)
  🪟 Ventanilla Cons-1 abierta (Consulta)

  📢 ¡Comienza la jornada! Llegada de estudiantes...

  🧑‍🎓 Estudiante-1 tomó ticket C-001 para Caja
  🧑‍🎓 Estudiante-2 tomó ticket I-001 para Informacion

  📺 TABLERO ─ Caja: Atendiendo C-001
    🪟 Caja-1 atendiendo C-001...
  🧑‍🎓 Estudiante-3 tomó ticket I-002 para Informacion

  📺 TABLERO ─ Informacion: Atendiendo I-001
    🪟 Info-1 atendiendo I-001...
  ...

  📋 Todos los estudiantes tienen ticket.
  ⏳ Esperando a que las ventanillas terminen...

    ✅ Caja-1 atendió C-001 (Caja)

  📺 TABLERO ─ Caja: Atendiendo C-002
    🪟 Caja-2 atendiendo C-002...
  ...

===============================================
  📊 REPORTE FINAL DEL DÍA
===============================================

  📊 Tickets emitidos por servicio:
     Caja           → 6 tickets
     Informacion    → 5 tickets
     Consulta       → 4 tickets

  Estudiantes atendidos: 15

===============================================
  ✅ ¡Jornada finalizada exitosamente!
===============================================
```

### 6.9 Diferencias clave con el sistema de restaurante

| Aspecto | Restaurante (sección 5) | Bienestar Universitario (sección 6) |
|---------|------------------------|-------------------------------------|
| **Colas** | Una sola cola compartida para todos los workers | Una cola independiente por tipo de servicio |
| **Tipo de worker** | Todos los cocineros hacen lo mismo | Cada ventanilla atiende SOLO su servicio |
| **Productores** | Pocos cajeros, muchos pedidos cada uno | Muchos estudiantes, un ticket cada uno |
| **Llegada** | Cajeros producen en ráfaga | Estudiantes llegan escalonados (más realista) |
| **Contadores** | Un solo contador global | Múltiples contadores (uno por servicio) |
| **Tiempo de servicio** | Constante (200-500ms todos) | Variable según el tipo de servicio |
| **Capacidad** | Cola limitada (controla flujo) | Aforo máximo por servicio |
| **Monitor** | No hay tablero | Tablero compartido visible para todos |

### 6.10 Ejercicios de extensión

1. **Agregar prioridad por tipo de estudiante** (Ej: estudiantes con discapacidad pasan al frente)
2. **Implementar uniones de colas** (si una ventanilla de Caja está libre, puede ayudar a Información)
3. **Agregar cancelación de tickets** (estudiante se cansa de esperar y se retira)
4. **Estadísticas en tiempo real** (tiempo de espera promedio por servicio)
5. **Notificación por SMS** cuando el ticket está próximo a ser atendido
6. **Múltiples pisos/edificios** con diferentes servicios y un dispensador central