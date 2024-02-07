// Copyright (c) 2024 Evan Overman (https://an-prata.it).
// Licensed under the MIT License.
// See LICENSE file in repository root for complete license text.

package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Creates an array of argv arrays for calling commands as specified by the
// given `Args` struct.
func ParseCommandFromArgs(args *Args) ([]*exec.Cmd, error) {
	var callbackArgv []string

	if args.splitThen {
		callbackArgv = SeparateThen(*args.command)
	} else if args.splitAnd {
		callbackArgv = SeparateAnd(*args.command)
	} else {
		callbackArgv = []string{*args.command}
	}

	for _, cmd := range callbackArgv {
		log.Println("watchman will call:", cmd)
	}

	enviornment := os.Environ()
	commands := []*exec.Cmd{}

	for _, command := range callbackArgv {
		args, err := ArgSplit(command)

		if err != nil {
			return nil, err
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = enviornment

		commands = append(commands, cmd)
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
