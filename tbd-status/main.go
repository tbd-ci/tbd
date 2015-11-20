package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Result struct {
	OverallSuccess bool `json:"overallSuccess"`
}

func main() {
	usage := func() string {
		return "Usage: tbd-status TREEISH"
	}

	if len(os.Args) != 2 {
		fmt.Println(usage())
	}
	treeish := os.Args[1]

	cmd := exec.Command("git", "notes", "--ref", "tbd", "show", treeish+"^{tree}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s", string(output))
	}

	result := Result{}
	if err := json.Unmarshal(output, &result); err != nil {
		log.Fatal("tbd-status - json: %s", err, err)
	}

	if result.OverallSuccess {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
