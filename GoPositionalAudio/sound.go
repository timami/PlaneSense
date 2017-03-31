package main

import "fmt"
import "math"

func main() {

	us := Plane{Lat: 48, Lng: 11., Alt: 20000, TrueCourse : 84.4}
	them:= Plane{Lat: 49, Lng: 10, Alt: 22000, TrueCourse : 30.4}

	playSound(us, them)
}

type Plane struct {

	Lat float32
	Lng float32
	Alt int32
	Tail string
	TrueCourse float32
}

type Point struct {
	px float32
	py float32
	pz float32
	pradius float32
	pnx float32
	pny float32
	pnz float32

func playSound(them Plane, me Plane) {

	fmt.Println("subconcious")

}

// The following code is based off cosinekitty.com/compass.html


func EarthRadiusInMeters(latRadians float32) {
	a:= 6378137.0
	b:= 6356752.3
	cos := math.Cos(latRadians)
	sin := math.Sin(latRadians)
	t1 := a * a * cos
	t2 := b * b * sin
	t3 := a * cos
	t4 := b * sin

	return math.Sqrt((t1*t1 + t2*t2) / (t3*t3 + t4*t4))
}

func GeocentricLatitude(lat float32) {
	e2 = 0.00669437999014
	clat = math.Atan((1.0 - e2) * Math.Tan(lat))
	return clat
}

func LocationToPoint(c Plane) {

	lat := c.Lat  * math.Pi / 180
	lng := c.Lng * math.Pi / 180
	radius = EarthRadiusInMeters(lat) // Earths radius
	clat = GeocentricLatitude(lat)

	cosLon := math.Cos(lon)
	sinLon := math.Sin(lon)
	cosLat := math.Cos(clat)
	sinLat := math.Sin(clat)

	x:= radius * cosLon * cosLat
	y := radius * sinLon * cosLat
    z := radius * sinLat
    
	cosGlat := math.Cos(lat)
    sinGlat := math.Sin(lat)
    //Normal Vector    
    nx := cosGlat * cosLon
    ny := cosGlat * sinLon
    nz := sinGlat
    
    x += c.elv * nx
    y += c.elv * ny
    z += c.elv * nz

	ret := Point(px:x, py:y, pz:z, pradius:radius, pnx:nx, pny:ny, pnz:nz)
	return ret
}

func distanceBetweenPoints (pA, pB)
{
	dx := pA.px - pB.px;
	dy := pA.py - pB.py;
	dz := pA.pz - pB.pz;
	distance = math.sqrt(dx*dx + dy*dy + dz*dz)
	return distance
}
