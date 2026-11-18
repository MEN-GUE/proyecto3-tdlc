package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuracion struct {
	Estados struct {
		Lista   []string `yaml:"q_list"`
		Inicial string   `yaml:"initial"`
		Final   string   `yaml:"final"`
	} `yaml:"q_states"`
	Alfabeto      []string `yaml:"alphabet"`
	AlfabetoCinta []string `yaml:"tape_alphabet"`
	Delta         []struct {
		Entrada struct {
			Estado  string `yaml:"initial_state"`
			Cache   string `yaml:"mem_cache_value"`
			Simbolo string `yaml:"tape_input"`
		} `yaml:"params"`
		Salida struct {
			Estado  string `yaml:"final_state"`
			Cache   string `yaml:"mem_cache_value"`
			Simbolo string `yaml:"tape_output"`
			Mover   string `yaml:"tape_displacement"`
		} `yaml:"output"`
	} `yaml:"delta"`
	Cadenas []string `yaml:"simulation_strings"`
}

// Un valor vacio en el yaml queda como "" y eso es el blanco.
func leerConfiguracion(ruta string) (*Configuracion, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %v", err)
	}
	defer archivo.Close()

	lector := yaml.NewDecoder(archivo)
	lector.KnownFields(true)

	var conf Configuracion
	if err := lector.Decode(&conf); err != nil {
		return nil, fmt.Errorf("el yaml tiene un error: %v", err)
	}
	return &conf, nil
}

func verSimbolo(s string) string {
	if s == "" {
		return "B"
	}
	return s
}
