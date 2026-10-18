package main

import (
	"log"
	"net/http"
	"time"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/widget"
)

// ==============================================================
// Config centraliza las referencias principales de la aplicación
// ==============================================================

// Crearem un struct amb totes les configuracions que necessiti la nostre App
type Config struct {
	UserConfig                         UserConfig      // [64 bytes] Mapeo de la estructura de configuración interna del usuario
	HTTPClient                         http.Client     // [32 bytes] Afegim la referència al client http sence necessitat de invocar la llibreria
	Registres                          [][]interface{} // [24 bytes] Per emmagatzemar el slice of slices en forma de interfície a on esta contingut les dades obtingudes de la bd
	Municipis                          []Municipio     // [24 bytes] Slice que contiene la lista estructurada de municipios
	App                                fyne.App        // [16 bytes] Definim que emprara Fyne per construir la GUI de l'App
	DB                                 Repository      // [16 bytes] Definim la referencia a la conexió a SQLite
	MainWindow                         fyne.Window     // [16 bytes] Aqui enmagatzemem la referencia a certes arees de la ui per controlar les actualitzacions de les mateixes.
	municipi                           string          // [16 bytes] Afegim la referencia a aquest valor de configuració
	apiKey                             string          // [16 bytes] Afegim la referencia a aquest valor de configuració
	InfoLog                            *log.Logger     // [8 bytes] Definim un Log d'accions
	ErrorLog                           *log.Logger     // [8 bytes] Definim un Log d'errors
	ClimaDadesContainer                *fyne.Container // [8 bytes] Guardem el contenidor de les dades del clima, referenciant el punter de memòria del contenidor de fyne.
	PronosticGraficContainer           *fyne.Container // [8 bytes] Definim un camp a on enmagatzem el contenidor del gràfic de clima, que ara sera de tipus contenidor fyne
	RegistresTable                     *widget.Table   // [8 bytes] Per emmagatzemar la referencia al punter que correspon el widget de la Taula.
	AfegirRegistresDataRegistreEntrada *widget.Entry   // [8 bytes] Afegim la referencia a la entrada del valor data registre per a nous registres que guardem en la bd
	AfegirRegistresPrecipitacioEntrada *widget.Entry   // [8 bytes] Afegim la referencia a la entrada del valor precipitacio per a nous registres que guardem en la bd
	AfegirRegistresTempMaximaEntrada   *widget.Entry   // [8 bytes] Afegim la referencia a la entrada del valor tempMaxima per a nous registres que guardem en la bd
	AfegirRegistresTempMinimaEntrada   *widget.Entry   // [8 bytes] Afegim la referencia a la entrada del valor tempMinima per a nous registres que guardem en la bd
	AfegirRegistresHumitatEntrada      *widget.Entry   // [8 bytes] Afegim la referencia a la entrada del valor humitat per a nous registres que guardem en la bd

	// CLEAN CODE: Estado centralizado de municipios cargados en RAM.
	// Evita el uso de variables globales sueltas y permite compartir la lista
	// entre pestañas de forma totalmente encapsulada y segura.

}

// Estructura per descodificar la primera resposta JSON de la API d'AEMET
type AemetRespuestaAPI struct {
	Descripcion string `json:"descripcion"` // [16 bytes] Cadena de texto informativa
	Datos       string `json:"datos"`       // [16 bytes] URL con los datos del reporte
	Metadatos   string `json:"metadatos"`   // [16 bytes] URL con la estructura de metadatos
	Estado      int    `json:"estado"`      // [8 bytes] Código de estado numérico de la API
}

type Municipio struct {
	ID        string `json:"id"`       // [16 bytes] Ej: "id08001"
	IDOld     string `json:"id_old"`   // [16 bytes] Ej: "08001" (Código INE de 5 dígitos)
	Nombre    string `json:"nombre"`   // [16 bytes] Ej: "Abrera"
	Longitude string `json:"longitud"` // [16 bytes] Coordenada de longitud en texto
	Latitude  string `json:"latitud"`  // [16 bytes] Coordenada de latitud en texto
	Altitude  string `json:"altitud"`  // [16 bytes] Altitud sobre el nivel del mar
	Provincia string `json:"-"`        // [16 bytes] Calculado según los 2 primeros dígitos de IDOld
	CCAA      string `json:"-"`        // [16 bytes] Calculado según la provincia
}

// UserConfig define las preferencias del usuario para config.json
type UserConfig struct {
	Municipios      []string `json:"municipios"`       // [24 bytes] Lista dinámica de favoritos guardados
	MunicipioCodigo string   `json:"municipio_codigo"` // [16 bytes] Código identificador único
	MunicipioNombre string   `json:"municipio_nombre"` // [16 bytes] Nombre legible del municipio
	Latitud         float64  `json:"latitud"`          // [8 bytes] Coordenada decimal
	Longitud        float64  `json:"longitud"`         // [8 bytes] Coordenada decimal
}

