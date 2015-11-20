package main

import (
	"encoding/json"
	"fmt"
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
		return
	}
	treeish := os.Args[1]

	cmd := exec.Command("git", "notes", "--ref", "tbd", "show", treeish+"^{tree}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", string(output))
		os.Exit(-1)
	}

	result := Result{}
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Printf("%s", err)
		os.Exit(-1)
	}

	if result.OverallSuccess {
		fmt.Print("success")
		os.Exit(0)
	} else {
		fmt.Print("failed")
		os.Exit(1)
	}
}
