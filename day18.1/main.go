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
var debug = true

const n = 0
const e = 1
const s = 2
const w = 3

type dig struct {
  dir int
  l int
  s string
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

type point struct {
  x, y int
}

func move(p point, d dig) (point, []point) {
  ret := point{p.x, p.y}
  rPath := []point{}
  switch d.dir{
  case n: {
    ret.y -= d.l
    for i := 0; i <= d.l; i++ {
      rPath = append(rPath, point{p.x, p.y-i})
    }
  }
  case s: {
    ret.y += d.l
    for i := 0; i <= d.l; i++ {
      rPath = append(rPath, point{p.x, p.y+i})
    }
  }
  case e: {
    ret.x += d.l
    for i := 0; i <= d.l; i++ {
      rPath = append(rPath, point{p.x+i, p.y})
    }
  }
  case w: {
    ret.x -= d.l
    for i := 0; i <= d.l; i++ {
      rPath = append(rPath, point{p.x-i, p.y})
    }
  }
  }
  return ret, rPath
}

func neigh(m *[][]int, p point) []point {
  c := []point{{p.x-1, p.y}, {p.x+1, p.y}, {p.x, p.y-1}, {p.x, p.y+1}}
  ret := []point{}
  for _, v := range c {
    if v.x >= 0 && v.x < len((*m)[0]) && v.y >= 0 && v.y < len(*m) {
      ret = append(ret, v)
    }
  }
  return ret
}

func flood(m *[][]int) {
  cur := point{0, 0}
  queue := []point{cur}
  visited := map[point]bool{{0,0}: true}
  for len(queue) > 0 {
    pop := queue[0]
    queue = queue[1:]
    (*m)[pop.y][pop.x] = -1
    n := neigh(m, pop)
    for _, v := range n {
      if visited[v] || (*m)[v.y][v.x] == 1 {
	continue
      }
      visited[v] = true
      queue = append(queue, v)
    }
  }
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

  digs := []dig{}
  minX := int(1e10)
  minY := int(1e10)
  cX := 0
  maxX := 0
  cY := 0
  maxY := 0
  for scan.Scan() {
    d := dig{}
    dString := ""
    fmt.Sscanf(scan.Text(), "%s %d %s", &dString, &(d.l), &(d.s))
    switch dString {
    case "R": {
      d.dir = e
      cX += d.l
      maxX = max(maxX, cX)
    }
    case "D": {
      d.dir = s
      cY += d.l
      maxY = max(maxY, cY)
    }
    case "L": {
      d.dir = w
      cX -= d.l
      minX = min(minX, cX)
    }
    case "U": {
      d.dir = n
      cY -= d.l
      minY = min(minY, cY)
    }
    }
    digs = append(digs, d)
  }
  if debug {
    fmt.Printf("minX: %v, minY: %v, maxX: %v, maxY: %v\n", minX, minY, maxX, maxY)
  }

  m := make([][]int, maxY - minY + 3, maxY - minY + 3)
  for i := range m {
    m[i] = make([]int, maxX - minX + 3, maxX - minX + 3)
  }
  cur := point{0, 0}
  for _, d := range digs {
    path := []point{}
    cur, path = move(cur, d)
    for _, p := range path {
      m[p.y-minY + 1][p.x-minX + 1] = 1
    }
  }

  flood(&m)
  for i := range m {
    for j := range m[i] {
      if m[i][j] == 0 {
	m[i][j] = 1
      }
    }
  }

  count := 0
  for i := range m {
    for j := range m[i] {
      if debug {
	if m[i][j] == 1 {
	  fmt.Print("#")
	} else if m[i][j] == -1 {
	  fmt.Print(".")
	} else {
	  fmt.Print("-")
	}
      }
      if m[i][j] == 1 {
	count++
      }
    }
    if debug {
      fmt.Println()
    }
  }

  fmt.Printf("%v\n", count)
}
