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

func parseLine(line string) [][]string {
	words := strings.Split(strings.TrimSpace(line), " ")

	commands := make([][]string, 0)
	cmd := make([]string, 0)

	for _, word := range words {
		trimed := strings.TrimSpace(word)

		if globRegexp.Match([]byte(trimed)) {
			expandeds, err := filepath.Glob(trimed)
			if err != nil {
				log.Fatal(err)
			}
			cmd = append(cmd, expandeds...)
		} else if trimed == "|" {
			commands = append(commands, cmd)
			cmd = make([]string, 0)
		} else {
			cmd = append(cmd, trimed)
		}
	}

	commands = append(commands, cmd)

	return commands
}

func createSubprocess(command []string, in *io.PipeReader, out *io.PipeWriter, ch chan<- bool) {
	go func() {
		cmd := exec.Command(command[0], command[1:]...)

		if in == nil {
			cmd.Stdin = os.Stdout
		} else {
			cmd.Stdin = in
		}

		if out == nil {
			cmd.Stdout = os.Stdout
		} else {
			cmd.Stdout = out
		}

		err := cmd.Start()
		if err != nil {
			log.Println(err)
		} else {
			err = cmd.Wait()
			if err != nil {
				log.Println(err)
			}

			if in != nil {
				in.Close()
			}

			if out != nil {
				out.Close()
			}
		}

		ch <- true
	}()
}

func waitSubprocess(processes int, ch <-chan bool) {
	for i := 0; i < processes; i++ {
		<-ch
	}
}

type DataPipes struct {
	in  *io.PipeReader
	out *io.PipeWriter
}

func makeSubprocessPipes(processes int) []*DataPipes {
	pipes := make([]*DataPipes, 0)

	for i := 0; i < processes; i++ {
		in, out := io.Pipe()

		data := &DataPipes{in, out}
		pipes = append(pipes, data)
	}

	return pipes
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

		commands := parseLine(string(line))
		processes := len(commands)

		pipes := makeSubprocessPipes(processes)

		ch := make(chan bool, len(commands))
		for i, command := range commands {
			var in *io.PipeReader
			var out *io.PipeWriter

			if i != 0 {
				in = pipes[i-1].in
			}

			if i != len(commands)-1 {
				out = pipes[i].out
			}

			createSubprocess(command, in, out, ch)
		}

		waitSubprocess(processes, ch)
	}
}
