# Simulador de Máquinas de Turing

Proyecto 3 de Teoría de la Computación. El programa lee la descripción de una máquina de Turing de una cinta desde un archivo yaml, arma la máquina y simula cada cadena del archivo. Para cada cadena muestra todas las descripciones instantáneas, si fue aceptada o rechazada y cómo quedó la cinta al final.

## Vídeo
[https://youtu.be/YeydxQvdJlI](https://youtu.be/YeydxQvdJlI)

## Requisitos

- Go 1.22 o más reciente.

## Cómo usarlo

```bash
go run . maquinas/reconocedora.yaml
go run . maquinas/alteradora.yaml
```

También se puede compilar una vez y usar el ejecutable:

```bash
go build
./proyecto3-tdlc maquinas/reconocedora.yaml
```

Como la salida puede ser larga, se puede guardar en un archivo:

```bash
go run . maquinas/alteradora.yaml > salida.txt
```

## Formato del archivo

El archivo sigue la estructura del enunciado: `q_states`, `alphabet`, `tape_alphabet`, `delta` y `simulation_strings`. Algunos detalles:

- El blanco se escribe dejando el campo vacío. En la salida el blanco se muestra como `B`, por eso la letra `B` no se puede usar como símbolo.
- El `#` va entre comillas (`"#"`), si no yaml lo toma como comentario.
- Cada símbolo es de un solo caracter y las cadenas se leen letra por letra.
- El alfabeto de la cinta siempre incluye al alfabeto de entrada, aunque no se repita en el archivo.
- Los desplazamientos son `L` (izquierda), `R` (derecha) y `S` (se queda en su lugar).
- Si hay un campo mal escrito, un estado que no existe, un símbolo que no está en el alfabeto, un desplazamiento inválido o dos transiciones para el mismo caso, el programa dice cuál es el problema y no simula nada.
- Si una cadena tiene símbolos que no están en el alfabeto de entrada, se marca como rechazada sin simularla.

## Cómo leer la salida

El estado de la máquina es un par `[estado, caché]`: el estado de siempre más un símbolo guardado en memoria. Así la máquina puede recordar un símbolo sin necesitar más estados. Las transiciones se muestran como:

```
delta([3,a], #) = ([4,a], #, R)
```

Que se lee: estando en el estado 3 con `a` en la caché y leyendo `#`, pasa al estado 4, deja `a` en la caché, escribe `#` y se mueve a la derecha.

En cada descripción instantánea el par `[estado, caché]` va justo antes del símbolo que está leyendo la cabeza. Los blancos de las orillas no se muestran, salvo el que esté bajo la cabeza. Por ejemplo `XaY[2,B]bZc` quiere decir que la máquina está en el estado 2, sin nada en la caché y leyendo la `b`.

- **Aceptada**: la máquina llegó al estado final. Ahí se detiene.
- **Rechazada**: no hay transición para el estado, la caché y el símbolo que se está leyendo, o la máquina pasó de 10000 pasos sin detenerse.

## Arquitectura

```
archivo yaml -> Configuracion -> Maquina -> por cada cadena: Cinta + simular -> salida
```

| Archivo | Qué hace |
|---|---|
| `config.go` | Lee el yaml a una estructura con los mismos campos del archivo. |
| `maquina.go` | Arma la máquina a partir de esa estructura y revisa que todo sea válido. Las transiciones se guardan en un mapa: la llave es (estado, caché, símbolo leído) y el valor es (nuevo estado, nueva caché, símbolo a escribir, desplazamiento). Así encontrar la transición de cada paso es directo. |
| `cinta.go` | La cinta es una lista de celdas y la posición de la cabeza. Si la cabeza se sale por la izquierda o por la derecha se agrega un blanco, así la cinta es infinita hacia los dos lados. También arma la descripción instantánea. |
| `simulacion.go` | El ciclo de la simulación: busca la transición, escribe, mueve la cabeza, cambia el estado y la caché, y guarda la descripción de cada paso. |
| `main.go` | Lee el archivo, muestra la máquina y simula cada cadena. |

## Máquina reconocedora: a^n b^n c^n

Archivo: `maquinas/reconocedora.yaml`. Reconoce cadenas con la misma cantidad de `a`, `b` y `c`, en ese orden y con n >= 1. Este lenguaje no es libre de contexto, así que un autómata de pila no lo puede reconocer. En cada vuelta la máquina cambia una `a` por `X`, una `b` por `Y` y una `c` por `Z`.

| Estado | Qué hace |
|---|---|
| 0 | Marca la primera `a` con `X`. Si ya no quedan `a`, pasa a revisar el final. |
| 1 | Avanza hasta la primera `b` sin marcar y la cambia por `Y`. |
| 2 | Avanza hasta la primera `c` sin marcar y la cambia por `Z`. |
| 3 | Regresa a la izquierda hasta la última `X`. |
| 4 | Revisa que solo queden `Y` y `Z` hasta el blanco. |
| 5 | Estado final. |

```
delta([0,B], a) = ([1,B], X, R)
delta([0,B], Y) = ([4,B], Y, R)
delta([1,B], a) = ([1,B], a, R)
delta([1,B], Y) = ([1,B], Y, R)
delta([1,B], b) = ([2,B], Y, R)
delta([2,B], b) = ([2,B], b, R)
delta([2,B], Z) = ([2,B], Z, R)
delta([2,B], c) = ([3,B], Z, L)
delta([3,B], a) = ([3,B], a, L)
delta([3,B], b) = ([3,B], b, L)
delta([3,B], Y) = ([3,B], Y, L)
delta([3,B], Z) = ([3,B], Z, L)
delta([3,B], X) = ([0,B], X, R)
delta([4,B], Y) = ([4,B], Y, R)
delta([4,B], Z) = ([4,B], Z, R)
delta([4,B], B) = ([5,B], B, S)
```

| Cadena | Resultado |
|---|---|
| `aabbcc` | Aceptada |
| `aaabbbccc` | Aceptada |
| `aabbbcc` | Rechazada, sobra una `b` |
| `aabcbc` | Rechazada, las letras no van en orden |

## Máquina alteradora: copiar w en w#w

Archivo: `maquinas/alteradora.yaml`. Recibe una cadena w de `a` y `b` y la deja como `w#w`. Aquí sí se usa la caché: cada letra se marca (`a` con `X`, `b` con `Y`), se guarda en la caché y se escribe al final de la copia. Al terminar, las marcas se regresan a su letra original.

| Estado | Qué hace |
|---|---|
| 0 | Va hasta el final de la cadena y escribe el `#`. |
| 1 | Regresa al inicio de la cadena. |
| 2 | Marca la siguiente letra y la guarda en la caché. Si llega al `#` ya terminó de copiar. |
| 3 | Con la letra en la caché, avanza hasta el `#`. |
| 4 | Avanza hasta el final de la copia y escribe la letra de la caché. |
| 5 | Regresa hasta la última `X` o `Y`. |
| 6 | Cambia las `X` por `a` y las `Y` por `b`. |
| 7 | Estado final. |

```
delta([0,B], a) = ([0,B], a, R)
delta([0,B], b) = ([0,B], b, R)
delta([0,B], B) = ([1,B], #, L)
delta([1,B], a) = ([1,B], a, L)
delta([1,B], b) = ([1,B], b, L)
delta([1,B], B) = ([2,B], B, R)
delta([2,B], a) = ([3,a], X, R)
delta([2,B], b) = ([3,b], Y, R)
delta([2,B], #) = ([6,B], #, L)
delta([3,a], a) = ([3,a], a, R)
delta([3,a], b) = ([3,a], b, R)
delta([3,a], #) = ([4,a], #, R)
delta([3,b], a) = ([3,b], a, R)
delta([3,b], b) = ([3,b], b, R)
delta([3,b], #) = ([4,b], #, R)
delta([4,a], a) = ([4,a], a, R)
delta([4,a], b) = ([4,a], b, R)
delta([4,a], B) = ([5,B], a, L)
delta([4,b], a) = ([4,b], a, R)
delta([4,b], b) = ([4,b], b, R)
delta([4,b], B) = ([5,B], b, L)
delta([5,B], a) = ([5,B], a, L)
delta([5,B], b) = ([5,B], b, L)
delta([5,B], #) = ([5,B], #, L)
delta([5,B], X) = ([2,B], X, R)
delta([5,B], Y) = ([2,B], Y, R)
delta([6,B], X) = ([6,B], a, L)
delta([6,B], Y) = ([6,B], b, L)
delta([6,B], B) = ([7,B], B, R)
```

| Cadena | Cinta final |
|---|---|
| `abbab` | `abbab#abbab` |
| `aaaaa` | `aaaaa#aaaaa` |
| `babbaa` | `babbaa#babbaa` |
| `bbababa` | `bbababa#bbababa` |
