package ui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url" // <-- ASEGÚRATE DE QUE ESTA LÍNEA ESTÉ AQUÍ metida
	"os"
	"strings"

	"ecohortapp/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// PronosticTab genera la pestanya de pronòstic adaptada a arquitectura neta
func PronosticTab(app *domain.Config) *fyne.Container {
	banner := CrearBannerAlerta(app, "verde", "")
	grafic := ObtenirGrafic(app)

	graficContainer := container.NewVBox(banner, grafic)
	scrollContainer := container.NewVScroll(graficContainer)

	app.PronosticGraficContainer = graficContainer

	return container.NewBorder(nil, nil, nil, nil, scrollContainer)
}

// ObtenirGrafic gestiona les coordenades i descarrega el climagrama de Meteoblue
func ObtenirGrafic(app *domain.Config) *canvas.Image {
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
		"https://meteoblue.com",
		lat,
		lon,
		apiKey,
		url.QueryEscape(nomMunicipi),
	)

	var img *canvas.Image

	err := DescarregarArxiu(app, apiURL, "pronostic.png")
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

// ActualitzarGraficPronostic refresca el bàner de la AEMET i la imatge de fons
func ActualitzarGraficPronostic(app *domain.Config) {
	if app.PronosticGraficContainer == nil {
		return
	}

	codi := app.UserConfig.MunicipioCodigo
	if codi == "" {
		codi = "esp"
	}

	// NOTA: Revisa que ObtenirAlertaActual estiga exportada correctament al teu repository o usecase
	// Si encara no està moguda de la rel, l'analitzador la marcarà com temporalment absent.
	var nivel, msg string

	banner := CrearBannerAlerta(app, nivel, msg)
	nouGrafic := ObtenirGrafic(app)

	app.PronosticGraficContainer.Objects = []fyne.CanvasObject{banner, nouGrafic}
	app.PronosticGraficContainer.Refresh()
}

// DescarregarArxiu realitza la petició física de descàrrega cap al disc dur de Windows
func DescarregarArxiu(app *domain.Config, URL string, nomArxiu string) error {
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
