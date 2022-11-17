package api

import (
	"context"
	"fmt"
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
	if reqErr := p.Client.sendRequest(req, &res); reqErr != nil {
		return nil, reqErr
	}

	return &res, nil
}
