package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/deckarep/golang-set"
)

type status string

type result struct {
	Success       bool
	CommitMessage string
}

var statuses = mapset.NewSet()

func main() {
	var treeIsh = flag.String("treeish", "HEAD", "TreeIsh Id")
	limit := flag.Int("limit", 100, "Depth of commits")
	flag.Parse()

	WalkTree(statuses, treeIsh, *limit)

	resultsMap := make(map[string]result)
	for _, s := range statuses.ToSlice() {
		var success bool
		if err := exec.Command("tbd-status", fmt.Sprint(s)).Run(); err == nil {
			success = true
		}
		resultsMap[s.(string)] = result{success, "<placeholder>"}
	}
	for k, v := range resultsMap {
		fmt.Printf("%s %t\t%s\n", k, v.Success, v.CommitMessage)
	}
}

func WalkTree(statuses mapset.Set, treeIsh *string, limit int) {
	combinedStd, err := logCmd(treeIsh).CombinedOutput()
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

func logCmd(id *string) (c *exec.Cmd) {
	c = exec.Command("git", "log", "--format='%p'", *id)
	return
}
