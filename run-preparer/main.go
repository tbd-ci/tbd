package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func CheckoutProject(treeId string) (tmpDir string) {
	tmpDir, err := ioutil.TempDir("/tmp", "tbd")
	if err != nil {
		panic(err)
	}

	if err := exec.Command("git", "--work-tree", tmpDir, "checkout", treeId, "--", ".").Run(); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	cmd2 := exec.Command("git", "add", "-A")
	cmd2.Dir = tmpDir
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	return
}

func main() {
	tmpDir := CheckoutProject("HEAD")
	fmt.Println(tmpDir)
}
