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
	grafic := app.obtenirGrafic()

	graficContainer := container.NewVBox(grafic)
	scrollContainer := container.NewVScroll(graficContainer)

	app.PronosticGraficContainer = graficContainer

	return container.NewBorder(nil, nil, nil, nil, scrollContainer)
}

func (app *Config) obtenirGrafic() *canvas.Image {
	nomMunicipi := app.UserConfig.MunicipioNombre
	if nomMunicipi == "" {
		nomMunicipi = "Abrera"
	}

	if idx := strings.Index(nomMunicipi, "("); idx != -1 {
		nomMunicipi = nomMunicipi[:idx]
	}
	nomMunicipi = strings.TrimSpace(nomMunicipi)

	cityEscaped := url.QueryEscape(nomMunicipi)

	// URL con formato gráfico alternativo
	apiURL := fmt.Sprintf("https://wttr.in/%s_2pn_lang=es.png", cityEscaped)

	var img *canvas.Image

	err := app.descarregarArxiu(apiURL, "pronostic.png")
	if err != nil {
		log.Println("Error descarregant el gràfic de temps:", err)
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

	nouGrafic := app.obtenirGrafic()
	app.PronosticGraficContainer.Objects = []fyne.CanvasObject{nouGrafic}
	app.PronosticGraficContainer.Refresh()
}

func (app *Config) descarregarArxiu(URL string, nomArxiu string) error {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	response, err := app.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error HTTP %d al descarregar la imatge d'URL: %s", response.StatusCode, URL)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(nomArxiu, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
