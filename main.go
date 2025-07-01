package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	staticDir := flag.String("static", "./static", "Ruta a los archivos estáticos (html, css, js)")
	port := flag.String("port", "8080", "Puerto para el servidor HTTP")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	// Canal para capturar señales
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Canal para indicar si el usuario quiere salir
	exitChan := make(chan bool, 1)

	go func() {
		for {
			sig := <-sigChan
			if sig == os.Interrupt {
				fmt.Println("\nPresionaste Ctrl+C. Si querés salir del programa, presioná Ctrl+D.")
				// No salimos, solo avisamos
			}
		}
	}()

	for {
		// Verificar si la ruta existe
		if _, err := os.Stat(*staticDir); os.IsNotExist(err) {
			fmt.Printf("La ruta %s no existe.\n", *staticDir)
			fmt.Print("Por favor, ingrese una ruta válida para los archivos estáticos (html, css, js): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				if err == os.ErrClosed || err.Error() == "EOF" {
					fmt.Println("\nSaliendo del programa por Ctrl+D.")
					exitChan <- true
					break
				}
				fmt.Println("Error leyendo la entrada:", err)
				continue
			}
			*staticDir = strings.TrimSpace(input)
			continue
		}

		// Verificar que exista index.html en la ruta
		indexPath := filepath.Join(*staticDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			fmt.Printf("No se encontró index.html en la ruta %s.\n", *staticDir)
			fmt.Print("Por favor, ingrese una ruta válida para los archivos estáticos (html, css, js): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				if err == os.ErrClosed || err.Error() == "EOF" {
					fmt.Println("\nSaliendo del programa por Ctrl+D.")
					exitChan <- true
					break
				}
				fmt.Println("Error leyendo la entrada:", err)
				continue
			}
			*staticDir = strings.TrimSpace(input)
			continue
		}

		break
	}

	select {
	case <-exitChan:
		os.Exit(0)
	default:
		fs := http.FileServer(http.Dir(*staticDir))
		http.Handle("/", fs)

		fmt.Printf("Sirviendo archivos estáticos desde %s en http://localhost:%s\n", *staticDir, *port)
		err := http.ListenAndServe(":"+*port, nil)
		if err != nil {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}
}
