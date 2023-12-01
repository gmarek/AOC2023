package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strconv"
  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

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
    for _, c := range(scan.Text()) {
      if !unicode.IsDigit(c) {
	continue
      }
      v, _ := strconv.Atoi(string(c))
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
