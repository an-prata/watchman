// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func main() {
	filePath := flag.String("file", "", "The file to watch for changes")
	callback := flag.String("command", "", "The command to run on file change events")
	splitAnd := flag.Bool("split-and", false, "Splits command string by the \"and\" operator (\"&&\"), using each new string as a command. Successive commands will only run upon the previous's success")
	splitThen := flag.Bool("split-then", false, "Splits command string by the \"then\" operator (\";\"), using each new string as a command. Successive commands will run regardless the previous's success")

	flag.Parse()

	if len(*filePath) < 1 {
		log.Fatal("Please given a file to watch.")
	}

	if len(*callback) < 1 {
		log.Fatal("Please give a command to call on file change events")
	}

	if *splitAnd && *splitThen {
		log.Fatal("Cannot split on both 'and' and 'then'")
	}

	callbackArgv := []string{*callback}

	if *splitThen {
		callbackArgv = SeparateThen(callbackArgv[0])
	}

	if *splitAnd {
		callbackArgv = SeparateAnd(callbackArgv[0])
	}

	for _, cmd := range callbackArgv {
		log.Println("watchman will call:", cmd)
	}

	commands := [][]string{}

	for _, command := range callbackArgv {
		args, err := ArgSplit(command)

		if err != nil {
			log.Fatal("Could not parse callback arguments:", err.Error())
		}

		commands = append(commands, args)
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal("Could not create watcher:", err.Error())
	}

	defer watcher.Close()
	err = watcher.Add(*filePath)

	if err != nil {
		log.Fatal("Could not add file for watching:", err.Error())
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal("Channel was not ok")
			}

			if !event.Has(fsnotify.Write) {
				break
			}

			log.Println("Got file write event: calling command ...")

			for _, c := range commands {
				cmd := exec.Command(c[0], c[1:]...)
				cmd.Env = os.Environ()

				err := cmd.Run()

				if err != nil {
					log.Println("Failed to run callback:", err.Error())

					// Only stop if requiring the previous command's success

					if *splitAnd {
						break
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				log.Fatal("Channel was not ok")
			}

			log.Println("Got error:", err)
		}
	}
}
