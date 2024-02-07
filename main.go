// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	args := GetArgs()
	commands, err := ParseCommandFromArgs(&args)

	if err != nil {
		log.Fatal("Could not parse callback arguments:", err.Error())
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal("Could not create watcher:", err.Error())
	}

	defer watcher.Close()
	err = watcher.Add(*args.file)

	if err != nil {
		log.Fatal("Could not add file for watching:", err.Error())
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal("Channel was not ok")
			}

			HandleEvent(event, commands, args.splitAnd)

		case err, ok := <-watcher.Errors:
			if !ok {
				log.Fatal("Channel was not ok")
			}

			log.Println("Got error:", err)
		}
	}
}
