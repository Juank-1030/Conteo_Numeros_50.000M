import java.util.Scanner;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Cuenta del 1 al N usando múltiples hilos en paralelo y mide el tiempo de
 * ejecución.
 * <p>
 * Divide el rango total en segmentos equitativos, asigna cada segmento a un
 * hilo
 * y utiliza un contador atómico compartido ({@link AtomicLong}) para sumar los
 * resultados sin condiciones de carrera.
 * </p>
 *
 * @since 1.0
 */
public class ContadorParalelo {

    /**
     * Solicita un número al usuario, lo lee desde la entrada estándar y valida
     * que esté dentro del rango especificado.
     *
     * @param teclado Escáner conectado a la entrada estándar.
     * @param mensaje Texto que se muestra al usuario para pedir el número.
     * @param minimo  Valor mínimo aceptado (incluido).
     * @param maximo  Valor máximo aceptado (incluido).
     * @return El número ingresado por el usuario, ya validado.
     * @since 1.0
     */
    static long pedirNumero(Scanner teclado, String mensaje, long minimo, long maximo) {
        System.out.print(mensaje);
        long valor = teclado.nextLong();

        if (valor < minimo || valor > maximo) {
            System.out.println("Error: el valor debe estar entre " + minimo + " y " + maximo + ".");
            System.exit(1);
        }

        return valor;
    }

    /**
     * Trabajo que ejecuta cada hilo: cuenta los números de su segmento asignado.
     * <p>
     * En modo detallado itera número por número llamando a
     * {@link AtomicLong#incrementAndGet()} e imprime cada valor. En modo resumen
     * suma todo el bloque de una sola vez con {@link AtomicLong#addAndGet(long)}.
     * </p>
     *
     * @param numeroHilo    Identificador del hilo (1-based).
     * @param inicio        Primer número del segmento (incluido).
     * @param fin           Último número del segmento (incluido).
     * @param modoDetallado {@code true} imprime cada número; {@code false} solo un
     *                      resumen.
     * @param contadorTotal Contador atómico compartido entre todos los hilos.
     * @since 1.0
     */
    static void contarSegmento(int numeroHilo, long inicio, long fin,
            boolean modoDetallado, AtomicLong contadorTotal) {

        long cantidadNumeros = fin - inicio + 1;

        if (modoDetallado) {
            // Print each number one by one
            for (long numero = inicio; numero <= fin; numero++) {
                contadorTotal.incrementAndGet();
                System.out.println("  [Hilo-" + numeroHilo + "] --> " + numero);
            }
        } else {
            // Add the complete block at once (fast mode)
            contadorTotal.addAndGet(cantidadNumeros);
            System.out.println("  Hilo-" + numeroHilo
                    + ": del " + inicio + " al " + fin
                    + "  (" + cantidadNumeros + " números)");
        }
    }

    /**
     * Punto de entrada del programa.
     * <p>
     * Solicita los parámetros al usuario (número final, cantidad de hilos y modo),
     * divide el trabajo en segmentos equitativos, lanza los hilos en paralelo,
     * espera a que finalicen y muestra el resultado con el tiempo total
     * transcurrido.
     * </p>
     *
     * @param args Argumentos de línea de comandos (no se utilizan).
     * @throws InterruptedException Si algún hilo es interrumpido mientras se espera
     *                              su finalización.
     * @since 1.0
     */
    public static void main(String[] args) throws InterruptedException {

        Scanner teclado = new Scanner(System.in);

        // 1. Read user parameters
        long numeroFinal = pedirNumero(teclado, "Número final (1 - 50.000.000.000): ", 1, 50_000_000_000L);
        int cantidadHilos = (int) pedirNumero(teclado, "Cantidad de hilos: ", 1, Integer.MAX_VALUE);
        boolean modoDetallado = pedirNumero(teclado, "Modo (1 = resumen | 2 = número a número): ", 1, 2) == 2;

        // 2. Calculate the size of each segment
        long tamanoSegmento = numeroFinal / cantidadHilos;

        // 3. Shared counter between all threads
        AtomicLong contadorTotal = new AtomicLong(0);

        System.out.println("\nContando del 1 al " + numeroFinal
                + " usando " + cantidadHilos + " hilo(s)...\n");

        // 4. Record the start time
        long tiempoInicio = System.currentTimeMillis();

        // 5. Create and start each thread
        Thread[] hilos = new Thread[cantidadHilos];

        for (int i = 0; i < cantidadHilos; i++) {
            final int numeroHilo = i + 1;
            final long inicio = i * tamanoSegmento + 1;
            final long fin = (i == cantidadHilos - 1) ? numeroFinal : (i + 1) * tamanoSegmento;

            hilos[i] = new Thread(
                    () -> contarSegmento(numeroHilo, inicio, fin, modoDetallado, contadorTotal),
                    "Hilo-" + numeroHilo);
            hilos[i].start();
        }

        // 6. Wait for all threads to finish
        for (Thread hilo : hilos) {
            hilo.join();
        }

        // 7. Calculate total elapsed time
        long tiempoFinal = System.currentTimeMillis();
        long milisegundos = tiempoFinal - tiempoInicio;
        double segundos = milisegundos / 1000.0;

        // 8. Show the final result
        System.out.println();
        System.out.println("--- Resultado ---");
        System.out.println("Total contado  : " + contadorTotal.get());
        System.out.printf("Tiempo de mora : %d ms (%.3f s)%n", milisegundos, segundos);
    }
}
