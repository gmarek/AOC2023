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

type vertex struct {
  name, left, right string
}

func main() {
  flag.Parse()

  filename := "test2"
  if *prod {
    filename = "input"
  }
  f, err := os.Open(filename)
  if err != nil {
    panic(err)
  }

  scan := bufio.NewScanner(bufio.NewReader(f))
  scan.Scan()

  turns := []string{}
  for _, c := range scan.Text() {
    turns = append(turns, string(c))
  }
  scan.Scan()

  graph := map[string]vertex{}
  for scan.Scan() {
    name := ""
    left := ""
    right := ""
    slice := strings.Split(scan.Text(), "=")
    name = strings.TrimSpace(slice[0])
    slice = strings.Split(slice[1], ",")
    left = strings.TrimSpace(slice[0])[1:]
    right = strings.TrimSpace(slice[1])
    right = right[:len(right)-1]

    graph[name] = vertex{name, left, right}
  }

  turn := 0
  current := "AAA"
  count := 0
  for current != "ZZZ" {
    count++
    if turns[turn] == "L" {
      current = graph[current].left
    } else {
      current = graph[current].right
    }
    turn = (turn + 1) % len(turns)
  }

  fmt.Printf("%v\n", count)
}
