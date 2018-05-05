# ff

Find and open a file quickly by sublime text style fuzzy match

## Installation

```
go install github.com/joshwillik/ff
```

## Usage

```
Find File (and open it in $EDITOR)

Usage:
  ff [-p | --print] [--ignore=<dir>]... <pattern>
  ff -h | --help

Options:
  -h --help   Show this screen
  -p --print  Print the path of the file instead of opening it
```

## Planned

- -r to do search from git root if inside a git project
- read ignore paths from .gitignore
