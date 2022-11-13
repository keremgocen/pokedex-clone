package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"pokedex-clone/internal/storage"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APICache   *storage.Store
}

func NewClient(url string, timeout time.Duration, cache *storage.Store) *Client {
	return &Client{
		BaseURL: url,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		APICache: cache,
	}
}

func (c *Client) GetSpecies(ctx context.Context, name string) (*PokemonSpecies, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.BaseURL, name), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res PokemonSpecies
	if res, ok := c.APICache.Load(name); ok {
		log.Println("returning cached", name)
		return res.(*PokemonSpecies), nil
	}

	if reqErr := c.sendRequest(req, &res); reqErr != nil {
		return nil, reqErr
	}

	if cacheErr := c.APICache.Save(name, &res); cacheErr != nil {
		log.Printf("failed to save %s in cache: [%v]", name, cacheErr.Error())
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
