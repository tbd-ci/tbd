package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	a := Runner("/work-space/work/rea/tbd-test-repo/")
	fmt.Println(a)
}

type success bool

type ResultStep struct {
	Success     success `json:"status"`
	StdErr      string  `json:"standardError"`
	StdOut      string  `json:standardOut`
	CombinedStd string  `json:combinedStd`
}

type Result struct {
	OverallSuccess success      `json:overallSuccess`
	Steps          []ResultStep `json:steps`
}

func Runner(path string) (res Result) {
	ciPath := path + "ci/"

	for _, step := range steps(ciPath) {
		resultStep := ResultStep{}
		var stdOut bytes.Buffer
		var stdErr bytes.Buffer

		rawCmd := ciPath + step.Name() + "/run"
		cmd := exec.Command(rawCmd)

		combinedStd, err := cmd.CombinedOutput()
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr

		if err != nil {
			res.OverallSuccess = false
			resultStep.Success = false
			resultStep.StdErr = stdErr.String()
			resultStep.CombinedStd = string(combinedStd)
		} else {
			fmt.Println(err)
			resultStep.Success = true
			resultStep.StdOut = stdOut.String()
			resultStep.CombinedStd = string(combinedStd)
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
