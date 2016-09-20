package main

import "os/exec"

func main() {

	f := ""
	err := exec.Command("go", "run", "mockaroo-temp.go", "user.go", f).Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("gofmt", "-w", "main-User.go").Run()
	if err != nil {
		panic(err)
	}
}
