# pokedex-clone

## How to Run it?

TBD

## Swagger

## Production API assumptions

TBD
- cache invalidation (timeout, LRU etc.)

Fun with Pokemons

The API has two main endpoints:

1. Return basic Pokemon information.
2. Return basic Pokemon information but with a ‘fun’ translation of the Pokemon description.
   Following are more detailed API requirements. Guidelines can be found on Page 3.

### API Requirements

#### Endpoint 1 - Basic Pokemon Information

Given a Pokemon name, returns standard Pokemon description and additional information. Example endpoint:

`/HTTP/GET /pokemon/<pokemon name>`

Example call (using httpie):
`http http://localhost:5000/pokemon/mewtwo`

The API response should contain a minimum of:

- Pokemon’s name
- Pokemon’s standard description
- Pokemon’s habitat
- Pokemon’s is_legendary status Example response:

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

1. If the Pokemon’s habitat is cave or it’s a legendary Pokemon then apply the Yoda translation.
2. For all other Pokemon, apply the Shakespeare translation.
3. If you can’t translate the Pokemon’s description (for whatever reason😉) then use the standard description

Example endpoint:
`HTTP/GET /pokemon/translated/<pokemon name>`

Example call (using httpie):
`http http://localhost:5000/pokemon/translated/mewtwo`

The API response should contain a minimum of:

- Pokemon name
- Translated Pokemon description
- Pokemon’s habitat
- Pokemon’s is_legendary status

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
