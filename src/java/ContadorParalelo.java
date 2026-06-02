import java.util.Scanner;
import java.util.concurrent.atomic.AtomicLong;

public class ContadorParalelo {

    // ---------------------------------------------------------------
    // Lectura y validación de un número ingresado por el usuario
    // ---------------------------------------------------------------
    static long pedirNumero(Scanner teclado, String mensaje, long minimo, long maximo) {
        System.out.print(mensaje);
        long valor = teclado.nextLong();

        if (valor < minimo || valor > maximo) {
            System.out.println("Error: el valor debe estar entre " + minimo + " y " + maximo + ".");
            System.exit(1);
        }

        return valor;
    }

    // ---------------------------------------------------------------
    // Trabajo que realiza cada hilo: contar su segmento asignado
    // ---------------------------------------------------------------
    static void contarSegmento(int numeroHilo, long inicio, long fin,
            boolean modoDetallado, AtomicLong contadorTotal) {

        long cantidadNumeros = fin - inicio + 1;

        if (modoDetallado) {
            // Imprime cada número uno a uno
            for (long numero = inicio; numero <= fin; numero++) {
                contadorTotal.incrementAndGet();
                System.out.println("  [Hilo-" + numeroHilo + "] --> " + numero);
            }
        } else {
            // Suma el bloque completo de una sola vez (modo rápido)
            contadorTotal.addAndGet(cantidadNumeros);
            System.out.println("  Hilo-" + numeroHilo
                    + ": del " + inicio + " al " + fin
                    + "  (" + cantidadNumeros + " números)");
        }
    }

    public static void main(String[] args) throws InterruptedException {

        Scanner teclado = new Scanner(System.in);

        // 1. Leer parámetros del usuario
        long numeroFinal = pedirNumero(teclado, "Número final (1 - 50.000.000.000): ", 1, 50_000_000_000L);
        int cantidadHilos = (int) pedirNumero(teclado, "Cantidad de hilos: ", 1, Integer.MAX_VALUE);
        boolean modoDetallado = pedirNumero(teclado, "Modo (1 = resumen | 2 = número a número): ", 1, 2) == 2;

        // 2. Calcular el tamaño de cada segmento
        long tamanoSegmento = numeroFinal / cantidadHilos;

        // 3. Contador compartido entre todos los hilos
        AtomicLong contadorTotal = new AtomicLong(0);

        System.out.println("\nContando del 1 al " + numeroFinal
                + " usando " + cantidadHilos + " hilo(s)...\n");

        // 4. Registrar el tiempo de inicio
        long tiempoInicio = System.currentTimeMillis();

        // 5. Crear y arrancar cada hilo
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

        // 6. Esperar a que todos los hilos terminen
        for (Thread hilo : hilos) {
            hilo.join();
        }

        // 7. Calcular el tiempo total transcurrido
        long tiempoFinal = System.currentTimeMillis();
        long milisegundos = tiempoFinal - tiempoInicio;
        double segundos = milisegundos / 1000.0;

        // 8. Mostrar el resultado final
        System.out.println();
        System.out.println("--- Resultado ---");
        System.out.println("Total contado  : " + contadorTotal.get());
        System.out.printf("Tiempo de mora : %d ms (%.3f s)%n", milisegundos, segundos);
    }
}
