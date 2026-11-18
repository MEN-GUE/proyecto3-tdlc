package main

import "strings"

// Las celdas vacias ("") son blancos. La cinta crece hacia donde haga falta.
type Cinta struct {
	celdas []string
	cabeza int
}

func nuevaCinta(cadena string) *Cinta {
	c := &Cinta{}
	for _, letra := range cadena {
		c.celdas = append(c.celdas, string(letra))
	}
	if len(c.celdas) == 0 {
		c.celdas = []string{""}
	}
	return c
}

func (c *Cinta) leer() string {
	return c.celdas[c.cabeza]
}

func (c *Cinta) escribir(s string) {
	c.celdas[c.cabeza] = s
}

func (c *Cinta) mover(direccion string) {
	switch direccion {
	case "R":
		c.cabeza++
		if c.cabeza == len(c.celdas) {
			c.celdas = append(c.celdas, "")
		}
	case "L":
		if c.cabeza == 0 {
			c.celdas = append([]string{""}, c.celdas...)
		} else {
			c.cabeza--
		}
	}
}

// Muestra la cinta sin los blancos de las orillas, con [estado,cache] antes del simbolo que se esta leyendo.
func (c *Cinta) descripcion(estado, cache string) string {
	inicio, fin := c.cabeza, c.cabeza
	for i, s := range c.celdas {
		if s != "" {
			inicio = min(inicio, i)
			fin = max(fin, i)
		}
	}

	var texto strings.Builder
	for i := inicio; i <= fin; i++ {
		if i == c.cabeza {
			texto.WriteString("[" + estado + "," + verSimbolo(cache) + "]")
		}
		texto.WriteString(verSimbolo(c.celdas[i]))
	}
	return texto.String()
}

func (c *Cinta) contenido() string {
	texto := ""
	for _, s := range c.celdas {
		texto += verSimbolo(s)
	}
	return strings.Trim(texto, "B")
}
