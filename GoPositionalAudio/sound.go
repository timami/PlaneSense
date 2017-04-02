package main

import "fmt"
import "math"

func main() {

	us := Plane{Lat: 48, Lng: 11., Alt: 20000, TrueCourse : 84.4}
	them:= Plane{Lat: 49, Lng: 10, Alt: 22000, TrueCourse : 30.4}

	playSound(us, them)
}

type Answer struct {
	azimuth float64
	distance float64
	altitude float64
}
type Plane struct {

	Lat float64
	Lng float64
	Alt float64
	Tail string
	TrueCourse float64
}

type Rotateglobestruct struct{

	px float64
	py float64
	pz float64
	pradius float64
}

type NormalDiff struct {
	ndx float64
	ndy float64
	ndz float64
	radius float64
}

type Point struct {
	px float64
	py float64
	pz float64
	pradius float64
	pnx float64
	pny float64
	pnz float64
}

type Hrir struct {
	x float64
	rx float64
	gamma float64
}

func playSound(me Plane, them Plane) {

	// aziumuth, distance, altitude
	ADA := calculate(me, them)

	// https://en.wikipedia.org/wiki/Spherical_coordinate_system
	// https://en.wikipedia.org/wiki/Azimuth

	// tested to the best of our knowlege
	r :=  ADA.distance // distance adduming this is ro
	beta :=  90 - ADA.altitude // elevation or theta based off wikipeida
	alpha := ADA.azimuth //azimuth or phi based off wikipedia

	// tested against her code and gives the same results
	hrirVal:= sphericalToHrir(alpha,beta,r)

}

// McMullen's function
func sphericalToHrir(alpha float64, beta float64, r float64) Hrir {
		x := math.Sin(alpha) * math.Cos(beta)
		rx := (math.Sqrt(1 - math.Pow(x/1,2.0)))
		gamma := math.Acos(math.Cos(alpha) / rx)

		// wrapping the javascript sign function for beta
		var sign float64 = 0;
		if(beta > 0) {
			sign = 1
		}

		if(beta < 0){
			sign = -1
		}

		gamma = sign * gamma

		//Boundaries
		if(gamma < -math.Pi * 3/4.0) {
			gamma = -math.Pi * 3/4.0
		}

		if(gamma > math.Pi * (3/4.0+0.03125)) {
			gamma = math.Pi * (3/4.0+0.03125)
		}

		return Hrir{ x: x, 	rx: rx,	gamma: gamma }
}

// The following code is based off cosinekitty.com/compass.html
func EarthRadiusInMeters(latRadians float64) float64 {
	a:= 6378137.0
	b:= 6356752.3
	cos := math.Cos(latRadians)
	sin := math.Sin(latRadians)
	t1 := a * a * cos
	t2 := b * b * sin
	t3 := a * cos
	t4 := b * sin

	mathlol:= (t1*t1 + t2*t2) / (t3*t3 + t4*t4)
	ret:= math.Sqrt(mathlol)
	return ret
}

func GeocentricLatitude(lat float64) float64 {
	e2 := 0.00669437999014
	clat := math.Atan((1.0 - e2) * math.Tan(lat))
	return clat
}

func RotateGlobe (b Plane, a Plane, bradius float64, aradius float64) Rotateglobestruct {

		// Get modified coordinates of 'b' by rotating the globe so that 'a' is at lat=0, lon=0.
		br := Plane{Lat: b.Lat, Lng: (b.Lng - a.Lng), Alt: b.Alt, TrueCourse : 0} // putting zero for true course because its not used
		brp := LocationToPoint(br)

		// Rotate brp cartesian coordinates around the z-axis by a.lon degrees,
		// then around the y-axis by a.lat degrees.
		// Though we are decreasing by a.lat degrees, as seen above the y-axis,
		// this is a positive (counterclockwise) rotation (if B's longitude is east of A's).
		// However, from this point of view the x-axis is pointing left.
		// So we will look the other way making the x-axis pointing right, the z-axis
		// pointing up, and the rotation treated as negative.

		var alat float64 = -a.Lat * math.Pi / 180
		alat = GeocentricLatitude(alat)

		acos := math.Cos(alat)
		asin := math.Sin(alat)

		bx := (brp.px * acos) - (brp.pz * asin)
		by := brp.py
		bz := (brp.px * asin) + (brp.pz * acos)

		return Rotateglobestruct{px: bx, py: by, pz: bz, pradius: bradius}
}

func LocationToPoint(c Plane) Point{

	lat := c.Lat  * math.Pi / 180
	lng := c.Lng * math.Pi / 180
	radius := EarthRadiusInMeters(lat) // Earths radius
	clat := GeocentricLatitude(lat)

	cosLon := math.Cos(lng)
	sinLon := math.Sin(lng)
	cosLat := math.Cos(clat)
	sinLat := math.Sin(clat)

	var x float64 = radius * cosLon * cosLat
	var y float64 = radius * sinLon * cosLat
	var z float64 = radius * sinLat

	cosGlat := math.Cos(lat)
    sinGlat := math.Sin(lat)
    //Normal Vector
    var nx float64 = cosGlat * cosLon
    var ny float64 = cosGlat * sinLon
    var nz float64 = sinGlat

    x += c.Alt * nx
    y += c.Alt * ny
    z += c.Alt * nz

	ret := Point{px:x, py:y, pz:z, pradius:radius, pnx:nx, pny:ny, pnz:nz}
	return ret
}

// returns distance in kilometers
func distanceBetweenPoints (pA Point, pB Point) float64{
	dx := pA.px - pB.px
	dy := pA.py - pB.py
	dz := pA.pz - pB.pz
	distance := math.Sqrt(dx*dx + dy*dy + dz*dz) / 1000
	return distance
}

func NormVectorDiff(b Point,a Point) NormalDiff {
	dx := b.px - a.px
	dy := b.py - a.py
	dz := b.pz - a.pz
	distance2 := dx*dx + dy*dy + dz*dz
	if (distance2 == 0) {
		ret:= NormalDiff{ndx: 0, ndy: 0, ndz: 0, radius: 0}
		return ret;
	}

	dist := math.Sqrt(distance2);
	ret := NormalDiff{ndx: (dx/dist), ndy: (dy/dist), ndz: (dz/dist), radius: 1.0}
	return ret
}

// takes in two lat/long/elev coordinates returns the azimuth, distance, and altitude
func calculate (a Plane, b Plane) Answer {
	ap := LocationToPoint(a)
	bp := LocationToPoint(b)

	var azimuth float64 = 0;
	var altitude float64 = 0;

	br := RotateGlobe (b, a, bp.pradius, ap.pradius)
    if (br.pz*br.pz + br.py*br.py > 1.0e-6) {
			theta := math.Atan2(br.pz, br.py) * 180.0 / math.Pi
			azimuth = 90.0 - theta
		if (azimuth < 0.0) {
			azimuth += 360.0
		}
		if (azimuth > 360.0) {
			azimuth -= 360.0
		}
	}
	bma := NormVectorDiff(bp, ap)
	if (0 != (bma.ndx+bma.ndy+bma.ndz)) {
		altitude = 90.0 - (180.0 / math.Pi)*math.Acos(bma.ndx*ap.pnx + bma.ndy*ap.pny + bma.ndz*ap.pnz)
	}
	ret:= Answer {azimuth: azimuth , distance: distanceBetweenPoints(ap, bp), altitude: altitude}
	return ret;
}
