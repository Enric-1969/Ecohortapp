package main

import (
	"database/sql"
	"ecohortapp/repository"
	"encoding/json" // <-- Para descodificar el JSON de AEMET
	"fmt"           // <-- Para formatear los errores (fmt.Errorf)
	"io"            // <-- Para leer el cuerpo de la respuesta (io.ReadAll)
	"log"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite"
)

// Crearem un struct amb totes les configuracions que necessiti la nostre App
type Config struct {
	App                                fyne.App              //Definim que emprara Fyne per construir la GUI de l'App
	InfoLog                            *log.Logger           //Definim un Log d'accions
	ErrorLog                           *log.Logger           //Definim un Log d'errors
	DB                                 repository.Repository //Definim la referencia a la conexió a SQLite
	MainWindow                         fyne.Window           //Aqui enmagatzemem la referencia a certes arees de la ui per controlar les actualitzacions de les mateixes.
	ClimaDadesContainer                *fyne.Container       //Guardem el contenidor de les dades del clima, referenciant el punter de memòria del contenidor de fyne.
	PronosticGraficContainer           *fyne.Container       //Definim un camp a on enmagatzem el contenidor del gràfic de clima, que ara sera de tipus contenidor fyne
	Registres                          [][]interface{}       //Per emmagatzemar el slice de slices en forma de interfície a on esta contingut les dades obtingudes de la bd
	RegistresTable                     *widget.Table         //Per emmagatzemar la referencia al punter que correspon el widget de la Taula.
	HTTPClient                         http.Client           //Afegim la referència al client http sence necessitat de invocar la llibreria
	AfegirRegistresDataRegistreEntrada *widget.Entry         //Afegim la referencia a la entrada del valor data registre per a nous registres que guardem en la bd
	AfegirRegistresPrecipitacioEntrada *widget.Entry         //Afegim la referencia a la entrada del valor precipitacio per a nous registres que guardem en la bd
	AfegirRegistresTempMaximaEntrada   *widget.Entry         //Afegim la referencia a la entrada del valor tempMaxima per a nous registres que guardem en la bd
	AfegirRegistresTempMinimaEntrada   *widget.Entry         //Afegim la referencia a la entrada del valor tempMinima per a nous registres que guardem en la bd
	AfegirRegistresHumitatEntrada      *widget.Entry         //Afegim la referencia a la entrada del valor humitat per a nous registres que guardem en la bd
	municipi                           string                //Afegim la referencia a aquest valor de configuració
	apiKey                             string                //Afegim la referencia a aquest valor de configuració
}

// Estructura per descodificar la primera resposta JSON de la API d'AEMET
type AemetRespuestaAPI struct {
	Descripcion string `json:"descripcion"`
	Estado      int    `json:"estado"`
	Datos       string `json:"datos"`
	Metadatos   string `json:"metadatos"`
}

