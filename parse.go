// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func MakeCommandsBuf(argvs [][]string, buf []*exec.Cmd) {
	// Checking length before hand is unecessary, as `slices.Grow` is a noop is no
	// growth is needed for `n` elements.

	slices.Grow[[]*exec.Cmd](buf, len(argvs))

	for i, command := range argvs {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Env = os.Environ()

		buf[i] = cmd
	}
}

func MakeCommands(argvs [][]string) []*exec.Cmd {
	commands := make([]*exec.Cmd, len(argvs))

	for i, command := range argvs {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Env = os.Environ()

		commands[i] = cmd
	}

	return commands
}

// Creates an array of argv arrays for calling commands as specified by the
// given `Args` struct.
func ParseArgvs(args *Args) ([][]string, error) {
	var callbackArgv []string

	if args.splitThen {
		callbackArgv = SeparateThen(*args.command)
	} else {
		callbackArgv = []string{*args.command}
	}

	for _, cmd := range callbackArgv {
		log.Println("watchman will call:", cmd)
	}

	commands := make([][]string, len(callbackArgv))

	for i, command := range callbackArgv {
		args, err := ArgSplit(command)

		if err != nil {
			return nil, err
		}

		commands[i] = args
	}

	return commands, nil
}

// Separates the argstring into a slice or strings by the "and" operator ('&&').
func SeparateAnd(argstr string) []string {
	split := strings.Split(argstr, "&&")

	for i, s := range split {
		split[i] = strings.TrimSpace(s)
	}

	return split
}

// Separates the argstring into a slice or strings by the "then" operator (';').
func SeparateThen(argstr string) []string {
	split := strings.Split(argstr, ";")

	for i, s := range split {
		split[i] = strings.TrimSpace(s)
	}

	return split
}

// Takes a single, concatenated, args string, and seperates it such that quoted
// regions remain whole, otherwise splitting by spaces.
func ArgSplit(argstr string) ([]string, error) {
	split := strings.Split(argstr, " ")

	// Check for quotes that are neither opening nor terminating

	for _, s := range split {
		if strings.ContainsAny(s, "\"'") && !isQuoted(s) {
			return nil, errors.New("Arguments contain stray quote")
		}
	}

	for arg := 0; arg < len(split)-1; arg++ {
		if !isNonTerminatedQuote(split[arg]) {
			continue
		}

		// Combine arguments then remove the latter

		split[arg] = split[arg] + split[arg+1]

		if arg < len(split)-2 {
			split = append(split[:arg+1], split[arg+2:]...)
		} else {
			split = split[:arg+1]
		}

		// Decrement so we revisit the concatinated elements again

		arg--
	}

	// If the loop failes to terminate a quote then the last arg will be
	// non-terminated

	if isNonTerminatedQuote(split[len(split)-1]) {
		return nil, errors.New("Arguments contain non-terminated quote")
	}

	return split, nil
}

// Returns true if the argument begins with a quote that is non-terminated.
func isNonTerminatedQuote(str string) bool {
	if str[0] != '\'' && str[0] != '"' {
		return false
	}

	return str[0] != str[len(str)-1]
}

// Returns true if the argument begins OR ends in quotes.
func isQuoted(str string) bool {
	return str[0] == '\'' || str[0] == '"' || str[len(str)-1] == '\'' || str[len(str)-1] == '"'
}
