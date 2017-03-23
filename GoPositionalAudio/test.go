// A simple example of using SoX libraries
package main

// Use this URL to import the go-sox library
import "github.com/krig/go-sox"
import "log"
import "flag"
import "strconv"

// Reads an input file, applies volume and flanger effects,
// plays to output device.
func main() {

	// ARGS FORMAT:
	// [0]: Sound file
	// [1]: Side string
	// [2]: Intensity for the side

	flag.Parse()

	// All libSoX applications must start by initializing the SoX library
	if !sox.Init() {
		log.Fatal("Failed to initialize SoX")
	}
	// Make sure to call Quit before terminating
	defer sox.Quit()

	in := sox.OpenRead(flag.Arg(0))
	if in == nil {
		log.Fatal("Failed to open input file")
	}
	// Close the file before exiting
	defer in.Release()

	// Open the output device: Specify the output signal characteristics.
	// Since we are using only simple effects, they are the same as the
	// input file characteristics.
	// Using "alsa" or "pulseaudio" should work for most files on Linux.
	// On other systems, other devices have to be used.
	out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
	if out == nil {
		out = sox.OpenWrite("default", in.Signal(), nil, "pulseaudio")
		if out == nil {
			log.Fatal("Failed to open output device")
		}
	}
	// Close the output device before exiting
	defer out.Release()

	// grab the side and intensity

	side:= flag.Arg(1)
	intensity:= flag.Arg(2)

	intensityVal, err:= strconv.ParseFloat(intensity, 64)

	if err != nil {
		log.Fatal("Some math went wrong on parsing intensity!")
	}

	if intensityVal <= 2.01 {
		log.Fatal("Volume must be above 2! I suggest trying 3")
	}

	// Create an effects chain: Some effects need to know about the
	// input or output encoding so we provide that information here.
	chain := sox.CreateEffectsChain(in.Encoding(), out.Encoding())
	// Make sure to clean up!
	defer chain.Release()

	// The first effect in the effect chain must be something that can
	// source samples; in this case, we use the built-in handler that
	// inputs data from an audio file.
	e := sox.CreateEffect(sox.FindEffect("input"))
	e.Options(in)
	// This becomes the first "effect" in the chain
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()

	// Create the `vol' effect, and initialise it with the desired parameters:
	e = sox.CreateEffect(sox.FindEffect("vol"))
	e.Options("3dB")
        // Add the effect to the end of the effects processing chain:
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()


	//*************************************************************************
	// THESE ARE THE TWO EFFECTS WE NEED, I DON"T KNOW WHY THE OTHERS ARE HERE, should try removing them.
	// Create the delay effect, and add it
	e = sox.CreateEffect(sox.FindEffect("delay"))

	// if left, change delay to be on the right
	if side == "left" {
		e.Options("00", "0.02")
	} else if side == "right" {
		e.Options("0.02", "00")
	}

	// How long to wait before playing the left sound, and right sound, respectively.
        // Add the effect to the end of the effects processing chain:
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()


	// Create the remix effect, and add it
	e = sox.CreateEffect(sox.FindEffect("remix"))


	// calc low intensity and make it into strings quick and dirty
	lowIntensityVal:= intensityVal - 2
	lowIntensityString:= strconv.FormatFloat(lowIntensityVal, 'f', -1, 64)

	// check sides
	if side == "left" {
		e.Options(("1v" + intensity), ("2v" + lowIntensityString)) // Left channel volume = 1v[volume], right channel volume = 2v[volume], volume 0 to 1.
	} else {
		e.Options(("1v" + lowIntensityString), ("2v" + intensity)) // Left channel volume = 1v[volume], right channel volume = 2v[volume], volume 0 to 1.
	}
  // Add the effect to the end of the effects processing chain:
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()
	//*************************************************************************

	// Create the `flanger' effect, and initialise it with default parameters:
	// e = sox.CreateEffect(sox.FindEffect("flanger"))
	// e.Options()
	// chain.Add(e, in.Signal(), in.Signal())
	// e.Release()

	// The last effect in the effect chain must be something that only consumes
	// samples; in this case, we use the built-in handler that outputs data.
	e = sox.CreateEffect(sox.FindEffect("output"))
	e.Options(out)
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()

	// Flow samples through the effects processing chain until EOF is reached.
	chain.Flow()
}
