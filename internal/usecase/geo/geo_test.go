package geo

import (
	"math"
	"testing"
)

func TestUtmToLatLon(t *testing.T) {
	// Coordenadas UTM Huso 30N conocidas para la Puerta del Sol (Madrid):
	// Easting ~ 440284, Northing ~ 4474241 -> Lat ~ 40.4167, Lon ~ -3.7037
	lat, lon := UtmToLatLon(440284.0, 4474241.0, 30)

	if math.Abs(lat-40.4167) > 0.01 {
		t.Errorf("Latitud UTM inesperada: se esperaba ~40.4167, se obtuvo %f", lat)
	}
	if math.Abs(lon-(-3.7037)) > 0.01 {
		t.Errorf("Longitud UTM inesperada: se esperaba ~-3.7037, se obtuvo %f", lon)
	}
}

func TestParseAEMETCoord(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		esperado float64
	}{
		{
			nombre:   "Formato DMS clásico (Grados Minutos Segundos Norte)",
			entrada:  "41° 23' 15\" N",
			esperado: 41.3875,
		},
		{
			nombre:   "Formato DMS Oeste en español ('O' -> Negativo)",
			entrada:  "03° 42' 13\" O",
			esperado: -3.703611,
		},
		{
			nombre:   "Formato Decimal directo con coma europea",
			entrada:  "41,3851",
			esperado: 41.3851,
		},
		{
			nombre:   "Formato Compacto AEMET DDMMSSW (6 dígitos + hemisferio)",
			entrada:  "034213W",
			esperado: -3.703611,
		},
		{
			nombre:   "Cadena vacía o inválida debe devolver 0.0 sin pánico",
			entrada:  "   ",
			esperado: 0.0,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			obtenido := ParseAEMETCoord(tc.entrada)
			if math.Abs(obtenido-tc.esperado) > 0.001 {
				t.Errorf("Para '%s': se esperaba %f, se obtuvo %f", tc.entrada, tc.esperado, obtenido)
			}
		})
	}
}
