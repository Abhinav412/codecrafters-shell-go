package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var COMMANDS map[string]func([]string)

func init() {
	COMMANDS = map[string]func([]string){
		"echo": echoCMD,
		"type": typeCMD,
		"exit": exitCMD,
	}
}

func exitCMD(input []string) {
	if len(input) == 1 && input[0] == "0" {
		os.Exit(0)
	}
	os.Exit(1)
}

func echoCMD(input []string) {
	fmt.Println(strings.Join(input, " "))
}

func typeCMD(input []string) {
	command := strings.Join(input, " ")
	_, ok := COMMANDS[command]
	if !ok {
		fmt.Println("%v: not found\n", command)
		return
	}
	fmt.Println("%v is a shell builtin\n", command)
}

func evaluate(input string) {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")
	command, optional := args[0], args[1:]

	output, ok := COMMANDS[command]
	if !ok {
		fmt.Println("%v: command not found\n", command)
		return
	}
	output(optional)
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		evaluate(command)
	}
}
