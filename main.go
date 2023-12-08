package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
  "sort"
  "strconv"
  "strings"
  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

func main() {
  flag.Parse()

  filename := "test"
  if *prod {
    filename = "input"
  }
  f, err := os.Open(filename)
  if err != nil {
    panic(err)
  }

  scan := bufio.NewScanner(bufio.NewReader(f))

  for scan.Scan() {
  }
}
