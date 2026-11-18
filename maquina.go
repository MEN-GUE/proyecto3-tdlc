package main

import (
	"fmt"
	"slices"
	"unicode/utf8"
)

type Clave struct {
	Estado string
	Cache  string
	Lee    string
}

type Accion struct {
	Estado  string
	Cache   string
	Escribe string
	Mover   string
}

type Maquina struct {
	Estados       []string
	Inicial       string
	Final         string
	Alfabeto      []string
	AlfabetoCinta []string
	Transiciones  map[Clave]Accion
	orden         []Clave
}

func construirMaquina(conf *Configuracion) (*Maquina, error) {
	m := &Maquina{
		Estados:      conf.Estados.Lista,
		Inicial:      conf.Estados.Inicial,
		Final:        conf.Estados.Final,
		Transiciones: map[Clave]Accion{},
	}

	if len(m.Estados) == 0 {
		return nil, fmt.Errorf("la lista de estados esta vacia")
	}
	if !slices.Contains(m.Estados, m.Inicial) {
		return nil, fmt.Errorf("el estado inicial %q no esta en la lista de estados", m.Inicial)
	}
	if !slices.Contains(m.Estados, m.Final) {
		return nil, fmt.Errorf("el estado final %q no esta en la lista de estados", m.Final)
	}

	if len(conf.Alfabeto) == 0 {
		return nil, fmt.Errorf("el alfabeto esta vacio")
	}
	for _, s := range conf.Alfabeto {
		if err := revisarSimbolo(s); err != nil {
			return nil, err
		}
		if !slices.Contains(m.Alfabeto, s) {
			m.Alfabeto = append(m.Alfabeto, s)
		}
	}

	// La cinta siempre incluye el alfabeto de entrada aunque el yaml no lo ponga
	m.AlfabetoCinta = slices.Clone(m.Alfabeto)
	for _, s := range conf.AlfabetoCinta {
		if s == "" || slices.Contains(m.AlfabetoCinta, s) {
			continue
		}
		if err := revisarSimbolo(s); err != nil {
			return nil, err
		}
		m.AlfabetoCinta = append(m.AlfabetoCinta, s)
	}

	for i, t := range conf.Delta {
		n := i + 1
		clave := Clave{t.Entrada.Estado, t.Entrada.Cache, t.Entrada.Simbolo}
		accion := Accion{t.Salida.Estado, t.Salida.Cache, t.Salida.Simbolo, t.Salida.Mover}

		if !slices.Contains(m.Estados, clave.Estado) {
			return nil, fmt.Errorf("transicion %d: el estado %q no existe", n, clave.Estado)
		}
		if !slices.Contains(m.Estados, accion.Estado) {
			return nil, fmt.Errorf("transicion %d: el estado %q no existe", n, accion.Estado)
		}

		for _, s := range []string{clave.Cache, clave.Lee, accion.Cache, accion.Escribe} {
			if s == "B" {
				return nil, fmt.Errorf("transicion %d: para el blanco hay que dejar el campo vacio, no poner B", n)
			}
			if s != "" && !slices.Contains(m.AlfabetoCinta, s) {
				return nil, fmt.Errorf("transicion %d: el simbolo %q no esta en el alfabeto de la cinta", n, s)
			}
		}

		if accion.Mover != "L" && accion.Mover != "R" && accion.Mover != "S" {
			return nil, fmt.Errorf("transicion %d: el desplazamiento debe ser L, R o S y no %q", n, accion.Mover)
		}

		if _, existe := m.Transiciones[clave]; existe {
			return nil, fmt.Errorf("transicion %d: ya hay otra transicion para (%s)", n, clave.texto())
		}
		m.Transiciones[clave] = accion
		m.orden = append(m.orden, clave)
	}

	return m, nil
}

func revisarSimbolo(s string) error {
	if utf8.RuneCountInString(s) != 1 {
		return fmt.Errorf("el simbolo %q debe ser de un solo caracter", s)
	}
	if s == "B" {
		return fmt.Errorf("el simbolo B no se puede usar porque asi se muestra el blanco")
	}
	return nil
}

func (c Clave) texto() string {
	return fmt.Sprintf("[%s,%s], %s", c.Estado, verSimbolo(c.Cache), verSimbolo(c.Lee))
}

func (a Accion) texto() string {
	return fmt.Sprintf("[%s,%s], %s, %s", a.Estado, verSimbolo(a.Cache), verSimbolo(a.Escribe), a.Mover)
}
