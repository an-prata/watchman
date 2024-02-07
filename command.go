// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func HandleEvent(event fsnotify.Event, commands []*exec.Cmd) {
	if !event.Has(fsnotify.Write) {
		return
	}

	log.Println("Got file write event: calling command ...")

	for _, c := range commands {
		go runCommand(c)
	}

}

func runCommand(command *exec.Cmd) {
	out, err := command.Output()

	if err != nil {
		log.Println("Failed to run callback:", err.Error())
	}

	outputLines := strings.Split(string(out), "\n")

	for _, line := range outputLines {
		log.Println("[Command Ouput]:", line)
	}
}
