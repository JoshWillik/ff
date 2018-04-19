package main
import (
  "fmt"
  "os"
  "syscall"
  "path/filepath"
)

func main() {
  debug := os.Getenv("DEBUG") != ""
  dir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  if len(os.Args) < 2 {
    fmt.Println("Usage: ff <pattern>")
    os.Exit(1)
  }
  pattern := os.Args[1]
  if debug {
    fmt.Fprintln(os.Stderr, "searching "+dir+"/"+pattern)
  }
  files, err := filepath.Glob(dir+"/"+pattern)
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
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
    fmt.Fprintf(os.Stderr, "running %s %v %s\n", editor, args, env)
  }
  if err := syscall.Exec(editor, args, env); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
