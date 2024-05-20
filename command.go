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

type AsyncEventHandler struct {
	lastHandled time.Time
	minimumGap  time.Duration
}

func NewAsyncEventHandler(minimumGap time.Duration) AsyncEventHandler {
	return AsyncEventHandler{time.Now(), minimumGap}
}

func (eh AsyncEventHandler) HandleEvent(event fsnotify.Event, commands []*exec.Cmd) {
	if !event.Has(fsnotify.Write) || time.Now().Sub(eh.lastHandled) < eh.minimumGap {
		return
	}

	eh.lastHandled = time.Now()
	log.Println("Got file write event: calling command ...")

	for _, c := range commands {
		go runCommand(c)
		log.Println("Started -", c.String())
	}
}

type SyncEventHandler struct {
	lastHandled time.Time
	minimumGap  time.Duration
}

func NewSyncEventHandler(minimumGap time.Duration) SyncEventHandler {
	return SyncEventHandler{time.Now(), minimumGap}
}

func (eh SyncEventHandler) HandleEvent(event fsnotify.Event, commands []*exec.Cmd) {
	if !event.Has(fsnotify.Write) || time.Now().Sub(eh.lastHandled) < eh.minimumGap {
		return
	}

	eh.lastHandled = time.Now()
	log.Println("Got file write event: calling command ...")

	for i, c := range commands {
		if i == len(commands)-1 {
			go runCommand(c)
		} else {
			runCommand(c)
		}

		log.Println("Started -", c.String())
	}
}

func runCommand(command *exec.Cmd) {

	// Not capturing output is deliberate, to clean it up and output it would add
	// about 2MB extra memory usage to the program in my tests. May consider making
	// output capture a flag for debugging configurations.

	err := command.Run()

	if err != nil {
		log.Println("Failed to run callback -", command.String(), ":", err.Error())
	}
}
