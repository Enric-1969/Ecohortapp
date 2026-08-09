package main

import "math"

// utmToLatLon convierte coordenadas UTM (Huso 30N España) a Lat/Lon WGS84
func utmToLatLon(easting, northing float64, zone int) (lat, lon float64) {
	a := 6378137.0
	f := 1.0 / 298.257223563
	b := a * (1.0 - f)
	e2 := (a*a - b*b) / (a * a)
	ePrime2 := (a*a - b*b) / (b * b)
	k0 := 0.9996

	x := easting - 500000.0
	y := northing

	m := y / k0
	mu := m / (a * (1.0 - e2/4.0 - 3.0*e2*e2/64.0 - 5.0*e2*e2*e2/256.0))

	e1 := (1.0 - math.Sqrt(1.0-e2)) / (1.0 + math.Sqrt(1.0-e2))

	j1 := 3.0*e1/2.0 - 27.0*math.Pow(e1, 3)/32.0
	j2 := 21.0*e1*e1/16.0 - 55.0*math.Pow(e1, 4)/32.0
	j3 := 151.0 * math.Pow(e1, 3) / 96.0

	fp := mu + j1*math.Sin(2.0*mu) + j2*math.Sin(4.0*mu) + j3*math.Sin(6.0*mu)

	c1 := ePrime2 * math.Pow(math.Cos(fp), 2)
	t1 := math.Pow(math.Tan(fp), 2)
	r1 := a * (1.0 - e2) / math.Pow(1.0-e2*math.Pow(math.Sin(fp), 2), 1.5)
	n1 := a / math.Sqrt(1.0-e2*math.Pow(math.Sin(fp), 2))
	d := x / (n1 * k0)

	fact1 := n1 * math.Tan(fp) / r1
	fact2 := d * d / 2.0
	fact3 := (5.0 + 3.0*t1 + 10.0*c1 - 4.0*c1*c1 - 9.0*ePrime2) * math.Pow(d, 4) / 24.0

	latitude := fp - fact1*(fact2-fact3)

	fact5 := d
	fact6 := (1.0 + 2.0*t1 + c1) * math.Pow(d, 3) / 6.0

	lon0 := float64(zone*6-183) * math.Pi / 180.0
	longitude := lon0 + (fact5-fact6)/math.Cos(fp)

	return latitude * 180.0 / math.Pi, longitude * 180.0 / math.Pi
}
