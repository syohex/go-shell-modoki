package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func parseLine(line string) []string {
	words := strings.Split(strings.TrimSpace(line), " ")

	command := make([]string, 0)
	for _, word := range words {
		trimed := strings.TrimSpace(word)
		command = append(command, trimed)
	}

	return command
}

func createSubprocess(command []string) chan bool {
	ch := make(chan bool)
	go func() {
		cmd := exec.Command(command[0], command[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			log.Println(err)
		} else {
			err = cmd.Wait()
			if err != nil {
				log.Println(err)
			}
		}

		ch <- true
		close(ch)
	}()

	return ch
}

func waitSubprocess(ch <-chan bool) {
	<-ch
}

func main() {
	bio := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("-> ")
		line, hasMoreLine, err := bio.ReadLine()
		if !hasMoreLine && err == io.EOF {
			fmt.Println("Bye")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		command := parseLine(string(line))
		ch := createSubprocess(command)

		waitSubprocess(ch)
	}
}
