package main

import (
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	drv, err := rtmididrv.New()
	must(err)

	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	if len(ins) < 1 {
		log.Println("No inputs!")
		os.Exit(1)
	}

	var bestIn drivers.In
	for i, in := range ins {
		log.Printf("in[%d] = %v", i, in)
		if bestIn != nil {
			continue
		}
		if strings.Contains(in.String(), "Midi Through") {
			continue
		}
		bestIn = in
	}

	// takes the first input
	in := bestIn
	if in == nil {
		in = ins[0]
	}

	log.Printf("opening MIDI Port %v", in)
	must(in.Open())

	defer in.Close()

	stop, err := midi.ListenTo(
		in,
		func(msg midi.Message, timestampms int32) {
			log.Printf("@%vms %v %v", timestampms, []byte(msg), msg)

			var bt []byte
			var ch, key, vel uint8
			switch {
			case msg.GetSysEx(&bt):
				log.Printf("got sysex: % X", bt)
			case msg.GetNoteStart(&ch, &key, &vel):
				log.Printf("starting note %s on channel %v with velocity %v",
					midi.Note(key), ch, vel)
			case msg.GetNoteEnd(&ch, &key):
				log.Printf("ending note %s on channel %v",
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
		log.Printf("ERROR: %s", err)
		return
	}
	defer stop()

	log.Println("interrupt to stop listening")
	for {
		time.Sleep(1 * time.Hour)
	}
}

func complain(err error) {
	log.Println("o hai, an error:", err)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
