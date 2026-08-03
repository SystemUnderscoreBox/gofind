# gofind

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](/LICENSE)

A simple, zero-dependency, file search script through folders and subfolders. Written in Go.
Gofind's time complexity follows a linear search: O(n). Meaning that it is recommended to NOT initiate
gofind from a starting point that includes a lot of data, e.g. from root '/'.

Further development of more efficient searching algorithm PoCs will be studied in the near-future.

To get started compile gofind with a go-compiler:
```bash
user@linux:~/path/to/gofind$ go build
user@linux:~/path/to/gofind$ ./gofind 'myfile.go'
```
You can make gofind into an executable:
```bash
user@linux:~/path/to/gofind$ go build && cd
user@linux:~$ mkdir bin
user@linux:~$ cd path/to/gofind
user@linux:~/path/to/gofind$ cp gofind ~/bin && cd
user@linux:~$ gofind '*.mp4' 
```

When calling gofind, write input with suffix (e.g. myfile.go) with citations for accurate findings. Use of citations heavily recommended when searching with regex pattern (*).

(No LLMs were utilized during the production of gofind)
