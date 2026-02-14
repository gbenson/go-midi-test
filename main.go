package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	drv, err := rtmididrv.New()
	must(err)

	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	if len(ins) < 1 {
		fmt.Println("No inputs!")
		os.Exit(1)
	}

	for i, in := range ins {
		fmt.Printf("in[%d] = %v\n", i, in)
	}

	// takes the first input
	in := ins[0]

	fmt.Printf("opening MIDI Port %v\n", in)
	must(in.Open())

	defer in.Close()

	stop, err := midi.ListenTo(
		in,
		func(msg midi.Message, timestampms int32) {
			fmt.Printf("[%v] %v\n", timestampms, msg)

			var bt []byte
			var ch, key, vel uint8
			switch {
			case msg.GetSysEx(&bt):
				fmt.Printf("got sysex: % X\n", bt)
			case msg.GetNoteStart(&ch, &key, &vel):
				fmt.Printf("starting note %s on channel %v with velocity %v\n",
					midi.Note(key), ch, vel)
			case msg.GetNoteEnd(&ch, &key):
				fmt.Printf("ending note %s on channel %v\n",
					midi.Note(key), ch)
			default:
				// ignore
			}
		},
		midi.UseSysEx(),
		midi.HandleError(complain),
		midi.UseActiveSense(),
		midi.UseTimeCode(),
	)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Println("sleeping now")
	time.Sleep(time.Second * 5)

	fmt.Println("stopping")
	stop()
	fmt.Printf("closing MIDI Port %v\n", in)
}

func complain(err error) {
	fmt.Println("o hai, an error:", err)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
