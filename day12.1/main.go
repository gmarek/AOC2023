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

type sequence []string

func options(p_in sequence, con_in []int, i int, force_true bool, force_false bool) int {
  p := make(sequence, len(p_in), len(p_in))
  con := make([]int, len(con_in), len(con_in))
  copy(p, p_in)
  copy(con, con_in)
  if debug {
    fmt.Printf("%v %v %v %v %v\n", p, con, i, force_true, force_false)
  }
  if i > len(p) - 1 {
    if len(con) == 0 {
      return 1
    }
    return 0
  }
  if len(con) == 0 {
    for ; i < len(p); i++ {
      if p[i] == "#" {
	return 0
      }
    }
    return 1
  }
  if force_true  {
    if p[i] == "." {
      return 0
    }
    p[i] = "#"
    con[0]--
    if con[0] == 0 {
      con = con[1:]
      return options(p, con, i+1, false, true)
    }
    return options(p, con, i+1, true, false)
  }
  if force_false {
    if p[i] == "#" {
      return 0
    }
    return options(p, con, i+1, false, false)
  }
  if p[i] == "." {
    return options(p, con, i+1, false, false)
  }
  if p[i] == "#" {
    con[0]--
    if con[0] == 0 {
      con = con[1:]
      return options(p, con, i+1, false, true)
    }
    p[i] = "#"
    return options(p, con, i+1, true, false)
  }
  no := options(p, con, i+1, false, false)
  yes := 0
  con[0]--
  if con[0] == 0 {
    con = con[1:]
    yes = options(p, con, i+1, false, true)
  } else {
    yes = options(p, con, i+1, true, false)
  }
  return yes + no
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

  sum := 0
  for scan.Scan() {
    pattern := sequence{}
    constraints := []int{}
    patternString := ""
    constraintsString := ""
    fmt.Sscanf(scan.Text(), "%s %s", &patternString, &constraintsString)
    for _, c := range patternString {
      pattern = append(pattern, string(c))
    }

    constraintsSlice := strings.Split(constraintsString, ",")
    for _, s := range constraintsSlice {
      v, _ := strconv.Atoi(s)
      constraints = append(constraints, v)
    }

    if debug {
      fmt.Printf("%v\n%v\n\n", pattern, constraints)
    }
    opts := options(pattern, constraints, 0, false, false)
    if debug {
      fmt.Printf("%v\n", opts)
    }
    sum += opts
  }

  fmt.Printf("%v\n", sum)
}
