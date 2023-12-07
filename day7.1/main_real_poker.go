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
var debug = true

type game struct {
  hand map[string]int
  points int
}

func strength(c string) int {
  if unicode.IsNumber(rune(c[0])) {
    v, _ := strconv.Atoi(c)
    return v
  }
  switch (c) {
  case "T": return 10
  case "J": return 11
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

func five(h map[string]int) (bool, handStr) {
  for k, v := range(h) {
    if v == 5 {
      return true, handStr{five: strength(k)}
    }
  }
  return false, handStr{}
}

func four(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v == 4 {
      str.four = strength(k)
    } else {
      str.one = append(str.one, strength(k))
    }
  }
  if str.four > 0 {
    return true, str
  }
  return false, handStr{}
}

func full(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v == 2 {
      str.two = append(str.two, strength(k))
    } else if v == 3 {
      str.three = strength(k)
    }
  }
  if len(str.two) > 0 && str.three > 0 {
    return true, str
  }
  return false, handStr{}
}

func three(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v == 2 {
      return false, handStr{}
    } else if v == 3 {
      str.three = strength(k)
    } else {
      str.one = append(str.one, strength(k))
    }
  }
  if str.three > 0 {
    return true, str
  }
  return false, handStr{}
}

func twoPairs(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v == 2 {
      str.two = append(str.two, strength(k))
    } else if v > 1 {
      return false, handStr{}
    }
    str.one = append(str.one, strength(k))
  }
  if len(str.two) != 2 {
    return false, handStr{}
  }
  sort.Ints(str.two)
  sort.Ints(str.one)
  return true, str
}

func pair(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v == 2 {
      str.two = append(str.two, strength(k))
    } else if v > 1 {
      return false, handStr{}
    } else {
      str.one = append(str.one, strength(k))
    }
  }
  if len(str.two) != 1 {
    return false, handStr{}
  }
  sort.Ints(str.two)
  sort.Ints(str.one)
  return true, str
}

func non(h map[string]int) (bool, handStr) {
  str := handStr{}
  for k, v := range(h) {
    if v > 1 {
      return false, handStr{}
    }
    str.one = append(str.one, strength(k))
  }
  sort.Ints(str.one)
  return true, str
}

type fullStr struct {
  typ int
  str handStr
}

func getHandStr(h map[string]int) fullStr {
  ok, v := five(h)
  if ok {
    if debug {
      fmt.Printf("five: %v\n", h)
    }
    return fullStr{6, v}
  }
  ok, v = four(h)
  if ok {
    if debug {
      fmt.Printf("four: %v\n", h)
    }
    return fullStr{5, v}
  }
  ok, v = full(h)
  if ok {
    if debug {
      fmt.Printf("full: %v\n", h)
    }
    return fullStr{4, v}
  }
  ok, v = three(h)
  if ok {
    if debug {
      fmt.Printf("three: %v\n", h)
    }
    return fullStr{3, v}
  }
  ok, v = twoPairs(h)
  if ok {
    if debug {
      fmt.Printf("two Pairs: %v\n", h)
    }
    return fullStr{2, v}
  }
  ok, v = pair(h)
  if ok {
    if debug {
      fmt.Printf("pair: %v\n", h)
    }
    return fullStr{1, v}
  }
  ok, v = non(h)
  if ok {
    if debug {
      fmt.Printf("non: %v\n", h)
    }
    return fullStr{0, v}
  }
  return fullStr{}
}

func weaker(lhs, rhs map[string]int) bool {
  if debug {
    fmt.Printf("parsing: %v vs %v\n", lhs, rhs)
  }

  lStr := getHandStr(lhs)
  rStr := getHandStr(rhs)

  if debug {
    fmt.Printf("compare: %v vs %v\n", lStr, rStr)
  }

  if lStr.typ < rStr.typ {
    return true
  }
  if rStr.typ < lStr.typ {
    return false
  }
  if lStr.str.five < rStr.str.five {
    return true
  }
  if rStr.str.five < lStr.str.five {
    return false
  }
  if lStr.str.four < rStr.str.four {
    return true
  }
  if rStr.str.four < lStr.str.four {
    return false
  }
  if lStr.str.three < rStr.str.three {
    return true
  }
  if rStr.str.three < lStr.str.three {
    return false
  }
  for i := range(lStr.str.two) {
    if lStr.str.two[len(lStr.str.two) - 1 - i] < rStr.str.two[len(lStr.str.two) - 1 - i] {
      return true
    }
    if rStr.str.two[len(lStr.str.two) - 1 - i] < lStr.str.two[len(lStr.str.two) - 1 - i] {
      return false
    }
  }
  for i := range(lStr.str.one) {
    if debug {
      fmt.Printf("%v vs %v\n", lStr.str.one, rStr.str.one)
    }
    if lStr.str.one[len(lStr.str.one) - 1 - i] < rStr.str.one[len(lStr.str.one) - 1 - i] {
      return true
    }
    if rStr.str.one[len(lStr.str.one) - 1 - i] < lStr.str.one[len(lStr.str.one) - 1 - i] {
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
    g := game{map[string]int{}, 0}
    for _, c := range tokens[0] {
      g.hand[string(c)]++
    }
    g.points, _ = strconv.Atoi(tokens[1])
    games = append(games, g)
  }

  sort.Slice(games, func(i, j int) bool {
    return weaker(games[i].hand, games[j].hand)
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
