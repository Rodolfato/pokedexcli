## Pokedexcli
Little Read–eval–print loop Pokedex built with Go using PokéAPI

### Prerequisites

* Go v1.22.2

### Installation
1. Clone the repo
   ```sh
   git clone https://github.com/Rodolfato/pokedexcli
   ```

2. Build and play around!
   ```sh
   go build && ./pokedexcli
   ```
   
### Usage

These are the available commands:

Progressively list all areas of the Pokemon map

```sh
  Pokedex > map
```

Go back to the last list of areas

```sh
  Pokedex > mapb
```

Explore the Pokemon in an area

```
  Pokedex > explore {area}
```

Try to catch a Pokemon

```
  Pokedex > catch {name}
```
  
 Inspect a Pokemon you've previously caught

```
  Pokedex > inspect {name}
```

List all Pokemon you've caught

```
  Pokedex > pokedex
```

List all available commands

```
  Pokedex > help
```

Exit the program

```
  Pokedex > exit
```


## Acknowledgments

* Project guided by [Boot.dev](https://www.boot.dev/)
* [PokéApi](https://pokeapi.co/) is super fun to work with
