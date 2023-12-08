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

type cycle struct {
  start string
  period int64
  offset int64
  ends []int64
  current int64
}

func gcd(a, b int64) int64 {
  for b != 0 {
    t := b
    b = a % b
    a = t
  }
  return a
}

func lcm(a, b int64) int64 {
  result := a * b / gcd(a, b)
  return result
}

func mergeCycles(c1, c2 cycle) cycle {
  c := cycle{
    period: lcm(c1.period, c2.period),
  }
  min := c1
  max := c2

  fmt.Printf("%v\n%v\n", c1, c2)
  if c2.period < c1.period {
    tmp := c1
    c1 = c2
    c2 = tmp
  }
  current := min.ends[0]
  step := min.period * ((max.period / min.period) + 1)
  visited := map[int64]bool{}
  for current % max.period != max.ends[0] {
    //fmt.Printf("%v\n", current % max.period)
    if visited[current % max.period] {
      panic("cycle!")
    }
    visited[current % max.period] = true
    current += step
  }

  c.ends = append(c.ends, current)
  return c
}

func main() {
  flag.Parse()

  filename := "test4"
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
  starts := []string{}
  ends := map[string]bool{}
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

    if string(name[2]) == "A" {
      starts = append(starts, name)
    }
    if string(name[2]) == "Z" {
      ends[name] = true
    }
    graph[name] = vertex{name, left, right}
  }


  cycles := []cycle{}
  for _, start := range(starts) {
    tracking := map[string]map[int]int{}
    for n := range(graph) {
      tracking[n] = map[int]int{}
    }

    turn := 0
    count := 0
    current := start
    tracking[start][0] = 0

    for {
      if turns[turn] == "L" {
	current = graph[current].left
      } else {
	current = graph[current].right
      }
      count += 1
      turn = (turn + 1) % len(turns)
      if v, ok := tracking[current][turn]; ok {
	c := cycle{}
	c.start = current
	c.period = int64(count - v + 1)
	c.offset = int64(turn)
	cycles = append(cycles, c)
	break
      }
      tracking[current][turn] = count
    }
  }

  for i, c := range(cycles) {
    turn := c.offset
    current := c.start
    ends := []int64{}
    for j := int64(0); j < c.period; j++ {
      if turns[turn] == "L" {
	current = graph[current].left
      } else {
	current = graph[current].right
      }
      turn = (turn + 1) % int64(len(turns))
      if string(current[2]) == "Z" {
	ends = append(ends, int64(j+1) + c.offset)
      }
    }
    cycles[i].ends = ends
  }

  for len(cycles) > 1 {
    c := cycle{}
    c.period = lcm(cycles[0].period, cycles[1].period)
    cycles = cycles[2:]
    cycles = append(cycles, c)
  }

  fmt.Printf("%v\n", cycles[0].period - 1)
}
