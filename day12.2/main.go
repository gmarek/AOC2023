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

type fingerprint struct {
  p string
  c string
  i int
  force_true, force_false bool
}

func getFingerprint(p sequence, c []int, i int, force_true, force_false bool) fingerprint {
  s := []string{}
  for _, v := range c {
    s = append(s, strconv.Itoa(v))
  }
  return fingerprint{strings.Join(p, ""), strings.Join(s, ""), i, force_true, force_false}
}

func options(p_in sequence, con_in []int, i int, force_true bool, force_false bool, cache map[fingerprint]int) int {
  fgp := getFingerprint(p_in, con_in, i, force_true, force_false)
  if v, ok := cache[fgp]; ok {
    return v
  }
  ret := 0
  defer func() {cache[fgp] = ret}()
  p := make(sequence, len(p_in), len(p_in))
  con := make([]int, len(con_in), len(con_in))
  copy(p, p_in)
  copy(con, con_in)
  if debug {
    fmt.Printf("%v %v %v %v %v\n", p, con, i, force_true, force_false)
  }
  if i > len(p) - 1 {
    if len(con) == 0 {
      ret = 1
      return ret
    }
    ret = 0
    return ret
  }
  if len(con) == 0 {
    for ; i < len(p); i++ {
      if p[i] == "#" {
	ret = 0
	return ret
      }
    }
    ret = 1
    return ret
  }
  if force_true  {
    if p[i] == "." {
      ret = 0
      return ret
    }
    p[i] = "#"
    con[0]--
    if con[0] == 0 {
      con = con[1:]
      ret = options(p, con, i+1, false, true, cache)
      return ret
    }
    ret = options(p, con, i+1, true, false, cache)
    return ret
  }
  if force_false {
    if p[i] == "#" {
      ret = 0
      return ret
    }
    ret = options(p, con, i+1, false, false, cache)
    return ret
  }
  if p[i] == "." {
    ret = options(p, con, i+1, false, false, cache)
    return ret
  }
  if p[i] == "#" {
    con[0]--
    if con[0] == 0 {
      con = con[1:]
      ret = options(p, con, i+1, false, true, cache)
      return ret
    }
    p[i] = "#"
    ret = options(p, con, i+1, true, false, cache)
    return ret
  }
  no := options(p, con, i+1, false, false, cache)
  yes := 0
  con[0]--
  if con[0] == 0 {
    con = con[1:]
    yes = options(p, con, i+1, false, true, cache)
  } else {
    yes = options(p, con, i+1, true, false, cache)
  }
  ret = yes + no
  return ret
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
  count := 0
  for scan.Scan() {
    fmt.Printf("%v\n", count)
    count++
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
    lenP := len(pattern)
    lenC := len(constraints)
    for i := 0; i < 4; i++ {
      pattern = append(pattern, "?")
      for j := 0; j < lenP; j++ {
	pattern = append(pattern, pattern[j])
      }
      for j:= 0; j < lenC; j++ {
	constraints = append(constraints, constraints[j])
      }
    }

    cache := map[fingerprint]int{}
    opts := options(pattern, constraints, 0, false, false, cache)
    fmt.Printf("%v\n", opts)
    sum += opts
  }

  fmt.Printf("%v\n", sum)
}
