package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var projectName string
	fmt.Print("Enter project name: ")
	fmt.Scan(&projectName)

	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		panic(err)
	}

	components := []string{"doc", "src", "job", "data", "ref", "eg"}

	for _, component := range components {
		err := os.Mkdir(projectName+"/"+component, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Create the files inside the directories
	err = os.WriteFile(projectName+"/doc/bib-file.bib", []byte("This is the bib file"), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(projectName+"/doc/README.md", []byte("This is the doc file"), 0644)
	if err != nil {
		panic(err)
	}

	// Create an example
	exampleName := "example1"
	exampleDir := projectName + "/eg/" + exampleName
	err = os.Mkdir(exampleDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, component := range components {
		err := os.Mkdir(exampleDir+"/"+component, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(exampleDir+"/doc/README.md", []byte("This is the doc file for the example"), 0644)
	if err != nil {
		panic(err)
	}

	// Create the small data directory
	err = os.Mkdir(projectName+"/data/large", os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Directory structure created successfully.")

	// Initialize Git repository
	cmd := exec.Command("git", "init", projectName)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
