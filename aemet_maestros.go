package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Mapa de relació entre els 2 primers dígits del codi INE i la Comunitat Autònoma
var MapaCCAA = map[string]string{
	"01": "País Vasco", "02": "Castilla-La Mancha", "03": "Comunitat Valenciana", "04": "Andalucía",
	"05": "Castilla y León", "06": "Extremadura", "07": "Illes Balears", "08": "Catalunya",
	"09": "Castilla y León", "10": "Extremadura", "11": "Andalucía", "12": "Comunitat Valenciana",
	"13": "Castilla-La Mancha", "14": "Andalucía", "15": "Galicia", "16": "Castilla-La Mancha",
	"17": "Catalunya", "18": "Andalucía", "19": "Castilla-La Mancha", "20": "País Vasco",
	"21": "Andalucía", "22": "Aragón", "23": "Andalucía", "24": "Castilla y León",
	"25": "Catalunya", "26": "La Rioja", "27": "Galicia", "28": "Comunidad de Madrid",
	"29": "Andalucía", "30": "Región de Murcia", "31": "Comunidad Foral de Navarra", "32": "Galicia",
	"33": "Principado de Asturias", "34": "Castilla y León", "35": "Canarias", "36": "Galicia",
	"37": "Castilla y León", "38": "Canarias", "39": "Cantabria", "40": "Castilla y León",
	"41": "Andalucía", "42": "Castilla y León", "43": "Catalunya", "44": "Aragón",
	"45": "Castilla-La Mancha", "46": "Comunitat Valenciana", "47": "Castilla y León", "48": "País Vasco",
	"49": "Castilla y León", "50": "Aragón", "51": "Ceuta", "52": "Melilla",
}

// ObtenirMunicipiosAEMET és un mètode de (*Config)
func (app *Config) ObtenirMunicipiosAEMET() ([]Municipio, error) {
	// 1. Cerca de la clave con fallback: memoria app -> preferencies Fyne -> variable d'entorn
	keyToUse := strings.TrimSpace(app.apiKey)
	if keyToUse == "" && app.App != nil {
		keyToUse = strings.TrimSpace(app.App.Preferences().StringWithFallback("aemet_api_key", ""))
	}
	if keyToUse == "" {
		keyToUse = strings.TrimSpace(os.Getenv("AEMET_API_KEY"))
	}

	if keyToUse == "" {
		return nil, fmt.Errorf("la API Key d'AEMET està buida. Configureu-la a la pestanya Mode PRO")
	}

	// 2. Paso 1: Petición al endpoint maestro
	urlMaestro := "https://opendata.aemet.es/opendata/api/maestro/municipios?api_key=" + keyToUse
	log.Println("[AEMET LOG] Paso 1 - Consultant endpoint mestre:", urlMaestro)

	resp, err := app.HTTPClient.Get(urlMaestro)
	if err != nil {
		return nil, fmt.Errorf("error de connexió al Paso 1: %w", err)
	}
	defer resp.Body.Close()

	log.Println("[AEMET LOG] Paso 1 - Status HTTP:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resposta HTTP no vàlida d'AEMET: %d", resp.StatusCode)
	}

	var apiResp AemetRespuestaAPI
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error descodificant JSON del Paso 1: %w", err)
	}

	if apiResp.Estado != 200 {
		return nil, fmt.Errorf("AEMET ha retornat l'estat %d: %s", apiResp.Estado, apiResp.Descripcion)
	}

	// 3. Paso 2: Descargar el JSON final desde la URL firmada
	log.Println("[AEMET LOG] Paso 2 - Descarregant municipis des de:", apiResp.Datos)

	datosResp, err := app.HTTPClient.Get(apiResp.Datos)
	if err != nil {
		return nil, fmt.Errorf("error de connexió al Paso 2: %w", err)
	}
	defer datosResp.Body.Close()

	log.Println("[AEMET LOG] Paso 2 - Status HTTP:", datosResp.StatusCode)

	if datosResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error descarregant fitxer de municipis (HTTP %d)", datosResp.StatusCode)
	}

	// Conversión de ISO-8859-1 (Latin-1) a UTF-8 para corregir acentos y caracteres especiales
	utf8Reader := transform.NewReader(datosResp.Body, charmap.ISO8859_1.NewDecoder())
	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("error llegint el cos de dades del Paso 2: %w", err)
	}

	var municipios []Municipio
	if err := json.Unmarshal(body, &municipios); err != nil {
		return nil, fmt.Errorf("error processant la llista de municipis: %w", err)
	}

	// 4. Mapear CCAA a cada municipio usando los 2 primeros dígitos de IDOld (Código INE)
	for i := range municipios {
		if len(municipios[i].IDOld) >= 2 {
			codiProv := municipios[i].IDOld[:2]
			if ccaa, ok := MapaCCAA[codiProv]; ok {
				municipios[i].CCAA = ccaa
			}
		}
	}

	return municipios, nil
}
