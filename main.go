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
		panic(err)
	}

	metaVer := getMetadataVersion("/Users/bas/Documents/Development/puppet/modules/node_red/metadata.json")
	gitTag:= getGitTag("/Users/bas/Documents/Development/puppet/modules/node_red")



	if gitTag != "" {
		GitVer, err = semver.Parse(gitTag)
	} else {
		GitVer = metaVer
		setGitTag("/Users/bas/Documents/Development/puppet/modules/node_red", metaVer.String())
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
		}
	} else if metaVer.Compare(GitVer) == 1 {
		fmt.Println("metadata version > git tag. we are doing stuff now")
        setGitTag("/Users/bas/Documents/Development/puppet/modules/node_red", metaVer.String())
	}

}


func getMetadataVersion(f string) semver.Version {

	file, _ := ioutil.ReadFile(f)
	json.Unmarshal(file, &MyJson)
	metaVer, err := semver.Parse(MyJson.Version)
	if err != nil {
		panic(err)
	}
	return metaVer

}


func getGitTag(d string) string {

	//var gitVersion semver.Version

	//gitVersion, err := semver.NewPRVersion("0.0.1")

	cmd := exec.Command("git", "tag")
	cmd.Dir = d
	out, _ := cmd.Output()

	//if out == []{
	//	gitVersion =
	//}
	//	gitVersion, _ =  semver.Parse(strings.TrimSuffix(string(out), "\n"))

	return strings.TrimRight(string(out), "\n")
}

func setGitTag(d string, ver string) string {

	//cmdOptions := fmt.Sprintf("tag -a" %v -m \"gitlab ci tag\"", ver)
    //fmt.Println(cmdOptions)
	cmd := exec.Command("git", "tag", "-a", ver, "-m", "\"gitlab ci tag\"")
	cmd.Dir = d
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	return string(out)
}