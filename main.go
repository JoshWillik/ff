package main

import (
  "fmt"
  "os"
  "syscall"
  "path/filepath"
  "github.com/manifoldco/promptui"
  "github.com/renstrom/fuzzysearch/fuzzy"
)

var debug bool

func files(dir string) []string {
  files := make([]string, 0, 30)
  filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
    files = append(files, path)
    return nil
  })
  return files
}

func mustGetwd() string {
  dir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  return dir
}

func fileMatches(pattern string) []string {
  dir := mustGetwd()+"/"
  if debug {
    fmt.Fprintln(os.Stderr, "searching "+pattern+" in "+dir)
  }
  all_files := files(dir)
  files := make([]string, 0, 20)
  for _, path := range all_files {
    rel_path := path[len(dir):]
    if fuzzy.Match(pattern, rel_path) {
      files = append(files, rel_path)
    }
  }
  return files
}

func chooseFile(files []string) string {
  prompt := promptui.Select{
    Label: "Choose file",
    Items: files,
    Size: 10,
    Searcher: func(input string, index int) bool {
      return fuzzy.Match(input, files[index])
    },
    StartInSearchMode: true,
  }
  _, file, err := prompt.Run()
  if err != nil {
    panic(err)
  }
  return file
}

func openFile(path string) {
  editor := os.Getenv("EDITOR")
  if editor == "" {
    fmt.Fprintf(os.Stderr, "$EDITOR not configured")
    fmt.Println(path)
  }
  args := []string{editor, path}
  env := os.Environ()
  if debug {
    fmt.Fprintf(os.Stderr, "running %s %v\n", editor, args)
  }
  if err := syscall.Exec(editor, args, env); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func main() {
  debug = os.Getenv("DEBUG") != ""
  if len(os.Args) < 2 {
    fmt.Println("Usage: ff <pattern>")
    os.Exit(1)
  }
  files := fileMatches(os.Args[1])
  if len(files) == 0 {
    fmt.Fprintln(os.Stderr, "no matches")
    os.Exit(1)
  }
  file := files[0]
  if len(files) > 1 {
    file = chooseFile(files)
  }
  if file == "" {
    os.Exit(1)
  }
  openFile(mustGetwd()+"/"+file)
}
