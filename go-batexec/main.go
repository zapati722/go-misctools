package main

import (
	"os/exec"
)

func testRun(cmd string, args ...string) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		println("NAY:")
		println(err.Error())
	} else {
		println("YAY:")
		println(string(out))
	}
}

func main() {
	testRun("sass", "-version")
	testRun("cmd", "/C", "sass", "-version")
}
