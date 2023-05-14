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

	exampleReadmeFilePath := filepath.Join(exampleDir, "eg", "README.md")
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
}

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

    \definecolor{codegreen}{rgb}{0,0.6,0}
    \definecolor{codegray}{rgb}{0.5,0.5,0.5}
    \definecolor{codepurple}{rgb}{0.58,0,0.82}
    \definecolor{backcolour}{rgb}{0.97,0.97,0.95}

    \lstdefinestyle{mystyle}{
        backgroundcolor=\color{backcolour},   
        commentstyle=\color{codegreen},
        keywordstyle=\color{blue},
        numberstyle=\tiny\color{codegray},
        stringstyle=\color{codepurple},
        basicstyle=\ttfamily\footnotesize,
        breakatwhitespace=false,         
        breaklines=true,                 
        captionpos=b,                    
        keepspaces=true,                 
        numbers=left,                    
        numbersep=5pt,                  
        showspaces=false,                
        showstringspaces=false,
        showtabs=false,                  
        tabsize=4
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

    \definecolor{codegreen}{rgb}{0,0.6,0}
    \definecolor{codegray}{rgb}{0.5,0.5,0.5}
    \definecolor{codepurple}{rgb}{0.58,0,0.82}
    \definecolor{backcolour}{rgb}{0.97,0.97,0.95}

    \lstdefinestyle{mystyle}{
        backgroundcolor=\color{backcolour},   
        commentstyle=\color{codegreen},
        keywordstyle=\color{blue},
        numberstyle=\tiny\color{codegray},
        stringstyle=\color{codepurple},
        basicstyle=\ttfamily\footnotesize,
        breakatwhitespace=false,         
        breaklines=true,                 
        captionpos=b,                    
        keepspaces=true,                 
        numbers=left,                    
        numbersep=5pt,                  
        showspaces=false,                
        showstringspaces=false,
        showtabs=false,                  
        tabsize=4
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

\begin{frame}
    \setcounter{footnote}{0}
    \frametitle{Contents}
    \tableofcontents
\end{frame}


\begin{frame}[plain]
    \frametitle{Title}
    \setcounter{footnote}{0}
    \setcounter{equation}{0}
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