// Estructures del fitxers clima-valors.go
type PreUrl struct {
	Url    string       `json:"datos"` // [16 bytes] Definim que el camp Url sera de tipus igual que el item "datos" del Json obtingut
	Client *http.Client // [8 bytes] Aquest camp ens servirà per quan realitzem testing, tindrem enmagatzemat aqui el tipus de crida que s'ha realitzat
}

type AutoGenerated []struct {
	Origen     Origen     `json:"origen"`     // [96 bytes]
	Prediccion Prediccion `json:"prediccion"` // [24 bytes]
	Elaborado  string     `json:"elaborado"`  // [16 bytes]
	Nombre     string     `json:"nombre"`     // [16 bytes]
	Provincia  string     `json:"provincia"`  // [16 bytes]
	ID         int        `json:"id"`         // [8 bytes]
	Version    float64    `json:"version"`    // [8 bytes]
}
type Origen struct {
	Productor string `json:"productor"` // [16 bytes]
	Web       string `json:"web"`       // [16 bytes]
	Enlace    string `json:"enlace"`    // [16 bytes]
	Language  string `json:"language"`  // [16 bytes]
	Copyright string `json:"copyright"` // [16 bytes]
	NotaLegal string `json:"notaLegal"` // [16 bytes]
}
type ProbPrecipitacion struct {
	Periodo string `json:"periodo"` // [16 bytes]
	Value   int    `json:"value"`   // [8 bytes]
}
type CotaNieveProv struct {
	Value   string `json:"value"`   // [16 bytes]
	Periodo string `json:"periodo"` // [16 bytes]
}
type EstadoCielo struct {
	Value       string `json:"value"`       // [16 bytes]
	Periodo     string `json:"periodo"`     // [16 bytes]
	Descripcion string `json:"descripcion"` // [16 bytes]
}
type Viento struct {
	Direccion string `json:"direccion"` // [16 bytes]
	Periodo   string `json:"periodo"`   // [16 bytes]
	Velocidad int    `json:"velocidad"` // [8 bytes]
}
type RachaMax struct {
	Value   string `json:"value"`   // [16 bytes]
	Periodo string `json:"periodo"` // [16 bytes]
}
type Dato struct {
	Value int `json:"value"` // [8 bytes]
	Hora  int `json:"hora"`  // [8 bytes]
}
type Temperatura struct {
	Dato   []Dato `json:"dato"`   // [24 bytes]
	Maxima int    `json:"maxima"` // [8 bytes]
	Minima int    `json:"minima"` // [8 bytes]
}
type SensTermica struct {
	Dato   []Dato `json:"dato"`   // [24 bytes]
	Maxima int    `json:"maxima"` // [8 bytes]
	Minima int    `json:"minima"` // [8 bytes]
}
type HumedadRelativa struct {
	Dato   []Dato `json:"dato"`   // [24 bytes]
	Maxima int    `json:"maxima"` // [8 bytes]
	Minima int    `json:"minima"` // [8 bytes]
}
type Dia struct {
	Temperatura       Temperatura         `json:"temperatura"`       // [40 bytes]
	SensTermica       SensTermica         `json:"sensTermica"`       // [40 bytes]
	HumedadRelativa   HumedadRelativa     `json:"humedadRelativa"`   // [40 bytes]
	ProbPrecipitacion []ProbPrecipitacion `json:"probPrecipitacion"` // [24 bytes]
	CotaNieveProv     []CotaNieveProv     `json:"cotaNieveProv"`     // [24 bytes]
	EstadoCielo       []EstadoCielo       `json:"estadoCielo"`       // [24 bytes]
	Viento            []Viento            `json:"viento"`            // [24 bytes]
	RachaMax          []RachaMax          `json:"rachaMax"`          // [24 bytes]
	Fecha             string              `json:"fecha"`             // [16 bytes]
	UvMax             int                 `json:"uvMax,omitempty"`   // [8 bytes]
}
type Prediccion struct {
	Dia []Dia `json:"dia"` // [24 bytes]
}
type Diaria struct {
	Time             time.Time `json:"-"`                 // [24 bytes]
	CodiIne          string    `json:"-"`                 // // [16 bytes] Nuevo campo para el código del municipio
	ProbPrecipitacio int       `json:"probPrecipitacion"` // [8 bytes]
	TemperaturaMax   int       `json:"maxima"`            // [8 bytes]
	TemperaturaMin   int       `json:"minima"`            // [8 bytes]
	HumitatRelativa  int       `json:"humedadRelativa"`   // [8 bytes]
}

// CodigoINE retorna el código INE oficial de 5 dígitos limpiando el prefijo "id" del campo ID
func (m Municipio) CodigoINE() string {
	if len(m.ID) > 2 && m.ID[:2] == "id" {
		return m.ID[2:]
	}
	if len(m.IDOld) == 5 {
		return m.IDOld
	}
	return m.ID
}
