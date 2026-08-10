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
		// Error de xarxa, corte de Wi-Fi o timeout
		return nil, fmt.Errorf("error de connexió a la API d'AEMET (comprovi el Wi-Fi): %w", err)
	}
	defer resp1.Body.Close()

	// Control directe de codis d'estat HTTP de la primera petició
	if resp1.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("AEMET ha aconseguit el límit de peticions (HTTP 429). Torni a intentar-ho en un minut")
	} else if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AEMET ha retornat el codi d'estat HTTP %d", resp1.StatusCode)
	}

	var apiRes AemetRespuestaAPI
	if err := json.NewDecoder(resp1.Body).Decode(&apiRes); err != nil {
		return nil, fmt.Errorf("error descodificant resposta JSON 1: %w", err)
	}

	// Comprovació de l'estat intern retornat en la resposta JSON d'AEMET
	if apiRes.Estado != 200 || apiRes.Datos == "" {
		if apiRes.Estado == 429 {
			return nil, fmt.Errorf("AEMET límit de peticions assolit (Estado 429): %s", apiRes.Descripcion)
		}
		return nil, fmt.Errorf("AEMET ha retornat l'estat %d: %s", apiRes.Estado, apiRes.Descripcion)
	}

	// 2. SEGONA PETICIÓ: Descarregar el JSON final des del CDN
	req2, err := http.NewRequest("GET", apiRes.Datos, nil)
	if err != nil {
		return nil, fmt.Errorf("error creant segona petició al CDN: %w", err)
	}

	req2.Header.Set("User-Agent", userAgent)

	resp2, err := app.HTTPClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("error descarregant les dades del CDN (error de xarxa): %w", err)
	}
	defer resp2.Body.Close()

	// Control de codi d'estat HTTP en la segona petició
	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el CDN d'AEMET ha retornat el codi HTTP %d", resp2.StatusCode)
	}

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("error llegint el cos del CDN: %w", err)
	}

	return body, nil
}
