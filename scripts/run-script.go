//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package cmd implements commands for the cobra CLI library

package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/exp/slog"
)

func RunScript(scriptName string, scriptPath string, args ...string) {
	fmt.Println("Executing script", c.C(scriptPath), "install", c.C(scriptName))

	bashArgs := []string{scriptPath}
	bashArgs = append(bashArgs, "install")
	bashArgs = append(bashArgs, args...)

	Exec("bash", bashArgs, nil)
}

// execute a command and restream Stdin & StdOut - return status
func Exec(command string, args []string, env map[string]string) int {
	exe := exec.Command(command, args...)

	// add in any env variables
	exe.Env = os.Environ()
	// append environemnt variables from the passed map
	for k, v := range env {
		exe.Env = append(exe.Env, k+"="+v)
	}

	// var stdout, stderr []byte
	var errStdout, errStderr error
	execStdOut, _ := exe.StdoutPipe()
	execStdErr, _ := exe.StderrPipe()

	// stdin is unconnected for now - to be debugged
	// exe.Stdin = bufio.NewReader(os.Stdin)
	// stdin, err := cmd.StdinPipe()

	err := exe.Start()
	if err != nil {
		slog.Error("FATAL cmd.Start() during scripts.Exec()")
		slog.Error("FATAL trying to execute %s\n", strings.Join(args, " "),"err",err)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	var exitCode = 0
	wg.Add(1)
	go func() {
		// defer stdin.Close()
		// io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")

		// stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		_, errStdout = rewriteStdout(os.Stdout, execStdOut)
		_, errStderr = rewriteStdout(os.Stderr, execStdErr)
		wg.Done()
	}()

	// stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	_, err = rewriteStdout(os.Stderr, execStdErr)
	if err != nil {
		slog.Error("Failed to rewrite script os.Stderr " + err.Error())
	}
	// wait for all the standard output to be rewritten
	wg.Wait()
	//wait for the process to exit
	exe.Wait()
	// grab the exit status of the process
	exitCode = exe.ProcessState.ExitCode()

	if errStdout != nil {
		slog.Warn("WARNING rewriting StdOut")
		slog.Warn("WARNING executing %s\n", strings.Join(args, " "), "err", errStdout)
	}
	if errStderr != nil {
		slog.Warn("WARNING rewriting Stderr")
		slog.Warn("WARNING executing %s\n", strings.Join(args, " "), "err",errStderr)
	}

	return exitCode
}
