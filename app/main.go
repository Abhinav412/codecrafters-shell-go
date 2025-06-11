package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}
		if command[:len(command)-1] == "exit 0" {
			break
		}
		if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5 : len(command)-1])
		} else if strings.HasPrefix(command, "type") {
			argument := command[5:]
			switch argument {
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "exit":
				fmt.Println("exit is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			default:
				fmt.Println(argument + ": not found")
			}
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
