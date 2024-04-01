package growsensoric

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"

	"github.com/growmpage/growhelper"
)

type Growobserver struct {
	mu        sync.Mutex
	observing Observing
}

type Observing struct {
	Minutes     int
	Temperature int
	Humidity    int
	Picture     bool
}

// func (g *Growobserver) withLockContextOnPi(fn func() Observing) {
// 	g.mu.Lock()
// 	defer g.mu.Unlock()
// 	if runtime.GOARCH == "arm" {
// 		fn()
// 	}
// }

func (g *Growobserver) Measure() Observing {
	if runtime.GOARCH != "arm" {
		return Observing{
			Minutes:     growhelper.Minutes(),
			Temperature: 25,
			Humidity:    60,
		}
	}
	
	var env physic.Env
	// g.withLockContextOnPi(func() Observing {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()
	if err := dev.Sense(&env); err != nil {
		log.Fatal(err)
	}
	dev.Halt()
	bus.Close()
	return Observing{
		Minutes:     growhelper.Minutes(),
		Temperature: int(math.Round(env.Temperature.Celsius())),
		Humidity:    int(math.Round(growhelper.ToFloat64(strings.Split(env.Humidity.String(), "%")[0]))),
	}
	// })

}

func (g *Growobserver) Picture() Observing {
	pictureMinutes := growhelper.Minutes()
	// base, _ := filepath.Abs(".")
	base := growhelper.Filename_pictures
	picturePath := base + growhelper.ToString(pictureMinutes) + ".png"

	if runtime.GOARCH == "arm" {
		exec.Command("raspistill", "-o", picturePath).Run()
		return Observing{
			Minutes: pictureMinutes,
			Picture: true,
		}
	} else {
		picturePathDumy := base + "toCopy.png"
		file, err := os.ReadFile(picturePathDumy)
		if err != nil {
			fmt.Println(err)
		}
		err = os.WriteFile(picturePath, file, 0644)
		if err != nil {
			fmt.Println(err)
		}
		return Observing{
			Minutes: pictureMinutes,
			Picture: true,
		}
	}
}
