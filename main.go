package main

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)


type Metadata struct {
	Name                   string      `json:"name"`
	Version                string      `json:"version"`
	Author                 string      `json:"author"`
	Summary                string      `json:"summary"`
	License                string      `json:"license"`
	Source                 string      `json:"source"`
	Dependencies           interface{} `json:"dependencies"`
	OperatingsystemSupport interface{} `json:"operatingsystem_support"`
	Requirements           interface{} `json:"requirements"`
	PdkVersion             string      `json:"pdk-version"`
	TemplateUrl            string      `json:"template-url"`
	TemplateRef            string      `json:"template-ref"`
}


var MyJson Metadata


var (
	cmdOut []byte
	err    error
)

var GitVer semver.Version


func main() {

	if _, err := os.Stat("./metadata.json"); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	//if _, err := os.Stat("./git/config"); os.IsNotExist(err) {
	//	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	//	os.Exit(1)
	//}


	metaVer := getMetadataVersion("./metadata.json")
	gitTag:= getGitTag("./")



	if gitTag != "" {
		GitVer, err = semver.Parse(gitTag)
	} else {
		GitVer = metaVer
		setGitTag("./", metaVer.String())
	}


	fmt.Printf("metaVersion is: %v\n", metaVer)
	fmt.Printf("gitTag is: %v\n", gitTag)
	fmt.Printf("gitVer is: %v\n", GitVer)


	//If the version number in metadata.json is higher then the git tag, the value is 1

	if metaVer.Compare(GitVer) <= 0 {
		if metaVer.Equals(GitVer) == true {
			fmt.Println("Probably WIP and/or no git tag. No actions taken")
		} else {
			fmt.Println("The metadata version is lower than the git tag. This must be fixed")
			os.Exit(1)
		}
	} else if metaVer.Compare(GitVer) == 1 {
		fmt.Println("metadata version > git tag. we are doing stuff now")
        setGitTag("./", metaVer.String())
	}

}


func getMetadataVersion(f string) semver.Version {

	file, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in metadata.json: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &MyJson)
	metaVer, err := semver.Parse(MyJson.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in metadata.json: %v\n", err)
		os.Exit(1)
	}
	return metaVer

}


func getGitTag(d string) string {

	cmd := exec.Command("git", "describe", "--abbrev=0")
	cmd.Dir = d
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting git tag: %v\n", err)
		os.Exit(1)
	}

	return strings.TrimRight(string(out), "\n")
}

func setGitTag(d string, ver string) string {

	cmdTag := exec.Command("git", "tag", "-a", ver, "-m", "\"gitlab ci tag\"")
	cmdTag.Dir = d
	_, err := cmdTag.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}


	cmdPushTag := exec.Command("/bin/bash", "-c", "git push --tags")
	cmdPushTag.Dir = d
	outPushTag, err := cmdPushTag.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in tag push: %v\n", err)
		os.Exit(1)
	}


	return string(outPushTag)
}