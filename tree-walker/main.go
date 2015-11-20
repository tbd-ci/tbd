package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/deckarep/golang-set"
)

type status string

var statuses = mapset.NewSet()

func main() {
	var treeIsh = flag.String("treeish", "HEAD", "TreeIsh Id")
	limit, _ := strconv.Atoi(*flag.String("limit", "100", "Depth of commits"))

	WalkTree(statuses, treeIsh, limit)
	for _, s := range statuses.ToSlice() {
		_, err := exec.Command("tbd-status", fmt.Sprint(s)).CombinedOutput()
		fmt.Println(err)
	}
}

func WalkTree(statuses mapset.Set, treeIsh *string, limit int) {
	combinedStd, err := cmd(treeIsh).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	commitsRaw := string(combinedStd)
	for _, s := range strings.Split(commitsRaw, "\n") {

		if len(statuses.ToSlice()) >= limit {
			break
		}

		s = strings.Replace(s, "'", "", -1)

		for _, t := range strings.Fields(s) {
			statuses.Add(t)
		}
	}
	return
}

func cmd(id *string) (c *exec.Cmd) {
	c = exec.Command("git", "log", "--format='%p'", *id)
	return
}
