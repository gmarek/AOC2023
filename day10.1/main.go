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
      return
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

  panic("zonk")
}
