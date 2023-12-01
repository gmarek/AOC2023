package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strings"
  "strconv"
  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

var mapping = map[string]int{
  "one": 1,
  "two": 2,
  "three": 3,
  "four": 4,
  "five": 5,
  "six": 6,
  "seven": 7,
  "eight": 8,
  "nine": 9,
}

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

  out := []int{}
  for scan.Scan() {
    first := -1
    last := -1
    for i, c := range(scan.Text()) {
      var v int
      if unicode.IsDigit(c) {
	v, _ = strconv.Atoi(string(c))
      } else {
	for k, val := range mapping {
	  if strings.HasPrefix(scan.Text()[i:], k) {
	    v = val
	    break
	  }
	}
      }
      if v == 0 {
	continue
      }
      if first == -1 {
	first = v
      } else {
	last = v
      }
    }
    if last == -1 {
      last = first
    }
    out = append(out, 10*first + last)
  }

  sum := 0
  for _, v := range out {
    sum += v
  }

  fmt.Printf("%v\n", sum)
}
