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

	maquina, err := construirMaquina(conf)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	mostrarMaquina(maquina)

	fmt.Println("Cadenas a simular:")
	for _, c := range conf.Cadenas {
		fmt.Println("  " + c)
	}
}

func mostrarMaquina(m *Maquina) {
	fmt.Println("Estados:", strings.Join(m.Estados, ", "))
	fmt.Println("Estado inicial:", m.Inicial)
	fmt.Println("Estado final:", m.Final)
	fmt.Println("Alfabeto:", strings.Join(m.Alfabeto, ", "))
	fmt.Println("Alfabeto de la cinta:", strings.Join(m.AlfabetoCinta, ", ")+", B")

	fmt.Println("Transiciones:")
	for _, c := range m.orden {
		fmt.Printf("  delta(%s) = (%s)\n", c.texto(), m.Transiciones[c].texto())
	}
}
