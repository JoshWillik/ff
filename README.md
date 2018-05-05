# ff

Find and open a file quickly by sublime text style fuzzy match

## Installation

```
go install github.com/joshwillik/ff
```

## Usage

```
Find File (and open it)

Usage:
  ff [-p|--print] [-r|--git-root] [--ignore=<dir>]... <pattern>
  ff -h|--help

Options:
  -h --help      Show this screen
  -p --print     Print the path of the file instead of opening it
  -r --git-root  Search for the file from the project git root
```

## Planned

- read ignore paths from .gitignore
