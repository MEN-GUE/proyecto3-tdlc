package main

import "fmt"

// Para no quedarnos esperando si la maquina nunca se detiene
const limitePasos = 10000

type Resultado struct {
	Descripciones []string
	Aceptada      bool
	Motivo        string
	Cinta         string
}

func (m *Maquina) simular(cadena string) Resultado {
	cinta := nuevaCinta(cadena)
	estado := m.Inicial
	cache := ""

	var r Resultado
	r.Descripciones = append(r.Descripciones, cinta.descripcion(estado, cache))

	for {
		if estado == m.Final {
			r.Aceptada = true
			r.Motivo = "llego al estado final " + m.Final
			break
		}
		if len(r.Descripciones) > limitePasos {
			r.Motivo = fmt.Sprintf("la maquina no se detuvo despues de %d pasos", limitePasos)
			break
		}

		clave := Clave{estado, cache, cinta.leer()}
		accion, existe := m.Transiciones[clave]
		if !existe {
			r.Motivo = fmt.Sprintf("no hay transicion para (%s)", clave.texto())
			break
		}

		cinta.escribir(accion.Escribe)
		cinta.mover(accion.Mover)
		estado = accion.Estado
		cache = accion.Cache
		r.Descripciones = append(r.Descripciones, cinta.descripcion(estado, cache))
	}

	r.Cinta = cinta.contenido()
	return r
}