func main() {
	var myApp Config //Creem una variable que sigui de tipus Config i aixi enmagatzemar la configuració de l'App

	// crearem una aplicació fyne
	fyneApp := app.NewWithID("cat.cibernarium.ecohortapp") //El definit el mètode New amb una id ens permet distribuir la nostre app en un MarketPlace
	myApp.App = fyneApp

	//crearem els nostres logs
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)        //Creem un Log per els registres informatius
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Lshortfile) //Creraem un log per els registres d'error

	//conexió amb la base de dades
	sqlDB, err := myApp.connectSQL() //Invoquem la funció d'establiment de la conexió
	if err != nil {
		//Recordem que Panic es l'equivalent a Print pero acompanyat a una crida a Panic
		log.Panic(err)
	}

	//crearem un repositori de base de dades
	myApp.setupDB(sqlDB)

	// INICIALIZAR EL CLIENTE HTTP CON TIMEOUT
	myApp.HTTPClient = http.Client{
		Timeout: 15 * time.Second,
	}

	//Definim la capacitat de que l'usuari modifiqui el municipi i la apiKey
	myApp.municipi = fyneApp.Preferences().StringWithFallback("municipi", "08001")

	//Definim la apiKey
	//os.Getenv("AEMET_API_KEY"):
	//En lugar de dejar un texto fijo con la clave antigua dentro del código Go, le decimos a Fyne:
	//"Si el usuario no tiene guardada una clave en sus preferencias, usa como valor por defecto la
	//clave que guardamos en el archivo .env".

	myApp.apiKey = fyneApp.Preferences().StringWithFallback("apiKey", os.Getenv("AEMET_API_KEY"))

	//crearem i definim el tamany de una pantalla de fyne
	myApp.MainWindow = fyneApp.NewWindow("Eco Hort App")
	myApp.MainWindow.Resize(fyne.NewSize(800, 700)) //Definim el tamany de la finestra
	myApp.MainWindow.SetFixedSize(true)             //Definim que tindra un tamany fixe
	myApp.MainWindow.SetMaster()                    //Indiquem que es la pantalla principal. Si tanquem aquesta pantalla la aplicacio finalitza

	myApp.makeUI() //Crearem una invocació a una funció externa que creara la interficié grafica a partir del contingut.

	//mostrar i executar l'aplicació
	myApp.MainWindow.ShowAndRun()
}

// Realitzarem una funció per invocar la conexió a la BD
func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	//Treballarem amb variables d'entorn per establir les configuracions i en aquest cas comprobem si DB_PATH (que fa referencia a la ruta de la db) te un valor o no.
	if os.Getenv("DB_PATH") != "" {
		//En cas de tenir valor, el recupera
		path = os.Getenv("DB_PATH")
	} else {
		//En cas contrari crearà aquest arxiu de Bases de dades dins de la ruta de l'aplicació
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		//Incloem un registre de control al log
		app.InfoLog.Println("db in:", path)
	}

	//A continuació establim la conexió i controlem possibles errors
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Creem un repositori per a la Base de Dades
func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}

// =========================================================================
// MÈTODE PER FER LA DOBLE PETICIÓ A AEMET CON USER-AGENT PERSONALITZAT
// =========================================================================
func (app *Config) ObtenirDadesAEMET(urlEndpoint string) ([]byte, error) {
	const userAgent = "EcoHortApp/1.0 (ecohortapp@cibernarium.cat)"

	// 1. PRIMERA PETICIÓ: Sol·licitar l'enllaç de dades a l'API
	req1, err := http.NewRequest("GET", urlEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creant primera petició: %w", err)
	}

	// Afegim l'API Key i les cabeceres necessàries
	q := req1.URL.Query()
	q.Add("api_key", app.apiKey)
	req1.URL.RawQuery = q.Encode()

	req1.Header.Set("User-Agent", userAgent)
	req1.Header.Set("Accept", "application/json")

	resp1, err := app.HTTPClient.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("error en la petició a la API d'AEMET: %w", err)
	}
	defer resp1.Body.Close()

	var apiRes AemetRespuestaAPI
	if err := json.NewDecoder(resp1.Body).Decode(&apiRes); err != nil {
		return nil, fmt.Errorf("error descodificant resposta 1: %w", err)
	}

	if apiRes.Estado != 200 || apiRes.Datos == "" {
		return nil, fmt.Errorf("AEMET ha retornat l'estat %d: %s", apiRes.Estado, apiRes.Descripcion)
	}

	// 2. SEGONA PETICIÓ: Descarregar el JSON final des del CDN
	req2, err := http.NewRequest("GET", apiRes.Datos, nil)
	if err != nil {
		return nil, fmt.Errorf("error creant segona petició: %w", err)
	}

	// TORNEM A AFECAR EL USER-AGENT AL CDN
	req2.Header.Set("User-Agent", userAgent)

	resp2, err := app.HTTPClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("error descarregant les dades del CDN: %w", err)
	}
	defer resp2.Body.Close()

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("error llegint el cos del CDN: %w", err)
	}

	return body, nil
}
