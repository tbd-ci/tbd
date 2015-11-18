package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Steps []string `yaml:"steps"`
}

func ParseConfig(filename string) (con Config, err error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return con, err
	}

	err = yaml.Unmarshal(source, &con)
	if err != nil {
		return con, err
	}

	return
}

func main() {
	testStuff()
}

func testStuff() {
	c, err := ParseConfig("test.yml")
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}
