package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "sort"
  "strings"
  "syscall"
  "github.com/docopt/docopt-go"
  "github.com/manifoldco/promptui"
  "github.com/renstrom/fuzzysearch/fuzzy"
)

var exit_ctrl_c = 130
var debug bool
// TODO josh: read this out of .gitignore if available
var ignorePatterns = []string{".min.js", ".git", "node_modules"}

func files(dir string, customIgnore []string) []string {
  files := make([]string, 0, 30)
  ignores := append(ignorePatterns, customIgnore...)
  // TODO josh: skip directories
  filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
    for _, pattern := range ignores {
      if strings.Contains(path, pattern) {
        if (file.IsDir()) {
          return filepath.SkipDir
        } else {
          return nil
        }
      }
    }
    if (!file.IsDir()) {
      files = append(files, path)
    }
    return nil
  })
  return files
}

func mustBaseDir(gitRoot bool) string {
  dir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  if (gitRoot) {
    output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
    if err != nil {
      fmt.Println()
      fmt.Fprintln(os.Stderr, "Could not find git root")
      os.Exit(1)
    }
    dir = strings.Trim(string(output), "\n")
  }
  return dir+"/"
}

func fileMatches(gitRoot bool, pattern string, ignore []string) []string {
  dir := mustBaseDir(gitRoot)
  if debug {
    fmt.Fprintln(os.Stderr, "searching "+pattern+" in "+dir)
  }
  rel_files := make([]string, 0, 20)
  for _, path := range files(dir, ignore) {
    rel_path := path[len(dir):]
    rel_files = append(rel_files, rel_path)
  }
  out_files := make([]string, 0, 20)
  ranked := fuzzy.RankFind(pattern, rel_files)
  sort.Slice(ranked, func(a, b int) bool {
    return ranked[a].Distance < ranked[b].Distance
  })
  for _, item := range ranked {
    if debug {
      fmt.Printf("distance=%d %s\n", item.Distance, item.Target)
    }
    out_files = append(out_files, item.Target)
  }
  return out_files
}

func chooseFile(files []string) string {
  prompt := promptui.Select{
    Label: "Choose file",
    Items: files,
    Size: 10,
    Searcher: func(input string, index int) bool {
      return fuzzy.Match(input, files[index])
    },
  }
  _, file, err := prompt.Run()
  if err == promptui.ErrInterrupt {
    os.Exit(exit_ctrl_c)
  } else if err != nil {
    fmt.Println(err)
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

type options struct{
  Print bool
  GitRoot bool
  Pattern string
  Ignore []string
}

func parseArgs() options {
  parsed , _ := docopt.ParseDoc(`Find File (and open it)

Usage:
  ff [-p|--print] [-r|--git-root] [--ignore=<dir>]... <pattern>
  ff -h|--help

Options:
  -h --help      Show this screen
  -p --print     Print the path of the file instead of opening it
  -r --git-root  Search for the file from the project git root
  `)
  opt := options{}
  if err := parsed.Bind(&opt); err != nil {
    panic(err)
  }
  return opt
}

func main() {
  args := parseArgs()
  debug = os.Getenv("DEBUG") != ""
  files := fileMatches(args.GitRoot, args.Pattern, args.Ignore)
  if len(files) == 0 {
    fmt.Fprintln(os.Stderr, "no matches")
    os.Exit(1)
  }
  file := files[0]
  if len(files) > 1 {
    if args.Print {
      for _, path := range files {
        fmt.Println(path)
      }
      os.Exit(0)
    } else {
      file = chooseFile(files)
    }
  }
  workDir := mustBaseDir(args.GitRoot)
  if args.Print {
    fmt.Println(workDir+file)
  } else {
    openFile(workDir+file)
  }
}
