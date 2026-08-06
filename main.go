package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	var myApp Config

	// 1. Carregar la configuració de l'usuari
	userCfg, err := LoadConfig()
	if err != nil {
		log.Println("No s'ha trobat config.json, usant configuració per defecte:", err)
		userCfg = UserConfig{
			MunicipioCodigo: "08001",
			Municipios:      []string{"08001"},
		}
	} else {
		log.Println("Configuració carregada amb èxit:", userCfg)
	}

	// 2. Crear l'aplicació Fyne
	fyneApp := app.NewWithID("cat.cibernarium.ecohortapp")
	myApp.App = fyneApp

	// 3. Crear logs de control
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Lshortfile)

	// 4. Connexió i configuració de la base de dades
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}
	myApp.setupDB(sqlDB)

	// 5. Inicialitzar el client HTTP amb timeout
	// Es el tiempo que espera la aplicación a que
	// la API responda antes de cancelar la conexión
	myApp.HTTPClient = http.Client{
		Timeout: 15 * time.Second,
	}

	// 6. Assignar configuració i API Key
	myApp.UserConfig = userCfg
	myApp.municipi = userCfg.MunicipioCodigo
	myApp.apiKey = fyneApp.Preferences().StringWithFallback("aemet_api_key", os.Getenv("AEMET_API_KEY"))

	// 7. Crear la finestra principal
	myApp.MainWindow = fyneApp.NewWindow("Eco Hort App")
	myApp.MainWindow.Resize(fyne.NewSize(800, 700))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	// 8. Executar l'aplicació
	myApp.MainWindow.ShowAndRun()
}
