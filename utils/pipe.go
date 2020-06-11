package utils

import (
    "github.com/alex-held/dev-env/manifest"
    scriptish "github.com/ganbarodigital/go_scriptish"
)

// The command's status code is stored in the pipeline.StatusCode.
//noinspection GoUnhandledErrorResult
func Com(commandFactory manifest.CommandFactory, args ...string) func(p *scriptish.Pipe) (int, error) {

    return func(p *scriptish.Pipe) (int, error) {

        expArgs := make([]string, len(args))
        for i := 0; i < len(args); i++ {
            expArgs[i] = p.Env.Expand(args[i])
        }
        cmd := commandFactory.Create(expArgs[0], expArgs[1:]...)

        // attach all of our inputs and outputs
        stdout := scriptish.NewDest()
        stderr := scriptish.NewDest()
        cmd.Stdin = p.Stdin
        cmd.Stdout = stdout
        cmd.Stderr = stderr

        // let's do it
        err := cmd.Start()
        if err != nil {
            return scriptish.StatusNotOkay, err
        }

        // wait for it to finish
        err = cmd.Wait()

        // copy the output to our pipe
        //
        // it's not ideal, because we can't preserve the original mixed
        // order of the command's output atm
        //
        // at some point, we'll need a new version of pipe that does
        // support preserving mixed order output!
        for line := range stdout.ReadLines() {
            scriptish.TracePipeStdout("%s", line)
            p.Stdout.WriteString(line)
            p.Stdout.WriteRune('\n')
        }
        for line := range stderr.ReadLines() {
            scriptish.TracePipeStderr("%s", line)
            p.Stderr.WriteString(line)
            p.Stderr.WriteRune('\n')
        }

        // we want the process's status code
        statusCode := cmd.ProcessState.ExitCode()

        // all done
        return statusCode, err
    }
}
