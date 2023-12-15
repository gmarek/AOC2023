package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
//  "strconv"
  "strings"
//  "unicode"
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

  tokens := []string{}
  sum := 0
  for scan.Scan() {
    tokens = strings.Split(scan.Text(), ",")
    for _, t := range tokens {
      val := 0
      for _, c := range t {
	val += int(c)
	val *= 17
	val = val % 256
      }
      sum += val
    }
  }

  fmt.Printf("%v\n", sum)
}
