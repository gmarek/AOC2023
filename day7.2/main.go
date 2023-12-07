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

type game struct {
  hand map[string]int
  orig []string
  points int
}

func strength(c string) int {
  if unicode.IsNumber(rune(c[0])) {
    v, _ := strconv.Atoi(c)
    return v
  }
  switch (c) {
  case "T": return 10
  case "J": return 1
  case "Q": return 12
  case "K": return 13
  case "A": return 14
  }
  panic(fmt.Sprintf("Unknown card: %v\n", c))
}

type handStr struct {
  five int
  four int
  three int
  two []int
  one []int
}

func five(h map[string]int) bool {
  for k, v := range(h) {
    if v == 5 || (k != "J" && v + h["J"] == 5) {
      return true
    }
  }
  return false
}

func four(h map[string]int) bool {
  for k, v := range(h) {
    if v == 4 || (k != "J" && v + h["J"] == 4) {
      return true
    }
  }
  return false
}

func full(h map[string]int) bool {
  if h["J"] == 0 {
    two := false
    three := false
    for _, v := range(h) {
      if v == 2 {
	two = true
      }
      if v == 3 {
	three = true
      }
    }
    if two && three {
      return true
    }
    return false
  } else if h["J"] == 1 {
    twos := 0
    for _, v := range(h) {
      if v == 2 {
	twos++
      }
    }
    if twos == 2 {
      return true
    }
    return false
  }
  return false
}

func three(h map[string]int) bool {
  for k, v := range(h) {
    if v == 3 || (k != "J" && v + h["J"] == 3) {
      return true
    }
  }
  return false
}

func twoPairs(h map[string]int) bool {
  if h["J"] == 0 {
    twos := 0
    for _, v := range(h) {
      if v == 2 {
	twos++
      }
    }
    if twos == 2 {
      return true
    }
  } else if h["J"] == 1 {
    twos := 0
    for _, v := range(h) {
      if v == 2 {
	twos++
      }
    }
    if twos == 1 {
      return true
    }
  }
  return false
}

func pair(h map[string]int) bool {
  for k, v := range(h) {
    if v == 2 || (k != "J" && v + h["J"] == 2) {
      return true
    }
  }
  return false
}

func non(h map[string]int) bool {
  str := handStr{}
  for k, v := range(h) {
    if v > 1 || k == "J"{
      panic("shouldn't happen")
    }
    str.one = append(str.one, strength(k))
  }
  return true
}

func getHandStr(h map[string]int) int {
  ok := five(h)
  if ok {
    if debug {
      fmt.Printf("five: %v\n", h)
    }
    return 6
  }
  ok = four(h)
  if ok {
    if debug {
      fmt.Printf("four: %v\n", h)
    }
    return 5
  }
  ok = full(h)
  if ok {
    if debug {
      fmt.Printf("full: %v\n", h)
    }
    return 4
  }
  ok = three(h)
  if ok {
    if debug {
      fmt.Printf("three: %v\n", h)
    }
    return 3
  }
  ok = twoPairs(h)
  if ok {
    if debug {
      fmt.Printf("two Pairs: %v\n", h)
    }
    return 2
  }
  ok = pair(h)
  if ok {
    if debug {
      fmt.Printf("pair: %v\n", h)
    }
    return 1
  }
  ok = non(h)
  if ok {
    if debug {
      fmt.Printf("non: %v\n", h)
    }
    return 0
  }
  return -1
}

func weaker(lhs, rhs game) bool {
  if debug {
    fmt.Printf("parsing: %v vs %v\n", lhs, rhs)
  }

  lStr := getHandStr(lhs.hand)
  rStr := getHandStr(rhs.hand)

  if lStr < rStr {
    return true
  }
  if rStr < lStr {
    return false
  }
  for i := range(lhs.orig) {
    if strength(lhs.orig[i]) < strength(rhs.orig[i]) {
      return true
    }
    if strength(rhs.orig[i]) < strength(lhs.orig[i]) {
      return false
    }
  }
  return false
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

  games := []game{}
  for scan.Scan() {
    tokens := strings.Split(scan.Text(), " ")
    g := game{map[string]int{}, []string{}, 0}
    for _, c := range tokens[0] {
      g.hand[string(c)]++
      g.orig = append(g.orig, string(c))
    }
    g.points, _ = strconv.Atoi(tokens[1])
    games = append(games, g)
  }

  sort.Slice(games, func(i, j int) bool {
    return weaker(games[i], games[j])
  })

  for _, g := range games {
    fmt.Printf("%v\n", g)
  }

  sum := int64(0)
  for i, v := range(games) {
    sum += int64(i+1) * int64(v.points)
  }

  fmt.Printf("%v\n", sum)
}
