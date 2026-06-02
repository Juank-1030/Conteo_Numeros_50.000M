@echo off
chcp 65001 > nul
title ContadorParalelo — Ejecutor

echo ============================================
echo   ContadorParalelo — Selector de lenguaje
echo ============================================
echo.
echo  [1]  Java   (Threads + AtomicLong)
echo  [2]  Go     (Goroutines + sync/atomic)
echo  [3]  Ambos  (Java primero, luego Go)
echo  [0]  Salir
echo.
set /p OPCION="Elige una opcion: "

if "%OPCION%"=="1" goto RUN_JAVA
if "%OPCION%"=="2" goto RUN_GO
if "%OPCION%"=="3" goto RUN_AMBOS
if "%OPCION%"=="0" goto FIN

echo Opcion no valida.
pause
exit /b 1

:: ---------------------------------------------------------------
:RUN_JAVA
echo.
echo [Java] Compilando...
javac src\java\ContadorParalelo.java -d out
if errorlevel 1 (
    echo Error al compilar Java. Asegurate de tener el JDK instalado.
    pause
    exit /b 1
)
echo [Java] Ejecutando...
echo.
java -cp out ContadorParalelo
echo.
pause
goto FIN

:: ---------------------------------------------------------------
:RUN_GO
echo.
echo [Go] Ejecutando...
echo.
go run src\go\contador_paralelo.go
if errorlevel 1 (
    echo Error al ejecutar Go. Asegurate de tener Go instalado.
    pause
    exit /b 1
)
echo.
pause
goto FIN

:: ---------------------------------------------------------------
:RUN_AMBOS
echo.
echo ============ [ Java ] ============
echo [Java] Compilando...
javac src\java\ContadorParalelo.java -d out
if errorlevel 1 (
    echo Error al compilar Java.
    pause
    exit /b 1
)
echo [Java] Ejecutando...
echo.
java -cp out ContadorParalelo

echo.
echo ============ [ Go ] ============
echo [Go] Ejecutando...
echo.
go run src\go\contador_paralelo.go
if errorlevel 1 (
    echo Error al ejecutar Go.
    pause
    exit /b 1
)
echo.
pause

:: ---------------------------------------------------------------
:FIN
exit /b 0
