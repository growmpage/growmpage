package growsensoric

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/smarthome-go/rpirf"

	"github.com/growmpage/growganizer"
	"github.com/growmpage/growhelper"
)

type Growswitcher struct {
	mu sync.Mutex
}

func (g *Growswitcher) withLockContextOnPi(fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if runtime.GOARCH == "arm" {
		fn()
	}
}

func (g *Growswitcher) Switch(plug growganizer.PlugControl) {
	fmt.Println("try to send plug: " + growhelper.ToString(plug.Code))
	g.withLockContextOnPi(func() {
		device, err := rpirf.NewRF(
			uint8(plug.PinNumberSend),
			uint8(plug.ProtocolIndex),
			uint8(plug.Repeat),
			uint16(plug.PulseLength),
			uint8(plug.Length),
		)
		if err != nil {
			fmt.Printf("err.Error() from device: %v\n", err.Error())
		}
		if err = device.Send(plug.Code); err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
		}
		fmt.Println("plug sent: " + growhelper.ToString(plug.Code))
		if err := device.Cleanup(); err != nil {
			fmt.Printf("Cleanup: %v\n", err.Error())
		}
	})
}

func (g *Growswitcher) SwitchSim(plugOn growganizer.PlugControl, plugOff growganizer.PlugControl) {
	if runtime.GOARCH == "arm" {
		g.Switch(plugOn)
		defer g.Switch(plugOff)
		time.Sleep(time.Second * 55)
		g.Switch(plugOff)
	} else {
		log.Println("simulate sim on")
		time.Sleep(time.Second * 1)
		log.Println("simulate sim off")
	}
}

// TODO: https://leangaurav.medium.com/golang-channels-vs-sync-once-for-one-time-execution-of-code-fafc81d2f54d
