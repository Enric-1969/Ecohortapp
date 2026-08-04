package main

import (
	"errors"
	"io"
	"os"

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
	apiURL := "https://my.meteoblue.com/visimage/meteogram_web_hd?look=KILOMETER_PER_HOUR%2CCELSIUS%2CMILLIMETER&apikey=5838a18e295d&temperature=C&windspeed=kmh&precipitationamount=mm&winddirection=3char&city=Abrera&iso2=es&lat=41.5168&lon=1.901&asl=111&tz=Europe%2FMadrid&lang=es&sig=b353aab637f77ab97ae54cbd760554f2"

	var img *canvas.Image

	err := app.descarregarArxiu(apiURL, "pronostic.png")
	if err != nil {
		img = canvas.NewImageFromResource(resourceNodisponiblePng)
	} else {
		img = canvas.NewImageFromFile("pronostic.png")
	}

	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(770, 480))

	return img
}

func (app *Config) descarregarArxiu(URL string, nomArxiu string) error {
	response, err := app.HTTPClient.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("rebem un codi de resposta erronia quan descarreguem la imatge")
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
