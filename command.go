// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

// For handling events and running commands based on write events.
type EventHandler struct {
	lastHandled time.Time
	minimumGap  time.Duration
	syncronous  bool
}

// Create a new handler for file events. The given duration will be used to
// ignore events made within that time from one another.
func NewEventHandler(minimumGap time.Duration, syncronous bool) EventHandler {
	return EventHandler{time.Now(), minimumGap, syncronous}
}

// Handle a file event, running the given set of commands if the event is a
// write and is beyond the minimum time duration between write events.
func (eh EventHandler) HandleEvent(event fsnotify.Event, commands []*exec.Cmd) {
	if !event.Has(fsnotify.Write) || time.Now().Sub(eh.lastHandled) < eh.minimumGap {
		return
	}

	eh.lastHandled = time.Now()
	log.Println("Got file write event: calling command ...")

	for i, c := range commands {
		errChan := runCommand(c)
		log.Println("Started -", c.String())

		// If its the last command it can be put in a go routine and allow the rest of
		// the program to run. This allows for deamon processes to run in this
		// position.
		if eh.syncronous && i != len(commands)-1 {
			<-errChan
		}

	}
}

func runCommand(command *exec.Cmd) chan error {

	// Not capturing output is deliberate, to clean it up and output it would add
	// about 2MB extra memory usage to the program in my tests. May consider making
	// output capture a flag for debugging configurations.

	errChan := make(chan error)

	go func() {
		err := command.Run()

		if err != nil {
			log.Println("Failed to run callback -", command.String(), ":", err.Error())
		}

		errChan <- err
	}()

	return errChan

}
