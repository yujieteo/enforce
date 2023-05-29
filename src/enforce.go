/*
This is a Go program that performs various operations on a project directory. Here's an overview of what the code does:

    It uses a dialog to allow the user to select a project directory.
    It validates that the project path exists.
    It creates several subdirectories within the project directory, such as "doc," "src," "job," "data," "ref," and "eg."
    It creates an example directory within the "doc" directory and further creates subdirectories within it.
    It creates a README file in the "eg" directory.
    It creates a "data/large" directory.
    It removes any empty subdirectories within the project directory.
    It checks if a Git repository already exists in the project directory. If not, it initializes a new Git repository.
    It creates additional subdirectories within the project directory.
    It generates a README file.
    It creates several subdirectories and files related to a LaTeX document.
    It creates a .gitignore file.
    It creates a Jupyter Notebook template.

Overall, the code sets up a project directory structure, initializes a Git repository (if necessary), and performs various file and directory operations within the project directory.
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

func main() {
	// Use dialog to select the project directory
	projectPath, err := dialog.Directory().Title("Select project directory").Browse()
	if err != nil {
		fmt.Println("Failed to select project directory:", err)
		return
	}

	// Validate the project path exists
	_, err = os.Stat(projectPath)
	if os.IsNotExist(err) {
		fmt.Println("Project path does not exist.")
		return
	}

	components := []string{"doc", "src", "job", "data", "ref", "eg"}

	for _, component := range components {
		componentPath := filepath.Join(projectPath, component)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Create an example
	exampleName := "report"
	exampleDir := filepath.Join(projectPath, "doc", exampleName)
	err = os.MkdirAll(exampleDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, component := range components {
		componentPath := filepath.Join(exampleDir, component)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	exampleReadmeFilePath := filepath.Join(projectPath, "eg", "README.md")
	err = os.WriteFile(exampleReadmeFilePath, []byte(`# List of Examples
	## Example 1
	## Example 2
	## Example 3
	`), 0644)
	if err != nil {
		panic(err)
	}

	// Create the large data directory
	dataLargeDir := filepath.Join(projectPath, "data", "large")
	err = os.MkdirAll(dataLargeDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	removeEmptySubdirectories(projectPath)

	for _, component := range components {
		componentPath := filepath.Join(projectPath, component)
		removeEmptySubdirectories(projectPath)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	errJob := os.MkdirAll(filepath.Join(projectPath, "job"), os.ModePerm)
	if errJob != nil {
		panic(err)
	}

	err = generateREADME(projectPath)
	if err != nil {
		fmt.Printf("Error generating README file: %v", err)
	}

	// Create the large data directory
	egStyDir := filepath.Join(projectPath, "doc", "report", "sty")
	err = os.MkdirAll(egStyDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = createColorThemeDSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createColorThemeSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createFontThemeSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createOuterThemeSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createInnerThemeSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createMainThemeSty(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createReportTex(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createGitignore(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createJupyterTemplate(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	// Check if Git repository already exists
	gitPath := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		// Initialize Git repository
		removeEmptySubdirectories(projectPath)

		sortFiles(projectPath)

		cmd := exec.Command("git", "-C", projectPath, "init")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("Git repository initialized.")
	} else {
		removeEmptySubdirectories(projectPath)
		fmt.Println("Git repository already exists.")
	}

}

// sortFiles recursively sorts files in a given folder path.
// It moves files to specific destination folders based on their file extensions.
// The sorting rules are as follows:
// - Files with extensions .pdf, .djvu, .epub, .html, .mkv, .mp4 are moved to the "ref" folder.
// - Files with extensions .docx, .md, .tex, .txt, .doc, .pptx are moved to the "doc" folder.
// - Files with extension .ipynb are moved to the "eg" folder.
// - Files with extensions .py, .go, .inp, .c, .m, .for, .cpp, .java are moved to the "src" folder.
// - All other files are moved to the "data" folder.
//
// Parameters:
// - folderPath: The path of the folder to sort.
//
// Example:
//
//	sortFiles("/path/to/folder")
func sortFiles(folderPath string) {

	// Check if the current folder is the .git repository
	if filepath.Base(folderPath) == ".git" {
		fmt.Printf("Skipping .git repository: '%s'\n", folderPath)
		return
	}

	subdirectories, err := getSubdirectories(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, subdir := range subdirectories {
		subdirPath := filepath.Join(folderPath, subdir)

		// Check if the current subdirectory is the .git directory
		if subdir == ".git" {
			fmt.Printf("Skipping .git repository: '%s'\n", subdirPath)
			continue
		}

		sortFiles(subdirPath)
	}

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		destFolder := ""

		switch extension {
		case ".pdf", ".djvu", ".epub", ".html", ".mkv", ".mp4":
			destFolder = filepath.Join("ref", strings.TrimSuffix(filepath.Base(path), extension))
		case ".docx", ".md", ".tex", ".txt", ".doc", ".pptx":
			destFolder = "doc"
		case ".ipynb":
			destFolder = filepath.Join("eg", strings.TrimSuffix(filepath.Base(path), extension))
		case ".py", ".go", ".inp", ".c", ".m", ".for", ".cpp", ".java":
			destFolder = filepath.Join("src", strings.TrimSuffix(filepath.Base(path), extension))
		default:
			destFolder = "data"
		}

		destFolderPath := filepath.Join(folderPath, destFolder)
		err = os.MkdirAll(destFolderPath, 0755)
		if err != nil {
			return err
		}

		destFilePath := filepath.Join(destFolderPath, filepath.Base(path))
		err = os.Rename(path, destFilePath)
		if err != nil {
			return err
		}

		fmt.Printf("Moved '%s' to '%s'\n", path, destFilePath)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// getSubdirectories retrieves a list of subdirectories in a given folder path.
//
// Parameters:
// - folderPath: The path of the folder to retrieve subdirectories from.
//
// Returns:
// - A slice of subdirectory names.
// - An error if any error occurs during the retrieval.
//
// Example:
//
//	subdirs, err := getSubdirectories("/path/to/folder")
//	if err != nil {
//	  log.Fatal(err)
//	}
//	for _, subdir := range subdirs {
//	  fmt.Println(subdir)
//	}
func getSubdirectories(folderPath string) ([]string, error) {
	var subdirectories []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != folderPath {
			subdirectories = append(subdirectories, strings.TrimPrefix(path, folderPath+string(os.PathSeparator)))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return subdirectories, nil
}

// removeEmptySubdirectories removes empty subdirectories in a given folder path.
//
// Parameters:
// - folderPath: The path of the folder to remove empty subdirectories from.
//
// Example:
//
//	removeEmptySubdirectories("/path/to/folder")
func removeEmptySubdirectories(folderPath string) {
	subdirectories := make([]string, 0)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != folderPath {
			subdirectories = append(subdirectories, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Remove empty subdirectories in reverse order
	for i := len(subdirectories) - 1; i >= 0; i-- {
		dir := subdirectories[i]
		isEmpty, err := isDirectoryEmpty(dir)
		if err != nil {
			log.Println(err)
			continue
		}

		if isEmpty {
			err := os.Remove(dir)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("Removed empty directory: %s\n", dir)
			}
		}
	}
}

// isDirectoryEmpty checks if a directory is empty.
//
// Parameters:
// - path: The path of the directory to check.
//
// Returns:
// - A boolean indicating whether the directory is empty.
// - An error if any error occurs during the check.
//
// Example:
//
//	empty, err := isDirectoryEmpty("/path/to/directory")
//	if err != nil {
//	  log.Fatal(err)
//	}
//	if empty {
//	  fmt.Println("Directory is empty")
//	} else {
//	  fmt.Println("Directory is not empty")
//	}
func isDirectoryEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		// Directory is not empty
		return false, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		// Directory does not exist
		return false, nil
	}

	if errors.Is(err, io.EOF) {
		// Directory is empty
		return true, nil
	}

	return false, err
}

// generateREADME generates a README.md file for the project at the specified path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the generation.
//
// Example:
//
//	err := generateREADME("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func generateREADME(projectPath string) error {
	readmePath := projectPath + "/README.md"

	// Check if the file already exists
	if _, err := os.Stat(readmePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	readmeContent := []byte(`# Project Name

One Paragraph of project description goes here

## Table of Contents

- [Project Name](#project-name)
  - [Table of Contents](#table-of-contents)
  - [About](#about)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [License](#license)
  - [Acknowledgements](#acknowledgements)

## About

Provide a brief introduction or overview of your project.

## Getting Started

Instructions on setting up and running the project.

### Prerequisites

List any software, libraries, or dependencies that need to be installed before running the project.

### Installation

Step-by-step instructions on how to install the project.

## Usage

Provide examples or instructions on how to use the project.

## Contributing

Explain how others can contribute to your project. Include guidelines for pull requests and code style.

## License

Mention the license under which the project is distributed (e.g., MIT License).

## Acknowledgements

Give credits to any external resources or individuals whose work has influenced your project.`)

	err := ioutil.WriteFile(readmePath, readmeContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write README file: %v", err)
	}

	return nil
}

func createColorThemeSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamercolorthemelazy.sty")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
\mode<presentation>

% Main colors
% ------------------
\definecolor{lblack}{HTML}{202124}
\definecolor{lblacktext}{HTML}{4D5156}
\definecolor{lwhite}{HTML}{FFFFFF}
\definecolor{lmain}{HTML}{E8EAED}
\definecolor{alinkblue}{HTML}{1A0DAB}

% Dark mode
%\definecolor{lblack}{HTML}{E8EAED}
%\definecolor{lblacktext}{HTML}{BCC0C3}
%\definecolor{lwhite}{HTML}{202124}
%\definecolor{lmain}{HTML}{202124}
%\definecolor{alinkblue}{HTML}{8AB4F8}

% Accented colors
\definecolor{solgreen}{HTML}{859900}
\definecolor{solblue}{HTML}{268BD2}
\definecolor{solred}{HTML}{DC322F}

% Structure dominant colors

\setbeamercolor*{background canvas}{bg=lwhite, fg=lblack}

\setbeamercolor*{palette primary}{bg=lmain, fg=lblack}
\setbeamercolor*{palette secondary}{bg=lmain, fg=lblack}
\setbeamercolor*{palette tertiary}{bg=lmain, fg=lblack}
\setbeamercolor*{frametitle}{bg=lmain, fg=lblack}

\setbeamercolor{title in head/foot}{bg=lmain, fg=lblack}
\setbeamercolor{section in head/foot}{parent=title in head/foot}
\setbeamercolor{subsection in head/foot}{parent=title in head/foot}

\setbeamercolor{headline}{bg=lmain, fg=lblack}
\setbeamercolor{title in headline}{parent=headline}
\setbeamercolor{author in headline}{parent=headline}
\setbeamercolor{institute in headline}{parent=headline}
\setbeamercolor{institute in footline}{parent=headline}

% Text dominant

\setbeamercolor*{title page header}{bg=lmain, fg=lblack}
\setbeamercolor*{author}{bg=lmain, fg=lblack}
\setbeamercolor*{date}{bg=lmain, fg=lblack}
\setbeamercolor*{structure}{bg=lmain, fg=lblack}
\setbeamercolor{subtitle}{bg=lmain, fg=lblack}

\setbeamercolor*{normal text}{fg=lblacktext}

\setbeamercolor*{titlelike}{bg=lmain, fg=lblack}
\setbeamercolor*{subtitle}{parent=title, fg=lblacktext}
\setbeamercolor*{author}{parent=title, fg=lblacktext}
\setbeamercolor*{date}{parent=title, fg=lblacktext}

\setbeamercolor*{block body}{bg=lwhite, fg=lblacktext}
\setbeamercolor*{block title}{bg=lwhite, fg=solblue}

\setbeamercolor{block body example}{bg=lwhite, fg=lblacktext}
\setbeamercolor{block title example}{bg=lwhite, fg=solgreen}

\setbeamercolor{block body alerted}{bg=lwhite, fg=lblacktext}
\setbeamercolor{block title alerted}{bg=lwhite, fg=solred}

\setbeamercolor{placeholder}{fg=, bg=}
\setbeamercovered{transparent=37}

\mode<all>

	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createColorThemeDSty creates the "beamercolorthemelazyd.sty" file for the color theme in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createColorThemeDSty("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createColorThemeDSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamercolorthemelazyd.sty")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
\mode<presentation>

% Main colors
% ------------------
%\definecolor{lblack}{HTML}{202124}
%\definecolor{lblacktext}{HTML}{4D5156}
%\definecolor{lwhite}{HTML}{FFFFFF}
%\definecolor{lmain}{HTML}{E8EAED}
%\definecolor{alinkblue}{HTML}{1A0DAB}

% Dark mode
\definecolor{lblack}{HTML}{E8EAED}
\definecolor{lblacktext}{HTML}{BCC0C3}
\definecolor{lwhite}{HTML}{202124}
\definecolor{lmain}{HTML}{202124}
\definecolor{alinkblue}{HTML}{8AB4F8}

% Accented colors
\definecolor{solgreen}{HTML}{859900}
\definecolor{solblue}{HTML}{268BD2}
\definecolor{solred}{HTML}{DC322F}

% Structure dominant colors

\setbeamercolor*{background canvas}{bg=lwhite, fg=lblack}

\setbeamercolor*{palette primary}{bg=lmain, fg=lblack}
\setbeamercolor*{palette secondary}{bg=lmain, fg=lblack}
\setbeamercolor*{palette tertiary}{bg=lmain, fg=lblack}
\setbeamercolor*{frametitle}{bg=lmain, fg=lblack}

\setbeamercolor{title in head/foot}{bg=lmain, fg=lblack}
\setbeamercolor{section in head/foot}{parent=title in head/foot}
\setbeamercolor{subsection in head/foot}{parent=title in head/foot}

\setbeamercolor{headline}{bg=lmain, fg=lblack}
\setbeamercolor{title in headline}{parent=headline}
\setbeamercolor{author in headline}{parent=headline}
\setbeamercolor{institute in headline}{parent=headline}
\setbeamercolor{institute in footline}{parent=headline}

% Text dominant

\setbeamercolor*{title page header}{bg=lmain, fg=lblack}
\setbeamercolor*{author}{bg=lmain, fg=lblack}
\setbeamercolor*{date}{bg=lmain, fg=lblack}
\setbeamercolor*{structure}{bg=lmain, fg=lblack}
\setbeamercolor{subtitle}{bg=lmain, fg=lblack}

\setbeamercolor*{normal text}{fg=lblacktext}

\setbeamercolor*{titlelike}{bg=lmain, fg=lblack}
\setbeamercolor*{subtitle}{parent=title, fg=lblacktext}
\setbeamercolor*{author}{parent=title, fg=lblacktext}
\setbeamercolor*{date}{parent=title, fg=lblacktext}

\setbeamercolor*{block body}{bg=lwhite, fg=lblacktext}
\setbeamercolor*{block title}{bg=lwhite, fg=solblue}

\setbeamercolor{block body example}{bg=lwhite, fg=lblacktext}
\setbeamercolor{block title example}{bg=lwhite, fg=solgreen}

\setbeamercolor{block body alerted}{bg=lwhite, fg=lblacktext}
\setbeamercolor{block title alerted}{bg=lwhite, fg=solred}

\setbeamercolor{placeholder}{fg=, bg=}
\setbeamercovered{transparent=37}

\mode<all>
	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createFontThemeSty creates the "beamerfontthemelazy.sty" file for the font theme in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createFontThemeSty("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createFontThemeSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamerfontthemelazy.sty")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
\mode<presentation>

\usefonttheme{professionalfonts}

\usepackage[T1]{fontenc}
\usepackage{newtxtext}
\usepackage{newtxmath}
\usepackage{courier}


%% Allow more stretching
\setlength{\emergencystretch}{3em}

\setbeamerfont{title}{size = \Large, series=\bfseries}
\setbeamerfont{subtitle}{size = \normalsize, series=\mdseries}
\setbeamerfont{author}{size=\small, series=\mdseries}
\setbeamerfont{date}{size=\small, series=\mdseries}
\setbeamerfont{footnote}{size=\tiny}
\setbeamerfont{frametitle}{size = \large, series=\upshape}
\setbeamerfont{block title}{size = \large, series=\upshape}

\mode<all>
	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createOuterThemeSty creates the "beamerouterthemelazy.sty" file for the outer theme in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createOuterThemeSty("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createOuterThemeSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamerouterthemelazy.sty")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
\mode<presentation>

% remove navigation symbols
\setbeamertemplate{navigation symbols}{}

\useoutertheme{tree}

\makeatletter
\setbeamertemplate{headline}
{%
    %\begin{beamercolorbox}[wd=\paperwidth,colsep=1.5pt]{upper separation line head}
    %\end{beamercolorbox}
    
    %\begin{beamercolorbox}[wd=\paperwidth,ht=2.5ex,dp=5ex,%
    %  leftskip=.3cm,rightskip=.3cm plus1fil]{title in head/foot}
    %  \usebeamerfont{title in head/foot}\insertshorttitle
    %\end{beamercolorbox}

    \setbeamertemplate{mini frames}[box]

    \begin{beamercolorbox}[wd=\paperwidth,ht=2.5ex,dp=6ex,%
      leftskip=.3cm,rightskip=.3cm plus1fil]{section in head/foot}
      \insertnavigation{0.6\paperwidth}
      \usebeamerfont{section in head/foot}%
      \ifbeamer@tree@showhooks
        \setbox\beamer@tempbox=\hbox{\insertsectionhead}%
        \ifdim\wd\beamer@tempbox>1pt%
          \hskip2pt\raise1.9pt\hbox{\vrule width0.4pt height1.875ex\vrule width 5pt height0.4pt}%
          \hskip1pt%
        \fi%
      \else%  
        \hskip6pt%
      \fi%
      \insertsectionhead
      \usebeamerfont{subsection in head/foot}%
      \ifbeamer@tree@showhooks
        \setbox\beamer@tempbox=\hbox{\insertsubsectionhead}%
        \ifdim\wd\beamer@tempbox>1pt%
          \ \raise1.9pt\hbox{\vrule width 5pt height0.4pt}%
          \hskip1pt%
        \fi%
      \else%  
        \hskip12pt%
      \fi%
      \insertsubsectionhead
      \hfill
    \end{beamercolorbox}
    \begin{beamercolorbox}[wd=\paperwidth,colsep=1.5pt]{lower separation line head}
    \end{beamercolorbox}
}
\makeatother


\makeatletter
\setbeamertemplate{footline}
{
  \leavevmode%
  \hbox{%
  \begin{beamercolorbox}[wd=.333333\paperwidth,ht=2.25ex,dp=2ex,center]{author in head/foot}%
    \usebeamerfont{author in
head/foot}%
  \insertshortauthor\hspace{1em}\beamer@ifempty{\insertshortinstitute}{}{(\insertshortinstitute)}
  \end{beamercolorbox}%
  \begin{beamercolorbox}[wd=.333333\paperwidth,ht=2.25ex,dp=2ex,center]{title in head/foot}%
    \usebeamerfont{title in head/foot}\insertshorttitle
  \end{beamercolorbox}%
  \begin{beamercolorbox}[wd=.333333\paperwidth,ht=2.25ex,dp=2ex,right]{date in head/foot}%
    \usebeamerfont{date in head/foot}\insertshortdate{}\hspace*{2em}
    \insertframenumber{} / \inserttotalframenumber\hspace*{2ex} 
  \end{beamercolorbox}}%
  \vskip0pt%
}
\makeatother

\mode<all>
	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createInnerThemeSty creates the "beamerinnerthemelazy.sty" file for the inner theme in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createInnerThemeSty("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createInnerThemeSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamerinnerthemelazy.sty")
	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
%% Vertical text alignment:
\DeclareOptionBeamer{c}{ \beamer@centeredtrue  }
\DeclareOptionBeamer{t}{ \beamer@centeredfalse }

%% Theorem numbers:
\DeclareOptionBeamer{unnumbered}{ \def \MATHtheorem {}          }
\DeclareOptionBeamer{numbered}  { \def \MATHtheorem {numbered}  }
\DeclareOptionBeamer{AMS}       { \def \MATHtheorem {ams style} }

\setbeamertemplate{title page}
{
    \AddToShipoutPictureFG*
    {
        \AtPageUpperLeft
        {
            \hspace{1.7 mm}
            \parbox[t][2cm][b]{\textwidth}
            {
                %\includegraphics[scale = 0.125]
                %{fig/logo.png}
            }
        }
    }

    \vbox to \textheight
    {
        \vspace{20 mm}

        \leftskip  = 1.7 mm
        \rightskip = 1.7 mm plus 2 cm

        \usebeamerfont{title}    \structure{\inserttitle}
        \\[0.1ex]
        \usebeamerfont{subtitle} \structure{\insertsubtitle}

        \vspace{5 mm}

        \usebeamerfont{author} \insertauthor
        \hfill
        \newlength{\datewidth}
        \settowidth{\datewidth}{\insertdate}
        \parbox{\datewidth}
        {
            \usebeamerfont{date} \insertdate
        }

        \vspace{5 mm}
        \usebeamerfont{institute} \insertinstitute
        \hfill
    }
}

\newcommand{\TitlePage}
{
    \begin{frame}[plain, noframenumbering]
        \titlepage
    \end{frame}
}

\setbeamertemplate{itemize items}[circle]
	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createMainThemeSty creates the "beamerthemelazy.sty" file for the main theme in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createMainThemeSty("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createMainThemeSty(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "sty", "beamerthemelazy.sty")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `\RequirePackage{tikz}
\RequirePackage{graphicx}
\RequirePackage{etoolbox}
\RequirePackage{xcolor}
\RequirePackage{calc}
\RequirePackage{eso-pic}
\RequirePackage{etoolbox}
\RequirePackage[LGR, T1]{fontenc}
\RequirePackage{thmtools}

\usepackage{sty/beamerinnerthemelazy}
\usepackage{sty/beamerouterthemelazy}
\usepackage{sty/beamerfontthemelazy}

\newif\if@dark
\@darkfalse
\DeclareOption{dark}{\@darktrue}
\newif\if@accent
\@accentfalse
\DeclareOption{accent}{\@accenttrue}
\ProcessOptions

\if@dark
\usepackage{sty/beamercolorthemelazyd}
\else\if@accent
\usepackage{sty/beamercolorthemelazy}
\else
\usepackage{sty/beamercolorthemelazy}
\fi\fi


\hypersetup{
  colorlinks=true,
  urlcolor=alinkblue,
  linkcolor=alinkblue,
}

\mode<all>
	`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createReportTex creates the "report.tex" file for the project report in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createReportTex("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createReportTex(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, "doc", "report", "report.tex")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
\documentclass[%
  beameroptions={ignorenonframetext,11pt,169},
  articleoptions={11pt},
  also={trans,handout,article},
  ]{beamerswitch}
\handoutlayout{nup=3plus,border=1pt}
\articlelayout{maketitle,frametitles=none}
\usepackage[british]{babel}
\mode<article>{
    \usepackage[hmargin=3cm,vmargin=2.5cm]{geometry}
    \usepackage{amsmath, amsthm, amssymb, amsfonts}
    \usepackage[T1]{fontenc}
    \usepackage{listings}
    \usepackage{color} 
    \usepackage{xcolor}  
    \usepackage{hyperref}
    \usepackage{tikz}
    \usepackage{float}
    \usepackage{courier}
    \usepackage{imakeidx}
    \usepackage{biblatex}
    \usepackage{pgfgantt}
    \addbibresource{ref.bib}
    
    \geometry{
    a4paper,
    total={170mm,257mm},
    left=20mm,
    top=20mm,
    }

    \definecolor{linkblue}{HTML}{1A0DAB}

    \newcommand\scalemath[2]{\scalebox{#1}{\mbox{\ensuremath{\displaystyle #2}}}}

    \hypersetup{
        colorlinks=true, 
        linktoc=all,    
        linkcolor=linkblue,  
    }

    \lstset{basicstyle=\footnotesize\ttfamily,breaklines=true}
    \lstset{framextopmargin=50pt,frame=bottomline}
    \lstset{basicstyle=\footnotesize\ttfamily,breaklines=true}
    \lstset{framextopmargin=50pt,frame=bottomline}

    
    \definecolor{solarized-base03}{HTML}{002B36}
    \definecolor{solarized-base02}{HTML}{073642}
    \definecolor{solarized-base01}{HTML}{586e75}
    \definecolor{solarized-base00}{HTML}{657b83}
    \definecolor{solarized-base0}{HTML}{839496}
    \definecolor{solarized-base1}{HTML}{93a1a1}
    \definecolor{solarized-base2}{HTML}{eee8d5}
    \definecolor{solarized-base3}{HTML}{fdf6e3}
    \definecolor{solarized-yellow}{HTML}{b58900}
    \definecolor{solarized-orange}{HTML}{cb4b16}
    \definecolor{solarized-red}{HTML}{dc322f}
    \definecolor{solarized-magenta}{HTML}{d33682}
    \definecolor{solarized-violet}{HTML}{6c71c4}
    \definecolor{solarized-blue}{HTML}{268bd2}
    \definecolor{solarized-cyan}{HTML}{2aa198}
    \definecolor{solarized-green}{HTML}{859900}
    
    \definecolor{backcolour}{HTML}{FFFFFF}

    \lstset{
        backgroundcolor=\color{backcolour},  
        basicstyle=\color{solarized-base00}\ttfamily,
        keywordstyle=\color{solarized-blue},
        stringstyle=\color{solarized-cyan},
        commentstyle=\color{solarized-green},
        numberstyle=\color{solarized-orange},
        identifierstyle=\color{solarized-violet},
        breaklines=true,                 
        captionpos=b,                    
        keepspaces=true,                 
        numbers=left,                    
        numbersep=5pt,                  
        showspaces=false,                
        showstringspaces=false,
        showtabs=true,                  
        tabsize=8,
    }
    \usepackage{newtxtext}
    \usepackage{newtxmath}
    \usepackage{courier}
}
\mode<presentation>{
    \usepackage[orientation=landscape,size=custom,width=16,height=9,scale=0.5,debug]{beamerposter} 
    \usepackage{hyperref}
    \usepackage{graphicx} % Allows including images
    \usepackage{booktabs}
    \usepackage[utf8]{inputenc} % 
    \usepackage{biblatex}
    \usepackage{pgfgantt}

    \usepackage{csquotes}      
    \usepackage{amsmath, amsthm, amssymb, amsfonts}        
    \usepackage{mathtools}    
    \usepackage[absolute, overlay]{textpos} 
    \setlength{\TPHorizModule}{\paperwidth}
    \setlength{\TPVertModule}{\paperheight}
    \usepackage{tikz}
    \usetikzlibrary{overlay-beamer-styles}
    \usepackage{listings}
    
    \usepackage{sty/beamerthemelazy}

    \lstset{basicstyle=\footnotesize\ttfamily,breaklines=true}
    \lstset{framextopmargin=50pt,frame=bottomline}

    
    \definecolor{solarized-base03}{HTML}{002B36}
    \definecolor{solarized-base02}{HTML}{073642}
    \definecolor{solarized-base01}{HTML}{586e75}
    \definecolor{solarized-base00}{HTML}{657b83}
    \definecolor{solarized-base0}{HTML}{839496}
    \definecolor{solarized-base1}{HTML}{93a1a1}
    \definecolor{solarized-base2}{HTML}{eee8d5}
    \definecolor{solarized-base3}{HTML}{fdf6e3}
    \definecolor{solarized-yellow}{HTML}{b58900}
    \definecolor{solarized-orange}{HTML}{cb4b16}
    \definecolor{solarized-red}{HTML}{dc322f}
    \definecolor{solarized-magenta}{HTML}{d33682}
    \definecolor{solarized-violet}{HTML}{6c71c4}
    \definecolor{solarized-blue}{HTML}{268bd2}
    \definecolor{solarized-cyan}{HTML}{2aa198}
    \definecolor{solarized-green}{HTML}{859900}
    
    \definecolor{backcolour}{HTML}{FFFFFF}

    \lstset{
        backgroundcolor=\color{backcolour},  
        basicstyle=\color{solarized-base00}\ttfamily,
        keywordstyle=\color{solarized-blue},
        stringstyle=\color{solarized-cyan},
        commentstyle=\color{solarized-green},
        numberstyle=\color{solarized-orange},
        identifierstyle=\color{solarized-violet},
        breaklines=true,                 
        captionpos=b,                    
        keepspaces=true,                 
        numbers=left,                    
        numbersep=5pt,                  
        showspaces=false,                
        showstringspaces=false,
        showtabs=true,                  
        tabsize=8,
    }

    % \addbibresource{ref.bib}
}
\mode<handout>{
    \usecolortheme{dove}
}

% The title
\title[Subtitle]{Title}

\author[]{Teo Yu Jie}

\institute[Institute]{School}

% Date, can be changed to a custom date
\date{\today}

\begin{document}

\maketitle

\section{Introduction}

\frame{\titlepage}

\tableofcontents


\begin{frame}[plain]
    \frametitle{Title}
    \setcounter{footnote}{0}
    \setcounter{equation}{0}
\end{frame}

\subsection{Background}

((blank)) (ADVICE: PROVIDE A BRIEF OVERVIEW OF THE RELEVANT 
LITERATURE). In retrospect, one can consider ((blank)) as a future 
direction (ADVICE: HIGHLIGHT AN AREA THAT NEEDS FURTHER 
INVESTIGATION OR EXPLORATION). Understanding the background of the 
study is crucial to grasp the context and motivation behind the 
current research. Previous studies have shown ((blank)) (ADVICE: 
SUMMARIZE KEY FINDINGS OR DISCOVERIES IN THE FIELD). However, there 
remains a gap in knowledge regarding ((blank)) (ADVICE: IDENTIFY A 
SPECIFIC GAP OR LIMITATION IN THE EXISTING RESEARCH). Therefore, 
this study aims to ((blank)) (ADVICE: STATE THE RESEARCH OBJECTIVES 
OR PURPOSE). By addressing this gap, the findings of this research 
can contribute to ((blank)) (ADVICE: DESCRIBE THE POTENTIAL IMPACT 
OR BENEFITS OF THE STUDY) and advance our understanding in the field.

\subsection{Problem Statement}

((blank)) (ADVICE: CLEARLY STATE THE PROBLEM). By identifying and 
addressing ((blank)) (ADVICE: SPECIFY THE PROBLEM), this study aims 
to contribute to the understanding and potential solutions for 
((blank)) (ADVICE: DESCRIBE THE IMPACT AND RELEVANCE OF THE 
PROBLEM). Furthermore, this research aims to bridge the existing gap 
in knowledge by ((blank)) (ADVICE: EXPLAIN HOW YOUR RESEARCH 
ADDRESSES THE GAP). Consequently, this investigation will provide 
valuable insights into ((blank)) (ADVICE: STATE THE BENEFITS OF 
SOLVING THE PROBLEM) and offer practical implications for ((blank)) 
(ADVICE: IDENTIFY THE RELEVANT STAKEHOLDERS). In retrospect, one can 
consider ((blank)) (ADVICE: HIGHLIGHT THE URGENCY OR TIMELINESS OF 
SOLVING THE PROBLEM) as a future direction. By addressing the issues 
highlighted in this study, we can pave the way for ((blank)) 
(ADVICE: DISCUSS THE POTENTIAL POSITIVE OUTCOMES).

To set the objectives of this study, ((blank)) (ADVICE: CLEARLY 
DEFINE THE RESEARCH SCOPE OR CONTEXT). The primary objective of this 
research is to ((blank)) (ADVICE: SPECIFY THE MAIN GOAL OR PURPOSE 
OF THE STUDY). By accomplishing this objective, we aim to ((blank)) 
(ADVICE: DESCRIBE THE INTENDED CONTRIBUTION OR OUTCOME). 
Additionally, the secondary objectives of this investigation are 
((blank)) (ADVICE: IDENTIFY THE SUBOBJECTIVES OR SPECIFIC ASPECTS TO 
BE ADDRESSED). These objectives will be pursued through ((blank)) 
(ADVICE: DISCUSS THE METHODOLOGY OR APPROACH). By achieving these 
objectives, we can provide ((blank)) (ADVICE: HIGHLIGHT THE VALUE OR 
BENEFITS OF ACHIEVING THE OBJECTIVES). Reflecting on the future, one 
can consider ((blank)) (ADVICE: SUGGEST POTENTIAL FUTURE DIRECTIONS 
TO BUILD UPON THE OBJECTIVES). Pursuing these future directions will 
enable us to ((blank)) (ADVICE: STATE THE POTENTIAL ENHANCEMENT OR 
EXPANSION OF THE RESEARCH).

\begin{frame}
  \frametitle{Introduction}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Background: Provides contextual information and sets the stage for the research.
    \item Problem Statement: Clearly states the research problem or question being addressed.
    \item Objectives: Outlines the specific goals and aims of the research.
  \end{itemize}
\end{frame}

\section{Literature}

\subsection{General survey}

A general survey of the literature reveals ((blank)) (ADVICE: 
IDENTIFY THE COMMON THEMES OR FINDINGS). Numerous studies have 
investigated ((blank)) (ADVICE: SPECIFY THE MAIN TOPICS OR RESEARCH 
AREAS). These studies have provided valuable insights into ((blank)) 
(ADVICE: HIGHLIGHT THE KNOWLEDGE AND UNDERSTANDING GAINED). 
Furthermore, it is evident that ((blank)) (ADVICE: DESCRIBE THE 
CURRENT STATE OR TRENDS IN THE FIELD). However, ((blank)) (ADVICE: 
POINT OUT THE GAPS OR LIMITATIONS IN THE EXISTING LITERATURE). 
Consequently, this research seeks to address these gaps and extend 
the current body of knowledge by ((blank)) (ADVICE: EXPLAIN THE 
NOVEL ASPECTS OR CONTRIBUTIONS OF YOUR STUDY). By doing so, we aim 
to provide a deeper understanding of ((blank)) (ADVICE: SPECIFY THE 
ASPECTS OR PHENOMENA UNDER INVESTIGATION). Moving forward, it is 
essential to ((blank)) (ADVICE: DISCUSS THE NEED FOR FUTURE RESEARCH 
OR DIRECTIONS). By considering these gaps and potential research 
areas, we can further advance our knowledge of ((blank)) (ADVICE: 
STATE THE RELEVANT TOPIC OR FIELD) and contribute to ((blank)) 
(ADVICE: HIGHLIGHT THE POTENTIAL BENEFITS OR IMPACT OF THE RESEARCH).

\subsection{Concepts and definitions}

In order to establish a strong foundation for this research, it is 
crucial to clarify the key concepts and provide precise definitions. 
((blank)) (ADVICE: INTRODUCE THE MAIN CONCEPTS OR TERMS). The 
concept of ((blank)) is central to this study and refers to 
((blank)) (ADVICE: PROVIDE A CLEAR AND CONCISE DEFINITION). 
Additionally, the notion of ((blank)) is relevant in understanding 
((blank)) (ADVICE: DEFINE ANOTHER KEY CONCEPT AND ITS RELATION TO 
THE RESEARCH). Furthermore, it is essential to define ((blank)) 
(ADVICE: IDENTIFY ANOTHER CONCEPT OR TERM) as it plays a significant 
role in this investigation. It is important to note that these 
definitions are not only limited to their traditional meanings, but 
also encompass ((blank)) (ADVICE: HIGHLIGHT ANY EXTENSIONS OR 
MODIFICATIONS OF THE CONCEPTS IN YOUR RESEARCH CONTEXT). By 
establishing clear definitions and conceptual frameworks, we can 
ensure a common understanding of the terminology used throughout 
this study. Additionally, these conceptual definitions will provide 
a basis for ((blank)) (ADVICE: HINT AT HOW THE CONCEPTS WILL BE 
APPLIED OR ANALYZED IN THE RESEARCH).

\begin{frame}
  \frametitle{Literature Review}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Overview of Relevant Literature: Summarizes the key 
    literature and theories related to the research topic.
    \item Concepts and definitions: Defines important terms and 
    concepts used in the research.
    \item Previous Research Findings: Highlights the main findings 
    of previous studies related to the research question.
    \item Current Knowledge Gap: Identifies areas where further 
    research is needed or where the current knowledge is limited.
  \end{itemize}
\end{frame}

\section{Methods}

\subsection{Research design}

The research design of this study is crucial for achieving the 
objectives and addressing the research questions. ((blank)) (ADVICE: 
CLEARLY STATE THE RESEARCH APPROACH OR STRATEGY). In this 
investigation, a ((blank)) (ADVICE: SPECIFY THE SPECIFIC RESEARCH 
DESIGN, e.g., experimental, qualitative, quantitative) approach will 
be employed to gather and analyze data. This design will enable us 
to ((blank)) (ADVICE: DESCRIBE HOW THE RESEARCH DESIGN WILL HELP IN 
ACHIEVING THE OBJECTIVES OR ADDRESSING THE RESEARCH QUESTIONS). To 
ensure the validity and reliability of the findings, ((blank)) 
(ADVICE: DISCUSS THE METHODOLOGICAL TECHNIQUES OR TOOLS THAT WILL BE 
UTILIZED). The data collection process will involve ((blank)) 
(ADVICE: EXPLAIN THE DATA COLLECTION METHODS OR SOURCES). 
Additionally, ((blank)) (ADVICE: MENTION ANY CONTROLS, VARIABLES, OR 
SAMPLING TECHNIQUES THAT ARE RELEVANT TO YOUR RESEARCH). The data 
will be analyzed through ((blank)) (ADVICE: IDENTIFY THE DATA 
ANALYSIS METHODS OR STATISTICAL TECHNIQUES TO BE APPLIED). By 
employing this research design, we aim to ((blank)) (ADVICE: STATE 
THE EXPECTED OUTCOMES OR CONTRIBUTION OF THE RESEARCH DESIGN). 
Looking ahead, the next step is to ((blank)) (ADVICE: INDICATE THE 
FUTURE STEPS IN THE RESEARCH PROCESS, SUCH AS PILOT TESTING OR 
IMPLEMENTATION OF THE RESEARCH DESIGN).

\subsection{Data collection and analysis}

The data collection process is a critical component of this study, 
ensuring the acquisition of reliable and relevant information. 
((blank)) (ADVICE: CLEARLY STATE THE PURPOSE OF DATA COLLECTION). In 
this research, data will be collected to ((blank)) (ADVICE: DESCRIBE 
THE SPECIFIC OBJECTIVES OF DATA COLLECTION). To obtain a 
comprehensive understanding of the phenomenon under investigation, a 
((blank)) (ADVICE: SPECIFY THE DATA COLLECTION METHOD, e.g., 
surveys, interviews, observations) approach will be employed. This 
method will enable us to ((blank)) (ADVICE: EXPLAIN HOW THE SELECTED 
METHOD WILL CAPTURE THE REQUIRED DATA). The sample population for 
data collection will consist of ((blank)) (ADVICE: INDICATE THE 
CHARACTERISTICS OR CRITERIA FOR SELECTING THE SAMPLE). The data 
collection instruments, such as ((blank)) (ADVICE: MENTION SPECIFIC 
TOOLS OR QUESTIONNAIRES), will be carefully designed to ensure 
clarity and comprehensiveness. Additionally, a pilot test will be 
conducted to ((blank)) (ADVICE: HIGHLIGHT THE IMPORTANCE OF THE 
PILOT TEST IN VALIDATING THE INSTRUMENTS OR METHODS). Furthermore, 
((blank)) (ADVICE: DISCUSS ANY ETHICAL CONSIDERATIONS OR APPROVALS 
REQUIRED FOR DATA COLLECTION). By adhering to rigorous data 
collection procedures, we aim to gather accurate and valid data that 
will serve as a foundation for robust analysis and meaningful 
findings.

The data analysis phase of this research is instrumental in deriving 
meaningful insights and drawing valid conclusions. ((blank)) 
(ADVICE: CLEARLY STATE THE PURPOSE OF DATA ANALYSIS). In this study, 
data will be analyzed to ((blank)) (ADVICE: SPECIFY THE OBJECTIVES 
OR RESEARCH QUESTIONS TO BE ADDRESSED THROUGH DATA ANALYSIS). The 
collected data will undergo a systematic process of ((blank)) 
(ADVICE: DESCRIBE THE DATA ANALYSIS METHOD OR APPROACH, e.g., 
qualitative content analysis, statistical analysis). This analysis 
will involve ((blank)) (ADVICE: MENTION THE SPECIFIC TECHNIQUES, 
TOOLS, OR SOFTWARE UTILIZED). The data will be examined for 
patterns, trends, and relationships, enabling us to ((blank)) 
(ADVICE: EXPLAIN HOW THE DATA ANALYSIS WILL UNCOVER INSIGHTS OR 
ANSWER THE RESEARCH QUESTIONS). Additionally, ((blank)) (ADVICE: 
DISCUSS ANY DATA TRANSFORMATION OR PREPROCESSING STEPS THAT WILL BE 
APPLIED). The findings obtained from the data analysis will be 
meticulously interpreted and synthesized to ((blank)) (ADVICE: 
INDICATE HOW THE FINDINGS WILL BE ORGANIZED AND PRESENTED). It is 
crucial to note that this research employs a ((blank)) (ADVICE: 
HIGHLIGHT THE INNOVATIVE ASPECTS OR NOVEL APPROACH IN DATA 
ANALYSIS). The outcomes of this data analysis will provide a 
comprehensive understanding of ((blank)) (ADVICE: SPECIFY THE 
PHENOMENON OR CONTEXT UNDER STUDY) and contribute to ((blank)) 
(ADVICE: STATE THE RELEVANT FIELD OR KNOWLEDGE DOMAIN) in a 
significant and impactful manner.

\subsection{Variables and measures}

In this section, we discuss the variables and measures employed in 
this research, as they are fundamental to understanding the 
phenomena under investigation. ((blank)) (ADVICE: INTRODUCE THE MAIN 
VARIABLES OF INTEREST). The primary variable in this study is 
((blank)) (ADVICE: DEFINE THE MAIN VARIABLE CLEARLY). It will be 
measured using ((blank)) (ADVICE: SPECIFY THE MEASUREMENT METHOD OR 
SCALE). Additionally, ((blank)) (ADVICE: IDENTIFY OTHER RELEVANT 
VARIABLES THAT WILL BE CONSIDERED). These variables, such as 
((blank)) (ADVICE: MENTION THE ADDITIONAL VARIABLES), will be 
measured through ((blank)) (ADVICE: DESCRIBE THE MEASUREMENT METHODS 
OR INDICATORS). It is important to note that the selection of 
appropriate measures is crucial for ensuring ((blank)) (ADVICE: 
DISCUSS THE VALIDITY AND RELIABILITY OF THE MEASURES). Furthermore, 
((blank)) (ADVICE: ADDRESS ANY CONTROL VARIABLES OR CONFOUNDING 
FACTORS THAT WILL BE ACCOUNTED FOR). By considering these variables 
and measures, we aim to capture a comprehensive picture of ((blank)) 
(ADVICE: STATE THE PHENOMENON OR RELATIONSHIPS UNDER STUDY). It is 
worth noting that the novel aspect of our research lies in ((blank)) 
(ADVICE: HIGHLIGHT THE INNOVATIVE OR UNIQUE ASPECTS OF THE VARIABLES 
OR MEASURES USED). These variables and measures will provide 
valuable insights and contribute to the advancement of knowledge in 
the field of ((blank)) (ADVICE: SPECIFY THE RELEVANT FIELD).

\begin{frame}
  \frametitle{Methods}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Research Design: Describes the overall approach and methodology employed in the research.
    \item Data Collection: Explains how data was gathered or collected for the study.
    \item Data Analysis: Describes the methods used to analyze the collected data.
    \item Variables and Measures: Specifies the variables studied and the measures used to assess them.
  \end{itemize}
\end{frame}


\section{Results}

\subsection{Presentation of findings}

The presentation of findings is a crucial component of this research, as it provides a comprehensive overview of the results obtained. ((blank)) (ADVICE: 
INTRODUCE THE FINDINGS SECTION AND SET THE CONTEXT). In this section, we present the key findings derived from the data analysis process. The 
findings will be organized and presented in a logical and coherent manner, focusing on ((blank)) (ADVICE: IDENTIFY THE MAIN THEMES, TRENDS, OR 
PATTERNS IN THE FINDINGS). Additionally, visual aids such as charts, graphs, and tables will be utilized to ((blank)) (ADVICE: EMPHASIZE THE 
IMPORTANCE OF VISUAL REPRESENTATION FOR CLARITY AND EASE OF UNDERSTANDING). The findings reveal ((blank)) (ADVICE: PROVIDE A SUMMARY OF THE 
MAIN FINDINGS). Moreover, it is worth noting that our research has uncovered ((blank)) (ADVICE: DISCUSS ANY UNIQUE OR UNEXPECTED FINDINGS 
THAT CONTRIBUTE TO THE NOVELTY OF YOUR WORK). These findings align with the research objectives and contribute to the existing body 
of knowledge in the field of ((blank)) (ADVICE: SPECIFY THE RELEVANT FIELD). Furthermore, ((blank)) (ADVICE: DISCUSS THE IMPLICATIONS 
OR SIGNIFICANCE OF THE FINDINGS IN RELATION TO THE RESEARCH QUESTIONS OR OBJECTIVES). Overall, the findings of this study provide 						valuable insights and lay the foundation for further discussion and analysis in subsequent sections.


\subsection{Data interpretation}

Data interpretation plays a crucial role in extracting meaningful 
insights from the collected data and understanding their 
implications within the context of the research. ((blank)) (ADVICE: 
INTRODUCE THE IMPORTANCE OF DATA INTERPRETATION). In this section, 
we analyze and interpret the findings derived from the data analysis 
process. The interpretation process involves a thorough examination 
of the data to identify patterns, trends, and relationships. By 
scrutinizing the data in depth, ((blank)) (ADVICE: DESCRIBE HOW DATA 
INTERPRETATION HELPS TO UNCOVER MEANINGFUL INSIGHTS OR ANSWER THE 
RESEARCH QUESTIONS). Additionally, we consider the theoretical 
frameworks and existing literature to provide a comprehensive 
understanding of the findings. The interpretation of the data 
enables us to ((blank)) (ADVICE: HIGHLIGHT THE SIGNIFICANCE OR 
IMPLICATIONS OF THE FINDINGS). Moreover, we examine any 
discrepancies or outliers that may arise and provide plausible 
explanations or potential factors contributing to these 
observations. This process facilitates the identification of key 
findings, underlying mechanisms, and potential areas for further 
investigation. ((blank)) (ADVICE: EMPHASIZE THE NOVEL ASPECTS OF 
YOUR INTERPRETATION THAT CONTRIBUTE TO THE OVERALL CONTRIBUTION OF 
YOUR WORK). The interpreted findings provide valuable insights that 
contribute to the existing knowledge base and address the research 
objectives. By presenting a comprehensive and well-grounded 
interpretation, this research contributes to the understanding of 
((blank)) (ADVICE: SPECIFY THE PHENOMENON, FIELD, OR CONTEXT UNDER 
STUDY) and offers practical implications for ((blank)) (ADVICE: 
IDENTIFY THE RELEVANT STAKEHOLDERS OR APPLICATIONS).

\subsection{Statistical analysis}

Statistical analysis is a crucial component of this research, as it 
provides a systematic approach to analyze and interpret the data 
collected. ((blank)) (ADVICE: INTRODUCE THE IMPORTANCE OF 
STATISTICAL ANALYSIS). In this section, we employ various 
statistical techniques to explore the relationships, patterns, and 
trends present in the data. The analysis begins with ((blank)) 
(ADVICE: IDENTIFY THE INITIAL STEPS OR PRELIMINARY ANALYSES, SUCH AS 
DESCRIPTIVE STATISTICS OR DATA CLEANING). Subsequently, we conduct 
((blank)) (ADVICE: SPECIFY THE SPECIFIC STATISTICAL TESTS, MODELS, 
OR PROCEDURES) to examine the associations and potential causality 
between variables. These statistical analyses enable us to ((blank)) 
(ADVICE: DESCRIBE HOW STATISTICAL ANALYSIS HELPS IN ANSWERING THE 
RESEARCH QUESTIONS OR OBJECTIVES). Furthermore, we assess the 
statistical significance of the findings and determine the strength 
of the relationships observed. The results are reported using 
appropriate statistical measures such as ((blank)) (ADVICE: MENTION 
THE RELEVANT STATISTICAL MEASURES, E.G., P-VALUES, EFFECT SIZES). 
Additionally, we explore potential confounding factors or 
interactions that may influence the outcomes. It is important to 
note that the statistical analysis conducted in this study is 
innovative in ((blank)) (ADVICE: HIGHLIGHT THE UNIQUE ASPECTS OR 
NOVEL APPLICATIONS OF YOUR STATISTICAL ANALYSIS). The statistical 
findings provide valuable insights into ((blank)) (ADVICE: SPECIFY 
THE PHENOMENON OR FIELD UNDER STUDY) and contribute to the overall 
understanding of ((blank)) (ADVICE: STATE THE RELEVANT FIELD OR 
KNOWLEDGE DOMAIN). By employing rigorous statistical analysis, we 
ensure the reliability and validity of the research findings, 
enhancing the robustness and impact of this study.

\begin{frame}
  \frametitle{Results and Analysis}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Presentation of Findings: Presents the results of the research in a clear and concise manner.
    \item Data Interpretation: Provides an explanation and interpretation of the research findings.
    \item Statistical Analysis: Describes any statistical tests or analyses conducted on the data.
  \end{itemize}
\end{frame}

\section{Discussion}

\subsection{Comparison with previous research}

In this section, we compare our research findings with those of 
previous studies to gain insights into the existing body of 
knowledge and identify novel contributions. ((blank)) (ADVICE: 
INTRODUCE THE IMPORTANCE OF COMPARISON WITH PREVIOUS RESEARCH). The 
comparison involves examining the similarities and differences 
between our findings and the results reported in prior literature. 
Through this comparative analysis, we aim to ((blank)) (ADVICE: 
STATE THE OBJECTIVES OF THE COMPARISON, SUCH AS VALIDATING OR 
EXTENDING PREVIOUS FINDINGS). Notably, our research offers a unique 
perspective by ((blank)) (ADVICE: HIGHLIGHT THE INNOVATIVE OR 
DISTINCT ASPECTS OF YOUR WORK). The comparison reveals that 
((blank)) (ADVICE: DESCRIBE THE KEY SIMILARITIES OR DIFFERENCES 
OBSERVED). Furthermore, it is important to consider the contextual 
factors that may account for any variations in findings. By 
critically analyzing and interpreting the similarities and 
discrepancies, we can provide a more comprehensive understanding of 
((blank)) (ADVICE: SPECIFY THE PHENOMENON, TOPIC, OR FIELD UNDER 
DISCUSSION). This comparative analysis not only helps us evaluate 
the consistency and generalizability of our results but also 
contributes to the advancement of knowledge in ((blank)) (ADVICE: 
STATE THE RELEVANT FIELD). The integration of our findings with 
existing research paves the way for further exploration and 
highlights the unique contributions of our study to the field.

\subsection{Implications and significance}

The implications and significance of this research are multifaceted 
and far-reaching, with important implications for both theory and 
practice. ((blank)) (ADVICE: INTRODUCE THE IMPORTANCE OF DISCUSSING 
IMPLICATIONS AND SIGNIFICANCE). Firstly, the findings of this study 
contribute to the theoretical understanding of ((blank)) (ADVICE: 
SPECIFY THE PHENOMENON OR FIELD UNDER STUDY) by ((blank)) (ADVICE: 
HIGHLIGHT THE NOVEL CONCEPTS, MODELS, OR INSIGHTS THAT ADVANCE 
THEORETICAL KNOWLEDGE). These findings challenge existing 
assumptions and provide new perspectives that extend the current 
body of knowledge. Secondly, this research has practical 
implications for ((blank)) (ADVICE: IDENTIFY THE RELEVANT 
STAKEHOLDERS, PRACTITIONERS, OR POLICYMAKERS). The insights gained 
from this study can inform decision-making processes and guide the 
development of effective strategies and interventions. ((blank)) 
(ADVICE: DISCUSS THE SPECIFIC WAYS IN WHICH THE FINDINGS CAN BE 
APPLIED OR ADD VALUE TO PRACTICE). Additionally, the innovative 
approaches and methodologies utilized in this research offer 
methodological contributions to the field of ((blank)) (ADVICE: 
SPECIFY THE RELEVANT FIELD OR RESEARCH DOMAIN). Lastly, this study 
opens up new avenues for future research by ((blank)) (ADVICE: 
HIGHLIGHT THE UNEXPLORED AREAS OR QUESTIONS THAT EMERGE FROM THE 
CURRENT RESEARCH). These future investigations can build upon our 
findings and delve deeper into the complexities of ((blank)) 
(ADVICE: STATE THE PHENOMENON OR TOPIC UNDER STUDY). In summary, the 
implications and significance of this research lie in its ability to 
advance theory, inform practice, contribute methodologically, and 
guide future research, ultimately making a valuable and lasting 
impact in the field of ((blank)) (ADVICE: SPECIFY THE RELEVANT FIELD 
OR KNOWLEDGE DOMAIN).

\subsection{Limitations and future research directions}

It is important to acknowledge the limitations inherent in this 
study, as they shape the scope and generalizability of the findings. 
((blank)) (ADVICE: INTRODUCE THE DISCUSSION OF LIMITATIONS). One 
limitation of this research is ((blank)) (ADVICE: IDENTIFY A 
SPECIFIC LIMITATION, E.G., SAMPLING BIAS, SMALL SAMPLE SIZE). This 
limitation may have influenced the representativeness of the sample 
and the generalizability of the results. Additionally, ((blank)) 
(ADVICE: MENTION ANOTHER LIMITATION, SUCH AS DATA COLLECTION 
CONSTRAINTS OR RESOURCE LIMITATIONS). These limitations could have 
impacted the comprehensiveness or accuracy of the data collected. 
Furthermore, ((blank)) (ADVICE: DISCUSS ANOTHER RELEVANT LIMITATION, 
E.G., POTENTIAL CONFOUNDING VARIABLES). The presence of confounding 
variables may have introduced bias or affected the internal validity 
of the study. It is also important to note that ((blank)) (ADVICE: 
HIGHLIGHT ANY SPECIFIC ASSUMPTIONS MADE OR SIMPLIFICATIONS ADOPTED). 
These assumptions or simplifications may have implications for the 
generalizability or applicability of the findings in real-world 
contexts. Despite these limitations, this research offers valuable 
insights and serves as a foundation for future investigations. By 
acknowledging these limitations, we foster transparency and ensure 
that readers can accurately interpret the scope and implications of 
our study.

The findings and implications of this study provide a foundation for 
future research in several promising directions. ((blank)) (ADVICE: 
INTRODUCE THE IMPORTANCE OF FUTURE RESEARCH DIRECTIONS). Firstly, 
further investigation is warranted to ((blank)) (ADVICE: STATE A 
SPECIFIC AREA OR ASPECT THAT REQUIRES FURTHER EXPLORATION). This 
includes delving deeper into ((blank)) (ADVICE: SPECIFY THE 
PHENOMENON, TOPIC, OR CONTEXT UNDER STUDY) to gain a more 
comprehensive understanding of its underlying mechanisms or 
dynamics. Additionally, future research could benefit from ((blank)) 
(ADVICE: IDENTIFY A METHOD OR APPROACH THAT COULD BE EMPLOYED TO 
EXTEND THE CURRENT STUDY). For instance, employing longitudinal or 
experimental designs may provide insights into causality or temporal 
relationships. Furthermore, it would be valuable to ((blank)) 
(ADVICE: SUGGEST AN AREA OR ASPECT THAT COULD BE EXPLORED FROM A 
DIFFERENT PERSPECTIVE OR USING ALTERNATIVE METHODS). This could 
involve interdisciplinary collaborations, exploring diverse 
populations, or integrating novel theoretical frameworks. It is also 
important to address the limitations identified in this study 
through ((blank)) (ADVICE: RECOMMEND STRATEGIES TO OVERCOME THE 
LIMITATIONS AND ENHANCE THE VALIDITY OR GENERALIZABILITY OF FUTURE 
RESEARCH). By addressing these limitations, future research can 
build upon the foundation laid by this study and expand our 
knowledge in ((blank)) (ADVICE: SPECIFY THE RELEVANT FIELD OR 
KNOWLEDGE DOMAIN). Overall, the findings of this study offer a 
springboard for future investigations that have the potential to 
advance theory, inform practice, and contribute to the existing body 
of knowledge in the field of ((blank)) (ADVICE: SPECIFY THE RELEVANT 
FIELD OR DOMAIN).

\begin{frame}
  \frametitle{Discussion}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Summary of Findings: Summarizes the main findings of the research.
    \item Comparison with Previous Research: Compares and contrasts the current findings with the results of previous studies.
    \item Implications and Significance: Discusses the implications and significance of the research findings.
    \item Limitations and Future Research Directions: Identifies limitations of the study and suggests areas for future research.
  \end{itemize}
\end{frame}

\section{Conclusion}

In conclusion, it is evident that this research has made significant 
contributions to the field of ((blank)). Through our rigorous 
analysis and interpretation of the data, we have gained valuable 
insights into ((blank)). The findings of this study have several 
implications for both theory and practice. ((blank)) (ADVICE: 

SPECIFY THE NOVEL ASPECTS OF YOUR WORK). These unique findings not 
only enhance our understanding of ((blank)) but also provide a fresh 
perspective on ((blank)). Furthermore, our research has identified 
potential areas for future investigation. ((blank)) (ADVICE: STATE 
THE RECOMMENDATIONS OR NEXT STEPS FOR FUTURE RESEARCH). By 
addressing these research gaps, researchers can further deepen their 
understanding of ((blank)). It is important to acknowledge the 
limitations of this study, such as ((blank)). (ADVICE: DESCRIBE A 
SPECIFIC LIMITATION). Nonetheless, these limitations provide 
opportunities for future studies to build upon and overcome these 
challenges. Overall, the findings of this research significantly 
contribute to the field and provide a solid foundation for further 
advancements in ((blank)).


\subsection{Summary of findings}

\begin{frame}
  \frametitle{Conclusion}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Summary of the Study: Summarizes the main points of the research paper.
    \item Contributions and Recommendations: Highlights the contributions of the research and provides recommendations for further action or study.
    \item Final Thoughts: Concludes the presentation with any final remarks or reflections.
  \end{itemize}
\end{frame}

\section*{References}

\begin{frame}
  \frametitle{References}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Lists the references cited in the research paper.
  \end{itemize}
\end{frame}

\section*{Appendices}

\begin{frame}
  \frametitle{Appendix (if applicable)}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Includes any supplementary materials or additional information that supports the research.
  \end{itemize}
\end{frame}

\begin{frame}[fragile]
  \frametitle{Appendix (if applicable)}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}

  \begin{lstlisting}[language=Python]
import math

def calculate_circle_area(radius):
  area = math.pi * radius**2
  return area

circle_radius = 3
area = calculate_circle_area(circle_radius)
print(f"The area of the circle is: {area}")\end{lstlisting}
\end{frame}

\end{document}
`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createGitignore creates the ".gitignore" file for the project in the specified project path.
//
// Parameters:
// - projectPath: The path of the project.
//
// Returns:
// - An error if any error occurs during the creation.
//
// Example:
//
//	err := createGitignore("/path/to/project")
//	if err != nil {
//	  log.Fatal(err)
//	}
func createGitignore(projectPath string) error {
	// Create the file path using filepath.Join
	filePath := filepath.Join(projectPath, ".gitignore")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
# Exclude sensitive information and credentials
*.env
*.pem
*.key
*.cer

# Exclude configuration files
config/
settings/
*.config

# Exclude log files
logs/
*.log

# Exclude temporary files and cache
tmp/
cache/
*.tmp

# Exclude build output
bin/
build/
dist/
*.exe
*.dll
*.o

# Exclude dependency files
node_modules/
vendor/
venv/

# Exclude operating system files
.DS_Store
Thumbs.db
.idea/

# Exclude personal user files
.bash_history
.vimrc
*.bak


# Exclude ANSYS result files
*.rst
*.db
*.dbb
*.err
*.esav
*.full
*.h3d
*.info
*.ldhi
*.ldpost
*.lff
*.load
*.mac
*.nas
*.odb
*.out
*.plo
*.plot
*.plt
*.pmacr
*.prm
*.prt
*.prtinfo
*.puz
*.read
*.readrst
*.results
*.resu
*.rfl
*.rstt
*.rth
*.rzh
*.ses
*.stat
*.stt
*.sum
*.tbin
*.tmc
*.trd
*.wdb
*.wrl
*.xy

## Core latex/pdflatex auxiliary files:
*.aux
*.lof
*.log
*.lot
*.fls
*.out
*.toc
*.fmt
*.fot
*.cb
*.cb2
.*.lb

## Intermediate documents:
*.dvi
*.xdv
*-converted-to.*
# these rules might exclude image files for figures etc.
# *.ps
# *.eps
# *.pdf

## Generated if empty string is given at "Please type another file name for output:"
.pdf

## Bibliography auxiliary files (bibtex/biblatex/biber):
*.bbl
*.bcf
*.blg
*-blx.aux
*-blx.bib
*.run.xml

## Build tool auxiliary files:
*.fdb_latexmk
*.synctex
*.synctex(busy)
*.synctex.gz
*.synctex.gz(busy)
*.pdfsync

## Build tool directories for auxiliary files
# latexrun
latex.out/

## Auxiliary and intermediate files from other packages:
# algorithms
*.alg
*.loa

# achemso
acs-*.bib

# amsthm
*.thm

# beamer
*.nav
*.pre
*.snm
*.vrb

# changes
*.soc

# comment
*.cut

# cprotect
*.cpt

# elsarticle (documentclass of Elsevier journals)
*.spl

# endnotes
*.ent

# fixme
*.lox

# feynmf/feynmp
*.mf
*.mp
*.t[1-9]
*.t[1-9][0-9]
*.tfm

#(r)(e)ledmac/(r)(e)ledpar
*.end
*.?end
*.[1-9]
*.[1-9][0-9]
*.[1-9][0-9][0-9]
*.[1-9]R
*.[1-9][0-9]R
*.[1-9][0-9][0-9]R
*.eledsec[1-9]
*.eledsec[1-9]R
*.eledsec[1-9][0-9]
*.eledsec[1-9][0-9]R
*.eledsec[1-9][0-9][0-9]
*.eledsec[1-9][0-9][0-9]R

# glossaries
*.acn
*.acr
*.glg
*.glo
*.gls
*.glsdefs
*.lzo
*.lzs
*.slg
*.slo
*.sls

# uncomment this for glossaries-extra (will ignore makeindex's style files!)
# *.ist

# gnuplot
*.gnuplot
*.table

# gnuplottex
*-gnuplottex-*

# gregoriotex
*.gaux
*.glog
*.gtex

# htlatex
*.4ct
*.4tc
*.idv
*.lg
*.trc
*.xref

# hyperref
*.brf

# knitr
*-concordance.tex
# TODO Uncomment the next line if you use knitr and want to ignore its generated tikz files
# *.tikz
*-tikzDictionary

# listings
*.lol

# luatexja-ruby
*.ltjruby

# makeidx
*.idx
*.ilg
*.ind

# minitoc
*.maf
*.mlf
*.mlt
*.mtc[0-9]*
*.slf[0-9]*
*.slt[0-9]*
*.stc[0-9]*

# minted
_minted*
*.pyg

# morewrites
*.mw

# newpax
*.newpax

# nomencl
*.nlg
*.nlo
*.nls

# pax
*.pax

# pdfpcnotes
*.pdfpc

# sagetex
*.sagetex.sage
*.sagetex.py
*.sagetex.scmd

# scrwfile
*.wrt

# svg
svg-inkscape/

# sympy
*.sout
*.sympy
sympy-plots-for-*.tex/

# pdfcomment
*.upa
*.upb

# pythontex
*.pytxcode
pythontex-files-*/

