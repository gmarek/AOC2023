package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

type sequence []int

func dif(in sequence) sequence {
  ret := sequence{}
  for i := 1; i < len(in); i++ {
    ret = append(ret, in[i] - in[i-1])
  }
  return ret
}

func isZero(in sequence) bool {
  for _, v := range in {
    if v != 0 {
      return false
    }
  }
  return true
}

func fullStack(in sequence) []sequence {
  ret := []sequence{}
  ret = append(ret, in)
  cur := in
  for !isZero(cur) {
    cur = dif(cur)
    ret = append(ret, cur)
  }
  fmt.Printf("%v\n", ret)
  return ret
}

func extrapolate(in []sequence) int {
  val := 0
  for i := len(in) - 2; i >= 0; i-- {
    val = in[i][0] - val
  }
  return val
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

  readings := []sequence{}
  for scan.Scan() {
    scn := sequence{}
    line := strings.Split(scan.Text(), " ")
    for _, v := range line {
      read, err := strconv.Atoi(v)
      if err != nil {
	panic("NAN")
      }
      scn = append(scn, read)
    }
    readings = append(readings, scn)
  }

  sum := 0
  for _, r := range(readings) {
    stack := fullStack(r)
    v := extrapolate(stack)
    fmt.Printf("%v\n", v)
    sum += v
  }

  fmt.Printf("%v\n", sum)
}
