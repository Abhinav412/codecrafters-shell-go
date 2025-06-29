package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	if len(input) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	command := input[0]
	_, ok := COMMANDS[command]
	if ok {
		fmt.Printf("%v is a shell builtin\n", command)
		return
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
			if fileInfo.Mode()&0111 != 0 {
				fmt.Printf("%v is %v\n", command, fullPath)
				return
			}
		}
		if strings.ToLower(filepath.Ext(command)) == "" {
			fullPathWithExt := filepath.Join(dir, command+".exe")
			if fileInfo, err := os.Stat(fullPathWithExt); err == nil && !fileInfo.IsDir() {
				if fileInfo.Mode()&0111 != 0 {
					fmt.Printf("%v is %v\n", command, fullPathWithExt)
					return
				}
			}
		}
	}

	fmt.Printf("%v: not found\n", command)
}

func evaluate(input string) {
	input = strings.TrimSpace(input)
	args := strings.Split(input, " ")
	command, optional := args[0], args[1:]

	output, ok := COMMANDS[command]
	if ok {
		output(optional)
		return
	}

	executablePath := findExecutable(command)
	if executablePath != "" {
		runExternalProgram(executablePath, optional)
		return
	}

	fmt.Printf("%v: command not found\n", command)
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

func findExecutable(command string) string {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
			if fileInfo.Mode()&0111 != 0 {
				return fullPath
			}
		}
		if strings.ToLower(filepath.Ext(command)) == "" {
			fullPathWithExt := filepath.Join(dir, command+".exe")
			if fileInfo, err := os.Stat(fullPathWithExt); err == nil && !fileInfo.IsDir() {
				if fileInfo.Mode()&0111 != 0 {
					return fullPathWithExt
				}
			}
		}
	}

	return ""
}

func runExternalProgram(executablePath string, args []string) {
	cmd := exec.Command(executablePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
