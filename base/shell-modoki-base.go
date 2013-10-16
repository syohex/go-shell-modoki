package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func createSubprocess(command string) chan bool {
	ch := make(chan bool)
	go func() {
		cmd := exec.Command(command)

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

		ch := createSubprocess(string(line))

		waitSubprocess(ch)
	}
}
