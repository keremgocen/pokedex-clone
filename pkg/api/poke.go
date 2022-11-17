package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

//go:generate mockgen -destination mocks/poke.go -package mocks -source poke.go

type PokeAPI interface {
	GetSpecies(ctx context.Context, name string) (*PokemonSpecies, error)
}

type Poke struct {
	Client *Client
	// retry/backoff
}

func (p Poke) GetSpecies(ctx context.Context, name string) (*PokemonSpecies, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", p.Client.BaseURL, name), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res PokemonSpecies
	if res, ok := p.Client.APICache.Load(name); ok {
		log.Println("returning cached", name)
		return res.(*PokemonSpecies), nil
	}

	if reqErr := p.Client.sendRequest(req, &res); reqErr != nil {
		return nil, reqErr
	}

	if cacheErr := p.Client.APICache.Save(name, &res); cacheErr != nil {
		log.Printf("failed to save %s in cache: [%v]", name, cacheErr.Error())
	}

	return &res, nil
}
