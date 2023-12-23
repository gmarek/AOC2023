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

const e = 0
const s = 1
const w = 2
const n = 3

type dig struct {
  dir int
  l int
}

func reflect(s *[]point, start, end int) {
  for start < end {
    tmp := (*s)[start]
    (*s)[start] = (*s)[end]
    (*s)[end] = tmp
    start++
    end--
  }
}

func makeFirst(s *[]point, i int) {
  if i == 0 {
    return
  }
  reflect(s, 0, len(*s) - 1)
  reflect(s, 0, len(*s) - i - 1)
  reflect(s, len(*s) - i, len(*s) - 1)
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

func abs(a int) int {
  if a < 0 {
    return -a
  }
  return a
}

func nextInd(i int, s *[]point) int {
  return (i+1) % len(*s)
}

func prevInd(i int, s *[]point) int {
  if i == 0 {
    return len(*s)-1
  }
  return i-1
}

func cutRectangle(s *[]point) int {
  if debug {
    fmt.Println("-----------------")
    fmt.Printf("in %v: %v\n", len(*s), *s)
  }
  minXInd := []int{}
  minX := int(1e10)
  nearMinX := int(1e10)
  for i, p := range *s {
    if p.x < minX {
      if minX != int(1e10) {
	nearMinX = minX
      }
      minX = p.x
      minXInd = []int{i}
    } else if p.x == minX {
      minXInd = append(minXInd, i)
    }
  }

  for _, i := range minXInd {
    (*s)[i].x = nearMinX
  }

  if (*s)[len(*s)-1].x == minX {
    minXInd = len(*s)-1
  }
  makeFirst(s, prevInd(minXInd, s))
  if debug {
    fmt.Printf("post rotate: %v\n", s)
  }
  if len(*s) == 4 {
    r1 := abs(nearMinX - minX) + 1
    r2 := abs((*s)[0].y - (*s)[3].y) + 1
    if debug {
      fmt.Printf("last: %v * %v = %v\n", r1, r2, r1*r2)
    }
    ret := r1 * r2
    *s = []point{}
    return ret
  }

  r1 := abs((*s)[1].y - (*s)[2].y) + 1
  r2 := abs(minX - min((*s)[0].x, (*s)[3].x))
  fmt.Printf("%v * %v = %v\n", r1, r2, r1 * r2)
  ret := r1 * r2
  newPoint := point{}
  newPoint.x = min((*s)[0].x, (*s)[3].x)
  if (*s)[0].x < (*s)[3].x {
    newPoint.y = (*s)[2].y
  } else {
    newPoint.y = (*s)[1].y
  }

  fmt.Printf("adding: %v\n", newPoint)
  remove(s, 2)
  remove(s, 1)

  fmt.Printf("\n%v\n", *s)
  suffix := make([]point, len(*s)-1, len(*s)-1)
  copy(suffix, (*s)[1:])
  *s = append((*s)[:1], newPoint)
  *s = append(*s, suffix...)
  fmt.Printf("%v\n\n", *s)
  if (*s)[0].x == (*s)[1].x {
    remove(s, 0)
    if (*s)[0].x == (*s)[1].x {
      remove(s, 0)
    }
  } else if (*s)[1].x == (*s)[2].x {
    remove(s, 2)
  }
  return ret
}

func remove(s *[]point, i int) {
  fmt.Printf("remove: %v (%v) from %v\n", i, (*s)[i], *s)
  if i == 0 {
    *s = (*s)[1:]
  } else if i == len(*s) - 1 {
    *s = (*s)[:len(*s)-1]
  } else {
    *s = append((*s)[:i], (*s)[i+1:]...)
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
  for scan.Scan() {
    d := dig{}
    dString := ""
    fmt.Sscanf(scan.Text(), "%s %d", &dString, &d.l)
    switch dString {
    case "R": d.dir = e
    case "D": d.dir = s
    case "L": d.dir = w
    case "U": d.dir = n
    }
    /*sTrash := ""
    iTrash := 0
    fmt.Sscanf(scan.Text(), "%s %d %s", &sTrash, &iTrash, &dString)
    lString := dString[2:len(dString) - 2]
    fmt.Sscanf(lString, "%x", &d.l)
    dirStr := dString[len(dString)-2:]
    fmt.Sscanf(dirStr, "%d", &d.dir)
    */
    digs = append(digs, d)
  }

  shape := []point{{0, 0}}
  for _, d := range digs {
    next := point{shape[len(shape) - 1].x, shape[len(shape) - 1].y}
    switch d.dir {
    case n: {
      next.y -= d.l
    }
    case e: {
      next.x += d.l
    }
    case s: {
      next.y += d.l
    }
    case w: {
      next.x -= d.l
    }
    }
    shape = append(shape, next)
  }

  shape = shape[:len(shape)-1]
  sum := 0
  for len(shape) > 0 {
    sum += cutRectangle(&shape)
  }

  fmt.Printf("%v\n", sum)
}
