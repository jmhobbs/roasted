package main

import (
	"flag"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/jmhobbs/roasted/pkg"
	"github.com/jmhobbs/roasted/sr700"
)

func main() {
	var (
		device = flag.String("device", "/dev/tty.wchusbserial1410", "SR700 serial device")
		debug  = flag.Bool("debug", false, "enable debug logging")
	)
	flag.Parse()

	recipe, err := pkg.LoadSimpleRecipe(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	// todo: validate program before running

	options := serial.OpenOptions{
		PortName:              *device,
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	roaster := sr700.New(port)
	roaster.SetDebug(*debug)

	pgm, err := roaster.Connect()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected.  Program: ", pgm)

	for i, step := range recipe.Steps {
		var (
			secondsRemaining int = step.Duration
			temp             sr700.Temperature
		)

		log.Printf("= Step %d - Heat: %v Fan: %v Cooling: %v Duration: %v", i, step.Heat, step.Fan, step.Cooling, step.Duration)
		for secondsRemaining > 0 {
			if step.Cooling {
				temp, err = roaster.Cool(step.Fan, time.Second*10)
			} else {
				temp, err = roaster.Roast(step.Fan, step.Heat, time.Second*10)
			}

			if err != nil {
				log.Println("error! ", err)
			} else {
				if temp == sr700.TemperatureBelow150F {
					log.Printf("<- -- F / %d s", secondsRemaining)
				}
				log.Printf("<- %v F / %d s", temp, secondsRemaining)
			}

			time.Sleep(time.Second * time.Duration(min(secondsRemaining, 5)))

			secondsRemaining = secondsRemaining - 5
		}
	}

	_, err = roaster.Stop()
	if err != nil {
		log.Printf("error stopping: %v", err)
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
