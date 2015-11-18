package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func createTempDir() (dir string) {
	dir, err := ioutil.TempDir("/tmp", "tbd")
	if err != nil {
		panic(err)
	}

	return
}

func CheckoutProject(projectDir string, treeId string) {
	tmpDir := createTempDir()
	cmd := fmt.Sprintf("git --git-dir=%s/.git --work-tree=%s checkout %s -f -q", projectDir, tmpDir, treeId)

	_, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	CheckoutProject("/work-space/work/rea/tbd-test-repo/", "HEAD")
}
