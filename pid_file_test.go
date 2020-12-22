package pidfile

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
)

func TestCrossCreateOrClearPidFile(t *testing.T) {
	path := "./pidfiletest.pid"
	pidFile, err := CreatePidFile(path)
	if err != nil {
		fmt.Printf("create pid file failed, error: %s\n", err.Error())
	}

	sc := make(chan os.Signal)
	isChildProcess := os.Getenv("CHILD_PROCESS") == "1"

	if !isChildProcess {
		go func() {
			cmd := exec.Command(os.Args[0], "-test.run=TestCrossCreateOrClearPidFile")
			cmd.Env = append(os.Environ(), "CHILD_PROCESS=1")

			err := cmd.Run()
			if err != nil {
				fmt.Printf("run father process failed, error: %s\n", err.Error())
			}
		}()

		// wait signal from child
		waitSignal(sc)

		// get pid of child
		tmpPid, err := ReadFromFile(NewFile("./pidfiletest.pid.tmp"))
		if err != nil || !tmpPid.ProcessExits() {
			fmt.Printf("read pid of child process failed. error: %s. isChildProcess: %v", err.Error(), isChildProcess)
		}

		clearPidFile(pidFile)

		// send signal to child
		sendSignal(tmpPid.Id, isChildProcess)
	} else {
		// send signal to father
		sendSignal(pidFile.Pid.Id, isChildProcess)

		// wait signal from father
		waitSignal(sc)

		clearPidFile(pidFile)
	}
}

func waitSignal(sc chan os.Signal) {
	signal.Notify(sc, syscall.SIGUSR1)
	<-sc
}

func sendSignal(pid int, isChildProcess bool) {
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("find process %d failed. error: %s. isChildProcess: %v\n", pid, err.Error(), isChildProcess)
	}

	err = process.Signal(syscall.SIGUSR1)
	if err != nil {
		fmt.Printf("send signal to process %d failed. error: %s\n. isChildProcess: %v\n", pid, err.Error(), isChildProcess)
	}
}

func clearPidFile(pidfile *PidFile) {
	err := ClearPidFile(pidfile)
	if err != nil {
		fmt.Printf("clear pid file %s failed, error: %s\n", pidfile.File.Path, err.Error())
	}
}
