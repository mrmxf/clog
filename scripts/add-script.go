//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package cmd implements commands for the cobra CLI library

package scripts

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

const KeywordUse = "clog"
const KeywordShort = "short"
const KeywordLong = "extra"

// Add a script from a given filename
// @todo - get script subcommands / options
func AddScript(rootCmd *cobra.Command, filePath string) {

	script := new(cobra.Command)
	scriptMeta, err := scriptMetadata(filePath)
	if err != nil {
		return
	}

	usage := strings.Split(scriptMeta[KeywordUse], " ")
	if len(usage[0]) == 0 {
		// we might have a legacy script with the word usage>
		usage = strings.Split(scriptMeta["usage"], " ")
		if len(usage[0]) == 0 {
			// no command to add - ignore this script
			return
		}
	}
	// we have a name, so add the script for processing
	name := usage[0]
	script.Use = name
	if _, shortExists := scriptMeta[KeywordShort]; shortExists {
		script.Short = scriptMeta[KeywordShort]
	}
	if _, longExists := scriptMeta[KeywordLong]; longExists {
		script.Long = scriptMeta[KeywordLong]
	}
	scriptMeta["filePath"] = filePath
	//preserve the metadata so that we know how to run the script later
	scriptsMap[name] = scriptMeta

	script.Run = func(cmd *cobra.Command, args []string) {
		//check that this is a valid script and nothing funky has happened
		if _, ok := scriptsMap[cmd.Name()]; !ok {
			slog.Error("command " + c.F(cmd.Name()) + c.E(" has no script associated with it. Aborting"))
			panic("Something catastrophic has happened inside clog")
		}

		scriptFilePath := scriptsMap[cmd.Name()]["filePath"]
		fmt.Println("Script(", c.C(cmd.Name()), ")", c.C(scriptFilePath))

		bashArgs := []string{scriptFilePath}
		bashArgs = append(bashArgs, args...)
		exe := exec.Command("bash", bashArgs...)

		// var stdout, stderr []byte
		// var errStdout, errStderr error
		stdoutIn, _ := exe.StdoutPipe()
		stderrIn, _ := exe.StderrPipe()

		//connect stdin for console IO
		exe.Stdin = bufio.NewReader(os.Stdin)
		err = exe.Start()
		if err != nil {
			slog.Error("cmd.Start() failed for "+cmd.Name(), "err", err.Error())
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			// stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
			_, err = rewriteStdout(os.Stdout, stdoutIn)
			if err != nil {
				slog.Error("Failed to rewrite script os.Stdout " + err.Error())
			}
			wg.Done()
		}()

		// stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
		_, err = rewriteStdout(os.Stderr, stderrIn)
		if err != nil {
			slog.Error("Failed to rewrite script os.StdErr " + err.Error())
		}

		wg.Wait()

		//output has finished, so disconnect StdIn to allow Command to complete
		// exe.Stdin = nil
		// fmt.Println("yo mama")
		// err = exe.Wait()
		if err != nil {
			slog.Error("cmd.Run() failed for %s with %s\n", cmd.Name(), err)
		}
		// if (errStdout != nil) || (errStderr != nil) {
		// 	log.Error("failed to capture stdout or stderr\n")
		// }
		exitStatus := exe.ProcessState.ExitCode()
		if err != nil || exitStatus > 0 {
			os.Exit(exitStatus)
		}
	}

	rootCmd.AddCommand(script)
}

func init() {
	slog.Debug("init scripts/add-script.go")
}
