package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ecohortapp/internal/domain"
	"ecohortapp/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	var myApp domain.Config

	// 1. Carregar la configuració de l'usuari desde la infraestructura
	userCfg, err := repository.LoadConfig()
	if err != nil {
		log.Println("No s'ha trobat config.json, usant configuració per defecte:", err)
		userCfg = domain.UserConfig{
			MunicipioCodigo: "08001",
			Municipios:      []string{"08001"},
		}
	} else {
		log.Println("Configuració carregada amb èxit:", userCfg)
	}

	// 2. Crear l'aplicació Fyne
	fyneApp := app.NewWithID("cat.cibernarium.ecohortapp")
	myApp.App = fyneApp

	// 3. Crear archivo de logs en disco i configuració de loggers des de la infraestructura
	logFile, infoLog, errorLog := repository.InicialitzarLoggers()
	myApp.InfoLog = infoLog
	myApp.ErrorLog = errorLog

	// 4. Connexió i inicialització de la base de dades local mitjançant infraestructura
	sqlDB, err := repository.ConnectSQL(fyneApp.Storage().RootURI().Path())
	if err != nil {
		myApp.ErrorLog.Panic(err)
	}

	// Acoblament del repositori de dades a la configuració del domini
	myApp.DB = repository.NewSQLiteRepository(sqlDB)
	if err := myApp.DB.Migrate(); err != nil {
		myApp.ErrorLog.Panic(err)
	}

	// 5. Inicialitzar el client HTTP amb timeout de seguretat contra talls de xarxa
	myApp.HTTPClient = http.Client{
		Timeout: 15 * time.Second,
	}

	// 6. Assignar preferències i claus d'accés API a les estructures públiques exportades
	myApp.UserConfig = userCfg
	myApp.Municipi = userCfg.MunicipioCodigo
	myApp.ApiKey = fyneApp.Preferences().StringWithFallback("aemet_api_key", os.Getenv("AEMET_API_KEY"))

	// 7. Càrrega silenciosa tolerant a fallades en segon pla per evitar bloquejos
	go func() {
		if m, err := repository.ObtenirMunicipiosAEMET(&myApp); err == nil {
			myApp.Municipis = m
			myApp.InfoLog.Println("[INIT] Municipis carregats automàticament a la memòria RAM.")
		} else {
			myApp.ErrorLog.Println("[INIT] No s'han pogut carregar els municipis en segon pla (reintentant automàticament al connectar):", err)
		}
	}()

	// 8. Crear la finestra principal de visualització de Fyne
	myApp.MainWindow = fyneApp.NewWindow("Eco Hort App")
	myApp.MainWindow.Resize(fyne.NewSize(800, 700))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	// Interfície visual provisional fins traslladar el fitxer ui.go final
	// myApp.makeUI()

	// 9. Executar l'aplicació en el fil principal del sistema operatiu
	myApp.MainWindow.ShowAndRun()
}