# tcolorbox
*.listing

# thmtools
*.loe

# TikZ & PGF
*.dpth
*.md5
*.auxlock

# titletoc
*.ptc

# todonotes
*.tdo

# vhistory
*.hst
*.ver

# easy-todo
*.lod

# xcolor
*.xcp

# xmpincl
*.xmpi

# xindy
*.xdy

# xypic precompiled matrices and outlines
*.xyc
*.xyd

# endfloat
*.ttt
*.fff

# Latexian
TSWLatexianTemp*

## Editors:
# WinEdt
*.bak
*.sav

# Texpad
.texpadtmp

# LyX
*.lyx~

# Kile
*.backup

# gummi
.*.swp

# KBibTeX
*~[0-9]*

# TeXnicCenter
*.tps

# auto folder when using emacs and auctex
./auto/*
*.el

# expex forward references with \gathertags
*-tags.tex

# standalone packages
*.sta

# Makeindex log files
*.lpz

# xwatermark package
*.xwm

# REVTeX puts footnotes in the bibliography by default, unless the nofootinbib
# option is specified. Footnotes are the stored in a file with suffix Notes.bib.
# Uncomment the next line to have this generated file ignored.
#*Notes.bib


# Exclude documentation and notes
docs/

# Exclude any other sensitive files or directories

`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}

// createJupyterTemplate creates a Jupyter Notebook template file in the specified project path.
// The template includes sections for introduction, data preparation, exploratory data analysis (EDA),
// results and conclusion, and references. Each section contains explanatory markdown cells and code cells
// that can be filled in with relevant content.
//
// Parameters:
//   - projectPath: The path to the project directory where the Jupyter Notebook template will be created.
//
// Returns:
//   - An error if there was a problem creating the template file or directory.
//
// Example:
//
//	err := createJupyterTemplate("/path/to/project")
//	if err != nil {
//	  fmt.Printf("Error creating Jupyter Notebook template: %v", err)
//	}
func createJupyterTemplate(projectPath string) error {
	// Create the file path using filepath.Join

	egStyDir := filepath.Join(projectPath, "eg", "notebook")
	err := os.MkdirAll(egStyDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(projectPath, "eg", "notebook", "notebook.ipynb")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write content to the file
	content := `
	{
		"cells": [
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"# Template Notebook"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"## 1. Introduction\n",
					"- Briefly explain the purpose of the study.\n",
					"- Describe the dataset and variables used in the analysis.\n",
					"- Provide an overview of the steps involved in the analysis."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# import numpy as np                      # Numerical computing library\n",
					"# import pandas as pd                     # Data manipulation and analysis library\n",
					"# import matplotlib.pyplot as plt        # Data visualization library\n",
					"# import seaborn as sns                   # Enhanced data visualization library\n",
					"# import scipy.stats as stats             # Statistical functions and tests\n",
					"# import sklearn                         # Machine learning library\n",
					"# import tensorflow as tf                # Deep learning library\n",
					"# import keras                           # Deep learning library\n",
					"# from keras.models import Sequential   # Sequential model for neural networks\n",
					"# from keras.layers import Dense        # Dense layer for neural networks\n",
					"# import statsmodels.api as sm          # Statistical models and tests\n",
					"# import plotly.express as px           # Interactive plotting library\n",
					"# import plotly.graph_objects as go     # Interactive plotting library\n",
					"# import networkx as nx                 # Network analysis library\n",
					"# import datetime                       # Date and time manipulation\n",
					"# import os                             # Operating system interaction"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"## 2. Data Preparation\n",
					"- Import necessary libraries and load the dataset.\n",
					"- Perform data cleaning and preprocessing, including handling missing values and outliers.\n",
					"- Split the dataset into training and testing sets if applicable.\n",
					"- Normalize or scale the variables if necessary."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# Data Preparation Code"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"## 3. Exploratory Data Analysis (EDA)\n",
					"- Visualize the dataset using appropriate plots and charts.\n",
					"- Calculate descriptive statistics such as mean, median, and standard deviation.\n",
					"- Explore the relationships between variables through correlation analysis or scatter plots.\n",
					"- Identify any patterns or trends in the data.\n"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# Exploratory Data Analysis Code"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"## 4. Operations\n",
					"\n",
					"- Hypothesis Testing: Statistical inference method to evaluate a claim about a population based on sample data."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"import scipy.stats as stats\n",
					"\n",
					"def perform_hypothesis_test(data, null_hypothesis, alternative_hypothesis, alpha=0.05):\n",
					"    \"\"\"\n",
					"    Perform a hypothesis test on the given data.\n",
					"\n",
					"    Args:\n",
					"        data (list or numpy.ndarray): The observed data.\n",
					"        null_hypothesis (float): The null hypothesis value to be tested.\n",
					"        alternative_hypothesis (str): The alternative hypothesis ('two-sided', 'greater', or 'less').\n",
					"        alpha (float, optional): The significance level (default is 0.05).\n",
					"\n",
					"    Returns:\n",
					"        tuple: A tuple containing the test statistic and the p-value.\n",
					"    \"\"\"\n",
					"    # Compute the test statistic and p-value based on the alternative hypothesis\n",
					"    if alternative_hypothesis == 'two-sided':\n",
					"        test_statistic, p_value = stats.ttest_1samp(data, null_hypothesis)\n",
					"    elif alternative_hypothesis == 'greater':\n",
					"        test_statistic, p_value = stats.ttest_1samp(data, null_hypothesis, alternative='greater')\n",
					"    elif alternative_hypothesis == 'less':\n",
					"        test_statistic, p_value = stats.ttest_1samp(data, null_hypothesis, alternative='less')\n",
					"    else:\n",
					"        raise ValueError(\"Invalid alternative hypothesis. Supported values are 'two-sided', 'greater', or 'less'.\")\n",
					"\n",
					"    # Determine if the null hypothesis should be rejected based on the p-value and significance level\n",
					"    if p_value < alpha:\n",
					"        result = \"Reject null hypothesis\"\n",
					"    else:\n",
					"        result = \"Fail to reject null hypothesis\"\n",
					"\n",
					"    return test_statistic, p_value, result\n",
					"\n",
					"\n",
					"perform_hypothesis_test([1, 1, 1, 1, 1], 0.2, 'two-sided')"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Regression Analysis: Statistical modeling technique to investigate the relationship between a dependent variable and one or more independent variables."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"import pandas as pd\n",
					"import statsmodels.api as sm\n",
					"\n",
					"def perform_multivariate_regression(X, y):\n",
					"    \"\"\"\n",
					"    Perform multivariate regression analysis on the given input and output variables.\n",
					"\n",
					"    Args:\n",
					"        X (numpy.ndarray): The input variables with shape (n_samples, n_features).\n",
					"        y (numpy.ndarray): The output variable with shape (n_samples,).\n",
					"\n",
					"    Returns:\n",
					"        statsmodels.regression.linear_model.RegressionResultsWrapper: The regression results object.\n",
					"    \"\"\"\n",
					"    # Add a constant term to the input variables\n",
					"    X = sm.add_constant(X)\n",
					"\n",
					"    # Perform the multivariate regression analysis\n",
					"    model = sm.OLS(y, X)\n",
					"    results = model.fit()\n",
					"\n",
					"    return results\n",
					"\n",
					"# Example input and output variables\n",
					"X = np.array([[1, 2], [3, 4], [5, 6]])\n",
					"y = np.array([3, 5, 7])\n",
					"\n",
					"# Perform multivariate regression analysis\n",
					"results = perform_multivariate_regression(X, y)\n",
					"\n",
					"# Print the regression results\n",
					"print(results.summary())"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Time Series Analysis: Analyzing and modeling data points collected over time to uncover patterns, trends, and make predictions."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"def perform_time_series_analysis(data, model='ARIMA', order=(1, 0, 0)):\n",
					"    \"\"\"\n",
					"    Perform time series analysis on the given data using the specified model.\n",
					"\n",
					"    Args:\n",
					"        data (pandas.Series): The time series data.\n",
					"        model (str, optional): The time series model to use ('ARIMA', 'SARIMAX', etc.).\n",
					"                              Default is 'ARIMA'.\n",
					"        order (tuple, optional): The order of the model (default is (1, 0, 0)).\n",
					"\n",
					"    Returns:\n",
					"        statsmodels.tsa.arima.model.ARIMAResults: The time series analysis results object.\n",
					"    \"\"\"\n",
					"    # Perform time series analysis using the specified model\n",
					"    if model == 'ARIMA':\n",
					"        model = sm.tsa.ARIMA(data, order=order)\n",
					"    elif model == 'SARIMAX':\n",
					"        model = sm.tsa.SARIMAX(data, order=order)\n",
					"    else:\n",
					"        raise ValueError(\"Invalid time series model. Supported models are 'ARIMA', 'SARIMAX', etc.\")\n",
					"\n",
					"    results = model.fit()\n",
					"\n",
					"    return results\n",
					"\n",
					"# Example time series data\n",
					"data = pd.Series([10, 15, 20, 25, 30])\n",
					"\n",
					"# Perform time series analysis\n",
					"results = perform_time_series_analysis(data, model='ARIMA', order=(1, 0, 0))\n",
					"\n",
					"# Print the time series analysis results\n",
					"print(results.summary())"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Classification: Predictive modeling technique that assigns input data points to predefined classes or categories based on their features."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"from sklearn.model_selection import train_test_split\n",
					"from sklearn.linear_model import LogisticRegression\n",
					"from sklearn.metrics import classification_report\n",
					"\n",
					"def perform_classification(X, y, test_size=0.2):\n",
					"    \"\"\"\n",
					"    Perform classification on the given input and output variables.\n",
					"\n",
					"    Args:\n",
					"        X (numpy.ndarray): The input variables with shape (n_samples, n_features).\n",
					"        y (numpy.ndarray): The output variable with shape (n_samples,).\n",
					"        test_size (float, optional): The proportion of the data to use for testing (default is 0.2).\n",
					"\n",
					"    Returns:\n",
					"        str: The classification report.\n",
					"    \"\"\"\n",
					"    # Split the data into training and testing sets\n",
					"    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=test_size, random_state=42)\n",
					"\n",
					"    # Create and train the classifier\n",
					"    classifier = LogisticRegression()\n",
					"    classifier.fit(X_train, y_train)\n",
					"\n",
					"    # Make predictions on the test set\n",
					"    y_pred = classifier.predict(X_test)\n",
					"\n",
					"    # Generate the classification report\n",
					"    report = classification_report(y_test, y_pred)\n",
					"\n",
					"    return report\n",
					"\n",
					"import numpy as np\n",
					"from sklearn.datasets import make_classification\n",
					"\n",
					"# Generate a synthetic classification dataset\n",
					"X, y = make_classification(\n",
					"    n_samples=1000,\n",
					"    n_features=10,\n",
					"    n_informative=5,\n",
					"    n_redundant=2,\n",
					"    random_state=42\n",
					")\n",
					"\n",
					"# Perform classification\n",
					"report = perform_classification(X, y)\n",
					"\n",
					"# Print the classification report\n",
					"print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Clustering: Unsupervised learning technique that groups similar data points together based on their characteristics or proximity.\n"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"from sklearn.cluster import KMeans\n",
					"\n",
					"def perform_clustering(X, n_clusters):\n",
					"    \"\"\"\n",
					"    Perform clustering on the given input variables.\n",
					"\n",
					"    Args:\n",
					"        X (numpy.ndarray): The input variables with shape (n_samples, n_features).\n",
					"        n_clusters (int): The number of clusters to form.\n",
					"\n",
					"    Returns:\n",
					"        numpy.ndarray: The cluster labels for each data point.\n",
					"    \"\"\"\n",
					"    # Create and fit the KMeans clustering model\n",
					"    kmeans = KMeans(n_clusters=n_clusters, random_state=42)\n",
					"    cluster_labels = kmeans.fit_predict(X)\n",
					"\n",
					"    return cluster_labels\n",
					"\n",
					"import numpy as np\n",
					"\n",
					"# Example input variables\n",
					"X = np.array([[1, 2], [1.5, 1.8], [5, 8], [8, 8], [1, 0.6], [9, 11]])\n",
					"\n",
					"# Perform clustering\n",
					"cluster_labels = perform_clustering(X, n_clusters=2)\n",
					"\n",
					"# Print the cluster labels\n",
					"print(cluster_labels)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Principal Component Analysis (PCA): Dimensionality reduction technique that transforms a high-dimensional dataset into a lower-dimensional space while preserving its most important information."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"import numpy as np\n",
					"import pandas as pd\n",
					"from sklearn.decomposition import PCA\n",
					"\n",
					"def generate_pca_report(X, n_components):\n",
					"    # Perform PCA\n",
					"    pca = PCA(n_components=n_components, random_state=42)\n",
					"    transformed_data = pca.fit_transform(X)\n",
					"\n",
					"    # Calculate explained variance ratio\n",
					"    explained_variance_ratio = pca.explained_variance_ratio_\n",
					"\n",
					"    # Calculate cumulative explained variance ratio\n",
					"    cumulative_explained_variance_ratio = np.cumsum(explained_variance_ratio)\n",
					"\n",
					"    # Create a report DataFrame\n",
					"    report_data = {\n",
					"        'Principal Component': range(1, n_components + 1),\n",
					"        'Explained Variance Ratio': explained_variance_ratio,\n",
					"        'Cumulative Explained Variance Ratio': cumulative_explained_variance_ratio\n",
					"    }\n",
					"    report = pd.DataFrame(report_data)\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example input variables\n",
					"X = np.array([[2, 3, 4], [5, 6, 7], [8, 9, 10], [11, 12, 13], [14, 15, 16]])\n",
					"\n",
					"# Generate PCA report\n",
					"pca_report = generate_pca_report(X, n_components=3)\n",
					"\n",
					"# Print the PCA report\n",
					"print(pca_report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Factor Analysis: Statistical technique used to uncover latent factors or constructs that explain the correlations among observed variables."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"import pandas as pd\n",
					"from sklearn.decomposition import FactorAnalysis\n",
					"\n",
					"def perform_factor_analysis(data, n_factors):\n",
					"    # Perform factor analysis\n",
					"    fa = FactorAnalysis(n_components=n_factors, random_state=42)\n",
					"    transformed_data = fa.fit_transform(data)\n",
					"    \n",
					"    # Get factor loadings\n",
					"    factor_loadings = pd.DataFrame(\n",
					"        fa.components_.T,\n",
					"        index=data.columns,\n",
					"        columns=[f\"Factor {i+1}\" for i in range(n_factors)]\n",
					"    )\n",
					"    \n",
					"    return transformed_data, factor_loadings\n",
					"\n",
					"# Example input variables\n",
					"data = pd.DataFrame({\n",
					"    'Item 1': [1, 4, 7, 10],\n",
					"    'Item 2': [2, 5, 8, 11],\n",
					"    'Item 3': [3, 6, 9, 12]\n",
					"})\n",
					"\n",
					"# Perform factor analysis\n",
					"transformed_data, factor_loadings = perform_factor_analysis(data, n_factors=2)\n",
					"\n",
					"# Print the transformed data and factor loadings\n",
					"print(\"Transformed Data:\")\n",
					"print(transformed_data)\n",
					"print(\"\\nFactor Loadings:\")\n",
					"print(factor_loadings)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Survival Analysis: Statistical method for analyzing time-to-event data, such as time until death or failure, to estimate survival probabilities and hazard rates."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"def generate_survival_analysis_report(data):\n",
					"    %pip install lifelines\n",
					"\n",
					"    import pandas as pd\n",
					"    from lifelines import KaplanMeierFitter\n",
					"\n",
					"    # Preprocessing\n",
					"    df = pd.DataFrame(data)\n",
					"    df['event_time'] = pd.to_datetime(df['event_time'])\n",
					"    df['event'] = df['event'].astype(bool)\n",
					"\n",
					"    # Survival Analysis\n",
					"    kmf = KaplanMeierFitter()\n",
					"    kmf.fit(df['event_time'], event_observed=df['event'])\n",
					"\n",
					"    # Generate Report\n",
					"    report = f\"Survival Analysis Report\\n\\n\"\n",
					"    report += f\"Number of Observations: {len(df)}\\n\"\n",
					"    report += f\"Number of Events: {df['event'].sum()}\\n\\n\"\n",
					"    report += f\"Survival Function:\\n\\n{kmf.survival_function_}\\n\\n\"\n",
					"    report += f\"Median Survival Time: {kmf.median_survival_time_}\\n\\n\"\n",
					"\n",
					"    return report\n",
					"\n",
					"# Example usage\n",
					"data = {\n",
					"    'event_time': ['2021-01-01', '2021-02-01', '2021-03-01', '2021-04-01', '2021-05-01'],\n",
					"    'event': [False, False, True, True, False]\n",
					"}\n",
					"\n",
					"# report = generate_survival_analysis_report(data)\n",
					"# print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Bayesian Analysis: Statistical approach that combines prior knowledge or beliefs with observed data to estimate posterior probabilities and make inferences."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scipy\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from scipy import stats\n",
					"\n",
					"def bayesian_analysis(data, prior_beliefs):\n",
					"    observed_data = data.values\n",
					"\n",
					"    # Compute posterior probabilities\n",
					"    posterior_probs = prior_beliefs * stats.norm.pdf(observed_data, loc=data.mean(), scale=data.std())\n",
					"\n",
					"    # Normalize posterior probabilities\n",
					"    normalized_probs = posterior_probs / np.sum(posterior_probs)\n",
					"\n",
					"    # Calculate summary statistics\n",
					"    mean = np.sum(normalized_probs * observed_data)\n",
					"    variance = np.sum(normalized_probs * ((observed_data - mean) ** 2))\n",
					"    standard_deviation = np.sqrt(variance)\n",
					"\n",
					"    # Construct report\n",
					"    report = pd.DataFrame({'Summary Statistic': ['Mean', 'Variance', 'Standard Deviation'],\n",
					"                           'Value': [mean, variance, standard_deviation]})\n",
					"\n",
					"    return report\n",
					"\n",
					"# Example usage\n",
					"# data = pd.read_csv('data.csv', header=None, names=['Observation'])\n",
					"# prior_beliefs = np.array([0.3, 0.5, 0.2])  # Example prior beliefs (weights)\n",
					"\n",
					"# report = bayesian_analysis(data, prior_beliefs)\n",
					"# print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Decision Trees: Non-parametric predictive model that partitions the data into hierarchical structures to make decisions or predictions based on feature values."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.tree import DecisionTreeClassifier\n",
					"from sklearn.metrics import classification_report\n",
					"\n",
					"def decision_tree_report(data, target):\n",
					"    # Split the data into features (X) and target variable (y)\n",
					"    X = data.drop(target, axis=1)\n",
					"    y = data[target]\n",
					"    \n",
					"    # Initialize and fit a decision tree classifier\n",
					"    dt_classifier = DecisionTreeClassifier()\n",
					"    dt_classifier.fit(X, y)\n",
					"    \n",
					"    # Generate predictions on the training data\n",
					"    y_pred = dt_classifier.predict(X)\n",
					"    \n",
					"    # Generate the classification report\n",
					"    report = classification_report(y, y_pred)\n",
					"    \n",
					"    return report\n",
					"\n",
					"# Example usage\n",
					"# data = pd.read_csv('your_data.csv')  # Load your dataset\n",
					"# target = 'target_variable'          # Specify the target variable column name\n",
					"# report = decision_tree_report(data, target)\n",
					"# print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Random Forests: Ensemble learning technique that combines multiple decision trees to improve prediction accuracy and handle complex relationships."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.ensemble import RandomForestClassifier\n",
					"from sklearn.metrics import classification_report\n",
					"\n",
					"def random_forest_report(X_train, y_train, X_test, y_test):\n",
					"    # Fit a random forest classifier\n",
					"    clf = RandomForestClassifier()\n",
					"    clf.fit(X_train, y_train)\n",
					"\n",
					"    # Make predictions on the test set\n",
					"    y_pred = clf.predict(X_test)\n",
					"\n",
					"    # Generate a classification report\n",
					"    report = classification_report(y_test, y_pred)\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"from sklearn.datasets import load_iris\n",
					"iris = load_iris()\n",
					"X = iris.data\n",
					"y = iris.target\n",
					"\n",
					"# Split the data into train and test sets\n",
					"from sklearn.model_selection import train_test_split\n",
					"X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)\n",
					"\n",
					"# Generate the report\n",
					"report = random_forest_report(X_train, y_train, X_test, y_test)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Support Vector Machines (SVM): Supervised learning algorithm that constructs a hyperplane or set of hyperplanes to separate data into different classes."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.svm import SVR\n",
					"from sklearn.datasets import fetch_california_housing\n",
					"from sklearn.metrics import mean_squared_error, r2_score\n",
					"\n",
					"def svr_report(X_train, y_train, X_test, y_test):\n",
					"    svr = SVR()\n",
					"    svr.fit(X_train, y_train)\n",
					"\n",
					"    y_pred = svr.predict(X_test)\n",
					"\n",
					"    mse = mean_squared_error(y_test, y_pred)\n",
					"    r2 = r2_score(y_test, y_pred)\n",
					"\n",
					"    report = f\"Mean Squared Error: {mse:.4f}\\nR-squared: {r2:.4f}\"\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example usage\n",
					"california_housing = fetch_california_housing()\n",
					"X = california_housing.data\n",
					"y = california_housing.target\n",
					"\n",
					"X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)\n",
					"\n",
					"report = svr_report(X_train, y_train, X_test, y_test)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Neural Networks: Computational models inspired by the structure and function of the human brain, used for pattern recognition, classification, and regression tasks."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scikit-learn tensorflow\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.model_selection import train_test_split\n",
					"from sklearn.preprocessing import StandardScaler\n",
					"from sklearn.metrics import classification_report\n",
					"import tensorflow as tf\n",
					"# from tensorflow.keras.models import Sequential\n",
					"# from tensorflow.keras.layers import Dense\n",
					"\n",
					"def neural_network_report(X_train, y_train, X_test, y_test):\n",
					"    # Standardize the input features\n",
					"    scaler = StandardScaler()\n",
					"    X_train = scaler.fit_transform(X_train)\n",
					"    X_test = scaler.transform(X_test)\n",
					"\n",
					"    # Build the neural network model\n",
					"    model = Sequential()\n",
					"    model.add(Dense(64, activation='relu', input_shape=(X_train.shape[1],)))\n",
					"    model.add(Dense(64, activation='relu'))\n",
					"    model.add(Dense(1, activation='sigmoid'))\n",
					"\n",
					"    # Compile the model\n",
					"    model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])\n",
					"\n",
					"    # Train the model\n",
					"    model.fit(X_train, y_train, epochs=10, batch_size=32, verbose=0)\n",
					"\n",
					"    # Make predictions on the test set\n",
					"    y_pred = model.predict_classes(X_test)\n",
					"\n",
					"    # Generate a classification report\n",
					"    report = classification_report(y_test, y_pred)\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"from sklearn.datasets import load_breast_cancer\n",
					"breast_cancer = load_breast_cancer()\n",
					"X = breast_cancer.data\n",
					"y = breast_cancer.target\n",
					"\n",
					"# Split the data into train and test sets\n",
					"X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)\n",
					"\n",
					"# Generate the report\n",
					"# report = neural_network_report(X_train, y_train, X_test, y_test)\n",
					"\n",
					"# print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Natural Language Processing (NLP): Field of study that focuses on the interaction between computers and human language, enabling machines to understand, interpret, and generate human language."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy scikit-learn tensorflow\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.model_selection import train_test_split\n",
					"from sklearn.feature_extraction.text import CountVectorizer\n",
					"from sklearn.linear_model import LogisticRegression\n",
					"from sklearn.metrics import classification_report\n",
					"import tensorflow as tf\n",
					"from tensorflow.keras.preprocessing.text import Tokenizer\n",
					"from tensorflow.keras.preprocessing.sequence import pad_sequences\n",
					"from tensorflow.keras.models import Sequential\n",
					"from tensorflow.keras.layers import Embedding, LSTM, Dense\n",
					"\n",
					"def nlp_report(X_train, y_train, X_test, y_test):\n",
					"    # Convert text data to numerical features using CountVectorizer\n",
					"    vectorizer = CountVectorizer()\n",
					"    X_train_vectorized = vectorizer.fit_transform(X_train)\n",
					"    X_test_vectorized = vectorizer.transform(X_test)\n",
					"\n",
					"    # Train logistic regression model\n",
					"    lr = LogisticRegression()\n",
					"    lr.fit(X_train_vectorized, y_train)\n",
					"\n",
					"    # Make predictions on the test set using logistic regression\n",
					"    y_pred_lr = lr.predict(X_test_vectorized)\n",
					"\n",
					"    # Generate classification report for logistic regression\n",
					"    report_lr = classification_report(y_test, y_pred_lr)\n",
					"\n",
					"    # Tokenize and pad the text data for LSTM\n",
					"    tokenizer = Tokenizer()\n",
					"    tokenizer.fit_on_texts(X_train)\n",
					"    X_train_tokenized = tokenizer.texts_to_sequences(X_train)\n",
					"    X_test_tokenized = tokenizer.texts_to_sequences(X_test)\n",
					"    vocab_size = len(tokenizer.word_index) + 1\n",
					"    max_length = max(len(sequence) for sequence in X_train_tokenized)\n",
					"    X_train_padded = pad_sequences(X_train_tokenized, maxlen=max_length, padding='post')\n",
					"    X_test_padded = pad_sequences(X_test_tokenized, maxlen=max_length, padding='post')\n",
					"\n",
					"    # Build and train LSTM model\n",
					"    model = Sequential()\n",
					"    model.add(Embedding(vocab_size, 100, input_length=max_length))\n",
					"    model.add(LSTM(64))\n",
					"    model.add(Dense(1, activation='sigmoid'))\n",
					"    model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])\n",
					"    model.fit(X_train_padded, y_train, epochs=10, batch_size=32, verbose=0)\n",
					"\n",
					"    # Make predictions on the test set using LSTM\n",
					"    y_pred_lstm = model.predict_classes(X_test_padded)\n",
					"\n",
					"    # Generate classification report for LSTM\n",
					"    report_lstm = classification_report(y_test, y_pred_lstm)\n",
					"\n",
					"    return report_lr, report_lstm\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"from sklearn.datasets import fetch_20newsgroups\n",
					"newsgroups = fetch_20newsgroups(subset='all', shuffle=True, random_state=42)\n",
					"\n",
					"X = newsgroups.data\n",
					"y = newsgroups.target\n",
					"\n",
					"# Split the data into train and test sets\n",
					"X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)\n",
					"\n",
					"# Generate the reports\n",
					"report_lr, report_lstm = nlp_report(X_train, y_train, X_test, y_test)\n",
					"\n",
					"print(\"Logistic Regression Report:\")\n",
					"print(report_lr)\n",
					"print()\n",
					"print(\"LSTM Report:\")\n",
					"print(report_lstm)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Association Rule Mining: Unsupervised learning technique that discovers interesting relationships or associations among variables in large datasets."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# %pip install pandas numpy mlxtend\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from mlxtend.frequent_patterns import apriori\n",
					"from mlxtend.frequent_patterns import association_rules\n",
					"\n",
					"def association_rule_mining_report(df, min_support, min_threshold):\n",
					"    # Perform one-hot encoding\n",
					"    encoded_df = pd.get_dummies(df['Item'])\n",
					"\n",
					"    # Concatenate Transaction ID column\n",
					"    encoded_df['Transaction ID'] = df['Transaction ID']\n",
					"\n",
					"    # Group items by transaction ID\n",
					"    grouped_df = encoded_df.groupby('Transaction ID').sum()\n",
					"\n",
					"    # Find frequent itemsets using Apriori algorithm\n",
					"    frequent_itemsets = apriori(grouped_df, min_support=min_support, use_colnames=True)\n",
					"\n",
					"    # Generate association rules\n",
					"    rules = association_rules(frequent_itemsets, metric=\"confidence\", min_threshold=min_threshold)\n",
					"\n",
					"    return rules\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"df = pd.DataFrame({\n",
					"    'Transaction ID': [1, 1, 2, 2, 2, 3, 3, 4, 4, 4, 4],\n",
					"    'Item': ['A', 'B', 'A', 'B', 'C', 'A', 'C', 'A', 'B', 'C', 'D']\n",
					"})\n",
					"\n",
					"# Generate the report\n",
					"report = association_rule_mining_report(df, min_support=0.3, min_threshold=0.7)\n",
					"\n",
					"print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Recommender Systems: Algorithms that provide personalized recommendations by predicting user preferences based on historical data and patterns."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"%pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.metrics.pairwise import cosine_similarity\n",
					"\n",
					"def recommender_system_report(ratings_df, user_id):\n",
					"    # Calculate item-item similarity matrix\n",
					"    item_similarity = cosine_similarity(ratings_df.T)\n",
					"\n",
					"    # Get user ratings and item similarities\n",
					"    user_ratings = ratings_df.loc[user_id].values.reshape(1, -1)\n",
					"    item_similarities = item_similarity * (1 - np.isnan(ratings_df.values))\n",
					"\n",
					"    # Calculate weighted average of item ratings\n",
					"    weighted_sum = np.dot(user_ratings, item_similarities)\n",
					"    weighted_avg = weighted_sum / np.sum(item_similarities, axis=1)\n",
					"\n",
					"    # Get recommended item indices sorted by predicted ratings\n",
					"    recommended_items = np.argsort(-weighted_avg)[0]\n",
					"\n",
					"    return recommended_items\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"ratings_data = {\n",
					"    'Item1': [4, 5, 2, 1],\n",
					"    'Item2': [3, 2, 5, 4],\n",
					"    'Item3': [1, 3, 4, 5],\n",
					"    'Item4': [5, 4, 3, 2]\n",
					"}\n",
					"\n",
					"ratings_df = pd.DataFrame(ratings_data, index=['User1', 'User2', 'User3', 'User4'])\n",
					"\n",
					"# Specify the user for recommendation\n",
					"user_id = 'User1'\n",
					"\n",
					"# Generate the report\n",
					"report = recommender_system_report(ratings_df, user_id)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Text Mining: Process of extracting useful information and knowledge from unstructured text data through techniques such as text classification, sentiment analysis, and topic modeling."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"%pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.feature_extraction.text import CountVectorizer\n",
					"from sklearn.decomposition import LatentDirichletAllocation\n",
					"\n",
					"def text_mining_report(documents, num_topics):\n",
					"    # Convert text documents into a matrix of token counts\n",
					"    vectorizer = CountVectorizer()\n",
					"    X = vectorizer.fit_transform(documents)\n",
					"\n",
					"    # Perform Latent Dirichlet Allocation (LDA) to extract topics\n",
					"    lda = LatentDirichletAllocation(n_components=num_topics, random_state=42)\n",
					"    lda.fit(X)\n",
					"\n",
					"    # Get the most important words for each topic\n",
					"    feature_names = vectorizer.get_feature_names_out()\n",
					"    top_words = []\n",
					"    for topic_idx, topic in enumerate(lda.components_):\n",
					"        top_words.append([feature_names[i] for i in topic.argsort()[:-6:-1]])\n",
					"\n",
					"    return top_words\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"documents = [\n",
					"    \"The cat sat on the mat.\",\n",
					"    \"The dog played in the garden.\",\n",
					"    \"The bird sang a beautiful song.\",\n",
					"    \"The cat and the dog chased each other.\"\n",
					"]\n",
					"\n",
					"# Specify the number of topics to extract\n",
					"num_topics = 2\n",
					"\n",
					"# Generate the report\n",
					"report = text_mining_report(documents, num_topics)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Anomaly Detection: Identifying rare or abnormal patterns or outliers in data that deviate significantly from the expected behavior."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"%pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.ensemble import IsolationForest\n",
					"\n",
					"def anomaly_detection_report(data, contamination):\n",
					"    # Reshape the data to 1-dimensional\n",
					"    data = np.array(data).reshape(-1, 1)\n",
					"\n",
					"    # Fit the Isolation Forest model\n",
					"    model = IsolationForest(contamination=contamination, random_state=42)\n",
					"    model.fit(data)\n",
					"\n",
					"    # Predict anomalies\n",
					"    predictions = model.predict(data)\n",
					"\n",
					"    # Convert predictions to boolean values\n",
					"    is_anomaly = np.where(predictions == -1, True, False)\n",
					"\n",
					"    # Create a report with anomaly labels and scores\n",
					"    report = pd.DataFrame({'Data': data.flatten(), 'Is Anomaly': is_anomaly, 'Anomaly Score': model.decision_function(data)})\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"data = [2.5, 1.2, 3.6, 4.2, 1.0, 1.9, 6.5, 4.8, 5.6]\n",
					"\n",
					"# Specify the contamination rate\n",
					"contamination = 0.2\n",
					"\n",
					"# Generate the report\n",
					"report = anomaly_detection_report(data, contamination)\n",
					"\n",
					"print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Ensemble Methods: Combining multiple models or predictions to improve overall performance and robustness."
				]
			},
			{
				"cell_type": "code",
				"execution_count": 52,
				"metadata": {},
				"outputs": [
					{
						"name": "stdout",
						"output_type": "stream",
						"text": [
							"Defaulting to user installation because normal site-packages is not writeable\n",
							"Requirement already satisfied: pandas in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (1.5.3)\n",
							"Requirement already satisfied: numpy in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (1.23.5)\n",
							"Requirement already satisfied: scikit-learn in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (1.2.2)\n",
							"Requirement already satisfied: python-dateutil>=2.8.1 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from pandas) (2.8.2)\n",
							"Requirement already satisfied: pytz>=2020.1 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from pandas) (2022.1)\n",
							"Requirement already satisfied: scipy>=1.3.2 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from scikit-learn) (1.10.1)\n",
							"Requirement already satisfied: joblib>=1.1.1 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from scikit-learn) (1.2.0)\n",
							"Requirement already satisfied: threadpoolctl>=2.0.0 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from scikit-learn) (3.1.0)\n",
							"Requirement already satisfied: six>=1.5 in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (from python-dateutil>=2.8.1->pandas) (1.16.0)\n",
							"Note: you may need to restart the kernel to use updated packages.\n",
							"              precision    recall  f1-score   support\n",
							"\n",
							"           0       1.00      1.00      1.00        10\n",
							"           1       1.00      1.00      1.00         9\n",
							"           2       1.00      1.00      1.00        11\n",
							"\n",
							"    accuracy                           1.00        30\n",
							"   macro avg       1.00      1.00      1.00        30\n",
							"weighted avg       1.00      1.00      1.00        30\n",
							"\n"
						]
					}
				],
				"source": [
					"%pip install pandas numpy scikit-learn\n",
					"\n",
					"import pandas as pd\n",
					"import numpy as np\n",
					"from sklearn.ensemble import RandomForestClassifier\n",
					"from sklearn.metrics import classification_report\n",
					"from sklearn.model_selection import train_test_split\n",
					"\n",
					"def ensemble_method_report(X, y, test_size=0.2, random_state=42):\n",
					"    # Split the data into training and testing sets\n",
					"    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=test_size, random_state=random_state)\n",
					"\n",
					"    # Initialize a Random Forest classifier\n",
					"    classifier = RandomForestClassifier(random_state=random_state)\n",
					"\n",
					"    # Fit the classifier on the training data\n",
					"    classifier.fit(X_train, y_train)\n",
					"\n",
					"    # Make predictions on the testing data\n",
					"    y_pred = classifier.predict(X_test)\n",
					"\n",
					"    # Generate the classification report\n",
					"    report = classification_report(y_test, y_pred)\n",
					"\n",
					"    return report\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Import example data\n",
					"from sklearn.datasets import load_iris\n",
					"\n",
					"# Load the Iris dataset\n",
					"data = load_iris()\n",
					"X = data.data\n",
					"y = data.target\n",
					"\n",
					"# Generate the report\n",
					"report = ensemble_method_report(X, y)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Genetic Algorithms: Optimization algorithms inspired by the process of natural selection, used to find near-optimal solutions to complex problems."
				]
			},
			{
				"cell_type": "code",
				"execution_count": 53,
				"metadata": {},
				"outputs": [
					{
						"name": "stdout",
						"output_type": "stream",
						"text": [
							"Defaulting to user installation because normal site-packages is not writeable\n",
							"Requirement already satisfied: numpy in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (1.23.5)\n",
							"Note: you may need to restart the kernel to use updated packages.\n",
							"HelNol Worldl\n"
						]
					}
				],
				"source": [
					"%pip install numpy\n",
					"\n",
					"import numpy as np\n",
					"from numpy.random import randint\n",
					"\n",
					"def genetic_algorithm_report(population_size, num_generations, target_string):\n",
					"    # Define the target string as an array of characters\n",
					"    target_array = np.array(list(target_string))\n",
					"\n",
					"    # Define the valid characters for the population\n",
					"    valid_chars = np.array(list(\"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ \"))\n",
					"\n",
					"    # Initialize the population with random strings\n",
					"    population = []\n",
					"    for _ in range(population_size):\n",
					"        random_string = ''.join(valid_chars[randint(0, len(valid_chars), len(target_string))])\n",
					"        population.append(random_string)\n",
					"\n",
					"    # Perform the genetic algorithm\n",
					"    for generation in range(num_generations):\n",
					"        fitness_scores = []\n",
					"        for individual in population:\n",
					"            # Calculate the fitness score as the number of matching characters\n",
					"            fitness_score = np.sum(np.array(list(individual)) == target_array)\n",
					"            fitness_scores.append(fitness_score)\n",
					"\n",
					"        # Select the fittest individuals to reproduce\n",
					"        fittest_indices = np.argsort(fitness_scores)[-population_size // 2:]\n",
					"        fittest_population = [population[i] for i in fittest_indices]\n",
					"\n",
					"        # Create the next generation through crossover and mutation\n",
					"        next_generation = []\n",
					"        while len(next_generation) < population_size:\n",
					"            parent1, parent2 = np.random.choice(fittest_population, size=2, replace=False)\n",
					"            crossover_point = np.random.randint(1, len(target_string))\n",
					"            child = parent1[:crossover_point] + parent2[crossover_point:]\n",
					"            mutation_point = np.random.randint(len(target_string))\n",
					"            child = child[:mutation_point] + valid_chars[randint(0, len(valid_chars))] + child[mutation_point+1:]\n",
					"            next_generation.append(child)\n",
					"\n",
					"        population = next_generation\n",
					"\n",
					"    # Find the best individual from the final population\n",
					"    best_individual = max(population, key=lambda x: np.sum(np.array(list(x)) == target_array))\n",
					"\n",
					"    return best_individual\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Specify the parameters\n",
					"population_size = 100\n",
					"num_generations = 200\n",
					"target_string = \"Hello, World!\"\n",
					"\n",
					"# Generate the report\n",
					"report = genetic_algorithm_report(population_size, num_generations, target_string)\n",
					"\n",
					"print(report)"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Markov Chains: Mathematical models that describe a sequence of events or states, where the probability of transitioning to the next state depends only on the current state."
				]
			},
			{
				"cell_type": "code",
				"execution_count": 55,
				"metadata": {},
				"outputs": [
					{
						"name": "stdout",
						"output_type": "stream",
						"text": [
							"Defaulting to user installation because normal site-packages is not writeable\n",
							"Requirement already satisfied: numpy in c:\\users\\yujie\\appdata\\roaming\\python\\python311\\site-packages (1.23.5)\n",
							"Note: you may need to restart the kernel to use updated packages.\n",
							"[0, 0, 0, 0, 1, 1, 2, 2, 2, 2, 2]\n"
						]
					}
				],
				"source": [
					"%pip install numpy\n",
					"\n",
					"import numpy as np\n",
					"\n",
					"def markov_chain_report(transition_matrix, initial_state, num_steps):\n",
					"    # Convert the transition matrix and initial state to numpy arrays\n",
					"    transition_matrix = np.array(transition_matrix)\n",
					"    initial_state = np.array(initial_state)\n",
					"\n",
					"    # Calculate the number of states\n",
					"    num_states = transition_matrix.shape[0]\n",
					"\n",
					"    # Perform the Markov chain simulation\n",
					"    current_state = np.random.choice(np.arange(num_states), p=initial_state)\n",
					"    states = [current_state]\n",
					"    for _ in range(num_steps):\n",
					"        # Generate a random number to determine the next state\n",
					"        random_number = np.random.rand()\n",
					"\n",
					"        # Find the next state based on the transition probabilities\n",
					"        next_state = np.random.choice(np.arange(num_states), p=transition_matrix[current_state])\n",
					"\n",
					"        # Update the current state\n",
					"        current_state = next_state\n",
					"\n",
					"        # Append the current state to the list of states\n",
					"        states.append(current_state)\n",
					"\n",
					"    return states\n",
					"\n",
					"\n",
					"# Example usage\n",
					"# Define the transition matrix and initial state\n",
					"transition_matrix = [\n",
					"    [0.7, 0.3, 0.0],\n",
					"    [0.0, 0.4, 0.6],\n",
					"    [0.5, 0.0, 0.5]\n",
					"]\n",
					"initial_state = [1.0, 0.0, 0.0]\n",
					"\n",
					"# Specify the number of steps for the Markov chain\n",
					"num_steps = 10\n",
					"\n",
					"# Generate the report\n",
					"report = markov_chain_report(transition_matrix, initial_state, num_steps)\n",
					"\n",
					"print(report)\n"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Hidden Markov Models (HMM): Statistical models used to model sequential data, where the underlying states are not directly observed but can be inferred.\n",
					"\n",
					"Find an actual working example!"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"## OTHER EXAMPLES"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"- Dimensionality Reduction: Techniques that reduce the number of variables or features while preserving important information and reducing noise.\n",
					"\n",
					"- Collaborative Filtering: Recommendation technique that predicts user preferences based on the preferences of similar users or items.\n",
					"\n",
					"- Network Analysis: Study of networks or graphs to understand and analyze relationships, connectivity, and patterns within complex systems.\n",
					"\n",
					"- Graph Mining: Analyzing and extracting useful information from large-scale graph structures, such as social networks or biological networks\n",
					"\n",
					"- Support Vector Regression (SVR): Regression technique that uses support vector machines to model and predict continuous variables.\n",
					"\n",
					"- Gradient Boosting Machines (GBM): Ensemble learning method that combines multiple weak prediction models (typically decision trees) to create a strong predictive model.\n",
					"\n",
					"- XGBoost: Gradient boosting library that provides optimized implementations of gradient boosting algorithms and is known for its efficiency and performance.\n",
					"\n",
					"- LightGBM: Gradient boosting framework that uses tree-based learning algorithms and is designed to be efficient with large-scale datasets.\n",
					"\n",
					"- CatBoost: Gradient boosting library that handles categorical features effectively and provides built-in handling of missing values.\n",
					"\n",
					"- K-Nearest Neighbors (KNN): Non-parametric classification algorithm that assigns labels to data points based on the majority vote of their nearest neighbors in the feature space.\n",
					"\n",
					"- Naive Bayes: Probabilistic classifier that applies Bayes' theorem with the assumption of independence among features.\n",
					"\n",
					"- Logistic Regression: Statistical regression model that predicts the probability of binary or categorical outcomes using a logistic function.\n",
					"\n",
					"- Poisson Regression: Regression technique used to model count data with an assumed Poisson distribution.\n",
					"\n",
					"- Lasso and Ridge Regression: Regularization techniques that introduce a penalty term to control the complexity of a regression model and prevent overfitting.\n",
					"\n",
					"- Elastic Net: Regression method that combines the penalties of Lasso and Ridge regression to handle high-dimensional datasets with correlated variables.\n",
					"\n",
					"- Quantile Regression: Regression technique that estimates conditional quantiles of the response variable, providing a more complete understanding of the relationship between variables.\n",
					"\n",
					"- K-Means Clustering: Partitioning method that aims to divide a dataset into K clusters based on the similarity of data points to the cluster centroids.\n",
					"\n",
					"- DBSCAN: Density-based clustering algorithm that groups data points into clusters based on their density and proximity.\n",
					"\n",
					"- Hierarchical Clustering: Agglomerative or divisive clustering method that builds a hierarchy of clusters by iteratively merging or splitting them based on proximity.\n",
					"\n",
					"- Gaussian Mixture Models (GMM): Probabilistic model that represents a dataset as a mixture of Gaussian distributions, often used for clustering or density estimation.\n",
					"\n",
					"- Hidden Markov Models for Time Series: Statistical models used to model and predict time series data, where the underlying states are not directly observable.\n",
					"\n",
					"- Long Short-Term Memory (LSTM): Recurrent neural network architecture that is capable of capturing long-term dependencies and has been widely used in sequence modeling tasks.\n",
					"\n",
					"- Convolutional Neural Networks (CNN): Neural network architecture designed to process structured grid-like data, such as images, using convolutional and pooling layers.\n",
					"\n",
					"- Recurrent Neural Networks (RNN): Neural network architecture that can handle sequential and time-dependent data by using feedback connections.\n",
					"\n",
					"- Transformers: Neural network architecture that utilizes self-attention mechanisms to capture relationships between different positions in the input sequence, often used in natural language processing and sequence-to-sequence tasks.\n",
					"\n",
					"- Word2Vec: Technique for learning word embeddings, representing words as dense vectors in a continuous space, often used in natural language processing tasks.\n",
					"\n",
					"- Latent Dirichlet Allocation (LDA): Generative probabilistic model used for topic modeling to uncover latent topics in a collection of documents.\n",
					"\n",
					"- Latent Semantic Analysis (LSA): Technique that analyzes relationships between documents and terms to uncover hidden semantic structures in a text corpus.\n",
					"\n",
					"- Singular Value Decomposition (SVD): Matrix factorization method that decomposes a matrix into three matrices to reveal its latent structure and reduce dimensionality.\n",
					"\n",
					"- Collaborative Filtering: Recommendation technique that predicts user preferences or item ratings based on the preferences of similar users or items.\n",
					"\n",
					"- Markov Chain Monte Carlo (MCMC): Method for sampling from complex probability distributions, often used in Bayesian inference to estimate posterior distributions of parameters.\n",
					"\n",
					"- Particle Swarm Optimization (PSO): Optimization algorithm inspired by the social behavior of bird flocking or fish schooling, used to find the optimal solution in a search space.\n",
					"\n",
					"- Reinforcement Learning: Branch of machine learning concerned with learning how to make decisions or take actions in an environment to maximize a reward signal.\n",
					"\n",
					"- Q-Learning: Model-free reinforcement learning algorithm that learns an optimal policy for an agent in a Markov decision process.\n",
					"\n",
					"- Deep Q-Networks (DQN): Deep reinforcement learning algorithm that combines deep neural networks with Q-learning to approximate the optimal action-value function.\n",
					"\n",
					"- Variational Autoencoders (VAE): Generative models that learn a latent representation of the input data and can generate new samples from the learned distribution.\n",
					"\n",
					"- Generative Adversarial Networks (GAN): Framework that consists of a generator and a discriminator network that are trained in an adversarial manner to generate realistic samples.\n",
					"\n",
					"- t-SNE (t-Distributed Stochastic Neighbor Embedding): Dimensionality reduction technique that maps high-dimensional data to a lower-dimensional space while preserving local structure.\n",
					"\n",
					"- UMAP (Uniform Manifold Approximation and Projection): Dimensionality reduction technique that preserves both local and global structure in the data and is known for its scalability.\n",
					"\n",
					"- Recurrent Neural Network (RNN): Neural network architecture designed to handle sequential and time-dependent data by using feedback connections between hidden units.\n",
					"\n",
					"- Gated Recurrent Unit (GRU): Variation of recurrent neural networks that uses gating mechanisms to better capture long-term dependencies and alleviate the vanishing gradient problem.\n",
					"\n",
					"- Transformer Networks: Neural network architecture that utilizes self-attention mechanisms to capture relationships between different positions in the input sequence, commonly used in natural language processing and machine translation tasks.\n",
					"\n",
					"- Deep Reinforcement Learning (DRL): Combining deep neural networks with reinforcement learning to train agents that can learn complex behaviors and make decisions in dynamic environments.\n",
					"\n",
					"- Self-Organizing Maps (SOM): Unsupervised learning technique that creates a low-dimensional representation of the input data, preserving the topological relationships between data points.\n",
					"\n",
					"- Non-negative Matrix Factorization (NMF): Matrix factorization technique that decomposes a non-negative matrix into two non-negative matrices, often used for dimensionality reduction or feature extraction.\n",
					"\n",
					"- Ordinal Regression: Regression technique used when the dependent variable has ordered categories or levels, providing predictions in the form of ordinal values.\n",
					"\n",
					"- Survival Regression: Regression technique used when the dependent variable represents the time until an event occurs, such as time until failure or time until a customer churns.\n",
					"\n",
					"- Hidden Semi-Markov Models (HSMM): Extension of Hidden Markov Models that allows for variable duration of states, often used in modeling sequential data with variable time intervals.\n",
					"\n",
					"- Imbalanced Data Techniques: Methods to handle imbalanced datasets, where the classes are not equally represented, including techniques like SMOTE (Synthetic Minority Over-sampling Technique) or ADASYN (Adaptive Synthetic Sampling).\n",
					"\n",
					"- Causal Inference Methods: Statistical techniques to determine causal relationships between variables and infer the effect of interventions or treatments on outcomes.\n",
					"\n",
					"- Synthetic Data Generation: Creating artificial data that mimics the characteristics of real data, often used for privacy protection, data augmentation, or simulating rare events.\n",
					"\n",
					"- Network Embedding: Mapping nodes in a network into low-dimensional vector representations, enabling various network analysis tasks such as link prediction or community detection.\n",
					"\n",
					"- Stacked Generalization (Stacking): Ensemble learning technique that trains multiple models and combines their predictions using another model to improve overall performance.\n",
					"\n",
					"- Gradient Boosting Decision Trees (GBDT): Ensemble learning method that combines multiple decision trees, trained in a stage-wise manner, to make accurate predictions.\n",
					"\n",
					"- Rule-based Models: Models that use a set of predefined rules or conditions to make predictions or decisions based on specific criteria.\n",
					"\n",
					"- Fuzzy Logic Systems: Mathematical framework that handles uncertainty and imprecision by assigning degrees of membership to variables, often used in decision-making systems.\n",
					"\n",
					"- Extreme Value Theory (EVT): Statistical theory that models the extreme values of a distribution, often used in risk management and predicting rare events.\n",
					"\n",
					"- Zero-Inflated Models: Statistical models used to analyze data with excessive zero values, such as count data with excess zeros or excessive non-response in surveys.\n",
					"\n",
					"- Dynamic Time Warping (DTW): Distance measure used to compare and align time series data that may vary in time or speed, often used in pattern recognition or speech recognition.\n",
					"\n",
					"- Transfer Learning: Technique that leverages knowledge learned from one task or domain to improve learning or performance on a different but related task or domain.\n",
					"\n",
					"- Active Learning: Process where an algorithm interacts with a human or an oracle to strategically select the most informative samples for labeling, reducing the labeling effort.\n",
					"\n",
					"- Autoencoders: Neural network architectures used for unsupervised learning and dimensionality reduction by learning to reconstruct the input data from a compressed representation.\n",
					"\n",
					"- Hyperparameter Optimization: Techniques to find the optimal hyperparameter values of a model or algorithm, often done through methods like grid search, random search, or Bayesian optimization.\n",
					"\n",
					"- Reinforcement Learning with Function Approximation: Combining reinforcement learning with function approximation methods, such as neural networks, to handle high-dimensional state spaces.\n",
					"\n",
					"- Multi-Task Learning: Learning paradigm where a model is trained to perform multiple related tasks simultaneously, leveraging shared information and improving generalization.\n",
					"\n",
					"- Markov Decision Processes (MDPs): Mathematical framework used to model decision-making processes under uncertainty, comprising states, actions, rewards, and transition probabilities.\n",
					"\n",
					"- Gaussian Processes: Probabilistic models that define distributions over functions, often used in regression problems and surrogate modeling.\n",
					"\n",
					"- Semi-Supervised Learning: Learning paradigm that combines labeled and unlabeled data to improve model performance, especially when labeled data is scarce or expensive to obtain.\n",
					"\n",
					"- Longitudinal Data Analysis: Statistical techniques for analyzing data collected over multiple time points from the same individuals, accounting for dependencies and temporal patterns.\n",
					"\n",
					"- Causal Graphical Models: Models that represent causal relationships among variables using directed acyclic graphs, facilitating causal inference and identification of causal mechanisms.\n",
					"\n",
					"- Bayesian Networks: Probabilistic graphical models that represent the probabilistic relationships among variables through directed acyclic graphs, enabling reasoning under uncertainty.\n",
					"\n",
					"- Deep Reinforcement Learning: Combining deep neural networks with reinforcement learning to train agents that can learn complex behaviors and make decisions in dynamic environments.\n",
					"\n",
					"- Federated Learning: Distributed machine learning approach where models are trained collaboratively across multiple devices or parties while keeping data decentralized and private.\n",
					"\n",
					"- Subspace Learning: Techniques that aim to learn a low-dimensional subspace that captures the most relevant information in high-dimensional data.\n",
					"\n",
					"- Ensemble Learning with Stacking: Combining predictions from multiple models by training a meta-model that learns to combine the outputs of the base models, often improving overall performance and generalization."
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"\n",
					"5. Results and Conclusion\n",
					"    - Summarize the findings from the analysis.\n",
					"    - Present the results using visualizations and tables.\n",
					"    - Discuss any limitations or assumptions of the study.\n",
					"    - Draw conclusions based on the results and their implications.\n",
					"    - Provide recommendations for future research or actions."
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"# Results and conclusion"
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"6. Conclusion\n",
					"    - Recap the main points of the study.\n",
					"    - Encourage further exploration and learning in the field of statistics."
				]
			},
			{
				"attachments": {},
				"cell_type": "markdown",
				"metadata": {},
				"source": [
					"\n",
					"6. References\n",
					"    - List any references or sources used in the study.\n",
					"\n"
				]
			}
		],
		"metadata": {
			"kernelspec": {
				"display_name": "Python 3",
				"language": "python",
				"name": "python3"
			},
			"language_info": {
				"codemirror_mode": {
					"name": "ipython",
					"version": 3
				},
				"file_extension": ".py",
				"mimetype": "text/x-python",
				"name": "python",
				"nbconvert_exporter": "python",
				"pygments_lexer": "ipython3",
				"version": "3.11.2"
			},
			"orig_nbformat": 4
		},
		"nbformat": 4,
		"nbformat_minor": 2
	}   
`
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("File created successfully.")
	return nil
}
