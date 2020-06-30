package ui

import "os/exec"

func StartCommand(cmd string) {
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Start()
}
