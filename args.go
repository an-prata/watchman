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
	splitAnd  bool
	splitThen bool
}

func GetArgs() Args {
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

	return Args{filePath, callback, *splitAnd, *splitThen}
}
