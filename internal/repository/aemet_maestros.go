package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"ecohortapp/internal/domain"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const municipiosCacheFile = "municipios_cache.json"

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

var MapaProvincias = map[string]string{
	"01": "Araba/Álava", "02": "Albacete", "03": "Alicante/Alacant", "04": "Almería",
	"05": "Ávila", "06": "Badajoz", "07": "Illes Balears", "08": "Barcelona",
	"09": "Burgos", "10": "Cáceres", "11": "Cádiz", "12": "Castellón/Castelló",
	"13": "Ciudad Real", "14": "Córdoba", "15": "A Coruña", "16": "Cuenca",
	"17": "Girona", "18": "Granada", "19": "Guadalajara", "20": "Gipuzkoa",
	"21": "Huelva", "22": "Huesca", "23": "Jaén", "24": "León",
	"25": "Lleida", "26": "La Rioja", "27": "Lugo", "28": "Madrid",
	"29": "Málaga", "30": "Murcia", "31": "Navarra", "32": "Ourense",
	"33": "Asturias", "34": "Palencia", "35": "Las Palmas", "36": "Pontevedra",
	"37": "Salamanca", "38": "Santa Cruz de Tenerife", "39": "Cantabria", "40": "Segovia",
	"41": "Sevilla", "42": "Soria", "43": "Tarragona", "44": "Teruel",
	"45": "Toledo", "46": "Valencia/València", "47": "Valladolid", "48": "Bizkaia",
	"49": "Zamora", "50": "Zaragoza", "51": "Ceuta", "52": "Melilla",
}

func PoblarCCAAiProvincia(municipis []domain.Municipio) {
	for i := range municipis {
		if len(municipis[i].IDOld) >= 2 {
			codiProv := municipis[i].IDOld[:2]
			if ccaa, ok := MapaCCAA[codiProv]; ok {
				municipis[i].CCAA = ccaa
			}
			if prov, ok := MapaProvincias[codiProv]; ok {
				municipis[i].Provincia = prov
			}
		}
	}
}

func ObtenirMunicipiosAEMET(app *domain.Config) ([]domain.Municipio, error) {
	if data, err := os.ReadFile(municipiosCacheFile); err == nil {
		var municipiosCache []domain.Municipio
		if err := json.Unmarshal(data, &municipiosCache); err == nil && len(municipiosCache) > 0 {
			PoblarCCAAiProvincia(municipiosCache)
			log.Println("[AEMET LOG] Carregats i poblats", len(municipiosCache), "municipis des del fitxer local en disc.")
			return municipiosCache, nil
		}
	}

	keyToUse := strings.TrimSpace(app.ApiKey)
	if keyToUse == "" && app.App != nil {
		keyToUse = strings.TrimSpace(app.App.Preferences().StringWithFallback("aemet_api_key", ""))
	}
	if keyToUse == "" {
		keyToUse = strings.TrimSpace(os.Getenv("AEMET_API_KEY"))
	}

	if keyToUse == "" {
		return nil, fmt.Errorf("la API Key d'AEMET està buida. Configureu-la a la pestanya Mode PRO")
	}

	urlMaestro := "https://aemet.es" + keyToUse
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

	var apiResp domain.AemetRespuestaAPI
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error descodificant JSON del Paso 1: %w", err)
	}

	if apiResp.Estado != 200 {
		return nil, fmt.Errorf("AEMET ha retornat l'estat %d: %s", apiResp.Estado, apiResp.Descripcion)
	}

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

	utf8Reader := transform.NewReader(datosResp.Body, charmap.ISO8859_1.NewDecoder())
	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("error llegint el cos de dades del Paso 2: %w", err)
	}

	var municipios []domain.Municipio
	if err := json.Unmarshal(body, &municipios); err != nil {
		return nil, fmt.Errorf("error processant la llista de municipis: %w", err)
	}

	PoblarCCAAiProvincia(municipios)

	if cacheData, err := json.Marshal(municipios); err == nil {
		_ = os.WriteFile(municipiosCacheFile, cacheData, 0644)
		log.Println("[AEMET LOG] Desats", len(municipios), "municipis a la memòria cau local en disc.")
	}

	return municipios, nil
}
