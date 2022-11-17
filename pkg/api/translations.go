package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -destination mocks/translations.go -package mocks -source translations.go

type TranslationsAPI interface {
	GetTranslation(ctx context.Context, name, text string, translationType TranslationType) (*TranslateAPIResponse, error)
}

type Translations struct {
	Client *Client
	// retry/backoff config
}

type TranslationType string

const (
	TTypeYoda        TranslationType = "yoda.json"
	TTypeShakespeare TranslationType = "shakespeare.json"
)

func (t Translations) GetTranslation(
	ctx context.Context,
	name, text string,
	translationType TranslationType,
) (*TranslateAPIResponse, error) {
	textToTranslate := &TranslationText{Text: text}
	b, err := json.Marshal(textToTranslate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, t.Client.BaseURL+string(translationType), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res TranslateAPIResponse

	if reqErr := t.Client.sendRequest(req, &res); reqErr != nil {
		return nil, reqErr
	}

	return &res, nil
}
