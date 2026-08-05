package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// ObtenirMunicipiosAEMET descarrega la llista de municipis i assigna la CCAA segons el codi INE
func ObtenirMunicipiosAEMET(apiKey string) ([]Municipio, error) {
	urlMaestro := "https://opendata.aemet.es/opendata/api/maestros/municipios?api_key=" + apiKey

	// Paso 1: Petición a la API maestro de AEMET
	resp, err := http.Get(urlMaestro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp AemetRespuestaAPI
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Estado != 200 {
		return nil, fmt.Errorf("error d'AEMET: %s", apiResp.Descripcion)
	}

	// Paso 2: Descarga del JSON real desde la URL temporal (apiResp.Datos)
	datosResp, err := http.Get(apiResp.Datos)
	if err != nil {
		return nil, err
	}
	defer datosResp.Body.Close()

	body, err := io.ReadAll(datosResp.Body)
	if err != nil {
		return nil, err
	}

	var municipios []Municipio
	if err := json.Unmarshal(body, &municipios); err != nil {
		return nil, err
	}

	// Enriquecer cada municipio con la CCAA basada en los 2 primeros dígitos del código INE
	for i := range municipios {
		codINE := municipios[i].IDOld
		if len(codINE) < 2 {
			codINE = strings.TrimPrefix(municipios[i].ID, "id")
			municipios[i].IDOld = codINE
		}

		if len(codINE) >= 2 {
			prefix := codINE[:2]
			if ccaa, ok := MapaCCAA[prefix]; ok {
				municipios[i].CCAA = ccaa
			}
		}
	}

	return municipios, nil
}
