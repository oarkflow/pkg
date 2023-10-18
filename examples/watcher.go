package main

import (
	"fmt"
	"log"

	"github.com/oarkflow/pkg/watcher"
)

// SetMaxEvents to 1 to allow at most 1 event's to be received
// on the Event channel per watching cycle.
//
// If SetMaxEvents is not set, the default is to send all events.
// w.SetMaxEvents(1)

// Only notify rename and move events.
// w.FilterOps(watcher.Rename, watcher.Move, watcher.Write, watcher.Create)

func main() {
	w, err := watcher.New(&watcher.Option{
		Path: []string{"./test_folder"},
	})
	if err != nil {
		panic(err)
	}
	w.OnAnyEvent(func(event watcher.Event) {
		fmt.Println(event)
	})
	if err := w.Start(); err != nil {
		log.Fatalln(err)
	}
}
