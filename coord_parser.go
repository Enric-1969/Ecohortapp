package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// parseAEMETCoord es el punto de entrada principal para convertir cualquier formato AEMET a WGS84
func parseAEMETCoord(coordStr string) float64 {
	raw := strings.TrimSpace(coordStr)
	if raw == "" {
		return 0.0
	}

	upper := strings.ToUpper(raw)
	isNegative := isSouthernOrWestern(upper)

	if val, ok := tryParseDMS(upper, isNegative); ok {
		fmt.Printf("📍 Coordenada DMS convertida: '%s' -> %f\n", coordStr, val)
		return val
	}

	if val, ok := tryParseDecimalOrUTM(upper, isNegative); ok {
		fmt.Printf("📍 Coordenada Decimal/UTM convertida: '%s' -> %f\n", coordStr, val)
		return val
	}

	if val, ok := tryParseCompactDMS(upper, isNegative); ok {
		fmt.Printf("📍 Coordenada DDMMSS convertida: '%s' -> %f\n", coordStr, val)
		return val
	}

	fmt.Printf("⚠️ Coordenada no reconocida ('%s').\n", coordStr)
	return 0.0
}

func isSouthernOrWestern(s string) bool {
	return strings.Contains(s, "S") ||
		strings.Contains(s, "W") ||
		strings.Contains(s, "O") ||
		strings.HasPrefix(s, "-")
}

func tryParseDMS(upper string, isNegative bool) (float64, bool) {
	r := strings.NewReplacer("º", " ", "°", " ", "'", " ", "\"", " ", "N", "", "S", "", "E", "", "W", "", "O", "")
	cleaned := strings.ReplaceAll(r.Replace(upper), ",", ".")
	fields := strings.Fields(cleaned)

	if len(fields) < 2 {
		return 0, false
	}

	deg, err1 := strconv.ParseFloat(fields[0], 64)
	min, err2 := strconv.ParseFloat(fields[1], 64)
	if err1 != nil || err2 != nil {
		return 0, false
	}

	var sec float64
	if len(fields) >= 3 {
		sec, _ = strconv.ParseFloat(fields[2], 64)
	}

	val := math.Abs(deg) + (min / 60.0) + (sec / 3600.0)
	if isNegative {
		val = -val
	}

	if val >= -180.0 && val <= 180.0 {
		return val, true
	}
	return 0, false
}

func tryParseDecimalOrUTM(upper string, isNegative bool) (float64, bool) {
	cleanFloatStr := strings.ReplaceAll(upper, ",", ".")
	cleanFloatStr = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			return r
		}
		return -1
	}, cleanFloatStr)

	val, err := strconv.ParseFloat(cleanFloatStr, 64)
	if err != nil || cleanFloatStr == "" {
		return 0, false
	}

	// 1. Coordenada decimal estándar WGS84
	if val >= -180.0 && val <= 180.0 {
		if isNegative && val > 0 {
			val = -val
		}
		return val, true
	}

	// 2. Coordenada UTM en metros
	if math.Abs(val) > 1000.0 {
		valAbs := math.Abs(val)
		var converted float64

		if valAbs > 1000000.0 { // Northing Y (Latitud)
			converted, _ = utmToLatLon(500000.0, valAbs, 30)
		} else { // Easting X (Longitud)
			_, converted = utmToLatLon(valAbs, 4400000.0, 30)
		}

		if isNegative && converted > 0 {
			converted = -converted
		}
		return converted, true
	}

	return 0, false
}

func tryParseCompactDMS(upper string, isNegative bool) (float64, bool) {
	var digits strings.Builder
	for _, ch := range upper {
		if ch >= '0' && ch <= '9' {
			digits.WriteRune(ch)
		}
	}
	cleanDigits := digits.String()

	if len(cleanDigits) < 5 || len(cleanDigits) > 7 {
		return 0, false
	}

	secStr := cleanDigits[len(cleanDigits)-2:]
	minStr := cleanDigits[len(cleanDigits)-4 : len(cleanDigits)-2]
	degStr := cleanDigits[:len(cleanDigits)-4]

	deg, _ := strconv.ParseFloat(degStr, 64)
	min, _ := strconv.ParseFloat(minStr, 64)
	sec, _ := strconv.ParseFloat(secStr, 64)

	val := deg + (min / 60.0) + (sec / 3600.0)
	if isNegative {
		val = -val
	}

	if val >= -180.0 && val <= 180.0 {
		return val, true
	}
	return 0, false
}
