package main

import (
	"time"
	"fmt"
	"strings"
	"os/exec"
	"os"
	"syscall"
)

type CheckResult struct {
	name string
	output []byte
	exitCode int
	err error
}

type Scheduler struct {
	checksByFrequency map[int][]Check
	resultsChan chan CheckResult
}

func NewScheduler(checks []Check, resultsChan chan CheckResult) Scheduler {
	checksByFrequency := make(map[int][]Check)

	for _, check := range checks {
		checksByFrequency[check.Frequency] = append(checksByFrequency[check.Frequency], check)
	}

	return Scheduler{checksByFrequency,resultsChan}
}

func (s *Scheduler) Run()  {
	tickChan := time.NewTicker(time.Second).C

	for {
		tick :=<- tickChan
		for frequency := range s.checksByFrequency {

			if tick.Unix() % int64(frequency) == 0 {
				for _, check := range s.checksByFrequency[frequency] {
					output, exitCode, err := s.execute(check)
					fmt.Printf("[%s] %s Output : \n %s\n", tick, check.Name, string(output))
					fmt.Printf("[%s] %s Exit Code: %d\n", tick, check.Name, exitCode)
					if err != nil {
						fmt.Printf("[%s] %s Error: %d\n", tick, check.Name, err)
					}
					s.resultsChan <- CheckResult{check.Name, output,exitCode, err}
				}
			}
		}
	}

}

func (s *Scheduler) execute(check Check) (output []byte, exitCode int, err error) {
	commandList := strings.Split(check.Command, " ")
	name := commandList[0]
	args := commandList[1:]

	cmd := exec.Command(name, args...)

	cmdEnv := make([]string, len(check.Environment))

	for envVar, value := range check.Environment {
		cmdEnv = append(cmdEnv, envVar + "=" + value)
	}

	cmd.Env = append(os.Environ(), cmdEnv...)

	output, err = cmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			return
		}
		ws := exitErr.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
		return output, exitCode, nil
	}
	return output, 0, nil

}

func mainScheduler(checks []Check, resultsChan chan CheckResult)  {
	scheduler := NewScheduler(checks, resultsChan)
	scheduler.Run()
}
