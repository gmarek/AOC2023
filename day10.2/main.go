package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
//  "strconv"
//  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

type point struct {
  x, y int
}

type connection struct {
  xp, xm, yp, ym bool
}

func strToCon(s string) connection {
  switch (s){
  case "-": return connection{xp: true, xm: true}
  case "|": return connection{yp: true, ym: true}
  case "L": return connection{xp: true, ym: true}
  case "F": return connection{xp: true, yp: true}
  case "7": return connection{xm: true, yp: true}
  case "J": return connection{xm: true, ym: true}
  case ".": return connection{}
  case "S": return connection{}
  }
  panic(s)
}

func printConnection(c connection) {
  str := ""
  if c.xp && c.xm {
    str = "-"
  }
  if c.yp && c.ym {
    str = "|"
  }
  if c.xp && c.ym {
    str = "L"
  }
  if c.xp && c.yp {
    str = "F"
  }
  if c.xm && c.yp {
    str = "7"
  }
  if c.xm && c.ym {
    str = "J"
  }
  fmt.Print(str)
}

func connectS(m [][]connection, start_x, start_y int) connection {
  ret := connection{}
  if start_x > 0 {
    fmt.Printf("%v\n", m[start_y][start_x-1])
    if m[start_y][start_x-1].xp {
      ret.xm = true
    }
  }
  if start_y > 0 {
    if m[start_y-1][start_x].yp {
      ret.ym = true
    }
  }
  if start_x < len(m[start_y]) {
    if m[start_y][start_x+1].xm {
      ret.xp = true
    }
  }
  if start_y < len(m) {
    if m[start_y+1][start_x].ym {
      ret.yp = true
    }
  }
  count := 0
  if ret.xp {
    count++
  }
  if ret.xm {
    count++
  }
  if ret.yp {
    count++
  }
  if ret.ym {
    count++
  }
  if count != 2 {
    panic(fmt.Sprintf("%v", ret))
  }
  return ret
}

func printMap(newMap [][]connection, vis map[point]bool) {
  for y := 0; y < len(newMap); y++ {
    for x := 0; x < len(newMap[0]); x++ {
      if newMap[y][x] == (connection{}) && !vis[point{x, y}] {
	fmt.Print(".")
      } else {
	if newMap[y][x] == (connection{}) {
	  fmt.Print("X")
	} else {
	  printConnection(newMap[y][x])
	}
      }
    }
    fmt.Println()
  }
}

func isTurn(c connection) bool {
  return ! ((c.xm && c.xp) || (c.ym && c.yp))
}

func countLeft(m [][]connection, p point) int {
  count := 0
  last_y := 0
  for x := p.x; x >= 0; x-- {
    c := m[p.y][x]
    if c.yp && c.ym {
      count++
      continue
    }
    if c.yp {
      if last_y == 0 {
	last_y = 1
	continue
      }
      if last_y == -1 {
	count++
      }
      last_y = 0
    }
    if c.ym {
      if last_y == 0 {
	last_y = -1
	continue
      }
      if last_y == 1 {
	count++
      }
      last_y = 0
    }
  }
  return count
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

  m := [][]connection{}
  start_x, start_y := 0, 0
  for scan.Scan() {
    row := []connection{}
    for _, c := range(scan.Text()) {
      row = append(row, strToCon(string(c)))
      if string(c) == "S" {
	start_x = len(row) - 1
	start_y = len(m)
      }
    }
    m = append(m, row)
  }
  m[start_y][start_x] = connectS(m, start_x, start_y)

  queue := []point{{start_x, start_y}, {-1, -1}}
  visited := map[point]bool{}
  step := 0
  for len(queue) > 0 {
    pop := queue[0]
    queue = queue[1:]
    if pop == (point{-1, -1}) {
      step++
      queue = append(queue, point{-1, -1})
      continue
    }
    if visited[pop] {
      fmt.Printf("%v\n", step)
      break
    }
    visited[pop] = true
    if m[pop.y][pop.x].xm {
      if !visited[point{pop.x-1, pop.y}] {
	queue = append(queue, point{pop.x-1, pop.y})
      }
    }
    if m[pop.y][pop.x].xp {
      if !visited[point{pop.x+1, pop.y}] {
	queue = append(queue, point{pop.x+1, pop.y})
      }
    }
    if m[pop.y][pop.x].ym {
      if !visited[point{pop.x, pop.y-1}] {
	queue = append(queue, point{pop.x, pop.y-1})
      }
    }
    if m[pop.y][pop.x].yp {
      if !visited[point{pop.x, pop.y+1}] {
	queue = append(queue, point{pop.x, pop.y+1})
      }
    }
  }

  printMap(m, map[point]bool{})
  fmt.Println()

  for y := 0; y < len(m); y++ {
    for x := 0; x < len(m[0]); x++ {
      if !visited[point{x, y}] {
	m[y][x] = connection{}
      }
    }
  }

  printMap(m, map[point]bool{})
  fmt.Println()

  count := 0
  for y := 0; y < len(m); y++ {
    for x := 0; x < len(m[0]); x++ {
      if m[y][x] == (connection{}) {
	if countLeft(m, point{x, y}) % 2 == 1 {
	  count++
	}
      }
    }
  }

  fmt.Printf("%v\n", count)
}
