package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . archivo.yaml")
		os.Exit(1)
	}

	conf, err := leerConfiguracion(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	mostrarConfiguracion(conf)
}

func mostrarConfiguracion(conf *Configuracion) {
	fmt.Println("Estados:", strings.Join(conf.Estados.Lista, ", "))
	fmt.Println("Estado inicial:", conf.Estados.Inicial)
	fmt.Println("Estado final:", conf.Estados.Final)
	fmt.Println("Alfabeto:", strings.Join(conf.Alfabeto, ", "))
	fmt.Println("Alfabeto de la cinta:", strings.Join(conf.AlfabetoCinta, ", "))

	fmt.Println("Transiciones:")
	for _, t := range conf.Delta {
		fmt.Printf("  delta([%s,%s], %s) = ([%s,%s], %s, %s)\n",
			t.Entrada.Estado, verSimbolo(t.Entrada.Cache), verSimbolo(t.Entrada.Simbolo),
			t.Salida.Estado, verSimbolo(t.Salida.Cache), verSimbolo(t.Salida.Simbolo), t.Salida.Mover)
	}

	fmt.Println("Cadenas a simular:")
	for _, c := range conf.Cadenas {
		fmt.Println("  " + c)
	}
}
