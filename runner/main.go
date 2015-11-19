package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"golang.org/x/crypto/ssh/terminal"
)

var workDir = flag.String("workdir", "", "directory to run script from")

func main() {
	flag.Parse()

	if !terminal.IsTerminal(0) {
		stdIn, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		os.Stdin.Close()
		*workDir = string(stdIn)
	} else {
		if *workDir == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	a := Runner(workDir)
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}

type ResultStep struct {
	Success bool   `json:"status"`
	Stage   string `json:"stage"`
	Output  string `json:"combinedStd"`
}

type Result struct {
	OverallSuccess bool         `json:"overallSuccess"`
	Steps          []ResultStep `json:"steps"`
}

func Runner(runDir *string) (res Result) {
	res.OverallSuccess = true

	ciPath := "ci/"

	for _, step := range steps(ciPath) {
		resultStep := ResultStep{}

		resultStep.Stage = step.Name()

		rawCmd := path.Join(ciPath + step.Name() + "/run")
		log.Println(rawCmd)
		cmd := exec.Command(rawCmd)
		cmd.Dir = *runDir

		combinedStd, err := cmd.CombinedOutput()
		if err != nil {
			res.OverallSuccess = false
			resultStep.Success = false
			resultStep.Output = string(combinedStd)
			log.Println(err)
		} else {
			resultStep.Success = true
			resultStep.Output = string(combinedStd)
		}

		res.Steps = append(res.Steps, resultStep)
	}
	return
}

func steps(ciPath string) (steps []os.FileInfo) {
	steps, err := ioutil.ReadDir(ciPath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return
}
