// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func HandleEvent(event fsnotify.Event, commands []*exec.Cmd) {
	if !event.Has(fsnotify.Write) {
		return
	}

	log.Println("Got file write event: calling command ...")

	for _, c := range commands {
		go runCommand(c)
		log.Println("Started")
	}

}

func runCommand(command *exec.Cmd) {

	// Not capturing output is deliberate, to clean it up and output it would add
	// about 2MB extra memory usage to the program in my tests. May consider making
	// output capture a flag for debugging configurations.

	err := command.Run()

	if err != nil {
		log.Println("Failed to run callback:", err.Error())
	}
}
