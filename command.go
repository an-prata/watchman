// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func HandleEvent(event fsnotify.Event, commands []*exec.Cmd, abortOnErr bool) {
	if !event.Has(fsnotify.Write) {
		return
	}

	log.Println("Got file write event: calling command ...")

	for _, c := range commands {
		err := c.Run()

		if err != nil {
			log.Println("Failed to run callback:", err.Error())

			if abortOnErr {
				break
			}
		}
	}

}
