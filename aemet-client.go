package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ObtenirDadesAEMET realitza la doble petició a l'API d'AEMET amb User-Agent personalitzat
func (app *Config) ObtenirDadesAEMET(urlEndpoint string) ([]byte, error) {
	const userAgent = "EcoHortApp/1.0 (ecohortapp@cibernarium.cat)"

	// 1. PRIMERA PETICIÓ: Sol·licitar l'enllaç de dades a l'API
	req1, err := http.NewRequest("GET", urlEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creant primera petició: %w", err)
	}

	q := req1.URL.Query()
	q.Add("api_key", app.apiKey)
	req1.URL.RawQuery = q.Encode()

	req1.Header.Set("User-Agent", userAgent)
	req1.Header.Set("Accept", "application/json")

	resp1, err := app.HTTPClient.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("error en la petició a la API d'AEMET: %w", err)
	}
	defer resp1.Body.Close()

	var apiRes AemetRespuestaAPI
	if err := json.NewDecoder(resp1.Body).Decode(&apiRes); err != nil {
		return nil, fmt.Errorf("error descodificant resposta 1: %w", err)
	}

	if apiRes.Estado != 200 || apiRes.Datos == "" {
		return nil, fmt.Errorf("AEMET ha retornat l'estat %d: %s", apiRes.Estado, apiRes.Descripcion)
	}

	// 2. SEGONA PETICIÓ: Descarregar el JSON final des del CDN
	req2, err := http.NewRequest("GET", apiRes.Datos, nil)
	if err != nil {
		return nil, fmt.Errorf("error creant segona petició: %w", err)
	}

	req2.Header.Set("User-Agent", userAgent)

	resp2, err := app.HTTPClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("error descarregant les dades del CDN: %w", err)
	}
	defer resp2.Body.Close()

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("error llegint el cos del CDN: %w", err)
	}

	return body, nil
}
