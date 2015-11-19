package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func main() {
	a := Runner()
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
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

func Runner() (res Result) {
	res.OverallSuccess = true

	ciPath := "ci/"

	for _, step := range steps(ciPath) {
		resultStep := ResultStep{}

		resultStep.Stage = step.Name()

		rawCmd := path.Join(ciPath + step.Name() + "/run")
		cmd := exec.Command(rawCmd)

		combinedStd, err := cmd.CombinedOutput()
		if err != nil {
			res.OverallSuccess = false
			resultStep.Success = false
			resultStep.Output = string(combinedStd)
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
