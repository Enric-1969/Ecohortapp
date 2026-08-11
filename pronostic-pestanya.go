package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (app *Config) pronosticTab() *fyne.Container {
	banner := app.crearBannerAlerta("verde", "")
	grafic := app.obtenirGrafic()

	graficContainer := container.NewVBox(banner, grafic)
	scrollContainer := container.NewVScroll(graficContainer)

	app.PronosticGraficContainer = graficContainer

	return container.NewBorder(nil, nil, nil, nil, scrollContainer)
}

func (app *Config) obtenirGrafic() *canvas.Image {
	apiKey := "zpPLdN8ijMAacbGX"

	lat := app.UserConfig.Latitud
	lon := app.UserConfig.Longitud

	if lat == 0 || lon == 0 {
		lat = 41.5161
		lon = 1.9021
	}

	nomMunicipi := app.UserConfig.MunicipioNombre
	if nomMunicipi == "" {
		nomMunicipi = "Abrera"
	}

	if idx := strings.Index(nomMunicipi, "("); idx != -1 {
		nomMunicipi = nomMunicipi[:idx]
	}
	nomMunicipi = strings.TrimSpace(nomMunicipi)

	apiURL := fmt.Sprintf(
		"https://my.meteoblue.com/images/meteogram?lat=%.4f&lon=%.4f&asl=100&tz=Europe%%2FMadrid&apikey=%s&format=png&dpi=72&lang=es&temperature_units=C&precipitation_units=mm&windspeed_units=kmh&location_name=%s",
		lat,
		lon,
		apiKey,
		url.QueryEscape(nomMunicipi),
	)

	var img *canvas.Image

	err := app.descarregarArxiu(apiURL, "pronostic.png")
	if err != nil {
		log.Println("Error descarregant el gràfic de Meteoblue:", err)
		img = canvas.NewImageFromResource(resourceNodisponiblePng)
	} else {
		img = canvas.NewImageFromFile("pronostic.png")
	}

	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(770, 480))

	return img
}

func (app *Config) actualitzarGraficPronostic() {
	if app.PronosticGraficContainer == nil {
		return
	}

	// 1. Leemos el código INE del municipio desde UserConfig
	codi := app.UserConfig.MunicipioCodigo
	if codi == "" {
		codi = "esp"
	}

	// 2. Consultamos la alerta para ese municipio concreto
	nivel, msg := app.ObtenirAlertaActual(codi)

	// 3. Generamos componentes y actualizamos pantalla
	banner := app.crearBannerAlerta(nivel, msg)
	nouGrafic := app.obtenirGrafic()

	app.PronosticGraficContainer.Objects = []fyne.CanvasObject{banner, nouGrafic}
	app.PronosticGraficContainer.Refresh()
}

func (app *Config) descarregarArxiu(URL string, nomArxiu string) error {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	response, err := app.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("error HTTP %d al descarregar la imatge. Resposta de Meteoblue: %s", response.StatusCode, string(b))
	}

	err = os.WriteFile(nomArxiu, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
