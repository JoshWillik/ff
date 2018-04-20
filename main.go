package main
import (
  "fmt"
  "os"
  "syscall"
  "path/filepath"
  "github.com/renstrom/fuzzysearch/fuzzy"
)

func main() {
  debug := os.Getenv("DEBUG") != ""
  dir, err := os.Getwd()
  dir = dir+"/"
  if err != nil {
    panic(err)
  }
  if len(os.Args) < 2 {
    fmt.Println("Usage: ff <pattern>")
    os.Exit(1)
  }
  pattern := os.Args[1]
  if debug {
    fmt.Fprintln(os.Stderr, "searching "+pattern+" in "+dir)
  }
  all_files, err := filepath.Glob(dir+"/**")
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
  }
  files := make([]string, 0, 20)
  for _, path := range all_files {
    rel_path := path[len(dir):]
    if fuzzy.Match(pattern, rel_path) {
      files = append(files, rel_path)
    }
  }
  if len(files) == 0 {
    fmt.Fprintln(os.Stderr, "no matches")
    os.Exit(1)
  }
  if len(files) > 1 {
    fmt.Println("More than 1 file match")
    for i, file := range files {
      fmt.Printf("%d: %s\n", i, file)
    }
    // TODO josh: allow picking which one eventually
    os.Exit(1)
  }
  file := files[0]
  if debug {
    for i, file := range files {
      fmt.Fprintf(os.Stderr, "%d: %s\n", i, file)
    }
  }
  editor := os.Getenv("EDITOR")
  if editor == "" {
    fmt.Fprintf(os.Stderr, "$EDITOR not configured")
    fmt.Println(file)
  }
  args := []string{file}
  env := os.Environ()
  if debug {
    fmt.Fprintf(os.Stderr, "running %s %v\n", editor, args)
  }
  if err := syscall.Exec(editor, args, env); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
