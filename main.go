package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	staticDir := flag.String("static", "./static", "Ruta a los archivos estáticos (html, css, js)")
	port := flag.String("port", "8080", "Puerto para el servidor HTTP")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Verificar si la ruta existe
		if _, err := os.Stat(*staticDir); os.IsNotExist(err) {
			fmt.Printf("La ruta %s no existe.\n", *staticDir)
			fmt.Print("Por favor, ingrese una ruta válida para los archivos estáticos (html, css, js): ")
			input, _ := reader.ReadString('\n')
			*staticDir = strings.TrimSpace(input)
			continue
		}

		// Verificar que exista index.html en la ruta
		indexPath := filepath.Join(*staticDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			fmt.Printf("No se encontró index.html en la ruta %s.\n", *staticDir)
			fmt.Print("Por favor, ingrese una ruta válida para los archivos estáticos (html, css, js): ")
			input, _ := reader.ReadString('\n')
			*staticDir = strings.TrimSpace(input)
			continue
		}

		break
	}

	fs := http.FileServer(http.Dir(*staticDir))
	http.Handle("/", fs)

	fmt.Printf("Sirviendo archivos estáticos desde %s en http://localhost:%s\n", *staticDir, *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
