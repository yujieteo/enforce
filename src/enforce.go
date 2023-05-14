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

	err = createGitignore(projectPath)
	if err != nil {
		fmt.Println(err)
	}

	err = createJupyterTemplate(projectPath)
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

\begin{frame}
  \frametitle{Literature Review}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Overview of Relevant Literature: Summarizes the key literature and theories related to the research topic.
    \item Key Concepts and Definitions: Defines important terms and concepts used in the research.
    \item Previous Research Findings: Highlights the main findings of previous studies related to the research question.
    \item Gaps in Existing Knowledge: Identifies areas where further research is needed or where the current knowledge is limited.
  \end{itemize}
\end{frame}

\begin{frame}
  \frametitle{Methodology}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Research Design: Describes the overall approach and methodology employed in the research.
    \item Data Collection: Explains how data was gathered or collected for the study.
    \item Data Analysis: Describes the methods used to analyze the collected data.
    \item Variables and Measures: Specifies the variables studied and the measures used to assess them.
  \end{itemize}
\end{frame}

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

\begin{frame}
  \frametitle{References}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Lists the references cited in the research paper.
  \end{itemize}
\end{frame}

\begin{frame}
  \frametitle{Appendix (if applicable)}
  \setcounter{footnote}{0}
  \setcounter{equation}{0}
  \begin{itemize}
    \item Includes any supplementary materials or additional information that supports the research.
  \end{itemize}
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
		  "execution_count": 11,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "import numpy as np                      # Numerical computing library\n",
		   "import pandas as pd                     # Data manipulation and analysis library\n",
		   "import matplotlib.pyplot as plt        # Data visualization library\n",
		   "import seaborn as sns                   # Enhanced data visualization library\n",
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
		   "import datetime                       # Date and time manipulation\n",
		   "import os                             # Operating system interaction"
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
		  "execution_count": 12,
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
		  "execution_count": 13,
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
		   "- Hypothesis Testing: Statistical inference method to evaluate a claim about a population based on sample data.\n",
		   "\n",
		   "- Regression Analysis: Statistical modeling technique to investigate the relationship between a dependent variable and one or more independent variables.\n",
		   "\n",
		   "- Time Series Analysis: Analyzing and modeling data points collected over time to uncover patterns, trends, and make predictions.\n",
		   "\n",
		   "- Classification: Predictive modeling technique that assigns input data points to predefined classes or categories based on their features.\n",
		   "\n",
		   "- Clustering: Unsupervised learning technique that groups similar data points together based on their characteristics or proximity.\n",
		   "\n",
		   "- Principal Component Analysis (PCA): Dimensionality reduction technique that transforms a high-dimensional dataset into a lower-dimensional space while preserving its most important information.\n",
		   "\n",
		   "- Factor Analysis: Statistical technique used to uncover latent factors or constructs that explain the correlations among observed variables.\n",
		   "\n",
		   "- Survival Analysis: Statistical method for analyzing time-to-event data, such as time until death or failure, to estimate survival probabilities and hazard rates.\n",
		   "\n",
		   "- Bayesian Analysis: Statistical approach that combines prior knowledge or beliefs with observed data to estimate posterior probabilities and make inferences.\n",
		   "\n",
		   "- Decision Trees: Non-parametric predictive model that partitions the data into hierarchical structures to make decisions or predictions based on feature values.\n",
		   "\n",
		   "- Random Forests: Ensemble learning technique that combines multiple decision trees to improve prediction accuracy and handle complex relationships.\n",
		   "\n",
		   "- Support Vector Machines (SVM): Supervised learning algorithm that constructs a hyperplane or set of hyperplanes to separate data into different classes.\n",
		   "\n",
		   "- Neural Networks: Computational models inspired by the structure and function of the human brain, used for pattern recognition, classification, and regression tasks.\n",
		   "\n",
		   "- Natural Language Processing (NLP): Field of study that focuses on the interaction between computers and human language, enabling machines to understand, interpret, and generate human language.\n",
		   "\n",
		   "- Deep Learning: Subset of machine learning that uses neural networks with multiple layers to learn hierarchical representations of data and solve complex tasks.\n",
		   "\n",
		   "- Association Rule Mining: Unsupervised learning technique that discovers interesting relationships or associations among variables in large datasets.\n",
		   "\n",
		   "- Recommender Systems: Algorithms that provide personalized recommendations by predicting user preferences based on historical data and patterns.\n",
		   "\n",
		   "- Text Mining: Process of extracting useful information and knowledge from unstructured text data through techniques such as text classification, sentiment analysis, and topic modeling.\n",
		   "\n",
		   "- Anomaly Detection: Identifying rare or abnormal patterns or outliers in data that deviate significantly from the expected behavior.\n",
		   "\n",
		   "- Ensemble Methods: Combining multiple models or predictions to improve overall performance and robustness.\n",
		   "\n",
		   "- Genetic Algorithms: Optimization algorithms inspired by the process of natural selection, used to find near-optimal solutions to complex problems.\n",
		   "\n",
		   "- Markov Chains: Mathematical models that describe a sequence of events or states, where the probability of transitioning to the next state depends only on the current state.\n",
		   "\n",
		   "- Hidden Markov Models (HMM): Statistical models used to model sequential data, where the underlying states are not directly observed but can be inferred.\n",
		   "\n",
		   "- Reinforcement Learning: Branch of machine learning concerned with learning how to make decisions or take actions in an environment to maximize a reward signal.\n",
		   "\n",
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
		  "execution_count": 14,
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
