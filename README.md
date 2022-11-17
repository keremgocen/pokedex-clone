# pokedex-clone

## How to Run it?

### Building the docker image

`-> docker build --tag pokedex-clone .`

### Running and attaching to the docker image locally

`-> docker run -p 5000:5000 --name test pokedex-clone`

## Assumptions

TBD

Fun with Pokemons

### Pokemon API

#### Endpoint 1 - Basic Pokemon Information

Given a Pokemon name, returns standard Pokemon description and additional information.

`/HTTP/GET /pokemon/<pokemon name>`

Example call (using curl):
`curl http://localhost:5000/pokemon/mewtwo`

Example API response:

```
{
 "name": "mewtwo",
 "description": "It was created by a scientist after years of horrific gene
 splicing and DNA engineering experiments.",
 "habitat": "rare",
 "isLegendary": true
}
```

#### Endpoint 2 - Translated Pokemon Description

Given a Pokemon name, return translated Pokemon description and other basic information using the following rules:

- If the Pokemon’s habitat is cave or it’s a legendary Pokemon then Yoda translation is applied.
- For all other Pokemon, Shakespeare translation is applied.
- Otherwise a standard description is returned.

`HTTP/GET /pokemon/translated/<pokemon name>`

Example call (using curl):
`curl http://localhost:5000/pokemon/translated/mewtwo`

Example response:

```
{
 "name": "mewtwo",
 "description": "Created by a scientist after years of horrific gene
 splicing and dna engineering experiments, it was.",
 "habitat": "rare",
 "isLegendary": true
}
```
