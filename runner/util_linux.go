package runner

import (
	"io"
	"os/exec"
	"syscall"
	"time"

	"github.com/creack/pty"
)

func (e *Engine) killCmd(cmd *exec.Cmd) (pid int, err error) {
	pid = cmd.Process.Pid

	if e.config.Build.SendInterrupt {
		// Sending a signal to make it clear to the process that it is time to turn off
		if err = syscall.Kill(-pid, syscall.SIGINT); err != nil {
			return
		}
		time.Sleep(e.config.killDelay())
	}

	// https://stackoverflow.com/questions/22470193/why-wont-go-kill-a-child-process-correctly
	err = syscall.Kill(-pid, syscall.SIGKILL)

	// Wait releases any resources associated with the Process.
	_, _ = cmd.Process.Wait()
	return
}

func (e *Engine) startCmd(cmd string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	c := exec.Command("/bin/sh", "-c", cmd)
	if e.config.Screen.NoPTY {
		stderr, err := c.StderrPipe()
		if err != nil {
			return nil, nil, nil, err
		}
		stdout, err := c.StdoutPipe()
		if err != nil {
			return nil, nil, nil, err
		}
		err = c.Start()
		if err != nil {
			return nil, nil, nil, err
		}
		return c, stdout, stderr, err
	}
	f, err := pty.Start(c)
	return c, f, f, err
}
