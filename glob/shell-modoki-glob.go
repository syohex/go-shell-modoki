package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var globRegexp = regexp.MustCompile(`\*`)

func parseLine(line string) ([]string, *os.File) {
	var output *os.File

	words := strings.Split(strings.TrimSpace(line), " ")

	commands := make([]string, 0)
	for _, word := range words {
		trimed := strings.TrimSpace(word)

		if globRegexp.Match([]byte(trimed)) {
			expandeds, err := filepath.Glob(trimed)
			if err != nil {
				log.Fatal(err)
			}
			commands = append(commands, expandeds...)
		} else {
			commands = append(commands, trimed)
		}
	}

	return commands, output
}

func createSubprocess(command []string, output *os.File) chan bool {
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

		command, output := parseLine(string(line))
		ch := createSubprocess(command, output)

		waitSubprocess(ch)
	}
}
