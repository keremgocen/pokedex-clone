package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pokedex-clone/pkg/storage"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APICache   *storage.Store
	// todo retry/backoff
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

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

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
