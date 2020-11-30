package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// not tested under OSes other than Linux.
func main() {
	shellPtr := flag.Bool("shell", false, "Use shell (/bin/sh) to execute program")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Missing program to be executed. Append it after parameters.")
		fmt.Printf("Example: %s ls -lh\n", os.Args[0])
		fmt.Printf("%s -shell echo '$PWD'\n", os.Args[0])
		os.Exit(-1)
	}
	program := flag.Args()

	// We use go-isatty manually, as the package fatih/color fails to consider one scenario:
	// stdout is redirected to a file (no color), but stderr still outputs to terminal (shall have color)
	// but fatih/color disables color globally when stdout is not terminal.
	// and when stderr is redirect to a file, but stdout connect to terminal.
	// fatih/color will output color to the file stderr connects.
	stdErrNoColor := os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stderr.Fd()) && !isatty.IsCygwinTerminal(os.Stderr.Fd()))
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	errYellowInstance := color.New(color.FgYellow)
	errRedInstance := color.New(color.FgRed)
	if stdErrNoColor {
		errYellowInstance.DisableColor()
		errRedInstance.DisableColor()
	} else {
		errYellowInstance.EnableColor()
		errRedInstance.EnableColor()
	}
	errYellow := errYellowInstance.SprintFunc()
	errRed := errRedInstance.SprintFunc()

	fmt.Printf("%s: Executing\n", yellow(program))
	var cmd *exec.Cmd = nil
	if *shellPtr {
		shellArgs := []string{"-c"}
		shellArgs = append(shellArgs, strings.Join(program, " "))
		cmd = exec.Command("sh", shellArgs...)
	} else {
		cmd = exec.Command(program[0], program[1:]...)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("%s cmd.StderrPipe(): %v", errYellow(program), err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("%s cmd.StdoutPipe(): %v", errYellow(program), err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("%s cmd.Start(): %v", errYellow(program), err)
	}
	pid := cmd.Process.Pid
	displayedName := yellow(program) + " " + red(pid)
	errDisplayedName := errYellow(program) + " " + errRed(pid)
	fmt.Printf("%s: PID=%s\n", yellow(program), red(pid))

	stdoutScanner := bufio.NewScanner(stdout)
	isStdoutFinished := make(chan bool)
	go func(fin chan bool) {
		for stdoutScanner.Scan() {
			fmt.Printf("%s stdout: %s\n", displayedName, stdoutScanner.Text())
		}
		fin <- true
	}(isStdoutFinished)

	stderrScanner := bufio.NewScanner(stderr)
	isStderrFinished := make(chan bool)
	go func(fin chan bool) {
		for stderrScanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s stderr: %s\n", errDisplayedName, stderrScanner.Text())
		}
		fin <- true
	}(isStderrFinished)

	<-isStdoutFinished
	<-isStderrFinished

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("%s exited with status %d", errDisplayedName, exitError.ExitCode())
		} else {
			log.Fatalf("%s cmd.Wait(): %v", errDisplayedName, err)
		}
	}
}

