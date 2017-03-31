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

type Rotatglobestruct struct{

	px: float32
	py: float32
	pz: float32
	pradius: float32
}

func playSound(them Plane, me Plane) {

	fmt.Println("subconcious")

}

// The following code is based off cosinekitty.com/compass.html
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

}

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
	clat = math.Atan((1.0 - e2) * math.Tan(lat))
	return clat
}

func RotateGlobe (b float32, a float32, bradius float32, aradius float32)
{
		// Get modified coordinates of 'b' by rotating the globe so that 'a' is at lat=0, lon=0.
		br := Plane{Lat: b.Lat, Lng: (b.Lng - a.Lng), Alt: b.Alt, TrueCourse : 0} // putting zero for true course because its not used
		brp := LocationToPoint(br, oblate)

		// Rotate brp cartesian coordinates around the z-axis by a.lon degrees,
		// then around the y-axis by a.lat degrees.
		// Though we are decreasing by a.lat degrees, as seen above the y-axis,
		// this is a positive (counterclockwise) rotation (if B's longitude is east of A's).
		// However, from this point of view the x-axis is pointing left.
		// So we will look the other way making the x-axis pointing right, the z-axis
		// pointing up, and the rotation treated as negative.

		alat := -a.lat * math.Pi / 180.0
		alat := GeocentricLatitude(alat)

		acos := math.Cos(alat)
		asin := math.Sin(alat)

		bx := (brp.x * acos) - (brp.z * asin)
		by := brp.y
		bz := (brp.x * asin) + (brp.z * acos)

		return Rotatglobestruct{px: bx, py: by, pz: bz, pradius: bradius}

}
