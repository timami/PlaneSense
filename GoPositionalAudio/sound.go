package main

import "fmt"

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


func playSound(them Plane, me Plane) {

	fmt.Println("subconcious")

}
