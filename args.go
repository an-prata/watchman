// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"flag"
	"log"
)

type Args struct {
	file      *string
	command   *string
	start     bool
	splitThen bool
	sync      bool
	msGap     int64
}

func GetArgs() Args {
	filePath := flag.String("file", "", "The file to watch for changes")
	callback := flag.String("command", "", "The command to run on file change events")
	start := flag.Bool("start", false, "Whether or not watchman should run the given command immediately first")
	splitThen := flag.Bool("split-then", false, "Splits command string by the \"then\" operator (\";\"), using each new string as a command. Successive commands will run regardless the previous's success")
	sync := flag.Bool("sync", false, "Runs commands one after the other, blocking further execution until all commands except the last has completed. This will not block watchman on the last command, meaning that it may a process not meant to exit.")
	millisecondGap := flag.Int64("ms-gap", 0, "Adds a minimum gap between events, any events received less than this many milliseconds after the previous event will be ignored. Useful if writes are infrequent but duplicate from some programs.")

	flag.Parse()

	if len(*filePath) < 1 {
		log.Fatal("Please given a file to watch.")
	}

	if len(*callback) < 1 {
		log.Fatal("Please give a command to call on file change events")
	}

	return Args{filePath, callback, *start, *splitThen, *sync, *millisecondGap}
}
