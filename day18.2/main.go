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
var details = false

const e = 0
const s = 1
const w = 2
const n = 3

const left = -1
const right = 1

type dig struct {
  dir int
  l int
  orientation int
}

func printMap(s []point) {
  if !details {
    return
  }
  for _, v := range s {
    fmt.Printf("(%v %v) ", v.x, v.y)
  }
  fmt.Println()
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
  from int
  inTurn bool
}

func abs(a int) int {
  if a < 0 {
    return -a
  }
  return a
}

func orientation(prev, cur int) int {
  turnMap := map[int]point{
    n: point{0, -1, 0, true},
    s: point{0, 1, 0, true},
    e: point{1, 0, 0, true},
    w: point{-1, 0, 0, true},
  }

  if turnMap[prev].x * turnMap[cur].y - turnMap[prev].y * turnMap[cur].x < 0 {
    return right
  }
  return left
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
  for scan.Scan() {
    d := dig{}
    dString := ""
    /*
    fmt.Sscanf(scan.Text(), "%s %d", &dString, &d.l)
    switch dString {
    case "R": d.dir = e
    case "D": d.dir = s
    case "L": d.dir = w
    case "U": d.dir = n
    }
    */
    sTrash := ""
    iTrash := 0
    fmt.Sscanf(scan.Text(), "%s %d %s", &sTrash, &iTrash, &dString)
    lString := dString[2:len(dString) - 2]
    fmt.Sscanf(lString, "%x", &d.l)
    dirStr := dString[len(dString)-2:]
    fmt.Sscanf(dirStr, "%d", &d.dir)
    digs = append(digs, d)
  }

  countLeft := 0
  for i, d := range digs {
    digs[i].orientation = orientation(d.dir, digs[(i+1)%len(digs)].dir)
    if digs[i].orientation == left {
      countLeft++
    } else {
      countLeft--
    }
  }

  globalTurn := 0
  if countLeft == 4 {
    globalTurn = left
  } else if countLeft == -4 {
    globalTurn = right
  } else {
    panic(fmt.Sprintf("countLeft: %v\n", countLeft))
  }

  shape := []point{{0, 0, digs[len(digs)-1].dir, digs[len(digs)-1].orientation == globalTurn}}
  for i, d := range digs {
    pInd := i-1
    if pInd < 0 {
      pInd = len(digs)-1
    }
    offset := 0
    if digs[pInd].orientation == d.orientation {
      if d.orientation == globalTurn {
	offset = 1
      } else {
	offset = -1
      }
    }
    next := point{shape[len(shape) - 1].x, shape[len(shape) - 1].y, d.dir, d.orientation == globalTurn}
    switch d.dir {
    case n: {
      next.y -= d.l + offset
    }
    case e: {
      next.x += d.l + offset
    }
    case s: {
      next.y += d.l + offset
    }
    case w: {
      next.x -= d.l + offset
    }
    }
    shape = append(shape, next)
  }

  minX := int(1e10)
  minY := int(1e10)
  for _, s := range shape {
    if s.x < minX {
      minX = s.x
    }
    if s.y < minY {
      minY = s.y
    }
  }

  shape[0].from = shape[len(shape)-1].from

  for i := range shape {
    shape[i].x -= minX
    shape[i].y -= minY
  }

  sum := 0
  for i := range shape {
    prev := i - 1
    if prev < 0 {
      prev = len(shape) - 1
    }
    switch shape[i].from {
    case n, s: continue
    case e, w: sum += (shape[i].y) * (shape[i].x - shape[prev].x)
    }
  }

  fmt.Printf("%v\n", abs(sum))
}
