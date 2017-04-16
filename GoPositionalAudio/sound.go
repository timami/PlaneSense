package main

import "math"
import "io/ioutil"
import "reflect"
import "fmt"
import "strings"
import "strconv"
import "github.com/krig/go-sox"
import "log"
import "container/heap"
import "github.com/parnurzeal/gorequest"
import "github.com/tidwall/gjson"

// global variable for our plane
var us Plane;
var APIKey string = "key1";

// start minheap code
type PlaneHeap []Plane

func (h PlaneHeap) Len() int { return len(h) }
func (h PlaneHeap) Less(i, j int) bool {

	adaplane1 := calculate(us, h[i])
	adaplane2 := calculate(us, h[j])

	return adaplane1.distance < adaplane2.distance
}

func (h PlaneHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PlaneHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *h = append(*h, x.(Plane))
}

func (h *PlaneHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
// end minheap code

func main() {

	if(strings.Compare(APIKey, "key1") == 0) {
		fmt.Println("Please replace API key with a valid one from Mashape!")
		return
	}

	// our coordinates which we will pull from google maps
	us = Plane{Lat: 29.63, Lng: -82.35, Alt: 10000}

	// initialize the PlaneHeap
	ph := &PlaneHeap{}
	heap.Init(ph)

	// make get request to grab all planes reporting
	url:="https://opensky-network.p.mashape.com/states/all"
	fmt.Println("URL:>", url)

	request := gorequest.New()
	resp, body, errs := request.Get(url).
	Set("X-Mashape-Key", "key1").
	Set("Accept", "application/json").
	End()

	if(errs != nil) {
		fmt.Println(errs)
		fmt.Println(resp)
		return;
	}

	// ** DEBUG **
	// fmt.Println("response Status:", errs)
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Body:", string(body))

	json := string(body)

	// get all states and loop through them and datamarshal them into a plane struct
	planeVectors:= gjson.Get(json, "states")
	planeVectors.ForEach(func(key, value gjson.Result) bool {
		// fmt.Println(value.Array()[5].Float()) // Long
		// fmt.Println(value.Array()[6].Float()) // Lat
		// fmt.Println(value.Array()[7].Float()) // Alt
		// fmt.Println(value.Array()[1].String()) // Alt


		// grab the plane data (more info can be found at https://opensky-network.org/mashape/#validity)
		planeArr:= value.Array()
		tempLng:= planeArr[5].Float() // Lng
		tempLat:= planeArr[6].Float() // Lat
		tempAlt:= planeArr[7].Float() // Alt
		tempTail:= planeArr[1].String() // Callsign

		// push into minheap
		if(tempLat != 0 && tempLng != 0 && tempAlt != 0) {
			tempPlane := Plane{Lat: tempLat, Lng: tempLng, Alt: tempAlt, Tail: tempTail}
			heap.Push(ph, tempPlane)
		}

		return true;
	})

	// Manual Sound Coordination for debug purposes

	// West
	// them1:= Plane{Lat: 29.63, Lng: -82.53, Alt: 20000, TrueCourse : 30.4}
	// East
	// them2:= Plane{Lat: 29.63, Lng: -82.11, Alt: 40000, TrueCourse : 30.4}
	// North
	// them3:= Plane{Lat: 29.66, Lng: -82.35, Alt: 40000, TrueCourse : 30.4}
	// //South
	// them4:= Plane{Lat: 29.60, Lng: -82.35, Alt: 20000, TrueCourse : 30.4}

	// heap.Push(ph, them1)
	// heap.Push(ph, them2)
	// heap.Push(ph, them3)
	// heap.Push(ph, them4)


	// ** DEBUG **
	// for ph.Len() > 0 {
	// fmt.Println(heap.Pop(ph))
	// }

	// print out closestPlane Info
	closestPlane := heap.Pop(ph).(Plane)
	fmt.Println(closestPlane)

	ADA := calculate(us, closestPlane)


	fmt.Println(ADA.azimuth);

	ADA.azimuth = float64(int(ADA.azimuth + us.TrueCourse) % 360)
		if ( ADA.altitude >= -30 && ADA.altitude <= 30 ) {

			// Play flat left or giht
			play_left_or_right("270_360/threepiovertwodirectleft.mp3","90_180/piovertwodirectright.mp3", "0_90/0pidirectinfront.mp3", "back/beep-09.mp3", ADA.azimuth);
		}

		if ( ADA.altitude > 30 && ADA.altitude <= 75 ) {
				// Play upper left or right
				play_left_or_right("270_360/topleft.mp3", "0_90/topright.mp3", "0_90/0pidirectinfront.mp3", "back/beep-09.mp3", ADA.azimuth);
				return
		}
		if ( ADA.altitude < -30 && ADA.altitude >= -75) {
				// Play bottom left or bottom right
				play_left_or_right("270_360/bottomleftfront.mp3","0_90/bottomright.mp3", "0_90/0pidirectinfront.mp3", "back/beep-09.mp3", ADA.azimuth);
				return
		}
		if ( ADA.altitude > 75 ) {

			// Play top
			top_bottom("/home/nicolas/Audio/0_90/top.mp3")
			return
		}
		if ( ADA.altitude < -75 ) {
			// Play bottom
			top_bottom("/home/nicolas/Audio/270_360/bottom.mp3")

			return
		}

	}

func play_left_or_right(Left string, Right string, front string, back string, azimuth float64) {

	if !sox.Init() {
		log.Fatal("Failed to initialize SoX")
	}
	// Make sure to call Quit before terminating
	defer sox.Quit()

	var in *sox.Format;

	fmt.Println(azimuth)

	moddedAzimuth  := int(azimuth) % 360.0
	fmt.Println(moddedAzimuth);

	if ( moddedAzimuth >= 45 && moddedAzimuth < 135 ) {
		// play right sound

		in = sox.OpenRead(Right)
		if in == nil {
			log.Fatal("Failed to open input file")
		}
		defer in.Release()


	} else if azimuth >= 315 || azimuth < 45  {
		in = sox.OpenRead(front)
		if in == nil {
			log.Fatal("Failed to open input file")
		}
		defer in.Release()
	}	else if azimuth >= 135 || azimuth < 225  {
		in = sox.OpenRead(back)
		if in == nil {
			log.Fatal("Failed to open input file")
		}
		defer in.Release()
	} else {
		// play left sound

		in = sox.OpenRead(Left)
		if in == nil {
			log.Fatal("Failed to open input file")
		}
		defer in.Release()

	}


		out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
		if out == nil {
			out = sox.OpenWrite("default", in.Signal(), nil, "pulseaudio")
			if out == nil {
				log.Fatal("Failed to open output device")
			}
		}

		chain := sox.CreateEffectsChain(in.Encoding(), out.Encoding())

		e := sox.CreateEffect(sox.FindEffect("input"))
		e.Options(in)
		// This becomes the first "effect" in the chain
		chain.Add(e, in.Signal(), in.Signal())
		e.Release()

		e = sox.CreateEffect(sox.FindEffect("vol"))
		e.Options("3dB")
		// Add the effect to the end of the effects processing chain:
		chain.Add(e, in.Signal(), in.Signal())
		e.Release()

		e = sox.CreateEffect(sox.FindEffect("output"))
		e.Options(out)
		chain.Add(e, in.Signal(), in.Signal())
		e.Release()

		chain.Flow()
		defer chain.Release()

		// Close the output device before exiting
		defer out.Release()
	}

func top_bottom(side string) {

	if !sox.Init() {
		log.Fatal("Failed to initialize SoX")
	}
	// Make sure to call Quit before terminating
	defer sox.Quit()


	in := sox.OpenRead(side)
	if in == nil {
		log.Fatal("Failed to open input file")
	}

	defer in.Release()

	out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
	if out == nil {
		out = sox.OpenWrite("default", in.Signal(), nil, "pulseaudio")
		if out == nil {
			log.Fatal("Failed to open output device")
		}
	}

	chain := sox.CreateEffectsChain(in.Encoding(), out.Encoding())

	chain.Flow()
	defer chain.Release()

	// Close the output device before exiting
	defer out.Release()
}


// Flow data from in to out via the samples buffer
func flow(in, out *sox.Format, samples []sox.Sample) {
	n := uint(len(samples))
	for number_read := in.Read(samples, n); number_read > 0; number_read = in.Read(samples, n) {
		out.Write(samples, uint(number_read))
	}
}

func loadHrtf(file string) Hrtf {

	rootPath:= "/home/nicolas/Desktop/hrtf-demo-from-site/hrir/"

	left:= rootPath + "/" + file + "_l.dat"
	right:= rootPath + "/" + file + "_r.dat"

	datL, errL := ioutil.ReadFile(left)
	datR, errR := ioutil.ReadFile(right)

	if(errL != nil || errR != nil) {
		fmt.Println(errL)
		fmt.Println(errR)
		return Hrtf{}
	}

	var hrirL [25][50][200]float64

	dataLString := string (datL)
	azimuthsL := strings.Split(dataLString, "\n")

	for azimuth := 0; azimuth < 25; azimuth++ {
		elevationsL := strings.Split(azimuthsL[azimuth], ",")

		for elevation := 0; elevation < 50; elevation++ {

			for k := 0; k < 200; k++ {
				parsedValue, err := strconv.ParseFloat((elevationsL[elevation + k*50]), 64);

				if(err != nil) {
					fmt.Println(err);
					return Hrtf{}
				}

				hrirL[azimuth][elevation][k] = parsedValue
			}
		}
	}

	var hrirR [25][50][200]float64

	dataRString := string (datR)
	azimuthsR := strings.Split(dataRString, "\n")

	for azimuth := 0; azimuth < 25; azimuth++ {
		elevationsR := strings.Split(azimuthsR[azimuth], ",")

		for elevation := 0; elevation < 50; elevation++ {

			for k := 0; k < 200; k++ {

				parsedValue, err := strconv.ParseFloat((elevationsR[elevation + k*50]), 64);

				if(err != nil) {
					fmt.Println(err);
					return Hrtf{}
				}

				hrirR[azimuth][elevation][k] = parsedValue;
			}
		}
	}

	hrtf := Hrtf{ hrirL: hrirL, hrirR: hrirR }
	// fmt.Println(hrtf)


	return hrtf
}

type Hrtf struct {

	// 25 rows of azimuth 50 cols of elevation, 200 deep of "Where's The Fun" based off the cipic database
	hrirL [25][50][200] float64
	hrirR [25][50][200] float64
}

// Goes into trafficSound.go
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

type Hrir_Buffer struct {
	left []float64
	right []float64

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
	hrirVal:= sphericalToHrir(alpha, beta, r)
	fmt.Println(hrirVal)

	// choose the hrtffile
	hrtfFile:= "3"
	hrtf:= loadHrtf(hrtfFile)

	// fmt.Println(hrtf)

	// go-sox loading sounds
	var samples [2048]sox.Sample

	if !sox.Init() {
		log.Fatal("Failed to initialize SoX")
	}
	defer sox.Quit()

	// Open the input file.
	in := sox.OpenRead("./planecloned.ogg");
	if in == nil {
		log.Fatal("Failed to open input file")
	}

	// Set up the memory buffer for writing
	buf := sox.NewMemstream()
	defer buf.Release()

	out := sox.OpenMemstreamWrite(buf, in.Signal(), nil, "sox")
	if out == nil {
		log.Fatal("Failed to open memory buffer")
	}

	flow(in, out, samples[:])
	out.Release()
	in.Release()

	// this should now be a buffer sound of the
	in = sox.OpenMemRead(buf)
	if in == nil {
		log.Fatal("Failed to open memory buffer for reading")
		}

	// print go deep and explore
	inType := reflect.TypeOf(in.Signal())
	for i:= 0; i < inType.NumMethod(); i++ {
		fmt.Println(inType.Method(i))
	}

	// update convolver and pass in theta and phi
	updateConvolver(alpha, beta, hrtf, in.Signal().Rate());

}

func updateConvolver(theta float64, phi float64, hrtf Hrtf, rate float64) {

	// linear interpolation taken from main.js line 440 from Online updating convolver

	// get sample rate
	var fs float64= 44100;
	var hrir_length = int(math.Ceil(200* rate / fs))
	var hrir_lengthsixtyfour = math.Ceil(200* rate / fs)


	var q float64
	var d float64
	var k_p float64
	var k_n float64

	// guess
	left:= make([]float64, hrir_length)
	right:= make([]float64, hrir_length)

// linear interpolation of the hrir to the sample rate
	var k float64 = 0;

	for ; k < hrir_lengthsixtyfour; k++ {
		q = k / (hrir_lengthsixtyfour) * 200;
    k_p = math.Floor(q);
    k_n = math.Ceil(q);
    d = q - k_p;

		fmt.Println("Theta: "+ strconv.FormatFloat(theta, 'f', -1, 64))
		fmt.Println("Phi: "+ strconv.FormatFloat(phi, 'f', -1, 64))
		fmt.Println("k_p: "+ strconv.FormatFloat(k_p, 'f', -1, 64))
		fmt.Println("index: " + strconv.FormatFloat(k, 'f', -1, 64))
		fmt.Println("arrayLen: " + strconv.FormatFloat(hrir_lengthsixtyfour, 'f', -1, 64))



    left[int(k)] = hrtf.hrirL[int(theta)][int(phi)][int(k_p)] * (1.0-d) + hrtf.hrirL[int(theta)][int(phi)][int(k_n)] * d;
    right[int(k)] = hrtf.hrirR[int(theta)][int(phi)][int(k_p)] * (1.0-d) + hrtf.hrirR[int(theta)][int(phi)][int(k_n)] * d;
	}



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
